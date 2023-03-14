package main

import (
	"fmt"
	"testing"
)

// 测试生成基本 NFA
func TestGenerateBasicNFA(t *testing.T) {
	g := generateBasicNFA(DriverNull, 114514)
	printNFA(g)
	if g.NumOfStates != 2 || g.EdgeTable[0].DriverID != 114514 || g.EdgeTable[0].DriverType != DriverNull {
		t.Error("generateBasicNFA failed")
	}
}

// 测试并运算
func TestUnionNFA(t *testing.T) {
	// 两个均无入边，无出边
	g1 := Graph{
		GraphId:     1,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 3, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 4, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g3 := unionNFA(&g1, &g2)
	printNFA(g3)
	// 结果：6 个状态，6 条边
	if g3.NumOfStates != 6 || len(g3.StateTable) != 6 || len(g3.EdgeTable) != 6 {
		t.Errorf("unionNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g3.NumOfStates, len(g3.StateTable), len(g3.EdgeTable))
	}
}

func TestProductNFA(t *testing.T) {
	_ = productNFA(&Graph{}, &Graph{})
}

func TestPlusClosureNFA(t *testing.T) {
	_ = plusClosureNFA(&Graph{})
}

func TestKleeneClosureNFA(t *testing.T) {
	_ = kleeneClosureNFA(&Graph{})
}

func TestZeroOrOneNFA(t *testing.T) {
	_ = zeroOrOneNFA(&Graph{})
}

func printNFA(g *Graph) {
	println("graph", g.GraphId)
	println("numOfStates", g.NumOfStates)
	for _, edge := range g.EdgeTable {
		fmt.Printf("Edge: from #%v, to #%v, driver %v, type %v\n", edge.FromState, edge.NextState, edge.DriverID, edge.DriverType)
	}
	for _, state := range g.StateTable {
		fmt.Printf("State: #%v, type %v, category %v\n", state.StateID, state.StateType, state.Category)
	}
}
