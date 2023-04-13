package lab2

import (
	"compiler-lab/lab1"
	"testing"
)

// https://cs.nyu.edu/~gottlieb/courses/2000s/2007-08-spring/compilers/lectures/lecture-07.html
var nyuMainExample struct {
	// E' \to E
	// E \to E + T | T
	// T \to T * F | F
	// F \to ( E ) | id
	E1, E, T, F                                      *NonTerminalSymbol
	Plus, Mul, LeftParenthesis, RightParenthesis, Id *TerminalSymbol
}

var nonSLRExample struct {
	// S' -> S
	// S -> A a | b A c | d c | b d a
	// A -> d
	S, A, S1   *NonTerminalSymbol
	a, b, c, d *TerminalSymbol
}

func initNYUMainExmaple() {
	// nonterminals
	nyuMainExample.E1 = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "E'",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
	}
	nyuMainExample.E = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "E",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
	}
	nyuMainExample.T = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "T",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
	}
	nyuMainExample.F = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "F",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
	}
	// terminals
	nyuMainExample.Plus = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "+",
			Type: Terminal,
		},
		Category: lab1.LexemeNumericOperator,
	}
	nyuMainExample.Mul = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "*",
			Type: Terminal,
		},
		Category: lab1.LexemeNumericOperator,
	}
	nyuMainExample.LeftParenthesis = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "(",
			Type: Terminal,
		},
	}
	nyuMainExample.RightParenthesis = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: ")",
			Type: Terminal,
		},
	}
	nyuMainExample.Id = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "id",
			Type: Terminal,
		},
	}

	// productions
	nyuMainExample.E1.ProductionTable = []*Production{
		{
			ID:         1,
			BodySize:   1,
			BodySymbol: []interface{}{nyuMainExample.E},
		},
	}
	nyuMainExample.E.ProductionTable = []*Production{
		{
			ID:         2,
			BodySize:   3,
			BodySymbol: []interface{}{nyuMainExample.E, nyuMainExample.Plus, nyuMainExample.T},
		},
		{
			ID:         3,
			BodySize:   1,
			BodySymbol: []interface{}{nyuMainExample.T},
		},
	}
	nyuMainExample.T.ProductionTable = []*Production{
		{
			ID:         4,
			BodySize:   3,
			BodySymbol: []interface{}{nyuMainExample.T, nyuMainExample.Mul, nyuMainExample.F},
		},
		{
			ID:         5,
			BodySize:   1,
			BodySymbol: []interface{}{nyuMainExample.F},
		},
	}
	nyuMainExample.F.ProductionTable = []*Production{
		{
			ID:         6,
			BodySize:   3,
			BodySymbol: []interface{}{nyuMainExample.LeftParenthesis, nyuMainExample.E, nyuMainExample.RightParenthesis},
		},
		{
			ID:         7,
			BodySize:   1,
			BodySymbol: []interface{}{nyuMainExample.Id},
		},
	}
}

func initNonSLRExmaple() {
	// nonterminals
	nonSLRExample.S1 = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "S'",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
	}
	nonSLRExample.S = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "S",
			Type: NonTerminal,
		},
		NumOfProduction: 4,
	}
	nonSLRExample.A = &NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "A",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
	}

	// terminals
	nonSLRExample.a = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "a",
			Type: Terminal,
		},
	}
	nonSLRExample.b = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "b",
			Type: Terminal,
		},
	}
	nonSLRExample.c = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "c",
			Type: Terminal,
		},
	}
	nonSLRExample.d = &TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "d",
			Type: Terminal,
		},
	}

	// productions
	// S' -> S
	nonSLRExample.S1.ProductionTable = []*Production{
		{
			ID:         0,
			BodySize:   1,
			BodySymbol: []interface{}{nonSLRExample.S},
		},
	}
	// S -> A a
	nonSLRExample.S.ProductionTable = []*Production{
		{
			ID:         1,
			BodySize:   2,
			BodySymbol: []interface{}{nonSLRExample.A, nonSLRExample.a},
		},
	}
	// S -> b A c
	nonSLRExample.S.ProductionTable = append(nonSLRExample.S.ProductionTable, &Production{
		ID:         2,
		BodySize:   3,
		BodySymbol: []interface{}{nonSLRExample.b, nonSLRExample.A, nonSLRExample.c},
	})
	// S -> d c
	nonSLRExample.S.ProductionTable = append(nonSLRExample.S.ProductionTable, &Production{
		ID:         3,
		BodySize:   2,
		BodySymbol: []interface{}{nonSLRExample.d, nonSLRExample.c},
	})
	// S -> b d a
	nonSLRExample.S.ProductionTable = append(nonSLRExample.S.ProductionTable, &Production{
		ID:         4,
		BodySize:   3,
		BodySymbol: []interface{}{nonSLRExample.b, nonSLRExample.d, nonSLRExample.a},
	})
	// A -> d
	nonSLRExample.A.ProductionTable = []*Production{
		{
			ID:         5,
			BodySize:   1,
			BodySymbol: []interface{}{nonSLRExample.d},
		},
	}
}

