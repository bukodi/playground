package tkeycloak

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type KeycloakContainer struct {
	AdminName          string
	AdminPass          string
	Image              string
	DockerTimeoutInSec *int
	cli                *client.Client
	id                 string
	ImportPath         string
}

func (kc *KeycloakContainer) checkDefaults() {
	if kc.AdminName == "" {
		kc.AdminName = "webadmin"
	}
	if kc.AdminPass == "" {
		kc.AdminPass = "Passw0rd"
	}
	if kc.Image == "" {
		kc.Image = "quay.io/keycloak/keycloak:latest"
		//kc.Image = "quay.io/keycloak/keycloak:24.0.4"
	}
	if kc.ImportPath != "" {
		if !filepath.IsAbs(kc.ImportPath) {
			if abs, err := filepath.Abs(kc.ImportPath); err == nil {
				if fi, err := os.Stat(abs); err == nil && fi.IsDir() {
					kc.ImportPath = abs
				}
			}
		}
	}
}

func (kc *KeycloakContainer) Start(ctx context.Context) error {
	kc.checkDefaults()
	if cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err != nil {
		return fmt.Errorf("can't connect to docker service: %w", err)
	} else {
		kc.cli = cli
	}

	containerConfig := &container.Config{
		Image: "quay.io/keycloak/keycloak:latest",
		Env: []string{
			"KEYCLOAK_ADMIN=webadmin",
			"KEYCLOAK_ADMIN_PASSWORD=Passw0rd",
		},
		Cmd: []string{
			"start-dev",
			"--import-realm",
			"--hostname-port=3999",
			"--http-port=3999",
			"--log-level=INFO",
		},
		ExposedPorts: nat.PortSet{
			"3999/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"3999/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "3999",
				},
			},
		},
		AutoRemove: true,
	}
	if kc.ImportPath != "" {
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: kc.ImportPath,
			Target: "/opt/keycloak/data/import",
		})
	}

	if exists, err := kc.checkImageExists(ctx); err != nil {
		return fmt.Errorf("image check failed: %w", err)
	} else if !exists {
		reader, err := kc.cli.ImagePull(ctx, kc.Image, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("image pull failed: %w", err)
		}
		io.Copy(os.Stdout, reader)
	} else {
		slog.Debug("Image exists", "Image", kc.Image)
	}

	if resp, err := kc.cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		"mrg-test-keycloak-from-go"); err != nil {
		return fmt.Errorf("container creation failed: %w", err)
	} else {
		kc.id = resp.ID
	}

	if err := kc.cli.ContainerStart(ctx, kc.id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("container start failed: %w", err)
	}

	slog.Info("Container started", "id", kc.id)
	return nil
}

func (kc *KeycloakContainer) Stop(ctx context.Context) error {
	if err := kc.cli.ContainerStop(ctx, kc.id, container.StopOptions{Timeout: kc.DockerTimeoutInSec}); err != nil {
		return fmt.Errorf("container stop failed: %w", err)
	}
	slog.Info("Container stopped", "id", kc.id)
	return nil
}

func (kc *KeycloakContainer) checkImageExists(ctx context.Context) (bool, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", kc.Image)

	if images, err := kc.cli.ImageList(ctx, types.ImageListOptions{
		Filters: filterArgs,
	}); err != nil {
		return false, err
	} else if len(images) == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
