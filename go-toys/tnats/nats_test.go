package tnats

import (
	"testing"
	"time"

	natssrv "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func TestEmbeddedNats(t *testing.T) {
	// Start embedded server
	opts := &natssrv.Options{
		Debug: true,
	}
	ns, err := natssrv.NewServer(opts)
	if err != nil {
		panic(err)
	}
	ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}
	clientUrl := ns.ClientURL()

	// Start subscriber
	ncSub, _ := nats.Connect(clientUrl)
	defer ncSub.Close()

	sub, _ := ncSub.Subscribe("greet.*", func(msg *nats.Msg) {
		t.Logf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
	})
	defer sub.Drain()

	// Connect as publisher
	ncPub, err := nats.Connect(clientUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer ncPub.Close()
	if err := ncPub.Publish("greet.alice", []byte("hello Alice")); err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 10)
	if err := ncPub.Publish("greet.bob", []byte("hello Bob")); err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 10)

	// Shutdown server
	ns.Shutdown()
	ns.WaitForShutdown()
}
