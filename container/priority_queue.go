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
type PriorityQueue struct {
	queue []Comparabler
	self  *PriorityQueue
}

// InitPq initial a priority queue instance
func InitPq(comps ...Comparabler) (pq *PriorityQueue) {
	pq = new(PriorityQueue)
	pq.queue = make([]Comparabler, 0, 16)
	pq.queue = append(pq.queue, comps...)
	pq.self = pq
	heap.Init(pq)
	return pq
}

// Len implement heap interface
func (pq *PriorityQueue) Len() int {
	return len(pq.queue)
}

// Less implement heap interface
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.queue[i].Less(pq.queue[j])
}

// Swap implement heap interface
func (pq *PriorityQueue) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
}

// Push implement heap interface
func (pq *PriorityQueue) Push(x interface{}) {
	pq.queue = append(pq.queue, x.(Comparabler))
}

// Pop implement heap interface
func (pq *PriorityQueue) Pop() interface{} {
	n := len(pq.queue)
	x := pq.queue[n-1]
	pq.queue[n-1] = nil
	pq.queue = pq.queue[:n-1]
	return x
}

// Pushx push a element into priority queue
func (pq *PriorityQueue) Pushx(x Comparabler) {
	heap.Push(pq.self, x)
}

// Popx pop a element from priority queue, removes and returns the top element
func (pq *PriorityQueue) Popx() (x Comparabler) {
	x = heap.Pop(pq.self).(Comparabler)
	return x
}

// Top returns the top element but don't remove
func (pq *PriorityQueue) Topx() (x Comparabler) {
	return pq.queue[len(pq.queue)-1]
}
