package graph

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"
)

type UndirectedGraph interface {
	Graph
	AddWithCostBoth(Edge) error
}

type UndirectedGraphImpl struct {
	sync.Mutex
	edges   []*Edge
	nodes   map[string]*Node
	edgeMap map[NodeAndNeighbor]*Edge
	*graphPath
}

var _ Graph = &UndirectedGraphImpl{}
var _ GraphPath = &UndirectedGraphImpl{}

func NewUndirectedGraph() *UndirectedGraphImpl {
	undirected := &UndirectedGraphImpl{
		nodes:   make(map[string]*Node),
		edgeMap: make(map[NodeAndNeighbor]*Edge),
	}
	undirected.graphPath = newGraphPath(undirected)

	return undirected
}

func (g *UndirectedGraphImpl) Order() int {
	return len(g.nodes)
}

func (g *UndirectedGraphImpl) Size() int {
	return len(g.edges)
}

func (g *UndirectedGraphImpl) Nodes() []string {
	nodes := make([]string, len(g.nodes))
	index := 0

	for name := range g.nodes {
		nodes[index] = name
		index++
	}

	return nodes
}

func (g *UndirectedGraphImpl) New() Graph {
	return NewUndirectedGraph()
}

func (g *UndirectedGraphImpl) String() string {
	g.Lock()
	defer g.Unlock()

	buffer := bytes.NewBufferString("")

	buffer.WriteString(fmt.Sprintf("Nodes: %d\n", len(g.nodes)))
	buffer.WriteString(fmt.Sprintf("Edges: %d\n", len(g.edges)))

	for name, n := range g.nodes {
		buffer.WriteString(fmt.Sprintf("Name: %s\n", name))
		buffer.WriteString(fmt.Sprintf("Incoming: %d\n", len(n.incoming)))
		buffer.WriteString(fmt.Sprintf("Outgoing: %d\n", len(n.outgoing)))
	}

	for _, e := range g.edges {
		buffer.WriteString(fmt.Sprintf("Edge: (%s->%s), Cost: %d\n", e.Node, e.Neighbor, e.Cost))
	}

	return buffer.String()
}

func (g *UndirectedGraphImpl) AddWithCost(edge Edge) error {
	g.Lock()
	defer g.Unlock()

	vertex := g.nodes[edge.Node]
	neighbor := g.nodes[edge.Neighbor]

	if edge.Node == edge.Neighbor {
		if vertex == nil {
			vertex = &Node{name: edge.Node}
			g.nodes[edge.Node] = vertex
		}

		return nil
	}

	if e, ok := g.edgeMap[NodeAndNeighbor{node: edge.Node, neighbor: edge.Neighbor}]; ok {
		e.Cost = edge.Cost

		return nil
	}

	if vertex == nil {
		vertex = &Node{name: edge.Node}
		g.nodes[edge.Node] = vertex
	}

	if neighbor == nil {
		neighbor = &Node{name: edge.Neighbor}
		g.nodes[edge.Neighbor] = neighbor
	}

	vertex.outgoing = append(vertex.outgoing, neighbor)
	neighbor.incoming = append(neighbor.incoming, vertex)

	e := &edge
	g.edgeMap[NodeAndNeighbor{node: edge.Node, neighbor: edge.Neighbor}] = e
	g.edges = append(g.edges, e)

	return nil
}

func (g *UndirectedGraphImpl) removeFromOutgoing(vertex, neighbor *Node) {
	for index, outgoing := range vertex.outgoing {
		if outgoing.name == neighbor.name {
			vertex.outgoing = append(vertex.outgoing[:index], vertex.outgoing[index+1:]...)
			break
		}
	}
}

func (g *UndirectedGraphImpl) removeFromIncoming(vertex, neighbor *Node) {
	for index, incoming := range vertex.incoming {
		if incoming.name == neighbor.name {
			vertex.incoming = append(vertex.incoming[:index], vertex.incoming[index+1:]...)
			break
		}
	}
}

