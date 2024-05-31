package tkeycloak

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const keycloakImage = "quay.io/keycloak/keycloak:latest"

func TestDockerStart(t *testing.T) {
	//slog.SetloSetDefault()SetlogLevel(slog.DEBUG)

	kc := &KeycloakContainer{}
	kc.DockerTimeout = time.Second * 20

	if err := kc.Start(context.Background()); err != nil {
		t.Fatal(fmt.Errorf("failed to start keycloak: %w", err))
	} else {
		defer kc.Stop(context.Background())
	}

	t.Logf("Keycloak container started: %s", kc.id)
}
