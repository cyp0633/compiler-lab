package lab1

import (
	"testing"
)

// 测试生成基本 NFA
func TestGenerateBasicNFA(t *testing.T) {
	g := generateBasicNFA(DriverNull, 114514)
	printGraph(g)
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
	printGraph(g3)
	// 结果：6 个状态，6 条边
	if g3.NumOfStates != 4 || len(g3.StateTable) != 4 || len(g3.EdgeTable) != 4 {
		t.Errorf("unionNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g3.NumOfStates, len(g3.StateTable), len(g3.EdgeTable))
	}
}

// NFA 并运算预处理 1 测试
func TestUnionNFAPreprocess(t *testing.T) {
	// 无入边，无出边
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
	unionNFAPreprocess(&g1)
	printGraph(&g1)
	// 仍然三个状态两条边
	if g1.NumOfStates != 3 || len(g1.StateTable) != 3 || len(g1.EdgeTable) != 2 {
		t.Errorf("unionNFAPreprocess failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g1.NumOfStates, len(g1.StateTable), len(g1.EdgeTable))
	}

	// 有入边，无出边
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 1, NextState: 0, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	unionNFAPreprocess(&g2)
	printGraph(&g2)
	// 前面多个状态
	if g2.NumOfStates != 4 || len(g2.StateTable) != 4 || len(g2.EdgeTable) != 4 {
		t.Errorf("unionNFAPreprocess failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g2.NumOfStates, len(g2.StateTable), len(g2.EdgeTable))
	}

	// 无入边，有出边
	g3 := Graph{
		GraphId:     3,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 2, NextState: 1, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	unionNFAPreprocess(&g3)
	printGraph(&g3)
	// 后面多个状态
	if g3.NumOfStates != 4 || len(g3.StateTable) != 4 || len(g3.EdgeTable) != 4 {
		t.Errorf("unionNFAPreprocess failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g3.NumOfStates, len(g3.StateTable), len(g3.EdgeTable))
	}
}

// 出入边计数测试
func TestInOutEdge(t *testing.T) {
	// 无入边，无出边
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
	in, out := g1.inOutEdge()
	if in || out {
		t.Errorf("inOutEdge failed, in: %v, out: %v\nExpected: no in, no out", in, out)
	}

	// 有入边，无出边
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 1, NextState: 0, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	in, out = g2.inOutEdge()
	if !in || out {
		t.Errorf("inOutEdge failed, in: %v, out: %v\nExpected: has in, no out", in, out)
	}

	// 无入边，有出边
	g3 := Graph{
		GraphId:     3,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 2, NextState: 1, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	in, out = g3.inOutEdge()
	if in || !out {
		t.Errorf("inOutEdge failed, in: %v, out: %v\nExpected: no in, has out", in, out)
	}
}

func TestUnionNFAPreprocess2(t *testing.T) {
	// 最后一个状态有 category 属性
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
			{StateID: 2, StateType: StateMatch, Category: LexemeIntegerConst},
		},
	}
	unionNFAPreprocess2(&g1)
	printGraph(&g1)
	if g1.NumOfStates != 4 || len(g1.StateTable) != 4 || len(g1.EdgeTable) != 3 {
		t.Errorf("unionNFAPreprocess2 failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g1.NumOfStates, len(g1.StateTable), len(g1.EdgeTable))
	}
}

func TestProductNFA(t *testing.T) {
	// 无入边，无出边
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
	// 有入边，无出边
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 1, NextState: 0, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	// 无入边，有出边
	g3 := Graph{
		GraphId:     3,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 2, NextState: 1, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}

	// 前一个有出，后一个有入
	g4 := productNFA(&g3, &g2)
	printGraph(g4)
	if g4.NumOfStates != 6 || len(g4.StateTable) != 6 || len(g4.EdgeTable) != 7 {
		t.Errorf("productNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g4.NumOfStates, len(g4.StateTable), len(g4.EdgeTable))
	}

	// 前一个有出，后一个无入
	g5 := productNFA(&g3, &g1)
	printGraph(g5)
	if g5.NumOfStates != 5 || len(g5.StateTable) != 5 || len(g5.EdgeTable) != 5 {
		t.Errorf("productNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g5.NumOfStates, len(g5.StateTable), len(g5.EdgeTable))
	}

	// 前一个无出，后一个有入
	g6 := productNFA(&g1, &g2)
	printGraph(g6)
	if g6.NumOfStates != 5 || len(g6.StateTable) != 5 || len(g6.EdgeTable) != 5 {
		t.Errorf("productNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g6.NumOfStates, len(g6.StateTable), len(g6.EdgeTable))
	}

	// 前一个无出，后一个无入
	g7 := productNFA(&g1, &g1)
	printGraph(g7)
	if g7.NumOfStates != 5 || len(g7.StateTable) != 5 || len(g7.EdgeTable) != 4 {
		t.Errorf("productNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g7.NumOfStates, len(g7.StateTable), len(g7.EdgeTable))
	}
}

// 正闭包测试
func TestPlusClosureNFA(t *testing.T) {
	// 无入边，无出边
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
	g11 := plusClosureNFA(&g1)
	printGraph(g11)
	// 多一条 epsilon 边，状态数不变
	if g11.NumOfStates != 3 || len(g11.StateTable) != 3 || len(g11.EdgeTable) != 3 {
		t.Errorf("plusClosureNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g11.NumOfStates, len(g11.StateTable), len(g11.EdgeTable))
	}

	// 有入边，无出边
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 1, NextState: 0, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g21 := plusClosureNFA(&g2)
	printGraph(g21)
	// 多两条 epsilon 边，前面加一个状态
	if g21.NumOfStates != 4 || len(g21.StateTable) != 4 || len(g21.EdgeTable) != 5 {
		t.Errorf("plusClosureNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g21.NumOfStates, len(g21.StateTable), len(g21.EdgeTable))
	}

	// 无入边，有出边
	g3 := Graph{
		GraphId:     3,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 2, NextState: 1, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g31 := plusClosureNFA(&g3)
	printGraph(g31)
	// 多两条 epsilon 边，后面加一个状态
	if g31.NumOfStates != 4 || len(g31.StateTable) != 4 || len(g31.EdgeTable) != 5 {
		t.Errorf("plusClosureNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g31.NumOfStates, len(g31.StateTable), len(g31.EdgeTable))
	}
}

// Kleene 闭包测试
func TestKleeneClosureNFA(t *testing.T) {
	// 无入边，无出边
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
	g11 := kleeneClosureNFA(&g1)
	printGraph(g11)
	// 多两条 epsilon 边，状态数不变
	if g11.NumOfStates != 3 || len(g11.StateTable) != 3 || len(g11.EdgeTable) != 4 {
		t.Errorf("kleeneClosureNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g11.NumOfStates, len(g11.StateTable), len(g11.EdgeTable))
	}

	// 有入边，无出边
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 1, NextState: 0, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g21 := kleeneClosureNFA(&g2)
	printGraph(g21)
	// 多三条 epsilon 边，前面加一个状态
	if g21.NumOfStates != 4 || len(g21.StateTable) != 4 || len(g21.EdgeTable) != 6 {
		t.Errorf("kleeneClosureNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g21.NumOfStates, len(g21.StateTable), len(g21.EdgeTable))
	}

	// 无入边，有出边
	g3 := Graph{
		GraphId:     3,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 2, NextState: 1, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g31 := kleeneClosureNFA(&g3)
	printGraph(g31)
	// 多三条 epsilon 边，后面加一个状态
	if g31.NumOfStates != 4 || len(g31.StateTable) != 4 || len(g31.EdgeTable) != 6 {
		t.Errorf("kleeneClosureNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g31.NumOfStates, len(g31.StateTable), len(g31.EdgeTable))
	}
}

func TestZeroOrOneNFA(t *testing.T) {
	// 无入边，无出边
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
	g11 := zeroOrOneNFA(&g1)
	printGraph(g11)
	// 多一条 epsilon 边，状态数不变
	if g11.NumOfStates != 3 || len(g11.StateTable) != 3 || len(g11.EdgeTable) != 3 {
		t.Errorf("zeroOrOneNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g11.NumOfStates, len(g11.StateTable), len(g11.EdgeTable))
	}

	// 有入边，无出边
	g2 := Graph{
		GraphId:     2,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 1, NextState: 0, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g21 := zeroOrOneNFA(&g2)
	printGraph(g21)
	// 多两条 epsilon 边，前面加一个状态
	if g21.NumOfStates != 4 || len(g21.StateTable) != 4 || len(g21.EdgeTable) != 5 {
		t.Errorf("zeroOrOneNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g21.NumOfStates, len(g21.StateTable), len(g21.EdgeTable))
	}

	// 无入边，有出边
	g3 := Graph{
		GraphId:     3,
		NumOfStates: 3,
		EdgeTable: []*Edge{
			{FromState: 0, NextState: 1, DriverID: 1, DriverType: DriverNull},
			{FromState: 1, NextState: 2, DriverID: 2, DriverType: DriverNull},
			{FromState: 2, NextState: 1, DriverID: 3, DriverType: DriverNull},
		},
		StateTable: []*State{
			{StateID: 0, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 1, StateType: StateUnmatch, Category: LexemeNull},
			{StateID: 2, StateType: StateMatch, Category: LexemeNull},
		},
	}
	g31 := zeroOrOneNFA(&g3)
	printGraph(g31)
	// 多两条 epsilon 边，后面加一个状态
	if g31.NumOfStates != 4 || len(g31.StateTable) != 4 || len(g31.EdgeTable) != 5 {
		t.Errorf("zeroOrOneNFA failed, NumOfStates: %v, len(StateTable): %v, len(EdgeTable): %v", g31.NumOfStates, len(g31.StateTable), len(g31.EdgeTable))
	}
}