func (g *UndirectedGraphImpl) RemoveEdge(nodeAndNeighbor NodeAndNeighbor) error {
	g.Lock()
	defer g.Unlock()

	edgeRef, ok := g.edgeMap[nodeAndNeighbor]
	if !ok {
		return fmt.Errorf("%w: no edge found for %s->%s", ErrNoEdgesInGraph,
			nodeAndNeighbor.node, nodeAndNeighbor.neighbor)
	}

	vertex := g.nodes[nodeAndNeighbor.node]
	neighbor := g.nodes[nodeAndNeighbor.neighbor]

	if vertex == nil {
		return fmt.Errorf("%w: no node %s in graph", ErrNoNodeInGraph, nodeAndNeighbor.node)
	}

	if neighbor == nil {
		return fmt.Errorf("%w: no neighbor node %s in graph", ErrNoNodeInGraph, nodeAndNeighbor.neighbor)
	}

	// remove the outgoing from vertex and incoming from neighbor and vice-versa
	g.removeFromOutgoing(vertex, neighbor)
	g.removeFromIncoming(neighbor, vertex)

	delete(g.edgeMap, nodeAndNeighbor)

	for index, edge := range g.edges {
		if edge == edgeRef {
			g.edges = append(g.edges[:index], g.edges[index+1:]...)
			break
		}
	}

	return nil
}

func (g *UndirectedGraphImpl) RemoveEdgeBoth(nodeAndNeighbor NodeAndNeighbor) error {
	err := g.RemoveEdge(nodeAndNeighbor)
	if err != nil {
		return err
	}

	err = g.RemoveEdge(NodeAndNeighbor{node: nodeAndNeighbor.neighbor, neighbor: nodeAndNeighbor.node})
	if err != nil {
		return err
	}

	return nil
}

func (g *UndirectedGraphImpl) AddWithCostBoth(edge Edge) error {
	err := g.AddWithCost(edge)
	if err != nil {
		return fmt.Errorf("error adding cost for first edge: %w", err)
	}

	err = g.AddWithCost(Edge{Node: edge.Neighbor, Neighbor: edge.Node, Cost: edge.Cost})
	if err != nil {
		return fmt.Errorf("error adding cost for second edge: %w", err)
	}

	return nil
}

func (g *UndirectedGraphImpl) visitNoLock(node string, visit func(w string, c uint) bool) error {
	vertex := g.nodes[node]
	if vertex == nil {
		return fmt.Errorf("%w: visit no lock no node %s in graph", ErrNoNodeInGraph, node)
	}

	for _, outgoing := range vertex.outgoing {
		edgeRef := g.edgeMap[NodeAndNeighbor{node: node, neighbor: outgoing.name}]

		if visit(edgeRef.Neighbor, edgeRef.Cost) {
			return nil
		}
	}

	return nil
}

func (g *UndirectedGraphImpl) Walk(node string, visit func(w string, c uint) bool) error {
	return g.visitNoLock(node, visit)
}

func (g *UndirectedGraphImpl) Visit(node string, visit func(w string, c uint) bool) error {
	g.Lock()
	defer g.Unlock()

	return g.visitNoLock(node, visit)
}

func (g *UndirectedGraphImpl) BFS(node string, visit func(v, w string, c uint) bool) ([]string, error) {
	g.Lock()
	defer g.Unlock()

	vertex := g.nodes[node]
	if vertex == nil {
		return nil, fmt.Errorf("%w: no node %s in graph", ErrNoNodeInGraph, node)
	}

	path := make([]string, 0, len(g.nodes))
	visitedMap := make(map[string]struct{})

	skip := false
	visitedMap[vertex.name] = struct{}{}

	for queue := []*Node{vertex}; len(queue) > 0 && skip == false; {
		vertex = queue[0]
		queue = queue[1:]
		path = append(path, vertex.name)

		err := g.visitNoLock(vertex.name, func(w string, c uint) bool {
			if _, ok := visitedMap[w]; ok {
				return false
			}

			visitedMap[w] = struct{}{}

			edgeRef := g.edgeMap[NodeAndNeighbor{node: vertex.name, neighbor: w}]
			if visit(edgeRef.Node, edgeRef.Neighbor, edgeRef.Cost) {
				skip = true

				return true
			}

			queue = append(queue, g.nodes[w])
			return false
		})

		if err != nil {
			return nil, err
		}
	}

	return path, nil
}

