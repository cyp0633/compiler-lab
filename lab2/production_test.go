package lab2

import (
	"compiler-lab/lab1"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	initializeTest1()
	code := m.Run()
	os.Exit(code)
}

var testData1 struct {
	T, E, E1, T1, F               NonTerminalSymbol
	plus, mul, lparen, rparen, id TerminalSymbol
}

func initializeTest1() {
	// 准备数据
	// E -> T E'
	// E' -> + T E' | epsilon
	// T -> F T'
	// T' -> * F T' | epsilon
	// F -> ( E ) | id

	// 非终结符：E E' T T' F
	testData1.T = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "T",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
		ProductionTable: []*Production{},
	}
	testData1.E = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "E",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
		ProductionTable: []*Production{},
	}
	testData1.E1 = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "E'",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
		ProductionTable: []*Production{},
	}
	testData1.T1 = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "T'",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
		ProductionTable: []*Production{},
	}
	testData1.F = NonTerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "F",
			Type: NonTerminal,
		},
		NumOfProduction: 1,
		ProductionTable: []*Production{},
	}

	// 终结符：+ * ( ) id
	testData1.plus = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "+",
			Type: Terminal,
		},
		Category: lab1.LexemeNumericOperator,
	}
	testData1.mul = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "*",
			Type: Terminal,
		},
		Category: lab1.LexemeNumericOperator,
	}
	testData1.lparen = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "(",
			Type: Terminal,
		},
		Category: lab1.LexemeNumericOperator,
	}
	testData1.rparen = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: ")",
			Type: Terminal,
		},
		Category: lab1.LexemeNumericOperator,
	}
	testData1.id = TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "id",
			Type: Terminal,
		},
		Category: lab1.LexemeStringConst,
	}

	// 产生式
	testData1.E.ProductionTable = append(testData1.E.ProductionTable, &Production{
		BodySize:   2,
		BodySymbol: []interface{}{&testData1.T, &testData1.E1},
	})
	testData1.E1.ProductionTable = append(testData1.E1.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&testData1.plus, &testData1.T, &testData1.E1},
	})
	testData1.E1.ProductionTable = append(testData1.E1.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&epsilonSymbol},
	})
	testData1.T.ProductionTable = append(testData1.T.ProductionTable, &Production{
		BodySize:   2,
		BodySymbol: []interface{}{&testData1.F, &testData1.T1},
	})
	testData1.T1.ProductionTable = append(testData1.T1.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&testData1.mul, &testData1.F, &testData1.T1},
	})
	testData1.T1.ProductionTable = append(testData1.T1.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&epsilonSymbol},
	})
	testData1.F.ProductionTable = append(testData1.F.ProductionTable, &Production{
		BodySize:   3,
		BodySymbol: []interface{}{&testData1.lparen, &testData1.E, &testData1.rparen},
	})
	testData1.F.ProductionTable = append(testData1.F.ProductionTable, &Production{
		BodySize:   1,
		BodySymbol: []interface{}{&testData1.id},
	})
}

// 测试类型断言
func TestTypeAssertion(t *testing.T) {
	a := NonTerminalSymbol{}
	s := []interface{}{}
	s = append(s, &a)
	for _, v := range s {
		if _, ok := v.(*NonTerminalSymbol); ok {
			t.Log("ok")
		} else {
			t.Error(reflect.TypeOf(v))
		}
	}
}

// 测试终结符的 First 函数
func TestTerminalFirst(t *testing.T) {
	a := TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "a",
			Type: Terminal,
		},
		Category: lab1.LexemeStringConst,
	}
	// 应为它自己
	f := a.First()
	if f == nil || len(f) != 1 || !f[a] {
		t.Error("Terminal First Error")
	}
}

// 测试产生式的 FIRST 函数
func TestProductionFirst(t *testing.T) {
	// E' -> epsilon
	prod1 := testData1.E1.ProductionTable[1]
	first1 := prod1.First()
	if len(first1) != 1 || first1[epsilonSymbol] != true {
		t.Error("Production First Error, first1:", first1)
	}
}

// 测试非终结符的 FIRST 函数
func TestNonTerminalFirst(t *testing.T) {
	// F -> ( E ) | id
	nt1 := testData1.F
	first1 := nt1.First()
	// { (, id }
	if len(first1) != 2 || first1[testData1.lparen] != true || first1[testData1.id] != true {
		t.Error("NonTerminal First Error, first1:", first1)
	}
}

// 测试两个 String() 函数
func TestPrettyPrint(t *testing.T) {
	// E -> T E'
	prod1 := testData1.E.ProductionTable[0]
	t.Log(prod1.String())
	// E'
	nt1 := testData1.E1
	t.Log(nt1.String())
}
