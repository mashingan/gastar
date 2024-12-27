# Gastar

Go A* path finding.  
Go implementation related to [graflib](https://github.com/mashingan/graflib).
While graflib intended to be full-fledged graph with its various utilities APIs, gastar is only intended to search
with A*.

# Install

Using `go mod`

```bash
go get github.com/mashingan/gastar
```

# How to use

Since this package doesn't all things about graph, so users need

1. Implement Hasher interface for its Node type.
2. Neighbors method which node's connections.
3. Optionally Cost method to define how costly the changes from Node1 to Node2
4. Optionally, Distance method to define how far the current node to next node in its connections.

## Example

Below is example of solving [water pouring problems](https://en.wikipedia.org/wiki/Water_pouring_puzzle)
with jug 1 has 3 liters and jugs 2 has 5 liters.
What are the paths on how to get 4 liters from jug 2

```go
package main

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
        // state of emptying jug 1
		next = append(next, jugState{
			jug1: jug{0, js.jug1.cap},
			jug2: jug{js.jug2.current, js.jug2.cap},
		})
	}
	if js.jug2.current > 0 {
        // state of emptying jug 2
		next = append(next, jugState{
			jug1: jug{js.jug1.current, js.jug1.cap},
			jug2: jug{0, js.jug2.cap},
		})
	}
	if js.jug1.current < js.jug1.cap {
        // state of filling jug 1 full
		next = append(next, jugState{
			jug1: jug{js.jug1.cap, js.jug1.cap},
			jug2: jug{js.jug2.current, js.jug2.cap},
		})
	}
	if js.jug2.current < js.jug2.cap {
        // state of filling jug 2 full
		next = append(next, jugState{
			jug1: jug{js.jug1.current, js.jug1.cap},
			jug2: jug{js.jug2.cap, js.jug2.cap},
		})
	}
	if js.jug1.current > 0 && js.jug2.current < js.jug2.cap {
        // state of pouring jug 1 to jug 2
		toPour := min(
			uint(math.Abs(float64(js.jug2.cap-js.jug2.current))),
			js.jug1.current)
		next = append(next, jugState{
			jug1: jug{max(0, js.jug1.current-toPour), js.jug1.cap},
			jug2: jug{min(js.jug2.cap, js.jug2.current+toPour), js.jug2.cap},
		})
	}
	if js.jug2.current > 0 && js.jug1.current < js.jug1.cap {
        // state of pouring jug 2 to jug 1
		toPour := min(
			uint(math.Abs(float64(js.jug1.cap-js.jug1.current))),
			js.jug2.current)
		next = append(next, jugState{
			jug1: jug{min(js.jug1.cap, js.jug1.current+toPour), js.jug1.cap},
			jug2: jug{max(js.jug2.cap, js.jug2.current-toPour), js.jug2.cap},
		})

	}
	return next
}

func main() {
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
		panic("Expected get result, found nothing")
	}
	fmt.Println(paths)
	fmt.Println("paths length:", len(paths))
}
```

In above example, we explicitly implement Neighbors method and Hasher interface for node type because both of
those are mandatory. We skipped Cost and Distance method
because both are optionals.  
In this case, no cost and distance for each state's transition
means we searched it by Breath-first Search (BFS).

There are several worked examples in [graflib example](https://github.com/mashingan/graflib/tree/master/examples)

# License

Apache 2.0