package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/queue"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"time"
)

type ServerNode struct {
	Config        Config `group:"autonode" namespace:"autonode"`
	Node          *noise.Node
	Overlay       *kademlia.Protocol
	Events        *queue.Queue
	EventHandler  func(event types.Event) (types.Result, error)
	ResultHandler func(event types.Result) error
}

func NewServerNode(config Config) *ServerNode {
	server := ServerNode{Config: config}
	node, err := noise.NewNode(noise.WithNodeAddress(server.Config.Host))
	util.Check(err)

	server.Node = node
	server.Node.Handle(server.Handle)

	server.Events = queue.NewQueue()

	server.Overlay = kademlia.New(kademlia.WithProtocolEvents(kademlia.Events{
		OnPeerAdmitted: func(id noise.ID) {
			fmt.Printf("New peer %s(%s).\n", id.Address, id.ID.String()[:printedLength])
		},
		OnPeerEvicted: func(id noise.ID) {
			fmt.Printf("Removed peer %s(%s).\n", id.Address, id.ID.String()[:printedLength])
		},
	}))
	server.Node.Bind(server.Overlay.Protocol())

	return &server
}

func (server *ServerNode) SetEventHandler(handler func(event types.Event) (types.Result, error)) {
	server.EventHandler = handler
}

func (server *ServerNode) SetResultHandler(handler func(event types.Result) error) {
	server.ResultHandler = handler
}

func (server *ServerNode) Start() {
	defer util.LogAndForget(server.Node.Close())
	util.Check(server.Node.Listen())
	bootstrap(server.Node, server.Config.Seeds...)
	discover(server.Overlay)
	for {
		event, err := server.Events.Items.DequeueOrWaitForNextElement()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		go server.ProcessEvent(event.(types.Event))
	}
}

func (server *ServerNode) ProcessEvent(event types.Event) {
	server.SendToNetworkSync(types.Ack{
		NodeId:    server.Node.ID().Marshal(),
		EventId:   event.Id,
		Timestamp: util.Now(),
	})
	result, err := server.EventHandler(event)
	util.LogAndForget(err)
	server.SendToNetworkBytes(event.NodeId, result)
}
