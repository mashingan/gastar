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
	// return js.jug1.current + 10*js.jug1.cap +
	// 	js.jug2.current + 10*js.jug2.cap
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
