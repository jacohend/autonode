package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/perlin-network/noise"
	"reflect"
)

func (server *ServerNode) Handle(ctx noise.HandlerContext) error {
	obj, err := ctx.DecodeMessage()
	if util.LogError(err) != nil {
		return err
	}
	fmt.Printf("Received %s Msg\n", reflect.TypeOf(obj))

	switch m := obj.(type) {
	case types.Event:
		util.LogAndForget(server.Events.PushItem(m))
	case types.Ack:
		fmt.Println("Received Ack Msg")
		util.LogAndForget(server.Events.RemoveItemById(m.EventId))
	case types.Result:
		fmt.Println("Received Result Msg")
		util.LogAndForget(server.ResultHandler(m))
	}

	return nil
}
