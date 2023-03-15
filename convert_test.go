package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
		t.Errorf("EpsilonClosure(0) failed, set: %v", ac1)
	}

	ac2 := g1.EpsilonClosureSet(map[int]bool{0: true, 4: true})
	// 含有 0 1 2 4
	if !ac2[0] || !ac2[1] || !ac2[2] || !ac2[4] || len(ac2) != 4 {
		t.Errorf("EpsilonClosureSet(0, 4) failed, set: %v", ac2)
	}
}

// Move 函数测试
func TestMove(t *testing.T) {
	g1 := Graph{
		NumOfStates: 5,
		EdgeTable: []*Edge{
			// Edges, some null and the other char
			{0, 1, 0, DriverNull},
			{0, 2, 0, DriverNull},
			{1, 3, 'a', DriverChar},
			{1, 4, 'b', DriverChar},
			{3, 0, 'a', DriverChar},
			{2, 4, 'b', DriverChar},
			{2, 4, 'a', DriverChar},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 3, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 4, StateType: StateMatch, Category: LexemeNull},
		},
	}

	// move 不含自己也不求 epsilon！切记！
	ac1 := g1.Move(1, 'a', DriverChar)
	// 含有 3
	if !ac1[3] || len(ac1) != 1 {
		t.Errorf("Move(1, 'a') failed, set: %v", ac1)
	}

	ac2 := g1.MoveSet(map[int]bool{1: true, 2: true}, 'a', DriverChar)
	// 含有 3 4
	if !ac2[3] || !ac2[4] || len(ac2) != 2 {
		t.Errorf("MoveSet(1, 4, 'a') failed, set: %v", ac2)
	}
}

// 子集构造法测试
func TestSubsetConstruction(t *testing.T) {
	// Dragon Book Fig 3.34
	g1 := Graph{
		GraphId:     0,
		NumOfStates: 11,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 0, DriverType: DriverNull},
			{FromState: 0, NextState: 7, DriverID: 0, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 0, DriverType: DriverNull},
			{FromState: 1, NextState: 4, DriverID: 0, DriverType: DriverNull},
			{FromState: 2, NextState: 3, DriverID: 'a', DriverType: DriverChar},
			{FromState: 3, NextState: 6, DriverID: 0, DriverType: DriverNull},
			{FromState: 4, NextState: 5, DriverID: 'b', DriverType: DriverChar},
			{FromState: 5, NextState: 6, DriverID: 0, DriverType: DriverNull},
			{FromState: 6, NextState: 1, DriverID: 0, DriverType: DriverNull},
			{FromState: 6, NextState: 7, DriverID: 0, DriverType: DriverNull},
			{FromState: 7, NextState: 8, DriverID: 'a', DriverType: DriverChar},
			{FromState: 8, NextState: 9, DriverID: 'b', DriverType: DriverChar},
			{FromState: 9, NextState: 10, DriverID: 'b', DriverType: DriverChar},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 3, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 4, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 5, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 6, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 7, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 8, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 9, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 10, StateType: StateMatch, Category: LexemeNull},
		},
	}
	solution := Graph{
		GraphId:     0,
		NumOfStates: 5,
		EdgeTable: []*Edge{
			// 0-a-1
			{FromState: 0, NextState: 1, DriverID: 'a', DriverType: DriverChar},
			// 0-b-2
			{FromState: 0, NextState: 2, DriverID: 'b', DriverType: DriverChar},
			// 1-a-1
			{FromState: 1, NextState: 1, DriverID: 'a', DriverType: DriverChar},
			// 1-b-3
			{FromState: 1, NextState: 3, DriverID: 'b', DriverType: DriverChar},
			// 2-a-1
			{FromState: 2, NextState: 1, DriverID: 'a', DriverType: DriverChar},
			// 2-b-2
			{FromState: 2, NextState: 2, DriverID: 'b', DriverType: DriverChar},
			// 3-a-1
			{FromState: 3, NextState: 1, DriverID: 'a', DriverType: DriverChar},
			// 3-b-4
			{FromState: 3, NextState: 4, DriverID: 'b', DriverType: DriverChar},
			// 4-a-1
			{FromState: 4, NextState: 1, DriverID: 'a', DriverType: DriverChar},
			// 4-b-2
			{FromState: 4, NextState: 2, DriverID: 'b', DriverType: DriverChar},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 3, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 4, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g2 := g1.SubsetConstruction()
	printNFA(g2)
	result := cmp.Equal(g2, &solution)
	if !result {
		t.Errorf("SubsetConstruction failed")
	}
}
