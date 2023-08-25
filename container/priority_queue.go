package container

import "container/heap"

// Interface types that can be compared
type Interface interface {
	Less(a any) bool
}

// Comparabler types that can be compared, with base types(todo)
type Comparabler interface {
	Interface
}

// PriorityQueue priority queue based container/heap
type PriorityQueue[T Comparabler] struct {
	queue []T
	self  *PriorityQueue[T]
}

// InitPq initial a priority queue instance
func NewPriorityqueue[T Comparabler](size int64) (pq *PriorityQueue[T]) {
	pq = new(PriorityQueue[T])
	pq.queue = make([]T, 0, size)
	pq.self = pq
	heap.Init(pq)
	return pq
}

// Len implement heap interface
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.queue)
}

// Less implement heap interface
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return pq.queue[i].Less(pq.queue[j])
}

// Swap implement heap interface
func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
}

// Push implement heap interface
func (pq *PriorityQueue[T]) Push(x interface{}) {
	pq.queue = append(pq.queue, x.(T))
}

// Pop implement heap interface
func (pq *PriorityQueue[T]) Pop() interface{} {
	n := len(pq.queue)
	x := pq.queue[n-1]
	// pq.queue[n-1] = nil
	pq.queue = pq.queue[:n-1]
	return x
}

// Pushx push a element into priority queue
func (pq *PriorityQueue[T]) Pushx(x T) {
	heap.Push(pq.self, x)
}

// Popx pop a element from priority queue, removes and returns the top element
func (pq *PriorityQueue[T]) Popx() (x T) {
	x = heap.Pop(pq.self).(T)
	return x
}

// Top returns the top element but don't remove
func (pq *PriorityQueue[T]) Topx() (x T) {
	return pq.queue[len(pq.queue)-1]
}

// Empty
func (pq *PriorityQueue[T]) Empty() bool {
	return len(pq.queue) == 0
}
