package autonode

import (
	"context"
	"fmt"
	leaderelection "github.com/chainpoint/leader-election"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"reflect"
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

func (server *ServerNode) SendToIdBytes(id []byte, msg noise.Serializable) {
	sendId, err := noise.UnmarshalID(id)
	if util.LogError(err) != nil {
		return
	}
	go server.SendToID(sendId, msg)
}

func (server *ServerNode) SendToNetwork(msg noise.Serializable) {
	if server.overlayCheck() {
		for _, id := range server.Overlay.Table().Peers() {
			fmt.Printf("Sending %s to ID %v\n", reflect.TypeOf(msg), id)
			go server.SendToID(id, msg)
		}
	} else {
		fmt.Println("Problem with Overlay- no peers")
	}
}

func (server *ServerNode) SendToNetworkSync(msg noise.Serializable) {
	if server.overlayCheck() {
		for _, id := range server.Overlay.Table().Peers() {
			fmt.Printf("Sending %s to ID %v\n", reflect.TypeOf(msg), id)
			server.SendToID(id, msg)
		}
	} else {
		fmt.Println("Problem with Overlay- no peers")
	}
}

func (server *ServerNode) DispatchRandom(msg types.Event) {
	if server.overlayCheck() {
		server.Events.NewEvent(msg, true)
		peers := server.Overlay.Table().Peers()
		result := leaderelection.ElectLeaders(peers, 1, time.Now().String()).([]noise.ID)
		if len(result) > 0 {
			fmt.Printf("Dispatching Event to %s\n", result[0])
			go server.SendToID(result[0], msg)
		}
	}
}

func (server *ServerNode) overlayCheck() bool {
	return server.Overlay != nil && server.Overlay.Table() != nil && server.Overlay.Table().Peers() != nil
}
