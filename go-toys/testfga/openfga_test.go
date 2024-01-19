package testfga

import (
	"context"
	"fmt"
	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	parser "github.com/openfga/language/pkg/go/transformer"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/server/commands"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/storage/memory"
	"github.com/openfga/openfga/pkg/typesystem"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestOpenFGA(t *testing.T) {
	var datastore storage.OpenFGADatastore = memory.New()
	defer datastore.Close()

	WriteAuthorizationModelTest(t, datastore)
	//server.NewMemoryStore()
}

func TestWriteAuthorizationModel(t *testing.T) {
	var datastore storage.OpenFGADatastore = memory.New()
	defer datastore.Close()
	ctx := context.TODO()

	srv, err := server.NewServerWithOpts(server.WithDatastore(datastore))
	require.NoError(t, err)

	storeID := ulid.Make().String()
	request := &openfgav1.WriteAuthorizationModelRequest{
		StoreId: storeID,
		TypeDefinitions: parser.MustTransformDSLToProto(`model
  schema 1.1
type user

type document
  relations
	define viewer: [user]`).TypeDefinitions,
	}

	cmd := commands.NewWriteAuthorizationModelCommand(datastore)
	response, err := cmd.Execute(ctx, request)
	stat, _ := status.FromError(err)
	t.Logf(stat.String())
	t.Logf("Response: %v", response)

	response2, err := srv.Check(ctx, &openfgav1.CheckRequest{})
	stat, _ = status.FromError(err)
	t.Logf(stat.String())
	t.Logf("Response2: %v", response2)
}

func WriteAuthorizationModelTest(t *testing.T, datastore storage.OpenFGADatastore) {
	storeID := ulid.Make().String()

	items := make([]*openfgav1.TypeDefinition, datastore.MaxTypesPerAuthorizationModel()+1)
	items[0] = &openfgav1.TypeDefinition{
		Type: "user",
	}
	for i := 1; i < datastore.MaxTypesPerAuthorizationModel(); i++ {
		items[i] = &openfgav1.TypeDefinition{
			Type: fmt.Sprintf("type%v", i),
			Relations: map[string]*openfgav1.Userset{
				"admin": {Userset: &openfgav1.Userset_This{}},
			},
			Metadata: &openfgav1.Metadata{
				Relations: map[string]*openfgav1.RelationMetadata{
					"admin": {
						DirectlyRelatedUserTypes: []*openfgav1.RelationReference{
							typesystem.DirectRelationReference("user", ""),
						},
					},
				},
			},
		}
	}
	var tests = []struct {
		name          string
		request       *openfgav1.WriteAuthorizationModelRequest
		allowSchema10 bool
		errCode       codes.Code
	}{
		{
			name: "direct_relationship_with_entrypoint",
			request: &openfgav1.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: parser.MustTransformDSLToProto(`model
  schema 1.1
type user

type document
  relations
	define viewer: [user]`).TypeDefinitions,
			},
		},
	}

	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := commands.NewWriteAuthorizationModelCommand(datastore)
			resp, err := cmd.Execute(ctx, test.request)
			status, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, test.errCode, status.Code())

			if err == nil {
				_, err = ulid.Parse(resp.AuthorizationModelId)
				require.NoError(t, err)
			}
		})
	}

}
