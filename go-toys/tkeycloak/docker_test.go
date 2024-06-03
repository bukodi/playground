package tkeycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

const keycloakImage = "quay.io/keycloak/keycloak:latest"

func TestDockerStart(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	kc := &KeycloakContainer{}
	kc.DockerTimeout = time.Second * 20

	if err := kc.Start(context.Background()); err != nil {
		t.Fatal(fmt.Errorf("failed to start keycloak: %w", err))
	}
	t.Logf("Keycloak container started: %s", kc.id)

	token, err := GetAccessToken()
	if err != nil {
		t.Fatal(fmt.Errorf("failed to get access token: %w", err))
	}
	r, err := GetRealms(token)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to get realms: %w", err))
	} else {
		t.Logf("Realms: %s", r)
	}

	if err := kc.Stop(context.Background()); err != nil {
		t.Fatal(fmt.Errorf("failed to stop keycloak: %w", err))
	} else {
		t.Logf("Keycloak container stopped: %s", kc.id)
	}

}

func GetRealms(token string) (any, error) {
	u := "http://0.0.0.0:3999/admin/realms/master"

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, body)
	}

	return body, nil
}

func GetAccessToken() (string, error) {
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "admin-cli")
	data.Set("username", "webadmin")
	data.Set("password", "Passw0rd")

	req, err := http.NewRequest("POST", "http://localhost:3999/realms/master/protocol/openid-connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
