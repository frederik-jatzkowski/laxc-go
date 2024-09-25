package graph

import "fmt"

type DiGraph[T comparable] struct {
	edges map[T][]T
}

func NewDiGraph[T comparable](nodes ...T) DiGraph[T] {
	result := DiGraph[T]{
		edges: make(map[T][]T, len(nodes)),
	}

	for _, node := range nodes {
		result.edges[node] = []T{}
	}

	return result
}

func (graph DiGraph[T]) AddEdge(from, to T) {
	graph.edges[from] = append(graph.edges[from], to)
	graph.edges[to] = graph.edges[to]
}

func (graph *DiGraph[T]) HasCycles() bool {
	visited := make(map[T]bool)
	recStack := make(map[T]bool) // Tracks nodes in the current recursion path

	// A helper function for DFS
	var dfs func(node T) bool
	dfs = func(node T) bool {
		if recStack[node] { // If node is already in recursion stack, we found a cycle
			return true
		}
		if visited[node] { // If already visited and not in recursion stack, no cycle from this node
			return false
		}

		// Mark the node as visited and add to recursion stack
		visited[node] = true
		recStack[node] = true

		// Visit all the adjacent nodes
		for _, neighbor := range graph.edges[node] {
			if dfs(neighbor) {
				return true
			}
		}

		// Remove the node from recursion stack before exiting
		recStack[node] = false
		return false
	}

	// Call dfs for each node in the graph
	for node := range graph.edges {
		if !visited[node] {
			if dfs(node) {
				return true
			}
		}
	}

	return false
}

func (graph *DiGraph[T]) TopologicalSort() ([]T, error) {
	var stack []T
	visited := make(map[T]bool)
	tempMarked := make(map[T]bool) // Temporary mark for nodes in the current DFS path

	// Recursive DFS function
	var dfs func(node T) bool
	dfs = func(node T) bool {
		if tempMarked[node] { // Cycle detected
			return true
		}
		if visited[node] { // Node already processed
			return false
		}

		// Mark the node temporarily
		tempMarked[node] = true

		// Visit all adjacent nodes
		for _, neighbor := range graph.edges[node] {
			if dfs(neighbor) {
				return true // Cycle detected, propagate upwards
			}
		}

		// Mark the node as permanently visited
		visited[node] = true
		tempMarked[node] = false
		stack = append([]T{node}, stack...) // Add to stack in reverse order

		return false
	}

	// Visit all nodes in the graph
	for node := range graph.edges {
		if !visited[node] {
			if dfs(node) { // If a cycle is detected, return an error
				return nil, fmt.Errorf("graph has cycles, topological sort not possible")
			}
		}
	}

	return stack, nil // Return the reverse of the collected order
}

func (graph *DiGraph[T]) Has(node T) bool {
	_, ok := graph.edges[node]

	return ok
}

func (graph *DiGraph[T]) Dfs(node T, visited map[T]bool) (last T) {
	if visited[node] {
		return node
	}

	visited[node] = true
	last = node
	for _, neighbor := range graph.edges[node] {
		last = graph.Dfs(neighbor, visited)
	}

	return last
}
