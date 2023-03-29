package lab2

import "testing"

var leftRecursionExample struct {
	S, T, F               NonTerminalSymbol
	plus, mul, lb, rb, id TerminalSymbol
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
