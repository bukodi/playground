package tnats

import (
	"fmt"
	"github.com/nats-io/nats-server/v2/server"
	"log/slog"
	"testing"
	"time"
)

const clusterName = "test-cluster"
const basePort = 5000

func TestCluster(t *testing.T) {
	srv1, err := startServer(1)
	if err != nil {
		t.Fatalf("%+v", err)
	} else {
		defer stopServer(srv1)
	}

	srv2, err := startServer(2)
	if err != nil {
		t.Fatalf("%+v", err)
	} else {
		defer stopServer(srv2)
	}

	t.Logf("srv1 peers: %+v", srv1.JetStreamClusterPeers())
	t.Logf("srv2 peers: %+v", srv2.JetStreamClusterPeers())
}

func startServer(srvIdx int) (*server.Server, error) {
	// Start embedded server
	srvName := fmt.Sprintf("srv%d", srvIdx)
	opts := &server.Options{
		ServerName: srvName,
		Port:       basePort + srvIdx,
		Debug:      true,
		RoutesStr:  fmt.Sprintf("nats://localhost:%d", basePort+100+1),
		Cluster: server.ClusterOpts{
			Name: clusterName,
			Port: basePort + 100 + srvIdx,
		},
		JetStream: true,
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("%s start failed: %w", srvName, err)
	}
	ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		return nil, fmt.Errorf("%s start timeout", srvName)
	} else {
		slog.Info("server started", "srvName", ns.Name())
		portInfo := ns.PortsInfo(100 * time.Millisecond)
		slog.Info("ports", "ports", portInfo)
	}
	return ns, nil
}

func stopServer(srv *server.Server) error {
	// Shutdown server
	srv.Shutdown()
	srv.WaitForShutdown()
	slog.Info("server shut down", "srvName", srv.Name())
	return nil
}
