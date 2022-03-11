package graph_test

import (
	"errors"
	"reflect"
	"sort"
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

	assert.Equal(t, g.Order(), 5, "Graph order should be 5")
	assert.Equal(t, g.Size(), 14, "Graph size should be 14")

	t.Log("Find shortest path", from, "to", to)

	path, cost, err := g.ShortestPathAndCost(from, to)
	assert.Nil(t, err, "error finding shortest path")

	shortestPath := strings.Join(path, "->")

	t.Log("shortest path from", from, "to", to, "=", shortestPath, "with cost", cost)

	assert.Equal(t, shortestPath, "a->d->e", "Shortest path from a to e")
	assert.Equal(t, cost, uint(11), "Shortest path from a to e should have cost 11")

	t.Log("Find all shortest paths", from, "to", to)

	allExpectedPaths := []string{"a->d->e", "a->d->c->e"}
	sort.Slice(allExpectedPaths, func(i, j int) bool {
		return allExpectedPaths[i] < allExpectedPaths[j]
	})

	paths, cost, err := g.FindAllShortestPathsAndCost(from, to)
	assert.Nil(t, err, "error finding all shortest paths")

	allPaths := make([]string, len(paths))

	for i, path := range paths {
		allPaths[i] = strings.Join(path, "->")
		t.Log("shortest path", from, "to", to, "=", allPaths[i], "cost", cost)
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return allPaths[i] < allPaths[j]
	})

	assert.Equal(t, reflect.DeepEqual(allPaths, allExpectedPaths), true, "All shortest paths with equal cost mismatch")
	assert.Equal(t, cost, uint(11), "All shortest paths cost mismatch")

	t.Log("BFS...")

	expectedNodes := []string{"a", "b", "c", "d", "e"}
	bfsNodes, err := g.BFS(from, func(vertex, neighbor string, cost uint) bool {
		t.Log("from", vertex, "to", neighbor, "cost", cost)
		return false
	})

	assert.Nil(t, err, "error doing a breadth first search")
	assert.Equal(t, reflect.DeepEqual(bfsNodes, expectedNodes), true, "BFS walk neighbor mismatch")

	t.Log("BFS nodes from", from, "=", bfsNodes)

	t.Log("DFS...")

	pathAndDepth, err := g.DFS()
	assert.Nil(t, err, "error doing depth first search")
	t.Log("DFS path and depth", pathAndDepth)
}

func TestUndirectedGraphBridges(t *testing.T) {
	g := graph.NewUndirectedGraph()
	g.AddWithCostBoth(graph.Edge{Node: "1", Neighbor: "0"})
	g.AddWithCostBoth(graph.Edge{Node: "0", Neighbor: "2"})
	g.AddWithCostBoth(graph.Edge{Node: "2", Neighbor: "1"})
	g.AddWithCostBoth(graph.Edge{Node: "0", Neighbor: "3"})
	g.AddWithCostBoth(graph.Edge{Node: "3", Neighbor: "4"})

	expectedEdges := []graph.Edge{
		{Node: "0", Neighbor: "3"},
		{Node: "3", Neighbor: "4"},
	}

	expectedMap := make(map[graph.Edge]struct{})

	for _, edge := range expectedEdges {
		expectedMap[edge] = struct{}{}
		expectedMap[graph.Edge{Node: edge.Neighbor, Neighbor: edge.Node}] = struct{}{}
	}

	res := g.GetBridges()
	t.Log("Get Bridges", res)
	assert.Equal(t, len(expectedEdges), len(res), "Get bridges result mismatch")

	for _, r := range res {
		_, ok := expectedMap[r]
		assert.Equal(t, ok, true, "Get bridges result mismatch")
	}

	g = graph.NewUndirectedGraph()
	g.AddWithCostBoth(graph.Edge{Node: "0", Neighbor: "1"})
	g.AddWithCostBoth(graph.Edge{Node: "1", Neighbor: "2"})
	g.AddWithCostBoth(graph.Edge{Node: "2", Neighbor: "3"})

	expectedEdges = []graph.Edge{
		{Node: "0", Neighbor: "1"},
		{Node: "1", Neighbor: "2"},
		{Node: "2", Neighbor: "3"},
	}

	expectedMap = make(map[graph.Edge]struct{})

	for _, edge := range expectedEdges {
		expectedMap[edge] = struct{}{}
		expectedMap[graph.Edge{Node: edge.Neighbor, Neighbor: edge.Node}] = struct{}{}
	}

	res = g.GetBridges()

	assert.Equal(t, len(expectedEdges), len(res), "Get bridges result mismatch")

	for _, r := range res {
		_, ok := expectedMap[r]
		assert.Equal(t, ok, true, "Get bridges result mismatch")
	}

	g = graph.NewUndirectedGraph()
	g.AddWithCostBoth(graph.Edge{Node: "0", Neighbor: "1"})
	g.AddWithCostBoth(graph.Edge{Node: "1", Neighbor: "2"})
	g.AddWithCostBoth(graph.Edge{Node: "2", Neighbor: "0"})
	g.AddWithCostBoth(graph.Edge{Node: "1", Neighbor: "3"})
	g.AddWithCostBoth(graph.Edge{Node: "1", Neighbor: "4"})
	g.AddWithCostBoth(graph.Edge{Node: "1", Neighbor: "6"})
	g.AddWithCostBoth(graph.Edge{Node: "3", Neighbor: "5"})
	g.AddWithCostBoth(graph.Edge{Node: "4", Neighbor: "5"})

	expectedMap = make(map[graph.Edge]struct{})
	expectedEdges = []graph.Edge{
		{Node: "1", Neighbor: "6"},
	}

	for _, edge := range expectedEdges {
		expectedMap[edge] = struct{}{}
		expectedMap[graph.Edge{Node: edge.Neighbor, Neighbor: edge.Node}] = struct{}{}
	}

	res = g.GetBridges()
	assert.Equal(t, len(expectedEdges), len(res), "Get bridges result mismatch")

	for _, r := range res {
		_, ok := expectedMap[r]
		assert.Equal(t, ok, true, "Get bridges result mismatch")
	}
}

