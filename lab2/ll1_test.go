package lab2

import (
	"fmt"
	"testing"
)

var leftRecursionExample struct {
	S, T, F               NonTerminalSymbol
	plus, mul, lb, rb, id TerminalSymbol
}

var leftFactorExample struct {
	S, E                  NonTerminalSymbol
	If, Then, Else, Other TerminalSymbol
}

var ll1Example struct {
	S, A NonTerminalSymbol
	a, b TerminalSymbol
}

func initLeftRecursionExample() {
	// 	S -> S + T | T
	// T -> T * F | F
	// F -> ( S ) | id

	// nonterminal: S T F
	leftRecursionExample.S = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "S",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
		ProductionTable: []*Production{},
	}
	leftRecursionExample.T = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "T",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
		ProductionTable: []*Production{},
	}
	leftRecursionExample.F = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "F",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
		ProductionTable: []*Production{},
	}

	// terminal: + * ( ) id
	leftRecursionExample.plus = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "+",
			Type: Terminal,
		},
	}
	leftRecursionExample.mul = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "*",
			Type: Terminal,
		},
	}
	leftRecursionExample.lb = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "(",
			Type: Terminal,
		},
	}
	leftRecursionExample.rb = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: ")",
			Type: Terminal,
		},
	}
	leftRecursionExample.id = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "id",
			Type: Terminal,
		},
	}

	// productions
	// S -> S + T | T
	leftRecursionExample.S.ProductionTable = append(leftRecursionExample.S.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&leftRecursionExample.S, &leftRecursionExample.plus, &leftRecursionExample.T},
	})
	leftRecursionExample.S.ProductionTable = append(leftRecursionExample.S.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&leftRecursionExample.T},
	})
	// T -> T * F | F
	leftRecursionExample.T.ProductionTable = append(leftRecursionExample.T.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&leftRecursionExample.T, &leftRecursionExample.mul, &leftRecursionExample.F},
	})
	leftRecursionExample.T.ProductionTable = append(leftRecursionExample.T.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&leftRecursionExample.F},
	})
	// F -> ( S ) | id
	leftRecursionExample.F.ProductionTable = append(leftRecursionExample.F.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&leftRecursionExample.lb, &leftRecursionExample.S, &leftRecursionExample.rb},
	})
	leftRecursionExample.F.ProductionTable = append(leftRecursionExample.F.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&leftRecursionExample.id},
	})
}

func initLeftFactorExample() {
	// S -> if E then S | if E then S else S | Other

	// nonterminal: S E
	leftFactorExample.S = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "S",
			Type: NonTerminal,
		},
		NumOfProduction: 3,
		ProductionTable: []*Production{},
	}
	leftFactorExample.E = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "E",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
		ProductionTable: []*Production{},
	}

	// terminal: if then else Other
	leftFactorExample.If = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "if",
			Type: Terminal,
		},
	}
	leftFactorExample.Then = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "then",
			Type: Terminal,
		},
	}
	leftFactorExample.Else = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "else",
			Type: Terminal,
		},
	}
	leftFactorExample.Other = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "Other",
			Type: Terminal,
		},
	}

	// productions
	// S -> if E then S | if E then S else S | Other
	leftFactorExample.S.ProductionTable = append(leftFactorExample.S.ProductionTable, &Production{
		ID:         0,
		BodySize:   4,
		BodySymbol: []interface{}{&leftFactorExample.If, &leftFactorExample.E, &leftFactorExample.Then, &leftFactorExample.S},
	})
	leftFactorExample.S.ProductionTable = append(leftFactorExample.S.ProductionTable, &Production{
		ID:         1,
		BodySize:   6,
		BodySymbol: []interface{}{&leftFactorExample.If, &leftFactorExample.E, &leftFactorExample.Then, &leftFactorExample.S, &leftFactorExample.Else, &leftFactorExample.S},
	})
	leftFactorExample.S.ProductionTable = append(leftFactorExample.S.ProductionTable, &Production{
		ID:         2,
		BodySize:   1,
		BodySymbol: []interface{}{&leftFactorExample.Other},
	})
}

func initLL1Example() {
	// S -> a A S
	// S -> b
	// A -> b A | \epsilon

	// nonterminal: S A
	ll1Example.S = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "S",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
		ProductionTable: []*Production{},
	}
	ll1Example.A = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "A",
			Type: NonTerminal,
		},
		NumOfProduction: 2,
		ProductionTable: []*Production{},
	}

	// terminal: a b
	ll1Example.a = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "a",
			Type: Terminal,
		},
	}
	ll1Example.b = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "b",
			Type: Terminal,
		},
	}

	// productions
	// S -> a A S
	ll1Example.S.ProductionTable = append(ll1Example.S.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&ll1Example.a, &ll1Example.A, &ll1Example.S},
	})
	// S -> b
	ll1Example.S.ProductionTable = append(ll1Example.S.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&ll1Example.b},
	})
	// A -> b A
	ll1Example.A.ProductionTable = append(ll1Example.A.ProductionTable, &Production{
		BodySize:   2,
		BodySymbol: []interface{}{&ll1Example.b, &ll1Example.A},
	})
	// A -> \epsilon
	ll1Example.A.ProductionTable = append(ll1Example.A.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&epsilonSymbol},
	})
}

