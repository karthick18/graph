package graph

import "errors"

var (
	ErrNoEdgesInGraph = errors.New("no-edges-in-graph")
	ErrNoNodesInHeap  = errors.New("no-nodes-in-heap")
	ErrNoPathInGraph  = errors.New("no-path-found-in-graph")
	ErrNoNodeInGraph  = errors.New("no-node-in-graph")
)
