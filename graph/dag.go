package graph

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type DAG interface {
	Graph
	TopologicalSort() ([]NodeAndDepth, error)
}

type DirectedGraphImpl struct {
	nodes   map[string]*Node
	edges   []*Edge
	edgeMap map[NodeAndNeighbor]*Edge
	*graphPath
	sync.Mutex
}

var _ Graph = &DirectedGraphImpl{}
var _ GraphPath = &DirectedGraphImpl{}

func NewDirectedGraph() *DirectedGraphImpl {
	dag := &DirectedGraphImpl{
		nodes:   make(map[string]*Node),
		edges:   make([]*Edge, 0),
		edgeMap: make(map[NodeAndNeighbor]*Edge),
	}
	dag.graphPath = newGraphPath(dag)

	return dag
}

func (d *DirectedGraphImpl) Clone() *DirectedGraphImpl {
	d.Lock()
	defer d.Unlock()

	c := &DirectedGraphImpl{
		nodes:   make(map[string]*Node),
		edges:   make([]*Edge, 0),
		edgeMap: make(map[NodeAndNeighbor]*Edge),
	}

	for name := range d.nodes {
		c.nodes[name] = &Node{name: name}
	}

	for k, e := range d.edgeMap {
		edgeCopy := &Edge{}
		*edgeCopy = *e
		c.edgeMap[k] = edgeCopy
	}

	for _, e := range d.edges {
		c.edges = append(c.edges, c.edgeMap[NodeAndNeighbor{e.Node, e.Neighbor}])
		c.nodes[e.Node].outgoing = append(c.nodes[e.Node].outgoing, c.nodes[e.Neighbor])
		c.nodes[e.Neighbor].incoming = append(c.nodes[e.Neighbor].incoming, c.nodes[e.Node])
	}

	return c
}

func (d *DirectedGraphImpl) New() Graph {
	return NewDirectedGraph()
}

func (d *DirectedGraphImpl) Order() int {
	return len(d.nodes)
}

func (d *DirectedGraphImpl) Size() int {
	return len(d.edges)
}

func (d *DirectedGraphImpl) Nodes() []string {
	nodes := make([]string, len(d.nodes))
	index := 0

	for name := range d.nodes {
		nodes[index] = name
		index++
	}

	return nodes
}

func (d *DirectedGraphImpl) String() string {
	buffer := bytes.NewBufferString("")

	buffer.WriteString(fmt.Sprintf("Nodes: %d\n", len(d.nodes)))
	buffer.WriteString(fmt.Sprintf("Edges: %d\n", len(d.edges)))

	for _, n := range d.nodes {
		buffer.WriteString(fmt.Sprintf("Name: %s\n", n.name))
		buffer.WriteString(fmt.Sprintf("Incoming: %d\n", len(n.incoming)))
		buffer.WriteString(fmt.Sprintf("Outgoing: %d\n", len(n.outgoing)))
	}

	for _, e := range d.edges {
		buffer.WriteString(fmt.Sprintf("Edge: (%s->%v, Cost: %d\n", e.Node, e.Neighbor, e.Cost))
	}

	return buffer.String()
}

func (n *Node) setWeight(weight uint, outgoing uint) {
	switch {
	case n.weight == 0 || outgoing == 0:
		n.weight = weight
	default:
		n.weight += weight
		n.weight = n.weight / outgoing
	}
}

func (d *DirectedGraphImpl) AddWithCost(edge Edge) error {
	d.Lock()
	defer d.Unlock()

	src := d.nodes[edge.Node]
	dst := d.nodes[edge.Neighbor]

	// add node with weight
	if edge.Node == edge.Neighbor {
		if src == nil {
			d.nodes[edge.Node] = &Node{name: edge.Node, weight: edge.Cost}
		} else {
			src.weight = edge.Cost
		}

		return nil
	}

	// if the edge already exists, just update the cost
	if e, ok := d.edgeMap[NodeAndNeighbor{edge.Node, edge.Neighbor}]; ok {
		e.Cost = edge.Cost
		return nil
	}

	if src != nil {
		path := make([]string, 0, len(d.nodes))
		checkForLoop := func(n *Node) error {
			if n.name == edge.Neighbor {
				return fmt.Errorf("cycle detected for directedGraph. Node %s is incoming to %s along path: %v", edge.Neighbor, edge.Node, path)
			}

			return nil
		}

		visited := make(map[string]struct{})
		if err := d.loopDetect(src, path, visited, checkForLoop); err != nil {
			return err
		}
		src.setWeight(edge.Cost, uint(len(src.outgoing)+1))
	} else {
		d.nodes[edge.Node] = &Node{name: edge.Node, weight: edge.Cost}
		src = d.nodes[edge.Node]
	}

	if dst == nil {
		d.nodes[edge.Neighbor] = &Node{name: edge.Neighbor}
		dst = d.nodes[edge.Neighbor]
	}

	src.outgoing = append(src.outgoing, dst)
	dst.incoming = append(dst.incoming, src)

	e := &edge
	d.edgeMap[NodeAndNeighbor{edge.Node, edge.Neighbor}] = e
	d.edges = append(d.edges, e)

	return nil
}

