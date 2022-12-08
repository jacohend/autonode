protos:
	protoc -I=. --go_out=. types/types.proto

mod:
	go mod tidy && go mod vendor

build:
	go build -o autonode github.com/jacohend/autonode/example
	sudo setcap 'cap_net_bind_service=+ep' autonode

pull:
	git pull

update: pull build