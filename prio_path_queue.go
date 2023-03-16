package graph

import (
	"container/heap"
)

type prioShortestPath struct {
	path []string
	cost uint
}

type prioShortestPaths []prioShortestPath

var _ heap.Interface = &prioShortestPaths{}

func NewPrioPathQueue() prioShortestPaths {
	return prioShortestPaths{}
}

func (p *prioShortestPaths) Len() int {
	return len(*p)
}

func (p *prioShortestPaths) Less(i, j int) bool {
	return (*p)[i].cost < (*p)[j].cost
}

func (p *prioShortestPaths) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *prioShortestPaths) Push(val interface{}) {
	sp := val.(prioShortestPath)
	*p = append(*p, sp)
}

func (p *prioShortestPaths) Pop() interface{} {
	l := p.Len()
	sp := (*p)[l-1]
	*p = (*p)[:l-1]

	return sp
}

func (p *prioShortestPaths) Add(sp prioShortestPath) {
	heap.Push(p, sp)
}

func (p *prioShortestPaths) Remove() (prioShortestPath, error) {
	if p.Len() == 0 {
		return prioShortestPath{}, ErrNoNodesInHeap
	}

	return heap.Pop(p).(prioShortestPath), nil
}
