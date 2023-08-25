package container

import (
	"sync/atomic"
	"time"
)

type Item struct {
	V          interface{}
	Expiration int64
}

func (i Item) Less(j any) bool {
	return i.Expiration < j.(Item).Expiration
}

// DelayQueue is an unbounded blocking queue of *Delayed* elements, in which
// an element can only be taken when its delay has expired. The head of the
// queue is the *Delayed* element whose delay expired furthest in the past.
type DelayQueue struct {
	C  chan interface{}
	pq PriorityQueue
	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan struct{}
	exitC    chan struct{}
}

// New creates an instance of delayQueue with the specified size.
func NewDelayQueue(size int64) *DelayQueue {
	return &DelayQueue{
		C:       make(chan interface{}),
		pq:      *NewPriorityqueue(size),
		wakeupC: make(chan struct{}),
	}
}

func (dq *DelayQueue) PeekAndShift(t int64) (item *Item, delta int64) {
	delta = 0
	item = new(Item)

	if dq.pq.Empty() {
		return nil, 0
	}

	*item = dq.pq.Topx().(Item)
	if item.Expiration < t {
		return item, t - item.Expiration
	}
	dq.pq.Popx()

	return item, 0
}

// Put inserts the element into the current queue.
func (dq *DelayQueue) Put(v interface{}, expiration int64) {
	item := &Item{V: v, Expiration: expiration}

	dq.pq.Pushx(item)

	if dq.pq.Topx().(Item) == *item {
		// A new item with the earliest expiration is added.
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

// TakeNoWait get a element from queue
// if no expired element, waits for an element to expire
func (dq *DelayQueue) Take() (item *Item) {
	for {
		t := time.Now().UnixMilli()
		var delta int64
		if item, delta = dq.PeekAndShift(t); delta == 0 {
			return item
		}
		select {
		case <-time.After(time.Duration(delta) * time.Millisecond):
			return item
		case <-dq.wakeupC:
			continue
		case <-dq.exitC:
			goto exit
		}
	}
exit:
	return
}

// TakeNoWait, Take get a element from queue, if no expired element, return nil
func (dq *DelayQueue) TakeNoWait() (item *Item) {
	t := time.Now().UnixMilli()
	var delta int64
	if item, delta = dq.PeekAndShift(t); delta == 0 {
		return item
	}
	return nil
}
