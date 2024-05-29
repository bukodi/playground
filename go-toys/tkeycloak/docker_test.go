package tkeycloak

import (
	"context"
	"fmt"
	"testing"
)

const keycloakImage = "quay.io/keycloak/keycloak:latest"

func TestDockerStart(t *testing.T) {

	kc := &KeycloakContainer{}

	if err := kc.Start(context.Background()); err != nil {
		t.Fatal(fmt.Errorf("failed to start keycloak: %w", err))
	} else {
		defer kc.Stop(context.Background())
	}

	t.Logf("Keycloak container started: %s", kc.id)
}
