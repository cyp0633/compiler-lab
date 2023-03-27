package lab2

import (
	"compiler-lab/lab1"
	"reflect"
	"testing"
)

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
