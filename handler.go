package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
)

func (server *ServerNode) Handle(ctx noise.HandlerContext) error {
	fmt.Printf("Raw Data: %v\n", ctx.Data())

	obj, err := ctx.DecodeMessage()
	if util.LogError(err) != nil {
		return err
	}
	fmt.Printf("Data: %#v\n", obj)

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