func (d *DirectedGraphImpl) loopDetect(src *Node, path []string, visited map[string]struct{}, checkForLoop func(n *Node) error) error {
	if _, ok := visited[src.name]; ok {
		return nil
	}

	visited[src.name] = struct{}{}
	path = append(path, src.name)

	for _, incoming := range src.incoming {
		if err := d.loopDetect(incoming, path, visited, checkForLoop); err != nil {
			return err
		}
	}

	if err := checkForLoop(src); err != nil {
		return err
	}

	path = path[:len(path)-1]
	return nil
}

func (d *DirectedGraphImpl) RemoveEdge(nodeAndNeighbor NodeAndNeighbor) error {
	edge, ok := d.edgeMap[nodeAndNeighbor]
	if !ok {
		return fmt.Errorf("%w: no edge found for %s->%s", ErrNoEdgesInGraph,
			nodeAndNeighbor.node, nodeAndNeighbor.neighbor)
	}

	delete(d.edgeMap, nodeAndNeighbor)

	for i, e := range d.edges {
		if e == edge {
			d.edges = append(d.edges[:i], d.edges[i+1:]...)
			break
		}
	}

	srcNode := d.nodes[nodeAndNeighbor.node]
	dstNode := d.nodes[nodeAndNeighbor.neighbor]

	for i, dst := range srcNode.outgoing {
		if dst.name == dstNode.name {
			srcNode.outgoing = append(srcNode.outgoing[:i], srcNode.outgoing[i+1:]...)
			break
		}
	}

	for i, src := range dstNode.incoming {
		if src.name == srcNode.name {
			dstNode.incoming = append(dstNode.incoming[:i], dstNode.incoming[i+1:]...)
			break
		}
	}

	return nil
}

func (d *DirectedGraphImpl) getIncomingDegrees() (map[string]int, error) {
	degree := make(map[string]int, len(d.nodes))

	for _, node := range d.nodes {
		degree[node.name] = len(node.incoming)
	}

	return degree, nil
}

func (d *DirectedGraphImpl) TopologicalSort() ([]NodeAndDepth, error) {
	d.Lock()
	defer d.Unlock()

	degreeMap, err := d.getIncomingDegrees()
	if err != nil {
		return nil, err
	}

	queue := list.New()
	for name, degree := range degreeMap {
		if degree == 0 {
			queue.PushBack(&nodeReferenceAndDepth{node: d.nodes[name], depth: 0, weight: d.nodes[name].weight})
		}
	}

	sortedQueue := list.New()

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)

		nodeRef := element.Value.(*nodeReferenceAndDepth)
		sortedQueue.PushBack(nodeRef)

		delete(degreeMap, nodeRef.node.name)
		depth := nodeRef.depth + 1

		for i := 0; i < len(nodeRef.node.outgoing); i++ {
			outgoing := nodeRef.node.outgoing[i]
			degreeMap[outgoing.name]--

			if degreeMap[outgoing.name] == 0 {
				queue.PushBack(&nodeReferenceAndDepth{node: outgoing, depth: depth, weight: outgoing.weight})
			}
		}
	}

	if len(degreeMap) > 0 {
		return nil, errors.New("not-a-directedGraph")
	}

	return refListToNodes(sortedQueue)
}

func (d *DirectedGraphImpl) dfs(vertex *Node, queue *list.List, data *DFSData, depth int) error {

	if data.NodeColor[vertex.name] == Black {
		return nil
	}

	if data.NodeColor[vertex.name] == Gray {
		return errors.New("not-a-directedGraph")
	}

	data.Time++
	data.NodeColor[vertex.name] = Gray
	data.Discover[vertex.name] = data.Time

	for _, outgoing := range vertex.outgoing {
		data.Prev[outgoing.name] = vertex.name
		if err := d.dfs(outgoing, queue, data, depth+1); err != nil {
			return err
		}
	}

	queue.PushFront(&nodeReferenceAndDepth{node: vertex, depth: depth, weight: vertex.weight})

	data.Time++
	data.NodeColor[vertex.name] = Black
	data.Finish[vertex.name] = data.Time

	return nil
}