func TestItemSetClosure(t *testing.T) {
	// 防止之前的测试影响
	ItemSetTable = []*ItemSet{}

	GrammarSymbolTable = []interface{}{&nyuMainExample.E1, &nyuMainExample.E, &nyuMainExample.T, &nyuMainExample.F, &nyuMainExample.Plus, &nyuMainExample.Mul, &nyuMainExample.LeftParenthesis, &nyuMainExample.RightParenthesis, &nyuMainExample.Id}
	// Closure({E' -> .E})
	set1 := &ItemSet{
		ID: maxItemSetID() + 1,
		ItemTable: map[LR0Item]struct{}{
			{NonTerminalSymbol: nyuMainExample.E1, Production: nyuMainExample.E1.ProductionTable[0], DotPosition: 0}: {},
		},
	}
	ItemSetTable = append(ItemSetTable, set1)
	closure1 := set1.Closure()
	// all 7 elements like ->.
	if len(closure1.ItemTable) != 7 {
		t.Error("closure1 should have 7 elements")
	}
	t.Log(closure1.String())
}

func TestExhaustTransition(t *testing.T) {
	// 防止之前的测试影响
	ItemSetTable = []*ItemSet{}

	GrammarSymbolTable = []interface{}{&nyuMainExample.E1, &nyuMainExample.E, &nyuMainExample.T, &nyuMainExample.F, &nyuMainExample.Plus, &nyuMainExample.Mul, &nyuMainExample.LeftParenthesis, &nyuMainExample.RightParenthesis, &nyuMainExample.Id}
	RootSymbol = nyuMainExample.E1
	// 初始化第一个状态
	// E' \to \cdot E 的闭包
	set1 := &ItemSet{
		ItemTable: map[LR0Item]struct{}{
			{NonTerminalSymbol: nyuMainExample.E1, Production: nyuMainExample.E1.ProductionTable[0], DotPosition: 0, Type: CoreItem}: {},
		},
	}
	set1 = set1.Closure()
	ItemSetTable = append(ItemSetTable, set1)

	// 每个状态的项目数
	// stateItems := []int{7, 2, 2, 1, 7, 1, 5, 3, 2, 2, 1, 1}
	stateItems := map[int]int{7: 2, 5: 1, 3: 1, 2: 4, 1: 4}

	for i := 0; i < len(ItemSetTable); i++ {
		t.Log(ItemSetTable[i].String())
		stateItems[len(ItemSetTable[i].ItemTable)]--
		ItemSetTable[i].ExhaustTransition()
	}
	for k, v := range stateItems {
		if v != 0 {
			t.Errorf("Missing %d states with %d items", v, k)
		}
	}
}

func TestLR0Goto(t *testing.T) {
	// 防止之前的测试影响
	ItemSetTable = []*ItemSet{}

	GrammarSymbolTable = []interface{}{&nyuMainExample.E1, &nyuMainExample.E, &nyuMainExample.T, &nyuMainExample.F, &nyuMainExample.Plus, &nyuMainExample.Mul, &nyuMainExample.LeftParenthesis, &nyuMainExample.RightParenthesis, &nyuMainExample.Id}
	RootSymbol = nyuMainExample.E1
	// 初始化第一个状态
	// E' \to \cdot E 的闭包
	set1 := &ItemSet{
		ItemTable: map[LR0Item]struct{}{
			{NonTerminalSymbol: nyuMainExample.E1, Production: nyuMainExample.E1.ProductionTable[0], DotPosition: 0, Type: CoreItem}: {},
		},
	}
	set1 = set1.Closure()
	goto1 := set1.Goto(nyuMainExample.E)

	t.Log(goto1.String())
	// E' \to E \cdot
	// E \to E \cdot + T
	if len(goto1.ItemTable) != 2 {
		t.Error("goto1 should have 2 elements")
	}
}