func TestDirectedGraph(t *testing.T) {

	dag := graph.NewDirectedGraph()

	dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "11", Cost: uint(3)})
	dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "7", Cost: uint(4)})
	dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "2", Cost: uint(5)})
	dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "9", Cost: uint(7)})
	dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "10", Cost: uint(10)})
	dag.AddWithCost(graph.Edge{Node: "7", Neighbor: "11", Cost: uint(1)})
	dag.AddWithCost(graph.Edge{Node: "7", Neighbor: "8", Cost: uint(2)})
	dag.AddWithCost(graph.Edge{Node: "8", Neighbor: "9", Cost: uint(4)})
	dag.AddWithCost(graph.Edge{Node: "3", Neighbor: "8", Cost: uint(6)})
	dag.AddWithCost(graph.Edge{Node: "3", Neighbor: "10", Cost: uint(2)})
	dag.AddWithCost(graph.Edge{Node: "9", Neighbor: "13", Cost: uint(3)})

	t.Log(dag)

	pathAndDepth, err := dag.DFS()
	assert.Nil(t, err, "error doing depth first search")

	t.Log("dfs of dag", pathAndDepth)

	expectedPathAndDepth := []graph.NodeAndDepth{
		{"5", 0}, {"3", 0}, {"7", 1}, {"8", 2}, {"11", 2}, {"2", 3}, {"10", 3}, {"9", 3}, {"13", 4},
	}
	pathAndDepth, err = dag.TopologicalSort()
	assert.Nil(t, err, "error doing topological sort")

	t.Log("topological sort of dag", pathAndDepth)
	assert.Equal(t, reflect.DeepEqual(pathAndDepth, expectedPathAndDepth), true, "DAG topological sort mismatch")

	dag.Visit("11", func(w string, cost uint) (skip bool) {
		t.Log("11", "->", w, "cost", cost)
		return
	})

	expectedPath := "5->11->9"
	path, cost, err := dag.ShortestPathAndCost("5", "9")
	assert.Nil(t, err, "error finding shortest path")

	actualPath := strings.Join(path, "->")
	t.Log("Shortest path from 5 -> 9 =", actualPath, "cost", cost)

	assert.Equal(t, actualPath, expectedPath, "Shortest path mismatch for DAG")
	assert.Equal(t, cost, uint(10), "Shortest path cost mismatch for DAG")

	expectedPaths := []string{"5->11->9->13", "5->7->8->9->13"}
	sort.Slice(expectedPaths, func(i, j int) bool {
		return expectedPaths[i] < expectedPaths[j]
	})

	paths, cost, err := dag.FindAllShortestPathsAndCost("5", "13")
	assert.Nil(t, err, "error finding all shortest paths and cost for DAG")

	results := make([]string, len(paths))

	for i, path := range paths {
		results[i] = strings.Join(path, "->")
		t.Log("Shortest path from 5->13 =", results[i], "cost", cost)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})

	assert.Equal(t, reflect.DeepEqual(results, expectedPaths), true, "Find all shortest paths mismatch")
	assert.Equal(t, cost, uint(13), "Find all shortest paths cost mismatch")

	dag = graph.NewDirectedGraph(graph.Unconnected)

	err = dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "11", Cost: uint(3)})
	assert.Nil(t, err, "DAG add failure for 5->11")

	err = dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "7", Cost: uint(4)})
	assert.Nil(t, err, "DAG add failure for 5->7")

	err = dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "2", Cost: uint(5)})
	assert.Nil(t, err, "DAG add failure for 11->2")

	t.Log(dag)

	err = dag.AddWithCost(graph.Edge{Node: "2", Neighbor: "5", Cost: uint(5)})
	if err != nil {
		t.Log("DAG add with cost 2->5 error", err)
	}

	assert.Equal(t, errors.Is(err, graph.ErrLoopInDag), true, "DAG add did not fail with loop detection for 2->5")
}

