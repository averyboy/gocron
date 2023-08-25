package container

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type Item struct {
	Callback func()
	Priority int64
	Index    int
}

func (i Item) Less(j any) bool {
	return i.Priority < j.(Item).Priority
}

// DelayQueue is an unbounded blocking queue of *Delayed* elements, in which
// an element can only be taken when its delay has expired. The head of the
// queue is the *Delayed* element whose delay expired furthest in the past.
type DelayQueue struct {
	C chan interface{}

	mu sync.Mutex
	pq PriorityQueue

	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan struct{}
}

// New creates an instance of delayQueue with the specified size.
func New(size int64) *DelayQueue {
	return &DelayQueue{
		C:       make(chan interface{}),
		pq:      *NewPriorityqueue(size),
		wakeupC: make(chan struct{}),
	}
}

// Offer inserts the element into the current queue.
func (dq *DelayQueue) Offer(f func(), expiration int64) {
	item := &Item{Callback: f, Priority: expiration}

	dq.mu.Lock()
	heap.Push(&dq.pq, item)
	index := item.Index
	dq.mu.Unlock()

	if index == 0 {
		// A new item with the earliest expiration is added.
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

func (dq *DelayQueue) PeekAndShift(now int64) (*Item, int64) {
	if dq.pq.Len() == 0 {
		return nil, 0
	}
	item := dq.pq.Topx()
	if item.(Item).Priority > now {
		return nil, item.(*Item).Priority - now
	}
	dq.pq.Popx()
	i := item.(Item)
	return &i, 0
}

// Poll starts an infinite loop, in which it continually waits for an element
// to expire and then send the expired element to the channel C.
func (dq *DelayQueue) Poll(exitC chan struct{}, nowF func() int64) {
	for {
		now := nowF()

		dq.mu.Lock()
		item, delta := dq.PeekAndShift(now)
		if item == nil {
			// No items left or at least one item is pending.

			// We must ensure the atomicity of the whole operation, which is
			// composed of the above PeekAndShift and the following StoreInt32,
			// to avoid possible race conditions between Offer and Poll.
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item != nil {
			if delta == 0 {
				// No items left.
				select {
				case <-dq.wakeupC:
					// Wait until a new item is added.
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 {
				// At least one item is pending.
				select {
				case <-dq.wakeupC:
					// A new item with an "earlier" expiration than the current "earliest" one is added.
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// The current "earliest" item expires.

					// Reset the sleeping state since there's no need to receive from wakeupC.
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		select {
		case dq.C <- item.Priority:
			// The expired element has been sent out successfully.
		case <-exitC:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
}