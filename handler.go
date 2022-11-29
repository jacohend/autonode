package autonode

import (
	"github.com/jacohend/autonode/types"
	"github.com/perlin-network/noise"
)

func handle(ctx noise.HandlerContext) error {
	if ctx.IsRequest() {
		return nil
	}

	obj, err := ctx.DecodeMessage()
	if err != nil {
		return err
	}

	switch obj.(type) {
	case types.Announce:

	case types.Time:

	case types.Event:

	case types.Ack:
	}
	return nil
}
