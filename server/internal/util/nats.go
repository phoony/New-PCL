package util

import (
	"fmt"
	"time"

	nats "github.com/nats-io/nats-server/v2/server"
)

func StartEmbeddedNATSServer(host string, port int) (*nats.Server, error) {
	opts := &nats.Options{
		Host: host,
		Port: port,
	}

	s, err := nats.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("error starting NATS server: %w", err)
	}

	go s.Start()

	// wait until the server is ready
	if !s.ReadyForConnections(10 * time.Second) {
		s.Shutdown()
		s.WaitForShutdown()
		return nil, fmt.Errorf("NATS Server not ready in time")
	}

	return s, nil
}
