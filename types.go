package graph

import (
	"container/list"
	"sort"
)

const (
	Infinity = ^uint(0)
)

type Color int

const (
	White Color = iota
	Gray
	Black
)

type Node struct {
	name     string
	weight   uint
	incoming []*Node
	outgoing []*Node
}

type NodeAndNeighbor struct {
	Node     string
	Neighbor string
}

type Edge struct {
	Node, Neighbor string
	Cost           uint
}

type NodeAndDepth struct {
	Node  string
	Depth int
}

type DFSData struct {
	Prev      map[string]string
	Discover  map[string]int
	Finish    map[string]int
	NodeColor map[string]Color
	Time      int
}

type nodeReferenceAndDepth struct {
	node   *Node
	depth  int
	weight uint
}

func refListToNodes(queue *list.List) ([]NodeAndDepth, error) {
	var refList []*nodeReferenceAndDepth

	for e := queue.Front(); e != nil; e = e.Next() {
		r := e.Value.(*nodeReferenceAndDepth)
		refList = append(refList, r)
	}

	sort.Slice(refList, func(i, j int) bool {
		if refList[i].depth == refList[j].depth {
			return refList[i].weight < refList[j].weight
		}

		return refList[i].depth < refList[j].depth
	})

	nodeAndDepth := []NodeAndDepth{}
	for _, ref := range refList {
		nodeAndDepth = append(nodeAndDepth, NodeAndDepth{Node: ref.node.name, Depth: ref.depth})
	}

	return nodeAndDepth, nil
}
