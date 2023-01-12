package graph

type Graph interface {
	Order() int
	Size() int
	Visit(vertex string, visit func(w string, c uint) bool) error
	Walk(vertex string, visit func(w string, c uint) bool) error
	BFS(vertex string, visit func(v, w string, c uint) bool) ([]string, error)
	DFS(...string) ([]NodeAndDepth, error)
	DFSWithData(...string) ([]NodeAndDepth, *DFSData, error)
	AddWithCost(Edge) error
	ShortestPath(v, w string) ([]string, error)
	ShortestPathAndCost(v, w string) ([]string, uint, error)
	ShortestPaths(v string) (map[string]uint, map[string]string, error)
	RemoveEdge(NodeAndNeighbor) error
}

func IsVisited(visited []string, child string) ([]string, bool) {
	for _, v := range visited {
		if v == child {
			return nil, true
		}
	}

	visited = append([]string{}, visited...)
	visited = append(visited, child)

	return visited, false
}
