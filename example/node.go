package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jacohend/autonode"
	"github.com/jacohend/autonode/types"
	ulid2 "github.com/jacohend/autonode/ulid"
	"github.com/jacohend/autonode/util"
	"github.com/jessevdk/go-flags"
	"net/http"
	"time"
)

var server *autonode.ServerNode

type Config struct {
	Config autonode.Config `group:"autonode" namespace:"autonode"`
	Host   string          `long:"addr" description:"host/port combo to listen in on"`
}

func main() {
	config := Config{}
	flagParser := flags.NewParser(&config, flags.IgnoreUnknown)
	_, err := flagParser.Parse()
	util.Check(err)

	server = autonode.NewServerNode(config.Config)
	server.SetEventHandler(ApiEventHandler)
	server.SetResultHandler(ApiResultHandler)

	fmt.Println("Starting autonode...")
	go server.Start()
	fmt.Println("Starting api server...")

	StartApi(config)
}

func StartApi(config Config) {
	r := mux.NewRouter()
	r.HandleFunc("/", ApiHandler)
	http.Handle("/", r)
	api := &http.Server{
		Handler:      r,
		Addr:         config.Host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	api.ListenAndServe()
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	ulidgen := ulid2.NewThreadSafeUlid()
	ulid, err := ulidgen.NewUlid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//dispatch an event to a random worker node for processing
	server.DispatchRandom(types.Event{
		NodeId:    server.Node.ID().Marshal(),
		Id:        ulid.Bytes(),
		Key:       "SAMPLE_EVENT",
		Value:     []byte("sample event"),
		Timestamp: util.Now(),
	})
	result := server.EventProcessor.WaitForResult(ulid.Bytes())
	if result != nil {
		resultbytes, _ := json.Marshal(result)
		w.Write(resultbytes)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// this runs on the worker
func ApiEventHandler(event types.Event) (types.Result, error) {
	// do some work here, then dispatch result
	return types.Result{
		NodeId:    server.Node.ID().Marshal(),
		EventId:   event.Id,
		Key:       "SAMPLE_RESULT",
		Value:     []byte("sample result"),
		Timestamp: util.Now(),
	}, nil
}

func ApiResultHandler(result types.Result) error {
	fmt.Printf("Result Handler %s: %#v", util.BytesToUlid(result.EventId), result)
	return nil
}
