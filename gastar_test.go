package gastar

import (
	"fmt"
	"math"
	"strings"
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
	Grapher[string, jugState, int]
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
	w.Grapher = NewDefault[string, jugState, int]()
	empty := jugState{
		jug1: jug{0, 3},
		jug2: jug{0, 5},
	}
	goal := jugState{
		jug1: jug{0, 3},
		jug2: jug{4, 5},
	}
	paths := PathFind[string, jugState, jugState](w, empty, goal)
	if len(paths) <= 0 {
		t.Error("Expected get result, found nothing")
	}
	t.Log(paths)
	t.Log("paths length:", len(paths))
}

func BenchmarkFindPath(b *testing.B) {
	w := waterjugs{}
	w.Grapher = NewDefault[string, jugState, int]()
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
		paths = PathFind[string, jugState, jugState](w, empty, goal)
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

func (ki knapsackItem) Hash() knapsackItem {
	return ki
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
	paths := PathFind[knapsackItem, knapsackItem, knapsackItem](k, empty, goal)
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
		paths = PathFind[knapsackItem, knapsackItem, knapsackItem](k, empty, goal)
	}
	_ = paths
}

const (
	tileWidth = 3
	tileTotal = 6
)

type (
	tileItem  int8
	theTiles  [tileTotal]tileItem
	tileGraph struct {
		Grapher[theTiles, theTiles, int]
		theTiles
		maxheight int
	}
)

func (tt theTiles) String() string {
	s := strings.Builder{}
	s.WriteByte('\n')
	for i, t := range tt {
		if i > 0 && i%tileWidth == 0 {
			s.WriteByte('\n')
		}
		s.WriteString(fmt.Sprintf(" %d ", t))
	}
	s.WriteByte('\n')
	return s.String()
}

func (ti theTiles) Hash() theTiles {
	return ti
}

func (tg tileGraph) Cost(t1, t2 theTiles) int {
	return 1
}

func divmod(a int) (int, int) {
	return a / tileWidth, a % tileWidth
}

func (tg tileGraph) Distance(t1, t2 theTiles) int {
	var r int
	for i, p2 := range t2 {
		a1, b1 := divmod(int(t1[i]))
		a2, b2 := divmod(int(p2))
		r += int(math.Abs(float64(a1-a2)) + math.Abs(float64(b1-b2)))
	}
	return r
}

func (tg tileGraph) Neighbors(tiles theTiles) []theTiles {
	restiles := []theTiles{}
	idx := -1
	for i, t := range tiles {
		if t == 0 {
			idx = i
			break
		}
	}
	if idx == -1 {
		return restiles
	}
	zrow, zcol := divmod(idx)
	addresult := func(row, col int) {
		newidx := row * tileWidth
		newidx += col
		if int(newidx) < len(tg.theTiles) && newidx >= 0 {
			newp := tiles
			newp[idx], newp[newidx] = newp[newidx], newp[idx]
			restiles = append(restiles, newp)
		}
	}
	if zcol-1 >= 0 {
		addresult(zrow, zcol-1)
	}
	if zcol+1 <= tileWidth {
		addresult(zrow, zcol+1)
	}
	if zrow-1 >= 0 {
		addresult(zrow-1, zcol)
	}
	if zrow+1 <= tg.maxheight {
		addresult(zrow+1, zcol)
	}
	return restiles
}

func TestSlider(t *testing.T) {
	start := [tileTotal]tileItem{3, 0, 5, 1, 4, 2}
	end := [tileTotal]tileItem{1, 2, 3, 4, 5, 0}
	tg := tileGraph{
		theTiles: [tileTotal]tileItem{3, 0, 5, 1, 4, 2},
	}
	tg.Grapher = NewDefault[theTiles, theTiles, int]()
	tg.maxheight = tileTotal / tileWidth
	if tileTotal%tileWidth != 0 {
		tg.maxheight++
	}
	t.Log("maxheight::", tg.maxheight)
	paths := PathFind[theTiles, theTiles, theTiles](tg, start, end)
	t.Log(paths)
	for _, tt := range paths {
		t.Log(tt)
	}
	path2 := [tileTotal]tileItem{1, 3, 5, 0, 4, 2}
	if paths[2] != path2 {
		t.Errorf("Expected %q, got %q", path2, paths[2])
	}
	if paths[len(paths)-1] != end {
		t.Errorf("Expected %q, got %q", end, paths[len(paths)-1])
	}
}