// 检测左递归测试
func TestCheckLeftRecursion(t *testing.T) {
	GrammarSymbolTable = []interface{}{&leftRecursionExample.S, &leftRecursionExample.T, &leftRecursionExample.F, &leftRecursionExample.plus, &leftRecursionExample.mul, &leftRecursionExample.lb, &leftRecursionExample.rb, &leftRecursionExample.id}
	RootSymbol = &leftRecursionExample.S
	test1 := CheckLeftRecursion()
	if !test1 {
		t.Error("TestCheckLeftRecursion failed")
	}

	GrammarSymbolTable = []interface{}{&testData1.E, &testData1.E1, &testData1.T, &testData1.T1, &testData1.F, &testData1.lparen, &testData1.rparen, &testData1.plus, &testData1.mul, &testData1.id}
	RootSymbol = &testData1.E
	test2 := CheckLeftRecursion()
	if test2 {
		t.Error("TestCheckLeftRecursion failed")
	}
}

func TestEliminateRecursion(t *testing.T) {
	GrammarSymbolTable = []interface{}{&leftRecursionExample.S, &leftRecursionExample.T, &leftRecursionExample.F, &leftRecursionExample.plus, &leftRecursionExample.mul, &leftRecursionExample.lb, &leftRecursionExample.rb, &leftRecursionExample.id}
	RootSymbol = &leftRecursionExample.S

	LeftRecursionElimination()

	if CheckLeftRecursion() {
		t.Error("TestEliminateRecursion failed")
	}
}

// 检测左因子测试
func TestCheckLeftFactor(t *testing.T) {
	test1 := leftFactorExample.S.LeftFactored()
	if !test1 {
		t.Error("TestCheckLeftFactor failed")
	}

	test2 := testData1.E.LeftFactored()
	if test2 {
		t.Error("TestCheckLeftFactor failed")
	}
}

func TestCheckLL1(t *testing.T) {
	GrammarSymbolTable = []interface{}{&ll1Example.S, &ll1Example.A, &ll1Example.a, &ll1Example.b}
	RootSymbol = &ll1Example.S

	if CheckLL1() {
		t.Error("TestCheckLL1 failed")
	}
}

// 测试生成 LL1 分析表
func TestBuildLL1AnalysisTable(t *testing.T) {
	// fill testData1
	GrammarSymbolTable = []interface{}{&testData1.E, &testData1.E1, &testData1.T, &testData1.T1, &testData1.F, &testData1.lparen, &testData1.rparen, &testData1.plus, &testData1.mul, &testData1.id}
	RootSymbol = &testData1.E
	Follow()

	BuildLL1AnalysisTable()
	for key, value := range LL1AnalysisTable {
		fmt.Printf("Symbol %v, input %v = production %v\n", key.NonTerminalSymbol.Name, key.string, value)
	}

	if len(LL1AnalysisTable) != 13 {
		t.Errorf("Length wanted 13, got %v", len(LL1AnalysisTable))
	}

	// E, id = E -> T E'
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.E, "id"}] != testData1.E.ProductionTable[0] {
		t.Error("E, id = E -> T E' not found")
	}
	// E, ( = E -> T E'
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.E, "("}] != testData1.E.ProductionTable[0] {
		t.Error("E, ( = E -> T E' not found")
	}
	// E', + = E' -> + T E'
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.E1, "+"}] != testData1.E1.ProductionTable[0] {
		t.Error("E', + = E' -> + T E' not found")
	}
	// E', ) = E' -> \epsilon
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.E1, ")"}] != testData1.E1.ProductionTable[1] {
		t.Error("E', ) = E' -> \\epsilon not found")
	}
	// E', # = E' -> \epsilon
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.E1, "#"}] != testData1.E1.ProductionTable[1] {
		t.Error("E', # = E' -> \\epsilon not found")
	}
	// T, id = T -> F T'
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.T, "id"}] != testData1.T.ProductionTable[0] {
		t.Error("T, id = T -> F T' not found")
	}
	// T, ( = T -> F T'
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.T, "("}] != testData1.T.ProductionTable[0] {
		t.Error("T, ( = T -> F T' not found")
	}
	// T', + = T' -> \epsilon
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.T1, "+"}] != testData1.T1.ProductionTable[1] {
		t.Error("T', + = T' -> \\epsilon not found")
	}
	// T', * = T' -> * F T'
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.T1, "*"}] != testData1.T1.ProductionTable[0] {
		t.Error("T', * = T' -> * F T' not found")
	}
	// T', ) = T' -> \epsilon
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.T1, ")"}] != testData1.T1.ProductionTable[1] {
		t.Error("T', ) = T' -> \\epsilon not found")
	}
	// T', # = T' -> \epsilon
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.T1, "#"}] != testData1.T1.ProductionTable[1] {
		t.Error("T', # = T' -> \\epsilon not found")
	}
	// F, id = F -> id
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.F, "id"}] != testData1.F.ProductionTable[1] {
		t.Error("F, id = F -> id not found")
	}
	// F, ( = F -> ( E )
	if LL1AnalysisTable[struct {
		*NonTerminalSymbol
		string
	}{&testData1.F, "("}] != testData1.F.ProductionTable[0] {
		t.Error("F, ( = F -> ( E ) not found")
	}
}

// 测试提取左因子
func TestExtractLeftFactor(t *testing.T) {
	GrammarSymbolTable = []interface{}{&leftFactorExample.S, &leftFactorExample.E, &leftFactorExample.If, &leftFactorExample.Else, &leftFactorExample.Other}
	RootSymbol = &leftFactorExample.S

	if !leftFactorExample.S.LeftFactored() {
		t.Error("TestExtractLeftFactor failed")
	}

	ExtractLeftFactor()

	for _, symbol := range GrammarSymbolTable {
		if symbol, ok := symbol.(*NonTerminalSymbol); ok {
			t.Log(symbol.String())
		}
	}

	if leftFactorExample.S.LeftFactored() {
		t.Error("TestExtractLeftFactor failed")
	}
}
