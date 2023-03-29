package lab2

import (
	"compiler-lab/lab1"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	initializeTest1()
	initLeftRecursionExample()
	initLeftFactorExample()
	initLL1Example()
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

	// E' -> + T E'
	prod2 := testData1.E1.ProductionTable[0]
	first2 := prod2.First()
	// { + }
	if len(first2) != 1 || first2[testData1.plus] != true {
		t.Error("Production First Error, first2:", first2)
	}

	// 这里就不放更多 case 了，[TestNonTerminalFirst] 里面有更多的 case
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

	// E -> T E'
	nt2 := testData1.E
	first2 := nt2.First()
	// { (, id }
	if len(first2) != 2 || first2[testData1.lparen] != true || first2[testData1.id] != true {
		t.Error("NonTerminal First Error, first2:", first2)
	}

	// E' -> + T E' | epsilon
	nt3 := testData1.E1
	first3 := nt3.First()
	// { +, epsilon }
	if len(first3) != 2 || first3[testData1.plus] != true || first3[epsilonSymbol] != true {
		t.Error("NonTerminal First Error, first3:", first3)
	}

	// T -> F T'
	nt4 := testData1.T
	first4 := nt4.First()
	// { (, id }
	if len(first4) != 2 || first4[testData1.lparen] != true || first4[testData1.id] != true {
		t.Error("NonTerminal First Error, first4:", first4)
	}

	// T' -> * F T' | epsilon
	nt5 := testData1.T1
	first5 := nt5.First()
	// 重复，试试会不会出错
	first51 := nt5.First()
	// { *, epsilon }
	if len(first5) != 2 || first5[testData1.mul] != true || first5[epsilonSymbol] != true {
		t.Error("NonTerminal First Error, first5:", first5)
	}
	if len(first51) != 2 || first51[testData1.mul] != true || first51[epsilonSymbol] != true {
		t.Error("NonTerminal First Error, first51:", first51)
	}

}

// 测试所有非终结符的 FOLLOW 函数
func TestAllFollow(t *testing.T) {
	GrammarSymbolTable = append(GrammarSymbolTable, &testData1.E, &testData1.E1, &testData1.T, &testData1.T1, &testData1.F)
	GrammarSymbolTable = append(GrammarSymbolTable, &testData1.lparen, &testData1.rparen, &testData1.plus, &testData1.mul, &testData1.id)
	RootSymbol = &testData1.E
	Follow()

	// FOLLOW(E) = { #, ) }
	follow1 := testData1.E.FollowSet
	if len(follow1) != 2 || follow1[testData1.rparen] != true {
		t.Error("Follow Error, follow1:", follow1)
	}

	// FOLLOW(E') = { #, ) }
	follow2 := testData1.E1.FollowSet
	if len(follow2) != 2 || follow2[testData1.rparen] != true {
		t.Error("Follow Error, follow2:", follow2)
	}

	// FOLLOW(T) = { +, #, ) }
	follow3 := testData1.T.FollowSet
	if len(follow3) != 3 || follow3[testData1.plus] != true || follow3[testData1.rparen] != true {
		t.Error("Follow Error, follow3:", follow3)
	}

	// FOLLOW(T') = { +, #, ) }
	follow4 := testData1.T1.FollowSet
	if len(follow4) != 3 || follow4[testData1.plus] != true || follow4[testData1.rparen] != true {
		t.Error("Follow Error, follow4:", follow4)
	}

	// FOLLOW(F) = { *, +, #, ) }
	follow5 := testData1.F.FollowSet
	if len(follow5) != 4 || follow5[testData1.mul] != true || follow5[testData1.plus] != true || follow5[testData1.rparen] != true {
		t.Error("Follow Error, follow5:", follow5)
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
