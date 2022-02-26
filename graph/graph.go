package graph

import "strings"

type Graph interface {
	Order() int
	Size() int
	Visit(vertex string, visit func(w string, c uint) bool) error
	Walk(vertex string, visit func(w string, c uint) bool) error
	BFS(vertex string, visit func(v, w string, c uint) bool) ([]string, error)
	DFS() ([]NodeAndDepth, error)
	DFSWithData() ([]NodeAndDepth, *DFSData, error)
	AddWithCost(Edge) error
	ShortestPath(v, w string) ([]string, error)
	ShortestPathAndCost(v, w string) ([]string, uint, error)
	ShortestPaths(v string) (map[string]uint, map[string]string, error)
	RemoveEdge(NodeAndNeighbor) error
}

type Visited string

func (visited Visited) Visit(parent, child string) (string, bool) {
	parts := strings.Split(parent, "/")

	for _, part := range parts {
		if part == child {
			return "", true
		}
	}

	childPath := parent + "/" + child

	return childPath, false
}
