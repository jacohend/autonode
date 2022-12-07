package types

import "github.com/golang/protobuf/proto"

func (a Announce) Marshal() []byte {
	result, _ := proto.Marshal(&a)
	return result
}

func AnnounceUnmarshal(m []byte) (Announce, error) {
	a := Announce{}
	err := proto.Unmarshal(m, &a)
	return a, err
}

func (a Time) Marshal() []byte {
	result, _ := proto.Marshal(&a)
	return result
}

func TimeUnmarshal(m []byte) (Time, error) {
	a := Time{}
	err := proto.Unmarshal(m, &a)
	return a, err
}

func (a Event) Marshal() []byte {
	result, _ := proto.Marshal(&a)
	return result
}

func EventUnmarshal(m []byte) (Event, error) {
	a := Event{}
	err := proto.Unmarshal(m, &a)
	return a, err
}

func (a Ack) Marshal() []byte {
	result, _ := proto.Marshal(&a)
	return result
}

func AckUnmarshal(m []byte) (Ack, error) {
	a := Ack{}
	err := proto.Unmarshal(m, &a)
	return a, err
}

func (a Result) Marshal() []byte {
	result, _ := proto.Marshal(&a)
	return result
}

func ResultUnmarshal(m []byte) (Result, error) {
	r := Result{}
	err := proto.Unmarshal(m, &r)
	return r, err
}
