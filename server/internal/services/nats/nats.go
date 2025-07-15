package nats

import (
	"log"
	"smart-fridge/internal/util"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type NATSService struct {
	server *server.Server
	Conn   *nats.Conn
	host   string
	port   int
}

func NewNATSService(host string, port int) *NATSService {
	return &NATSService{
		host: host,
		port: port,
	}
}

func (n *NATSService) Start() error {
	server, err := util.StartEmbeddedNATSServer(n.host, n.port)
	if err != nil {
		return err
	}
	n.server = server
	log.Printf("NATS server started on %s:%d", n.host, n.port)
	return nil
}

// stop gracefully shuts down NATS server and connection
func (n *NATSService) Stop() {
	if n.server != nil && n.server.ReadyForConnections(0) {
		n.server.Shutdown()
		n.server.WaitForShutdown()
		log.Println("NATS server stopped")
	}
}
