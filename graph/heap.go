package graph

import "fmt"

type heap struct {
	nodes    []string
	indexMap map[string]int
	less     func(n1, n2 string) bool
}

func NewHeap(nodes []string, less func(n1, n2 string) bool) *heap {
	h := &heap{
		nodes:    nodes,
		indexMap: make(map[string]int),
		less:     less,
	}

	for i, node := range nodes {
		h.indexMap[node] = i
	}

	// heapify
	for i := (len(nodes) - 1) / 2; i >= 0; i-- {
		h.down(i, len(nodes))
	}

	return h
}

func (h *heap) lessIndex(i, j int) bool {
	return h.less(h.nodes[i], h.nodes[j])
}

func (h *heap) down(i0, n int) bool {
	root := i0
	for {
		child := 2*root + 1
		if child >= n {
			break
		}

		if child+1 < n && h.lessIndex(child+1, child) {
			child++
		}

		if h.lessIndex(root, child) {
			break
		}

		h.swap(root, child)
		root = child
	}

	return root != i0
}

func (h *heap) up(i0 int) {
	child := i0
	for {
		parent := (child - 1) / 2
		if parent == child {
			break
		}

		if h.lessIndex(parent, child) {
			break
		}

		h.swap(parent, child)
		child = parent
	}
}

func (h *heap) swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
	h.indexMap[h.nodes[i]] = i
	h.indexMap[h.nodes[j]] = j
}

// the cost has changed for node. fix the heap
func (h *heap) Fix(node string) {
	if index, ok := h.indexMap[node]; ok {
		if !h.down(index, len(h.nodes)) {
			h.up(index)
		}
	}
}

func (h *heap) PushOrFix(node string) {
	if _, ok := h.indexMap[node]; !ok {
		h.Push(node)
	} else {
		h.Fix(node)
	}
}

func (h *heap) Pop() (string, error) {
	if len(h.nodes) == 0 {
		return "", fmt.Errorf("%w: no nodes in heap", ErrNoNodesInHeap)
	}

	top := h.nodes[0]
	h.swap(0, len(h.nodes)-1)
	h.nodes = h.nodes[:len(h.nodes)-1]
	delete(h.indexMap, top)

	h.down(0, len(h.nodes))

	return top, nil
}

func (h *heap) Push(node string) {
	h.nodes = append(h.nodes, node)
	h.indexMap[node] = len(h.nodes) - 1
	h.up(len(h.nodes) - 1)
}

func (h *heap) Len() int {
	return len(h.nodes)
}
