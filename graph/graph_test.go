package graph_test

import (
	"strings"
	"testing"

	"github.com/karthick18/go-graph/graph"
	"github.com/stretchr/testify/assert"
)

func TestUndirectedGraph(t *testing.T) {
	from, to := "a", "e"

	g := graph.NewUndirectedGraph()

	g.AddWithCostBoth(graph.Edge{Node: "a", Neighbor: "b", Cost: uint(3)})
	g.AddWithCostBoth(graph.Edge{Node: "b", Neighbor: "c", Cost: uint(5)})
	g.AddWithCostBoth(graph.Edge{Node: "a", Neighbor: "c", Cost: uint(8)})
	g.AddWithCostBoth(graph.Edge{Node: "a", Neighbor: "d", Cost: uint(1)})
	g.AddWithCostBoth(graph.Edge{Node: "d", Neighbor: "e", Cost: uint(10)})
	g.AddWithCostBoth(graph.Edge{Node: "e", Neighbor: "c", Cost: uint(4)})
	g.AddWithCostBoth(graph.Edge{Node: "c", Neighbor: "d", Cost: uint(6)})

	t.Log("Undirected Graph", g)
	t.Log("Find shortest path", from, "to", to)

	path, cost, err := g.ShortestPathAndCost(from, to)
	assert.Nil(t, err, "error finding shortest path")

	t.Log("shortest path from", from, "to", to, "=", strings.Join(path, "->"), "with cost", cost)

	t.Log("Find all shortest paths", from, "to", to)

	paths, cost, err := g.FindAllShortestPathsAndCost(from, to)
	assert.Nil(t, err, "error finding all shortest paths")

	for _, path := range paths {
		t.Log("shortest path", from, "to", to, "=", strings.Join(path, "->"), "cost", cost)
	}

	t.Log("BFS...")
	bfsNodes, err := g.BFS(from, func(vertex, neighbor string, cost uint) bool {
		t.Log("from", vertex, "to", neighbor, "cost", cost)
		return false
	})

	assert.Nil(t, err, "error doing a breadth first search")

	t.Log("BFS nodes from", from, "=", bfsNodes)

	t.Log("DFS...")
	pathAndDepth, err := g.DFS()
	assert.Nil(t, err, "error doing depth first search")

	t.Log("DFS path and depth", pathAndDepth)
}

func TestDirectedGraph(t *testing.T) {

	dag := graph.NewDirectedGraph()

	dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "11", Cost: uint(3)})
	dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "7", Cost: uint(5)})
	dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "2", Cost: uint(5)})
	dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "9", Cost: uint(7)})
	dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "10", Cost: uint(10)})
	dag.AddWithCost(graph.Edge{Node: "7", Neighbor: "11", Cost: uint(1)})
	dag.AddWithCost(graph.Edge{Node: "7", Neighbor: "8", Cost: uint(2)})
	dag.AddWithCost(graph.Edge{Node: "8", Neighbor: "9", Cost: uint(4)})
	dag.AddWithCost(graph.Edge{Node: "3", Neighbor: "8", Cost: uint(6)})
	dag.AddWithCost(graph.Edge{Node: "3", Neighbor: "10", Cost: uint(2)})

	t.Log(dag)

	pathAndDepth, err := dag.DFS()
	assert.Nil(t, err, "error doing depth first search")

	t.Log("dfs of dag", pathAndDepth)

	pathAndDepth, err = dag.TopologicalSort()
	assert.Nil(t, err, "error doing topological sort")

	t.Log("topological sort of dag", pathAndDepth)

	dag.Visit("11", func(w string, cost uint) (skip bool) {
		t.Log("11", "->", w, "cost", cost)
		return
	})

	path, cost, err := dag.ShortestPathAndCost("5", "9")
	assert.Nil(t, err, "error finding shortest path")

	t.Log("Shortest path from 5 -> 9 =", strings.Join(path, "->"), "cost", cost)
}
