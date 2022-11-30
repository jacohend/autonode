package queue

import (
	"github.com/enriquebris/goconcurrentqueue"
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
	queue.Lock.Unlock()
	return nil
}