func (d *DirectedGraphImpl) DFSWithData() ([]NodeAndDepth, *DFSData, error) {
	queue := list.New()
	data := &DFSData{
		Prev:      make(map[string]string, len(d.nodes)),
		Discover:  make(map[string]int, len(d.nodes)),
		NodeColor: make(map[string]Color, len(d.nodes)),
		Finish:    make(map[string]int, len(d.nodes)),
	}

	d.Lock()
	defer d.Unlock()

	for _, node := range d.nodes {
		if len(node.incoming) > 0 {
			continue
		}

		if err := d.dfs(node, queue, data, 0); err != nil {
			return nil, nil, err
		}
	}

	nodeAndDepth, err := refListToNodes(queue)
	if err != nil {
		return nil, nil, err
	}

	return nodeAndDepth, data, nil
}

func (d *DirectedGraphImpl) DFS() ([]NodeAndDepth, error) {
	nd, _, err := d.DFSWithData()
	if err != nil {
		return nil, err
	}

	return nd, nil
}

func (d *DirectedGraphImpl) visitNoLock(vertex string, visit func(w string, c uint) bool) error {
	node, ok := d.nodes[vertex]
	if !ok {
		return fmt.Errorf("%w: no node %s in graph", ErrNoNodeInGraph, vertex)
	}

	for _, outgoing := range node.outgoing {
		e := d.edgeMap[NodeAndNeighbor{vertex, outgoing.name}]
		if visit(e.Neighbor, e.Cost) {
			return nil
		}
	}

	return nil
}

func (d *DirectedGraphImpl) Visit(vertex string, visit func(w string, c uint) bool) error {
	d.Lock()
	defer d.Unlock()

	return d.visitNoLock(vertex, visit)
}

func (d *DirectedGraphImpl) Walk(vertex string, visit func(w string, c uint) bool) error {
	return d.visitNoLock(vertex, visit)
}

func (d *DirectedGraphImpl) BFS(vertex string, visit func(v, w string, c uint) bool) ([]string, error) {
	d.Lock()
	defer d.Unlock()

	if _, ok := d.nodes[vertex]; !ok {
		return nil, fmt.Errorf("%w: no node %s in graph", ErrNoNodeInGraph, vertex)
	}

	path := make([]string, 0, len(d.nodes))
	skip := false

	node := d.nodes[vertex]

	visitedMap := make(map[string]struct{})
	visitedMap[node.name] = struct{}{}

	for queue := []*Node{node}; len(queue) > 0 && skip == false; {
		n := queue[0]
		queue = queue[1:]

		path = append(path, n.name)
		for _, outgoing := range n.outgoing {
			if _, ok := visitedMap[outgoing.name]; ok {
				continue
			}
			visitedMap[outgoing.name] = struct{}{}
			e := d.edgeMap[NodeAndNeighbor{n.name, outgoing.name}]

			if visit(e.Node, e.Neighbor, e.Cost) {
				skip = true
				break
			}

			queue = append(queue, outgoing)
		}
	}

	return path, nil
}

func (d *DirectedGraphImpl) ShortestPathAndCost(v, w string) ([]string, uint, error) {
	d.Lock()
	defer d.Unlock()

	cost, parents, err := d.shortestPaths(v)
	if err != nil {
		return nil, 0, err
	}

	if cost[w] == Infinity || parents[w] == "" {
		return nil, 0, fmt.Errorf("%w: no path in graph from %s->%s", ErrNoPathInGraph, v, w)
	}

	paths := []string{}

	for p := w; p != ""; p = parents[p] {
		paths = append(paths, p)
	}

	for i, j := 0, len(paths)-1; i < j; i, j = i+1, j-1 {
		paths[i], paths[j] = paths[j], paths[i]
	}

	return paths, cost[w], nil
}

func (d *DirectedGraphImpl) ShortestPath(v, w string) ([]string, error) {
	path, _, err := d.ShortestPathAndCost(v, w)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (d *DirectedGraphImpl) ShortestPaths(from string) (map[string]uint, map[string]string, error) {
	d.Lock()
	defer d.Unlock()

	return d.shortestPaths(from)
}

func (d *DirectedGraphImpl) FindAllShortestPathsAndCost(from, to string) ([][]string, uint, error) {
	d.Lock()
	defer d.Unlock()

	return d.findAllShortestPathsAndCost(from, to)
}

func (d *DirectedGraphImpl) FindAllShortestPaths(from, to string) ([][]string, error) {
	paths, _, err := d.FindAllShortestPathsAndCost(from, to)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func (d *DirectedGraphImpl) FindAllShortestPathsAndCostBFS(from, to string) ([][]string, uint, error) {
	d.Lock()
	defer d.Unlock()

	return d.findAllShortestPathsAndCostBFS(from, to)
}

func (d *DirectedGraphImpl) FindAllShortestPathsBFS(from, to string) ([][]string, error) {
	paths, _, err := d.FindAllShortestPathsAndCostBFS(from, to)
	if err != nil {
		return nil, err
	}

	return paths, nil
}
