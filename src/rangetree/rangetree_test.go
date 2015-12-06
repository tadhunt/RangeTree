package rangetree

import (
	"log"
	"testing"
)

func TestNewRangeTree(t *testing.T) {
	rt := NewRangeTree()

	if rt == nil {
		t.Fatalf("NewRangeTree returned nil")
	}
}

type RangeTest struct {
	Input    Range
	Expected []Range
}

func TestAddRange(t *testing.T) {
	tests := []RangeTest{
		{Input: Range{500, 1000}, Expected: []Range{{500, 1000}}},

		{Input: Range{300, 400}, Expected: []Range{{300, 400}, {500, 1000}}},  // disjoint left
		{Input: Range{499, 501}, Expected: []Range{{300, 400}, {499, 1000}}},  // extend left
		{Input: Range{498, 1001}, Expected: []Range{{300, 400}, {498, 1001}}}, // extend left and right

		{Input: Range{402, 410}, Expected: []Range{{300, 400}, {402, 410}, {498, 1001}}},  // disjoint right
		{Input: Range{498, 1001}, Expected: []Range{{300, 400}, {402, 410}, {498, 1001}}}, // fully subsumed (noop)
		{Input: Range{498, 1010}, Expected: []Range{{300, 400}, {402, 410}, {498, 1010}}}, // extend right

		{Input: Range{0, 10000}, Expected: []Range{{0, 10000}}}, // overlap all

	}

	rt := NewRangeTree()
	for i, v := range tests {
		log.Printf("AddRange (%d, %d)\n", v.Input.Start, v.Input.End)
		rt.AddRange(v.Input.Start, v.Input.End)
		rt.Dump(false)
		log.Printf("Run check %v\n", v)
		if !rt.Check(v.Expected) {
			log.Printf("FAILED\n")
			t.Fatalf("Test %d failed", i)
		}
	}

	rt.Dump(true)
}

func TestAdjacent(t *testing.T) {
	tests := []RangeTest{
		/*
		 * make a bunch of disjoint regions
		 */
		{Input: Range{50, 60}, Expected: []Range{{50, 60}}},
		{Input: Range{30, 40}, Expected: []Range{{30, 40}, {50, 60}}},
		{Input: Range{70, 80}, Expected: []Range{{30, 40}, {50, 60}, {70, 80}}},
		{Input: Range{10, 20}, Expected: []Range{{10, 20}, {30, 40}, {50, 60}, {70, 80}}},
		{Input: Range{95, 99}, Expected: []Range{{10, 20}, {30, 40}, {50, 60}, {70, 80}, {95, 99}}},
		{Input: Range{90, 93}, Expected: []Range{{10, 20}, {30, 40}, {50, 60}, {70, 80}, {90, 93}, {95, 99}}},

		/*
		 * Now cause them to be connected together by filling in non-overalapping holes
		 */
		{Input: Range{41, 49}, Expected: []Range{{10, 20}, {30, 60}, {70, 80}, {90, 93}, {95, 99}}},
		{Input: Range{61, 69}, Expected: []Range{{10, 20}, {30, 80}, {90, 93}, {95, 99}}},
	}

	rt := NewRangeTree()
	for i, v := range tests {
		log.Printf("AddRange Test %d (%d, %d)\n", i, v.Input.Start, v.Input.End)
		rt.AddRange(v.Input.Start, v.Input.End)
		rt.Dump(false)
		log.Printf("Run check %v\n", v)
		if !rt.Check(v.Expected) {
			log.Printf("FAILED\n")
			t.Fatalf("Test %d failed", i)
		}
	}

	rt.Dump(true)
}

func TestHasRange(t *testing.T) {
	tests := []RangeTest{
		/*
		 * make a bunch of disjoint regions
		 */
		{Input: Range{50, 60}, Expected: []Range{{50, 60}}},
		{Input: Range{30, 40}, Expected: []Range{{30, 40}, {50, 60}}},
		{Input: Range{70, 80}, Expected: []Range{{30, 40}, {50, 60}, {70, 80}}},
		{Input: Range{10, 20}, Expected: []Range{{10, 20}, {30, 40}, {50, 60}, {70, 80}}},
		{Input: Range{95, 99}, Expected: []Range{{10, 20}, {30, 40}, {50, 60}, {70, 80}, {95, 99}}},
		{Input: Range{90, 93}, Expected: []Range{{10, 20}, {30, 40}, {50, 60}, {70, 80}, {90, 93}, {95, 99}}},
		{Input: Range{41, 49}, Expected: []Range{{10, 20}, {30, 60}, {70, 80}, {90, 93}, {95, 99}}},
		{Input: Range{61, 69}, Expected: []Range{{10, 20}, {30, 80}, {90, 93}, {95, 99}}},
	}

	rt := NewRangeTree()
	for _, v := range tests {
		rt.AddRange(v.Input.Start, v.Input.End)
	}

	if rt.HasRange(0, 9) {
		t.Errorf("shouldn't have range")
	}
	if !rt.HasRange(30, 80) {
		t.Errorf("missing range")
	}

	if rt.HasRange(27,30) {
		t.Errorf("shouldn't have range")
	}
	if rt.HasRange(80,85) {
		t.Errorf("shouldn't have range")
	}

	if rt.HasRange(100,100) {
		t.Errorf("shouldn't have range")
	}
}