// 测试构造 DFA
func TestBuildLR0DFA(t *testing.T) {
	// 防止之前的测试影响
	ItemSetTable = []*ItemSet{}

	GrammarSymbolTable = []interface{}{&nyuMainExample.E1, &nyuMainExample.E, &nyuMainExample.T, &nyuMainExample.F, &nyuMainExample.Plus, &nyuMainExample.Mul, &nyuMainExample.LeftParenthesis, &nyuMainExample.RightParenthesis, &nyuMainExample.Id}
	RootSymbol = nyuMainExample.E1
	// 初始化第一个状态
	// E' \to \cdot E 的闭包
	set1 := &ItemSet{
		ItemTable: map[LR0Item]struct{}{
			{NonTerminalSymbol: nyuMainExample.E1, Production: nyuMainExample.E1.ProductionTable[0], DotPosition: 0, Type: CoreItem}: {},
		},
	}
	set1 = set1.Closure()
	ItemSetTable = append(ItemSetTable, set1)
	for i := 0; i < len(ItemSetTable); i++ {
		t.Log(ItemSetTable[i].String())
		ItemSetTable[i].ExhaustTransition()
	}

	BuildDFA()
	t.Log("Start state:", DFA.StartItemSet.ID)
	for key, value := range DFA.EdgeSet {
		t.Logf("%v -> %v", key, value)
	}
}

// 测试 SLR(1) 检查
func TestCheckSLR1(t *testing.T) {
	// 测试 1：SLR

	// 防止之前的测试影响
	ItemSetTable = []*ItemSet{}

	GrammarSymbolTable = []interface{}{nyuMainExample.E1, nyuMainExample.E, nyuMainExample.T, nyuMainExample.F, nyuMainExample.Plus, nyuMainExample.Mul, nyuMainExample.LeftParenthesis, nyuMainExample.RightParenthesis, nyuMainExample.Id}
	RootSymbol = nyuMainExample.E1
	// 初始化第一个状态
	// E' \to \cdot E 的闭包
	set1 := &ItemSet{
		ItemTable: map[LR0Item]struct{}{
			{NonTerminalSymbol: nyuMainExample.E1, Production: nyuMainExample.E1.ProductionTable[0], DotPosition: 0, Type: CoreItem}: {},
		},
	}
	set1 = set1.Closure()
	ItemSetTable = append(ItemSetTable, set1)
	for i := 0; i < len(ItemSetTable); i++ {
		ItemSetTable[i].ExhaustTransition()
	}

	check1 := CheckSLR1()
	if !check1 {
		t.Error("Example 1 should be SLR1")
	}

	// 测试 2：不是 SLR
	// 防止之前的测试影响
	ItemSetTable = []*ItemSet{}

	GrammarSymbolTable = []interface{}{nonSLRExample.S1, nonSLRExample.S, nonSLRExample.A, nonSLRExample.a, nonSLRExample.b, nonSLRExample.c, nonSLRExample.d}
	RootSymbol = nonSLRExample.S1
	// 初始化第一个状态
	set2 := &ItemSet{
		ItemTable: map[LR0Item]struct{}{
			{NonTerminalSymbol: nonSLRExample.S1, Production: nonSLRExample.S1.ProductionTable[0], DotPosition: 0, Type: CoreItem}: {},
		},
	}
	set2 = set2.Closure()
	ItemSetTable = append(ItemSetTable, set2)
	for i := 0; i < len(ItemSetTable); i++ {
		t.Log(ItemSetTable[i].String())
		ItemSetTable[i].ExhaustTransition()
	}

	check2 := CheckSLR1()
	if check2 {
		t.Error("Example 2 should not be SLR1")
	}

}
