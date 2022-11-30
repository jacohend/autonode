package autonode

import (
	"context"
	"fmt"
	"github.com/jacohend/autonode/types"
	"github.com/perlin-network/noise"
	"time"
)

func (server *ServerNode) Handle(ctx noise.HandlerContext) error {
	if !ctx.IsRequest() {
		return nil
	}

	obj, err := ctx.DecodeMessage()
	if err != nil {
		return err
	}

	switch m := obj.(type) {
	case types.Announce:

	case types.Time:

	case types.Event:
		server.Events.PushItem(m)
	case types.Ack:

	default:
		for _, id := range server.Overlay.Table().Peers() {
			go server.SendToID(id, obj)
		}
	}
	return nil
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
