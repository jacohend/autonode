package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/queue"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/oklog/ulid/v2"
	"sync"
	"time"
)

type Processor struct {
	State         map[string]*EventStateMachine
	Events        *queue.Queue
	Process       func(event types.Event)
	EventHandler  func(event types.Event) (types.Result, error)
	ResultHandler func(event types.Result) error
	Standalone    bool
	Lock          sync.Mutex //required to protect queue during removal
}

type EventStateMachine struct {
	Dispatcher bool //is this node the event's creator?
	Event      *types.Event
	Ack        *types.Ack
	Result     *types.Result
	ResultSub  chan types.Result
}

func NewEventProcessor() *Processor {
	return &Processor{
		State:      make(map[string]*EventStateMachine),
		Events:     queue.NewQueue(),
		Standalone: true,
	}
}

func (processor *Processor) NewEvent(event types.Event, dispatching bool) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, _, exists := processor.GetEvent(event.Id); !exists {
		fmt.Printf("Creating new event %s\n", id.String())
		processor.State[id.String()] = &EventStateMachine{
			Dispatcher: dispatching,
			Event:      &event,
			Ack:        nil,
			Result:     nil,
			ResultSub:  make(chan types.Result),
		}
		processor.Events.PushItem(event)
	}
}

func (processor *Processor) AcknowledgeEvent(ack types.Ack) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, s, exists := processor.GetEvent(ack.EventId); exists {
		if !s.Dispatcher {
			fmt.Printf("Deleting event id %s\n", id.String())
			delete(processor.State, id.String())
			go processor.Events.RemoveItemById(ack.EventId)
		} else {
			fmt.Printf("Storing ack for event id %s\n", id.String())
			s.Ack = &ack
		}
	}
}

func (processor *Processor) AddResult(result types.Result) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, _, exists := processor.GetEvent(result.EventId); exists {
		processor.State[id.String()].Result = &result
		processor.State[id.String()].ResultSub <- result
		fmt.Printf("AddResult %s: %#v\n", util.BytesToUlid(result.EventId), result)
	}
}

func (processor *Processor) WaitForResult(idbytes []byte) *types.Result {
	defer processor.CompleteEvent(idbytes)
	if _, s, exists := processor.GetEvent(idbytes); exists {
		select {
		case result := <-s.ResultSub:
			fmt.Println("Found Result")
			return &result
		case <-time.After(10 * time.Second):
			fmt.Println("Timed out waiting for Result")
			return nil
		}
	}
	return nil
}

func (processor *Processor) CompleteEvent(idbytes []byte) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	fmt.Println("Completing Event")
	if id, _, exists := processor.GetEvent(idbytes); exists {
		fmt.Println("Deleting event...")
		delete(processor.State, id.String())
		fmt.Println("Deletion finished")
	}
}

func (processor *Processor) GetEvent(idbytes []byte) (ulid.ULID, *EventStateMachine, bool) {
	id := util.BytesToUlid(idbytes)
	if _, exists := processor.State[id.String()]; exists {
		return id, processor.State[id.String()], true
	}
	return id, nil, false
}

func (processor *Processor) Start() {
	for {
		event, err := processor.Events.PopItem()
		fmt.Printf("Processing Event: %v\n", event)

		_, s, exists := processor.GetEvent(event.Id)

		if exists && (s.Dispatcher && !processor.Standalone) {
			fmt.Printf("We're the dispatcher and we have workers; skipping self-assignment\n")
			processor.Events.PushItem(event)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if util.LogError(err) != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		go processor.Process(event)
	}
}

func (processor *Processor) Prune() {

}
