package graph

import (
	"fmt"
)

type GraphPath interface {
	New() Graph
	Nodes() []string
	Walk(v string, visit func(w string, c uint) bool) error
}

type graphPath struct {
	graphRef GraphPath
}

func newGraphPath(graphRef GraphPath) *graphPath {
	return &graphPath{graphRef: graphRef}
}

func (g *graphPath) shortestPaths(from string) (map[string]uint, map[string]string, error) {
	nodes := g.graphRef.Nodes()
	cost := make(map[string]uint, len(nodes))
	parents := make(map[string]string, len(nodes))

	cost[from] = 0
	parents[from] = ""

	prioQueue := NewHeap([]string{}, func(n1, n2 string) bool {
		c1, ok := cost[n1]
		if !ok {
			c1 = Infinity
		}

		c2, ok := cost[n2]
		if !ok {
			c2 = Infinity
		}

		return c1 < c2
	})

	prioQueue.Push(from)

	for prioQueue.Len() > 0 {
		vertex, err := prioQueue.Pop()
		if err != nil {
			return nil, nil, err
		}

		g.graphRef.Walk(vertex, func(w string, c uint) (skip bool) {
			costW := cost[vertex] + c
			if _, ok := cost[w]; !ok {
				cost[w] = costW
				parents[w] = vertex
				prioQueue.Push(w)
			} else if costW < cost[w] {
				cost[w] = costW
				parents[w] = vertex
				prioQueue.Fix(w)
			}

			return
		})
	}

	return cost, parents, nil
}

func (g *graphPath) findAllShortestPathsAndCost(from, to string) ([][]string, uint, error) {
	subGraph, costMap, err := g.shortestPathsAll(from)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: error finding all shortest paths", err)
	}

	cost, ok := costMap[to]
	if !ok {
		return nil, 0, fmt.Errorf("%w: no path found in graph for %s->%s", ErrNoPathInGraph, from, to)
	}

	paths := [][]string{}

	visitedPaths, _ := IsVisited([]string{}, from)

	subGraph.Walk(from, func(w string, c uint) bool {
		subpaths := g.findAllPaths(subGraph, w, to, visitedPaths)
		if len(subpaths) > 0 {
			for _, subpath := range subpaths {
				subpath = append([]string{from}, subpath...)
				paths = append(paths, subpath)
			}
		}

		return false
	})

	return paths, cost, nil
}

func (g *graphPath) findAllShortestPathsAndCostBFS(from, to string) ([][]string, uint, error) {
	subGraph, costMap, err := g.shortestPathsAll(from)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: error finding all shortest paths", err)
	}

	cost, ok := costMap[to]
	if !ok {
		return nil, 0, fmt.Errorf("%w: no path found in graph for %s->%s", ErrNoPathInGraph, from, to)
	}

	paths := g.findAllPathsBFS(subGraph, from, to)

	return paths, cost, nil
}

func (g *graphPath) findAllPaths(subGraph Graph, from, to string, visitedPaths []string) [][]string {
	paths := [][]string{}

	if from == to {
		paths = append(paths, []string{from})

		return paths
	}

	visitedPaths, visited := IsVisited(visitedPaths, from)
	if visited {
		return paths
	}

	subGraph.Walk(from, func(w string, c uint) bool {
		subpaths := g.findAllPaths(subGraph, w, to, visitedPaths)
		if len(subpaths) > 0 {
			// prefix ourselves to each of the subpath
			for _, subpath := range subpaths {
				subpath = append([]string{from}, subpath...)
				paths = append(paths, subpath)
			}
		}

		return false
	})

	return paths
}

func (g *graphPath) findAllPathsBFS(subGraph Graph, from, to string) [][]string {
	paths := [][]string{}

	if from == to {
		paths = append(paths, []string{from})

		return paths
	}

	for queue := []interface{}{from, []string{}}; len(queue) > 0; {
		entry := queue[0].(string)
		parents := queue[1].([]string)
		queue = queue[2:]
		path := append([]string{}, parents...)
		path = append(path, entry)

		subGraph.Walk(entry, func(w string, c uint) bool {
			if w == to {
				result := append([]string{}, path...)
				result = append(result, w)
				paths = append(paths, result)

				return false
			}

			if _, visited := IsVisited(path, w); visited {
				return false
			}

			queue = append(queue, []interface{}{w, path}...)

			return false
		})
	}

	return paths
}

func (g *graphPath) shortestPathsAll(v string) (Graph, map[string]uint, error) {
	subGraph := g.graphRef.New()

	cost := make(map[string]uint, len(g.graphRef.Nodes()))

	edgesMap := make(map[string][]NodeAndNeighbor)

	prioQueue := NewHeap([]string{}, func(n1, n2 string) bool {
		c1, ok := cost[n1]
		if !ok {
			c1 = Infinity
		}

		c2, ok := cost[n2]
		if !ok {
			c2 = Infinity
		}

		return c1 < c2
	})

	cost[v] = 0

	prioQueue.Push(v)

	for prioQueue.Len() > 0 {
		v, err := prioQueue.Pop()
		if err != nil {
			return nil, nil, err
		}

		g.graphRef.Walk(v, func(w string, c uint) bool {
			costW := cost[v] + c
			currentCost, ok := cost[w]
			if !ok {
				cost[w] = costW
				err = subGraph.AddWithCost(Edge{Node: v, Neighbor: w, Cost: costW})
				if err != nil {
					return true
				}

				edgesMap[w] = append(edgesMap[w], NodeAndNeighbor{node: v, neighbor: w})
				prioQueue.Push(w)
			} else if costW <= currentCost {

				if costW < currentCost {
					cost[w] = costW
					prioQueue.Fix(w)

					removeEdges(subGraph, edgesMap[w])
					edgesMap[w] = []NodeAndNeighbor{}
				}

				err = subGraph.AddWithCost(Edge{Node: v, Neighbor: w, Cost: costW})
				if err != nil {
					return true
				}

				edgesMap[w] = append(edgesMap[w], NodeAndNeighbor{node: v, neighbor: w})
			}

			return false
		})

		if err != nil {
			return nil, nil, err
		}
	}

	return subGraph, cost, nil
}

func removeEdges(graph Graph, edges []NodeAndNeighbor) {
	for _, e := range edges {
		graph.RemoveEdge(e)
	}
}
