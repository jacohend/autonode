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
	State         map[ulid.ULID]EventStateMachine
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
		State:  make(map[ulid.ULID]EventStateMachine),
		Events: queue.NewQueue(),
	}
}

func (processor *Processor) NewEvent(event types.Event, dispatching bool) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	id := util.BytesToUlid(event.Id)
	if _, exists := processor.State[id]; !exists {
		processor.Events.PushItem(event)
		processor.State[id] = EventStateMachine{
			Dispatcher: dispatching,
			Event:      &event,
			Ack:        nil,
			Result:     nil,
		}
	}
}

func (processor *Processor) AcknowledgeEvent(ack types.Ack) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	id := util.BytesToUlid(ack.EventId)
	if s, exists := processor.State[id]; exists {
		if !s.Dispatcher {
			delete(processor.State, id)
			processor.Events.RemoveItemById(ack.EventId)
		} else {
			s.Ack = &ack
		}
	}
}

func (processor *Processor) AddResult(result types.Result) {
	processor.Lock.Lock()
	defer processor.Lock.Unlock()
	id := util.BytesToUlid(result.EventId)
	if s, exists := processor.State[id]; exists {
		s.Result = &result
	}
}

func (processor *Processor) WaitForResult(idbytes []byte) *types.Result {
	defer processor.CompleteEvent(idbytes)
	id := util.BytesToUlid(idbytes)
	t := time.Now()
	timeout := t.Add(10 * time.Second)
	for !t.After(timeout) {
		if s, exists := processor.State[id]; exists && s.Result != nil {
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
	id := util.BytesToUlid(idbytes)
	if _, exists := processor.State[id]; exists {
		delete(processor.State, id)
		processor.Events.RemoveItemById(idbytes)
	}
}

func (processor *Processor) Start() {
	for {
		event, err := processor.Events.PopItem()
		fmt.Println("Received Event")

		if util.LogError(err) != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		go processor.Process(event.(types.Event))
	}
}

func (processor *Processor) Prune() {

}
