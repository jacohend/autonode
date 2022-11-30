package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/queue"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

type ServerNode struct {
	Config  Config
	Node    *noise.Node
	Overlay *kademlia.Protocol
	Events  *queue.Queue
}

func (server *ServerNode) Start() {
	node, err := noise.NewNode(noise.WithNodeAddress(server.Config.Host))
	util.Check(err)

	server.Node = node
	defer server.Node.Close()

	server.Events = queue.NewQueue()
	server.Node.Handle(server.Handle)

	server.Overlay = kademlia.New(kademlia.WithProtocolEvents(kademlia.Events{
		OnPeerAdmitted: func(id noise.ID) {
			fmt.Printf("New peer %s(%s).\n", id.Address, id.ID.String()[:printedLength])
		},
		OnPeerEvicted: func(id noise.ID) {
			fmt.Printf("Removed peer %s(%s).\n", id.Address, id.ID.String()[:printedLength])
		},
	}))
	server.Node.Bind(server.Overlay.Protocol())
	util.Check(server.Node.Listen())
	bootstrap(server.Node, server.Config.Seeds...)
	discover(server.Overlay)
}

func (server *ServerNode) Send(event types.Event) error {
	return nil
}

func (server *ServerNode) SendSync(event types.Event) (types.Ack, error) {
	return types.Ack{}, nil
}