func TestDirectedGraphConnected(t *testing.T) {
	dag := graph.NewDirectedGraph(graph.Unconnected)

	err := dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "11", Cost: uint(3)})
	assert.Nil(t, err, "DAG add failure for 5->11")

	err = dag.AddWithCost(graph.Edge{Node: "5", Neighbor: "7", Cost: uint(4)})
	assert.Nil(t, err, "DAG add failure for 5->7")

	err = dag.AddWithCost(graph.Edge{Node: "11", Neighbor: "2", Cost: uint(5)})
	assert.Nil(t, err, "DAG add failure for 11->2")

	t.Log(dag)

	err = dag.AddWithCost(graph.Edge{Node: "2", Neighbor: "5", Cost: uint(5)})
	if err != nil {
		t.Log("DAG add with cost 2->5 error", err)
	}

	assert.Equal(t, errors.Is(err, graph.ErrLoopInDag), true, "DAG add did not fail with loop detection for 2->5")

	dag = graph.NewDirectedGraph(graph.Connected)
	err = dag.AddWithCost(graph.Edge{Node: "0", Neighbor: "1"})
	assert.Nil(t, err, "DAG add failure for 0->1")

	err = dag.AddWithCost(graph.Edge{Node: "1", Neighbor: "2"})
	assert.Nil(t, err, "DAG add failure for 1->2")

	err = dag.AddWithCost(graph.Edge{Node: "2", Neighbor: "3"})
	assert.Nil(t, err, "DAG add failure for 2->3")

	isConnected, err := dag.IsStronglyConnected()
	assert.Nil(t, err, "DAG is strongly connected failed")

	assert.Equal(t, isConnected, false, "DAG should not be strongly connected")

	err = dag.AddWithCost(graph.Edge{Node: "3", Neighbor: "0"})
	assert.Nil(t, err, "DAG add failure for 3->0")

	isConnected, err = dag.IsStronglyConnected()
	assert.Nil(t, err, "DAG is strongly connected failed")

	assert.Equal(t, isConnected, true, "DAG should be strongly connected ")

	t.Log("Directed connected graph", dag)
	t.Log("Connection status", isConnected)

	scc, err := dag.GetStronglyConnectedComponents()
	assert.Nil(t, err, "DAG get strongly connected failed")

	actual := []string{}

	for _, cc := range scc {
		t.Log("connected-component", cc)

		sort.Slice(cc, func(i, j int) bool {
			return cc[i] < cc[j]
		})

		actual = append(actual, strings.Join(cc, ","))
	}

	expected := []string{"0,1,2,3"}

	assert.Equal(t, reflect.DeepEqual(expected, actual), true, "DAG strongly connected mismatch")

	dag = graph.NewDirectedGraph(graph.Connected)

	dag.AddWithCost(graph.Edge{Node: "0", Neighbor: "2"})
	dag.AddWithCost(graph.Edge{Node: "2", Neighbor: "1"})
	dag.AddWithCost(graph.Edge{Node: "1", Neighbor: "0"})
	dag.AddWithCost(graph.Edge{Node: "0", Neighbor: "3"})
	dag.AddWithCost(graph.Edge{Node: "3", Neighbor: "4"})

	isConnected, err = dag.IsStronglyConnected()
	assert.Nil(t, err, "DAG is strongly connected failed")
	assert.Equal(t, isConnected, false, "DAG should not be strongly connected")

	sccComponents, err := dag.GetStronglyConnectedComponents()
	assert.Nil(t, err, "DAG get strongly connected components failed")

	expected = []string{"0,1,2", "3", "4"}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i] < expected[j]
	})

	actual = []string{}

	for _, comp := range sccComponents {
		t.Log("connected-component", comp)
		sort.Slice(comp, func(i, j int) bool {
			return comp[i] < comp[j]
		})

		actual = append(actual, strings.Join(comp, ","))
	}

	assert.Equal(t, reflect.DeepEqual(expected, actual), true, "DAG strongly connected components mismatch")
}
