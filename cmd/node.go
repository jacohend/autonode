package main

import (
	"github.com/jacohend/autonode"
	"github.com/jessevdk/go-flags"
)

func main() {
	config := autonode.Config{}
	flagParser := flags.NewParser(&config, flags.IgnoreUnknown)
	if _, err := flagParser.Parse(); err != nil {
		panic(err)
	}
	server := autonode.NewServerNode(config)
	//server.SetCallback()
	go server.Start()
}
