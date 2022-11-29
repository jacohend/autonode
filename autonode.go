package autonode

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var config Config
	flagParser := flags.NewParser(&config, flags.IgnoreUnknown)
	if _, err := flagParser.Parse(); err != nil {
		panic(err)
	}
	server := ServerNode{Config: config}
	server.Start()
}

type ServerNode struct {
	Config  Config
	Node    *noise.Node
	Overlay *kademlia.Protocol
}

func (server *ServerNode) Start() {
	node, err := noise.NewNode(noise.WithNodeAddress(server.Config.Host))
	server.Node = node
	check(err)
	defer server.Node.Close()
	server.Node.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}
		ctx.ID()

		fmt.Printf("Got a message: '%s'\n", string(ctx.Data()))

		return ctx.Send([]byte("Hi!"))
	})

	server.Overlay = kademlia.New()
	server.Node.Bind(server.Overlay.Protocol())
	check(server.Node.Listen())
	bootstrap(server.Node, server.Config.Seeds...)
	discover(server.Overlay)
}
