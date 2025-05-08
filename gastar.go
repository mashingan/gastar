// Copyright (c) 2025 Rahmatullah
// This library is licensed with Apache license which can be found in LICENSE

/*
gastar is fork of graflib library (https://github.com/mashingan/graflib)
which implemented in Nim.
gastar specifically implements A* path finding that can be
customizable like graflib.
While graflib offer more integrated graph solution with
additional edges and vertices, gastar only gives solution
of path finding search. In order gastar can search the solution,
any graph instance will need to implement method
- Neighbors, which acts as edges connector for each Node.
- Cost, between two Nodes
- Distance, between two Nodes
*/
package gastar

import (
	"cmp"
	"container/heap"
	"slices"
)

type Hasher[C comparable] interface {
	Hash() C
}

type Grapher[C comparable, K Hasher[C], H cmp.Ordered] interface {
	Neighbors(K) []K
	Cost(k1, k2 K) H
	Distance(k1, k2 K) H
}

type DefaultGrapher[C comparable, K Hasher[C], H cmp.Ordered] struct{}

func (*DefaultGrapher[C, K, H]) Cost(n1, n2 K) H {
	var cost H
	return cost
}

func (*DefaultGrapher[C, K, H]) Distance(n1, n2 K) H {
	var distance H
	return distance
}

func NewDefault[C comparable, K Hasher[C], H cmp.Ordered]() Grapher[C, K, H] {
	return &DefaultGrapher[C, K, H]{}
}

func (d *DefaultGrapher[C, K, H]) Neighbors(node K) []K {
	return []K{}
}

func PathFind[C comparable, K Hasher[C], V any, H cmp.Ordered](g Grapher[C, K, H], start, goal K) []K {
	var (
		costSoFar = make(map[C]H)
		visited   = make(map[C]K)
		visiting  = priorityQueue[C, K, H]{}
		thecost   H
	)
	heap.Init(&visiting)
	costSoFar[start.Hash()] = thecost
	visited[start.Hash()] = start
	heap.Push(&visiting, queueNode[C, K, H]{node: start, cost: thecost})
	for visiting.Len() > 0 {
		next := heap.Pop(&visiting).(queueNode[C, K, H])
		node := next.node
		if node.Hash() == goal.Hash() {
			break
		}
		neighbors := g.Neighbors(node)
		for _, neighbor := range neighbors {
			thecost := costSoFar[neighbor.Hash()] + g.Cost(node, neighbor)
			nextcost, ok := costSoFar[neighbor.Hash()]
			if !ok || thecost < nextcost {
				priority := thecost + g.Distance(node, neighbor)
				costSoFar[neighbor.Hash()] = thecost
				heap.Push(&visiting, queueNode[C, K, H]{node: neighbor, cost: priority})
				visited[neighbor.Hash()] = node
			}
		}
	}
	current := goal
	paths := []K{}
	exists := false
	for {
		paths = append(paths, current)
		if current.Hash() == start.Hash() {
			break
		}
		current, exists = visited[current.Hash()]
		if !exists {
			return []K{}
		}
	}
	slices.Reverse(paths)
	return paths
}

type queueNode[C comparable, K Hasher[C], H cmp.Ordered] struct {
	node K
	cost H
}

type priorityQueue[C comparable, K Hasher[C], H cmp.Ordered] []queueNode[C, K, H]

func (pq priorityQueue[C, K, H]) Len() int           { return len(pq) }
func (pq priorityQueue[C, K, H]) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq priorityQueue[C, K, H]) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue[C, K, H]) Push(v any) {
	vv := v.(queueNode[C, K, H])
	*pq = append(*pq, vv)
}

func (pq *priorityQueue[C, K, H]) Pop() any {
	length := len(*pq)
	old := *pq
	v := old[length-1]
	*pq = old[0 : length-1]
	return v
}
