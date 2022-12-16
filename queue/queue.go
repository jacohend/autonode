package queue

import (
	"bytes"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/jacohend/autonode/types"
	"sync"
)

type Queue struct {
	Items goconcurrentqueue.Queue
	Lock  sync.Mutex //required to protect queue during removal
}

func NewQueue() *Queue {
	return &Queue{Items: goconcurrentqueue.NewFIFO()}
}

func (queue *Queue) PushItem(msg types.Event) error {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	return queue.Items.Enqueue(msg)
}

func (queue *Queue) PopItem() (types.Event, error) {
	result, err := queue.Items.DequeueOrWaitForNextElement()
	return result.(types.Event), err
}

func (queue *Queue) RemoveItem(msg types.Event) error {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	finder := goconcurrentqueue.NewFIFO()
	i := queue.Items.GetLen()
	for i >= 0 {
		item, err := queue.Items.Dequeue()
		if err != nil {
			return err
		}
		if bytes.Compare(item.(types.Event).Id, msg.Id) == 0 {
			finder.Enqueue(item)
		}
		i--
	}
	queue.Items = finder
	return nil
}

func (queue *Queue) RemoveItemById(id []byte) error {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	finder := goconcurrentqueue.NewFIFO()
	i := queue.Items.GetLen()
	for i >= 0 {
		item, err := queue.Items.Dequeue()
		if err != nil {
			return err
		}
		if bytes.Compare(item.(types.Event).Id, id) == 0 {
			finder.Enqueue(item)
		}
		i--
	}
	queue.Items = finder
	return nil
}
