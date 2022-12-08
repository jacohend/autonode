package autonode

import (
	"context"
	"fmt"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"strings"
	"time"
)

func bootstrap(node *noise.Node, addresses ...string) {
	for _, addr := range addresses {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := node.Ping(ctx, addr)
		cancel()

		if err != nil {
			fmt.Printf("Failed to ping bootstrap node (%s). Skipping... [error: %s]\n", addr, err)
			continue
		}
	}
}

// discover uses Kademlia to discover new peers from nodes we already are aware of.
func discover(overlay *kademlia.Protocol) {
	ids := overlay.Discover()

	var str []string
	for _, id := range ids {
		str = append(str, fmt.Sprintf("%s(%s)", id.Address, id.ID.String()[:printedLength]))
	}

	if len(ids) > 0 {
		fmt.Printf("Discovered %d peer(s): [%v]\n", len(ids), strings.Join(str, ", "))
	} else {
		fmt.Printf("Did not discover any peers.\n")
	}
}

func (server *ServerNode) SendToID(id noise.ID, msg noise.Serializable) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := server.Node.SendMessage(ctx, id.Address, msg)
	cancel()
	if err != nil {
		fmt.Printf("Failed to send message to %s(%s). Skipping... [error: %s]\n",
			id.Address,
			id.ID.String()[:printedLength],
			err,
		)
		return err
	}
	return nil
}

func (server *ServerNode) SendToNetworkBytes(id []byte, msg noise.Serializable) {
	sendId, err := noise.UnmarshalID(id)
	if util.LogError(err) != nil {
		return
	}
	go server.SendToID(sendId, msg)
}

func (server *ServerNode) SendToNetwork(msg noise.Serializable) {
	if server.Overlay != nil && server.Overlay.Table() != nil && server.Overlay.Table().Peers() != nil {
		for _, id := range server.Overlay.Table().Peers() {
			go server.SendToID(id, msg)
		}
	}
}

func (server *ServerNode) SendToNetworkSync(msg noise.Serializable) {
	if server.Overlay != nil && server.Overlay.Table() != nil && server.Overlay.Table().Peers() != nil {
		for _, id := range server.Overlay.Table().Peers() {
			server.SendToID(id, msg)
		}
	}
}
