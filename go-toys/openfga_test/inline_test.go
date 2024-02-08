package openfga_test

import (
	"context"
	"github.com/openfga/go-sdk/client"
	"os"
	"testing"
)

func TestEmbeddedOpenFGA(t *testing.T) {
	t.Logf("test")
	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiScheme:            os.Getenv("FGA_API_SCHEME"), // optional. Can be "http" or "https". Defaults to "https"
		ApiHost:              os.Getenv("FGA_API_HOST"),   // required, define without the scheme (e.g. api.fga.example instead of https://api.fga.example)
		StoreId:              os.Getenv("FGA_STORE_ID"),   // optional, not needed for \`CreateStore\` and \`ListStores\`, required before calling for all other methods
		AuthorizationModelId: os.Getenv("FGA_MODEL_ID"),   // optional, can be overridden per request
	})

	if err != nil {
		// .. Handle error
	}

	resp, err := fgaClient.CreateStore(context.Background()).Body(client.ClientCreateStoreRequest{Name: "FGA Demo"}).Execute()
	if err != nil {
		// .. Handle error
	}
	t.Logf("CreateStore response: %+v", resp)
}
