package autonode

import (
	"fmt"
	"github.com/jacohend/autonode/queue"
	"github.com/jacohend/autonode/types"
	"github.com/jacohend/autonode/util"
	"github.com/oklog/ulid/v2"
	"os"
	"sync"
	"time"
)

type Processor struct {
	State         map[ulid.ULID]*EventStateMachine
	Events        *queue.Queue
	Process       func(event types.Event)
	EventHandler  func(event types.Event) (types.Result, error)
	ResultHandler func(event types.Result) error
	Lock          sync.Mutex //required to protect queue during removal
}

type EventStateMachine struct {
	Dispatcher bool //is this node the event's creator?
	Event      *types.Event
	Ack        *types.Ack
	Result     *types.Result
}

func NewEventProcessor() *Processor {
	return &Processor{
		State:  make(map[ulid.ULID]*EventStateMachine),
		Events: queue.NewQueue(),
	}
}

func (processor *Processor) NewEvent(event types.Event, dispatching bool) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, _, exists := processor.GetEvent(event.Id); !exists {
		processor.State[id] = &EventStateMachine{
			Dispatcher: dispatching,
			Event:      &event,
			Ack:        nil,
			Result:     nil,
		}
		processor.Events.PushItem(event)
	}
}

func (processor *Processor) AcknowledgeEvent(ack types.Ack) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, s, exists := processor.GetEvent(ack.EventId); exists {
		if !s.Dispatcher {
			fmt.Printf("Deleting event id %s", id.String())
			delete(processor.State, id)
			processor.Events.RemoveItemById(ack.EventId)
		} else {
			fmt.Printf("Storing ack for event id %s", id.String())
			s.Ack = &ack
		}
	}
}

func (processor *Processor) AddResult(result types.Result) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, _, exists := processor.GetEvent(result.EventId); exists {
		processor.State[id].Result = &result
		fmt.Printf("AddResult %s: %#v", util.BytesToUlid(result.EventId), result)
	}
}

func (processor *Processor) WaitForResult(idbytes []byte) *types.Result {
	defer processor.CompleteEvent(idbytes)
	id := util.BytesToUlid(idbytes)
	t := time.Now()
	timeout := t.Add(10 * time.Second)
	for !t.After(timeout) {
		os.Stdout.Write([]byte(fmt.Sprintf("State: %#v\n", processor.State[id])))
		if _, s, exists := processor.GetEvent(idbytes); exists && s.Result != nil {
			return s.Result
		}
		t = time.Now()
		time.Sleep(time.Millisecond * 60)
	}
	return nil
}

func (processor *Processor) CompleteEvent(idbytes []byte) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	if id, _, exists := processor.GetEvent(idbytes); exists {
		delete(processor.State, id)
		processor.Events.RemoveItemById(idbytes)
	}
}

func (processor *Processor) GetEvent(idbytes []byte) (ulid.ULID, *EventStateMachine, bool) {
	id := util.BytesToUlid(idbytes)
	if _, exists := processor.State[id]; exists {
		return id, processor.State[id], true
	}
	return id, nil, false
}

func (processor *Processor) Start() {
	for {
		event, err := processor.Events.PopItem()
		fmt.Printf("Processing Event: %v\n", event)

		if util.LogError(err) != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		go processor.Process(event)
	}
}

func (processor *Processor) Prune() {

}
