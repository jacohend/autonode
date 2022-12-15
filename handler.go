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
		server.Events.NewEvent(m, false)
	case types.Ack:
		server.Events.AcknowledgeEvent(m)
	case types.Result:
		server.Events.ResultHandler(m)
		server.Events.AddResult(m)
	}

	return nil
}
