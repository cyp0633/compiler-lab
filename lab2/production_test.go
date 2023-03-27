package lab2

import (
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
