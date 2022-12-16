package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"net"
	"strconv"
)

type ServerNode struct {
	Config         Config
	Node           *noise.Node
	Overlay        *kademlia.Protocol
	EventProcessor *Processor
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

	server.EventProcessor = NewEventProcessor()
	server.SetEventProcessor(server.ProcessEvent)

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

//Internal method
func (server *ServerNode) SetEventProcessor(handler func(event types.Event)) {
	server.EventProcessor.Process = handler
}

//External method- for devs
func (server *ServerNode) SetEventHandler(handler func(event types.Event) (types.Result, error)) {
	server.EventProcessor.EventHandler = handler
}

//External method- for devs
func (server *ServerNode) SetResultHandler(handler func(event types.Result) error) {
	server.EventProcessor.ResultHandler = handler
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
	server.Overlay.Discover()
	fmt.Println("Server started. Listening to events.")
	server.EventProcessor.Start()
}

func (server *ServerNode) ProcessEvent(event types.Event) {
	numPeers, _ := server.Peers()
	//skip processing this if we dispatched it?
	if _, s, exists := server.EventProcessor.GetEvent(event.Id); exists && (s.Dispatcher && numPeers > 0) {
		fmt.Printf("We're the dispatcher and we have workers; skipping self-assignment\n")
		return
	}
	ack := types.Ack{
		NodeId:    server.Node.ID().Marshal(),
		EventId:   event.Id,
		Timestamp: util.Now(),
	}
	defer server.EventProcessor.AcknowledgeEvent(ack)
	server.SendToNetworkSync(ack)
	result, err := server.EventProcessor.EventHandler(event)
	util.LogAndForget(err)
	server.SendToIdBytes(event.NodeId, result)
}
