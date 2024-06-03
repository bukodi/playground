package tkeycloak

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type KeycloakContainer struct {
	TcpPort       int
	AdminName     string
	AdminPass     string
	Image         string
	ContainerName string
	DockerTimeout time.Duration
	cli           *client.Client
	id            string
	ImportPath    string
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
	if kc.ContainerName == "" {
		kc.ContainerName = "mrg-test-keycloak-from-go"
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
	if kc.TcpPort == 0 {
		kc.TcpPort = 3999
	}
}

func (kc *KeycloakContainer) tcpPort() string {
	return fmt.Sprintf("%d", kc.TcpPort)
}

func (kc *KeycloakContainer) Start(ctx context.Context) error {
	kc.checkDefaults()
	if cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err != nil {
		return fmt.Errorf("can't connect to docker service: %w", err)
	} else {
		kc.cli = cli
	}

	containerConfig := &container.Config{
		Image: kc.Image,
		Env: []string{
			"KEYCLOAK_ADMIN=" + kc.AdminName,
			"KEYCLOAK_ADMIN_PASSWORD=" + kc.AdminPass,
		},
		Cmd: []string{
			"start-dev",
			"--import-realm",
			"--hostname-port=" + kc.tcpPort(),
			"--http-port=" + kc.tcpPort(),
			"--log-level=INFO",
		},
		ExposedPorts: nat.PortSet{
			nat.Port(kc.tcpPort() + "/tcp"): struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(kc.tcpPort() + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: kc.tcpPort(),
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

	if err := kc.removeExistingContainers(ctx); err != nil {
		return fmt.Errorf("can't remove existing containers: %w", err)
	}

	if resp, err := kc.cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		kc.ContainerName); errdefs.IsConflict(err) {
		slog.Warn("Container already exists, removing it", "id", resp.ID)
	} else if err != nil {
		return fmt.Errorf("container creation failed: %w", err)
	} else {
		kc.id = resp.ID
	}

	if err := kc.cli.ContainerStart(ctx, kc.id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("container start failed: %w", err)
	}

	// Replace "container_id" with your container's ID
	out, err := kc.cli.ContainerLogs(ctx, kc.id, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		panic(err)
	}

	err = scanContainerLog(ctx, out, func(line string) error {
		slog.Debug(line, "type", "out")
		if strings.Contains(line, "[io.quarkus] (main) Keycloak") &&
			strings.Contains(line, "started") {
			return fmt.Errorf("")
		} else {
			return nil
		}
	}, func(line string) error {
		slog.Error(line, "type", "err")
		return fmt.Errorf("ERROR: %s", line)
	})
	if err.Error() == "" {
		slog.Info("KeycloakContainer started", "id", kc.id)
		return nil
	} else {
		return fmt.Errorf("container start failed: %w", err)
	}
}

func scanContainerLog(ctx context.Context, out io.ReadCloser, fnOutLine func(string) error, fnErrLine func(string) error) error {
	if r := recover(); r != nil {
		fmt.Println("main select:", r)
	}

	outPipeR, outPipeW := io.Pipe()
	errPipeR, errPipeW := io.Pipe()
	//defer func() {
	//	outPipeW.Close()
	//	errPipeW.Close()
	//	outPipeR.Close()
	//	errPipeR.Close()
	//}()

	outLinesCh := make(chan string)
	errLinesCh := make(chan string)

	// Copy the output to standard output
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("split docker pipe:", r)
			}
		}()

		_, err := stdcopy.StdCopy(outPipeW, errPipeW, out)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("outPipeR: Recovered from panic:", r)
			}
		}()

		stdScanner := bufio.NewScanner(outPipeR)
		for stdScanner.Scan() {
			if stdScanner.Err() == nil {
				outLinesCh <- stdScanner.Text()
			}
		}
		close(outLinesCh)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("errPipeR: Recovered from panic:", r)
			}
		}()

		errScanner := bufio.NewScanner(errPipeR)
		for errScanner.Scan() {
			if errScanner.Err() == nil {
				errLinesCh <- errScanner.Text()
			}
		}
		close(errLinesCh)
	}()

	for {
		select {
		case outLine, ok := <-outLinesCh:
			if !ok {
				return fmt.Errorf("stdout closed")
			}
			if err := fnOutLine(outLine); err != nil {
				return err
			}
		case errLine, ok := <-errLinesCh:
			if !ok {
				return fmt.Errorf("stderr closed")
			}
			if err := fnErrLine(errLine); err != nil {
				return err
			}
		case <-ctx.Done():
			return fmt.Errorf("context done")
		}
	}
}

func (kc *KeycloakContainer) Stop(ctx context.Context) error {
	var timeoutInSec *int
	if kc.DockerTimeout == 0 {
		timeoutInSec = nil
	} else {
		x := int(kc.DockerTimeout.Seconds())
		timeoutInSec = &x
	}
	if err := kc.cli.ContainerStop(ctx, kc.id, container.StopOptions{Timeout: timeoutInSec}); err != nil {
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

func (kc *KeycloakContainer) removeExistingContainers(ctx context.Context) (retError error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", kc.ContainerName)

	if containers, err := kc.cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filterArgs,
	}); err != nil {
		return fmt.Errorf("can't list existing containers: %w", err)
	} else {
		for _, c := range containers {
			if err = kc.cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{
				Force: true,
			}); err != nil {
				retError = errors.Join(retError, fmt.Errorf("error removing container: %w", err))
			} else {
				slog.Debug("Container removed", "id", c.ID)
			}
		}
		return retError
	}
}
