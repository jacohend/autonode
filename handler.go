package autonode

import (
	"github.com/jacohend/autonode/types"
	"github.com/perlin-network/noise"
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
		//TODO: send this message when we begin queue processing this item
		/*		go ctx.SendMessage(types.Ack{
				NodeId:    server.Node.ID().String(),
				EventId:   m.Id,
				Key:       "",
				Value:     nil,
				Timestamp: util.Now(),
			})*/
	case types.Ack:
		server.Events.RemoveItemById(m.EventId)
	case types.Result:

	default:
		server.SendToNetwork(m)
	}
	return nil
}
