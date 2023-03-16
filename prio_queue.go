package graph

import (
	"container/heap"
)

type prioQueueImpl struct {
	nodes    []string
	indexMap map[string]int
	less     func(n1, n2 string) bool
}

var _ heap.Interface = &prioQueueImpl{}

func NewPrioQueue(nodes []string, less func(n1, n2 string) bool) *prioQueueImpl {
	h := &prioQueueImpl{
		nodes:    nodes,
		indexMap: make(map[string]int),
		less:     less,
	}

	for i, node := range nodes {
		h.indexMap[node] = i
	}

	heap.Init(h)

	return h
}

func (h *prioQueueImpl) Len() int {
	return len(h.nodes)
}

func (h *prioQueueImpl) Less(i, j int) bool {
	return h.less(h.nodes[i], h.nodes[j])
}

func (h *prioQueueImpl) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
	h.indexMap[h.nodes[i]] = i
	h.indexMap[h.nodes[j]] = j
}

// the cost has changed for node. fix the heap
func (h *prioQueueImpl) Fix(node string) {
	if index, ok := h.indexMap[node]; ok {
		heap.Fix(h, index)
	}
}

func (h *prioQueueImpl) PushOrFix(node string) {
	if _, ok := h.indexMap[node]; !ok {
		heap.Push(h, node)
	} else {
		h.Fix(node)
	}
}

func (h *prioQueueImpl) Push(val interface{}) {
	node := val.(string)
	h.nodes = append(h.nodes, node)
	h.indexMap[node] = len(h.nodes) - 1
}

func (h *prioQueueImpl) Pop() interface{} {
	l := len(h.nodes)
	node := h.nodes[l-1]
	h.nodes = h.nodes[:l-1]
	delete(h.indexMap, node)

	return node
}

func (h *prioQueueImpl) Add(node string) {
	heap.Push(h, node)
}

func (h *prioQueueImpl) Remove() (string, error) {
	if len(h.nodes) == 0 {
		return "", ErrNoNodesInHeap
	}

	return heap.Pop(h).(string), nil
}
