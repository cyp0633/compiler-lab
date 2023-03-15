package main

import "testing"

func TestEpsilonClosure(t *testing.T) {
	g1 := Graph{
		NumOfStates: 5,
		EdgeTable: []*Edge{
			// Edges, some null and the other char
			{0, 1, 0, DriverNull},
			{0, 2, 0, DriverNull},
			{1, 3, 0, DriverChar},
			{2, 4, 0, DriverChar},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 3, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 4, StateType: StateMatch, Category: LexemeNull},
		},
	}
	ac1 := g1.EpsilonClosure(0)
	// 含有 0 1 2
	if !ac1[0] || !ac1[1] || !ac1[2] || len(ac1) != 3 {
		t.Errorf("EpsilonClosure(0) failed")
	}

	ac2 := g1.EpsilonClosureSet(map[int]bool{0: true, 4: true})
	// 含有 0 1 2 4
	if !ac2[0] || !ac2[1] || !ac2[2] || !ac2[4] || len(ac2) != 4 {
		t.Errorf("EpsilonClosureSet(0, 4) failed")
	}
}
