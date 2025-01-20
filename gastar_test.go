package gastar

import (
	"fmt"
	"math"
	"testing"
)

type jug struct {
	current uint
	cap     uint
}

type jugState struct {
	jug1, jug2 jug
}

func (js jugState) Hash() string {
	return fmt.Sprintf("(%d,%d)(%d,%d)",
		js.jug1.current, js.jug1.cap, js.jug2.current,
		js.jug2.cap)
}

type waterjugs struct {
	Grapher[jugState, jugState, int]
}

func (w waterjugs) Neighbors(js jugState) []jugState {
	next := []jugState{}
	if js.jug1.current > 0 {
		next = append(next, jugState{
			jug1: jug{0, js.jug1.cap},
			jug2: jug{js.jug2.current, js.jug2.cap},
		})
	}
	if js.jug2.current > 0 {
		next = append(next, jugState{
			jug1: jug{js.jug1.current, js.jug1.cap},
			jug2: jug{0, js.jug2.cap},
		})
	}
	if js.jug1.current < js.jug1.cap {
		next = append(next, jugState{
			jug1: jug{js.jug1.cap, js.jug1.cap},
			jug2: jug{js.jug2.current, js.jug2.cap},
		})
	}
	if js.jug2.current < js.jug2.cap {
		next = append(next, jugState{
			jug1: jug{js.jug1.current, js.jug1.cap},
			jug2: jug{js.jug2.cap, js.jug2.cap},
		})
	}
	if js.jug1.current > 0 && js.jug2.current < js.jug2.cap {
		toPour := min(
			uint(math.Abs(float64(js.jug2.cap-js.jug2.current))),
			js.jug1.current)
		next = append(next, jugState{
			jug1: jug{max(0, js.jug1.current-toPour), js.jug1.cap},
			jug2: jug{min(js.jug2.cap, js.jug2.current+toPour), js.jug2.cap},
		})
	}
	if js.jug2.current > 0 && js.jug1.current < js.jug1.cap {
		toPour := min(
			uint(math.Abs(float64(js.jug1.cap-js.jug1.current))),
			js.jug2.current)
		next = append(next, jugState{
			jug1: jug{max(js.jug1.cap, js.jug1.current+toPour), js.jug1.cap},
			jug2: jug{min(js.jug2.cap, js.jug2.current-toPour), js.jug2.cap},
		})

	}
	return next
}

func TestUseGastar(t *testing.T) {
	w := waterjugs{}
	w.Grapher = NewDefault[jugState, jugState, int]()
	empty := jugState{
		jug1: jug{0, 3},
		jug2: jug{0, 5},
	}
	goal := jugState{
		jug1: jug{0, 3},
		jug2: jug{4, 5},
	}
	paths := PathFind[jugState, jugState](w, empty, goal)
	if len(paths) <= 0 {
		t.Error("Expected get result, found nothing")
	}
	t.Log(paths)
	t.Log("paths length:", len(paths))
}

func BenchmarkFindPath(b *testing.B) {
	w := waterjugs{}
	w.Grapher = NewDefault[jugState, jugState, int]()
	empty := jugState{
		jug1: jug{0, 3},
		jug2: jug{0, 5},
	}
	goal := jugState{
		jug1: jug{0, 3},
		jug2: jug{4, 5},
	}
	b.ResetTimer()
	var paths []jugState
	for i := 0; i < b.N; i++ {
		paths = PathFind[jugState, jugState](w, empty, goal)
	}
	_ = paths
}

type knapsackGraph struct {
	Grapher[knapsackItem, knapsackItem, int]
	items []knapsackItem
}

type knapsackItem struct {
	name          string
	capacity      uint
	currentWeight uint
	weight        uint
	value         uint
	currentValue  uint
}

func (ki knapsackItem) Hash() string {
	return fmt.Sprintf("(%s,%d,%d,%d,%d,%d)", ki.name, ki.weight, ki.value,
		ki.capacity, ki.currentWeight, ki.currentValue)
}

const knapsackCap = 50

func (knapsackGraph) Cost(a, b knapsackItem) int {
	if b.weight == 0 {
		return 0
	}
	return int(-(a.currentValue + (b.value / b.weight)))
}

func (kg knapsackGraph) Neighbors(item knapsackItem) []knapsackItem {
	kitems := []knapsackItem{}
	for _, inv := range kg.items {
		if item.currentWeight+inv.weight > item.capacity {
			continue
		}
		newitem := inv
		newitem.currentWeight = item.currentWeight + newitem.weight
		newitem.currentValue = newitem.value + item.currentValue
		kitems = append(kitems, newitem)

	}
	if len(kitems) == 0 {
		kitems = append(kitems, knapsackItem{name: "full", weight: knapsackCap})
	}
	return kitems
}

func TestKnapsack(t *testing.T) {
	k := knapsackGraph{items: []knapsackItem{
		{name: "ransom", capacity: knapsackCap, weight: 10, value: 30},
		{name: "health-kit", capacity: knapsackCap, weight: 20, value: 100},
		{name: "elixir", capacity: knapsackCap, weight: 30, value: 120},
	},
	}

	k.Grapher = NewDefault[knapsackItem, knapsackItem, int]()
	empty := knapsackItem{name: "empty", capacity: knapsackCap}
	goal := knapsackItem{name: "full", weight: knapsackCap}
	paths := PathFind[knapsackItem, knapsackItem](k, empty, goal)
	if len(paths) < 5 {
		t.Errorf("Expected get 5 paths exact, found less: (%d)\n", len(paths))
	}
	const hk = "health-kit"
	if paths[1].name != hk || paths[2].name != hk {
		t.Errorf("Expecting health-kit item for 2 first, got other: (%s, %s)\n",
			paths[1].name, paths[2].name)
	}
	if paths[3].name != "ransom" {
		t.Errorf("Expecting last item is ransom, got other: (%s)\n",
			paths[3].name)
	}
	t.Log(paths)
	t.Log("paths length:", len(paths))
}

func BenchmarkKnapsack(b *testing.B) {
	k := knapsackGraph{items: []knapsackItem{
		{name: "ransom", capacity: knapsackCap, weight: 10, value: 30},
		{name: "health-kit", capacity: knapsackCap, weight: 20, value: 100},
		{name: "elixir", capacity: knapsackCap, weight: 30, value: 120},
	},
	}

	k.Grapher = NewDefault[knapsackItem, knapsackItem, int]()
	empty := knapsackItem{name: "empty", capacity: knapsackCap}
	goal := knapsackItem{name: "full", weight: knapsackCap}
	var paths []knapsackItem
	for i := 0; i < b.N; i++ {
		paths = PathFind[knapsackItem, knapsackItem](k, empty, goal)
	}
	_ = paths
}
