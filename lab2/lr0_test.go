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

func TestItemSetClosure(t *testing.T) {
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
