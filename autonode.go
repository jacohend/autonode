package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/queue"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"net"
	"strconv"
	"time"
)

type ServerNode struct {
	Config        Config
	Node          *noise.Node
	Overlay       *kademlia.Protocol
	Events        *queue.Queue
	EventHandler  func(event types.Event) (types.Result, error)
	ResultHandler func(event types.Result) error
}

func NewServerNode(config Config) *ServerNode {
	server := ServerNode{Config: config}
	host, port, _ := net.SplitHostPort(server.Config.Host)
	ip, _ := net.ResolveIPAddr("ip", host)

	portInt, err := strconv.ParseUint(port, 10, 16)
	util.Check(err)

	node, err := noise.NewNode(noise.WithNodeBindHost(ip.IP),
		noise.WithNodeBindPort(uint16(portInt)),
		noise.WithNodeAddress(server.Config.Host))
	util.Check(err)

	server.Node = node
	server.Node.RegisterMessage(types.Announce{}, types.AnnounceUnmarshal)
	server.Node.RegisterMessage(types.Time{}, types.TimeUnmarshal)
	server.Node.RegisterMessage(types.Event{}, types.EventUnmarshal)
	server.Node.RegisterMessage(types.Ack{}, types.AckUnmarshal)
	server.Node.RegisterMessage(types.Result{}, types.ResultUnmarshal)
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

	return &server
}

func (server *ServerNode) SetEventHandler(handler func(event types.Event) (types.Result, error)) {
	server.EventHandler = handler
}

func (server *ServerNode) SetResultHandler(handler func(event types.Result) error) {
	server.ResultHandler = handler
}

func (server *ServerNode) Start() {
	fmt.Println("Starting...")
	defer server.Node.Close()
	fmt.Println("Binding overlay")
	server.Node.Bind(server.Overlay.Protocol())
	fmt.Println("Listening in on specified interface")
	util.Check(server.Node.Listen())
	fmt.Println("Bootstrapping from seeds...")
	bootstrap(server.Node, server.Config.Seeds...)
	fmt.Println("Discovering peers...")
	discover(server.Overlay)
	fmt.Println("Server started. Listening to events.")
	for {
		event, err := server.Events.PopItem()
		fmt.Printf("Received Msg: %v\n", event)
		if util.LogError(err) != nil {
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
