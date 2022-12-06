package queue

import (
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/jacohend/autonode/types"
	"sync"
)

type Queue struct {
	Items goconcurrentqueue.Queue
	Lock  sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{Items: goconcurrentqueue.NewFIFO()}
}

func (queue *Queue) PushItem(T any) error {
	return queue.Items.Enqueue(T)
}

func (queue *Queue) PopItem() (any, error) {
	result, err := queue.Items.Dequeue()
	return result.(any), err
}

func (queue *Queue) RemoveItem(T any) error {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	finder := goconcurrentqueue.NewFIFO()
	i := queue.Items.GetLen()
	for i >= 0 {
		item, err := queue.PopItem()
		if err != nil {
			return err
		}
		if item != T {
			finder.Enqueue(item)
		}
		i--
	}
	queue.Items = finder
	return nil
}

func (queue *Queue) RemoveItemById(id string) error {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	finder := goconcurrentqueue.NewFIFO()
	i := queue.Items.GetLen()
	for i >= 0 {
		item, err := queue.PopItem()
		if err != nil {
			return err
		}
		if item.(types.Event).Id != id {
			finder.Enqueue(item)
		}
		i--
	}
	queue.Items = finder
	return nil
}
