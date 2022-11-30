protos:
	protoc -I=. --go_out=. types/types.proto

mod:
	go mod tidy && go mod vendor

build:
	go build cmd/main.go
