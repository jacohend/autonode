package util

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func Now() *timestamppb.Timestamp {
	ts := timestamppb.New(time.Now())
	return ts
}
