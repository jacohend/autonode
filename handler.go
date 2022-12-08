package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/types"
	"github.com/perlin-network/noise"
)

func (server *ServerNode) Handle(ctx noise.HandlerContext) error {
	fmt.Printf("Data: %v\n", ctx.Data())

	obj, err := ctx.DecodeMessage()
	if err != nil {
		return err
	}

	switch m := obj.(type) {
	case types.Event:
		server.Events.PushItem(m)
	case types.Ack:
		server.Events.RemoveItemById(m.EventId)
	case types.Result:
		server.ResultHandler(m)
	}
	return nil
}