func (g *UndirectedGraphImpl) dfs(vertex *Node, queue *list.List, data *DFSData, depth int) error {
	if data.NodeColor[vertex.name] != White {
		return nil
	}

	data.Time++
	data.NodeColor[vertex.name] = Gray
	data.Discover[vertex.name] = data.Time

	for _, node := range vertex.outgoing {
		if err := g.dfs(node, queue, data, depth+1); err != nil {
			return err
		}
	}

	data.Time++
	data.NodeColor[vertex.name] = Black
	data.Finish[vertex.name] = data.Time

	queue.PushFront(&nodeReferenceAndDepth{node: vertex, depth: depth})

	return nil
}

func (g *UndirectedGraphImpl) DFSWithData() ([]NodeAndDepth, *DFSData, error) {
	queue := list.New()
	data := &DFSData{
		Prev:      make(map[string]string),
		Finish:    make(map[string]int),
		Discover:  make(map[string]int),
		NodeColor: make(map[string]Color),
	}

	g.Lock()
	defer g.Unlock()

	for _, node := range g.nodes {
		if err := g.dfs(node, queue, data, 0); err != nil {
			return nil, nil, err
		}
	}

	nodeAndDepth, err := refListToNodes(queue)
	if err != nil {
		return nil, nil, err
	}

	return nodeAndDepth, data, nil
}

func (g *UndirectedGraphImpl) DFS() ([]NodeAndDepth, error) {
	nodeAndDepth, _, err := g.DFSWithData()
	if err != nil {
		return nil, err
	}

	return nodeAndDepth, nil
}

func (g *UndirectedGraphImpl) ShortestPathAndCost(v, w string) ([]string, uint, error) {
	g.Lock()
	defer g.Unlock()

	cost, parents, err := g.shortestPaths(v)
	if err != nil {
		return nil, 0, err
	}

	if parents[w] == "" {
		return nil, 0, fmt.Errorf("%w: no path in graph for path: %s->%s", ErrNoPathInGraph, v, w)
	}

	queue := list.New()

	for p := w; p != ""; p = parents[p] {
		queue.PushFront(p)
	}

	path := make([]string, queue.Len())

	for index, e := 0, queue.Front(); e != nil; index, e = index+1, e.Next() {
		path[index] = e.Value.(string)
	}

	return path, cost[w], nil
}

func (g *UndirectedGraphImpl) ShortestPath(v, w string) ([]string, error) {
	path, _, err := g.ShortestPathAndCost(v, w)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (g *UndirectedGraphImpl) ShortestPaths(v string) (map[string]uint, map[string]string, error) {
	g.Lock()
	defer g.Unlock()

	return g.shortestPaths(v)
}

func (g *UndirectedGraphImpl) FindAllShortestPathsAndCost(from, to string) ([][]string, uint, error) {
	g.Lock()
	defer g.Unlock()

	return g.findAllShortestPathsAndCost(from, to)
}

func (g *UndirectedGraphImpl) FindAllShortestPaths(from, to string) ([][]string, error) {
	paths, _, err := g.FindAllShortestPathsAndCost(from, to)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func (g *UndirectedGraphImpl) FindAllShortestPathsAndCostBFS(from, to string) ([][]string, uint, error) {
	g.Lock()
	defer g.Unlock()

	return g.findAllShortestPathsAndCostBFS(from, to)
}

func (g *UndirectedGraphImpl) FindAllShortestPathsBFS(from, to string) ([][]string, error) {
	paths, _, err := g.FindAllShortestPathsAndCostBFS(from, to)
	if err != nil {
		return nil, err
	}

	return paths, nil
}
