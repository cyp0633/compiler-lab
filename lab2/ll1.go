package lab2

import (
	"github.com/google/go-cmp/cmp"
)

// 消除左递归
//
// 使用之前应当先检测左递归
func LeftRecursionElimination() {
	// 遍历非终结符
	for i := 0; i < len(GrammarSymbolTable); i++ {
		ai, ok := GrammarSymbolTable[i].(*NonTerminalSymbol)
		if !ok {
			continue
		}
		for j := 0; j < i; j++ {
			aj, ok := GrammarSymbolTable[j].(*NonTerminalSymbol)
			if !ok {
				continue
			}
			// 找到 ai -> aj \gamma
			for index, production := range ai.ProductionTable {
				if len(production.BodySymbol) > 0 && cmp.Equal(production.BodySymbol[0], aj) {
					// 对于每个 aj -> \delta
					for _, production2 := range aj.ProductionTable {
						// 加入 ai -> \delta \gamma
						newProduction := &Production{
							BodySymbol: append(production2.BodySymbol, production.BodySymbol[1:]...),
						}
						newProduction.BodySize = len(newProduction.BodySymbol)
						ai.ProductionTable = append(ai.ProductionTable, newProduction)
					}
					// 删除 ai -> aj \gamma
					ai.ProductionTable = append(ai.ProductionTable[:index], ai.ProductionTable[index+1:]...)
				}
			}
		}

		// 对 ai 消除直接左递归
		// 没有左递归，跳过
		if !ai.LeftRecursive() {
			continue
		}
		newSymbol := &NonTerminalSymbol{
			GrammarSymbol: GrammarSymbol{
				Name: ai.Name + "'",
				Type: NonTerminal,
			},
			// 自带一个 \epsilon 产生式
			ProductionTable: []*Production{
				{BodySymbol: []interface{}{&epsilonSymbol}, BodySize: 1},
			},
		}
		// 有左递归，加入新的非终结符
		GrammarSymbolTable = append(GrammarSymbolTable, newSymbol)
		// 遍历 ai 的产生式
		// 有删除行为，所以不能使用 range
		for index := 0; index < len(ai.ProductionTable); index++ {
			prod := ai.ProductionTable[index]
			// 对于有左递归的产生式 ai -> ai \alpha，删除，加入 ai' -> \alpha ai'
			if len(prod.BodySymbol) > 0 && cmp.Equal(prod.BodySymbol[0], ai) {
				ai.ProductionTable = append(ai.ProductionTable[:index], ai.ProductionTable[index+1:]...)
				newProduction := &Production{
					BodySymbol: append(prod.BodySymbol[1:], newSymbol),
				}
				newProduction.BodySize = len(newProduction.BodySymbol)
				newSymbol.ProductionTable = append(newSymbol.ProductionTable, newProduction)
				index--
			} else {
				// 对于没有左递归的产生式 ai -> \beta，替换为 ai -> \beta ai'
				newProduction := &Production{
					BodySymbol: append(prod.BodySymbol, newSymbol),
				}
				newProduction.BodySize = len(newProduction.BodySymbol)
				ai.ProductionTable[index] = newProduction
			}
		}
	}
}

// 对单个非终结符检测左递归
func (s *NonTerminalSymbol) LeftRecursive() bool {
	for _, prod := range s.ProductionTable {
		if len(prod.BodySymbol) > 0 && cmp.Equal(prod.BodySymbol[0], s) {
			return true
		}
	}
	return false
}

// 对所有非终结符检测左递归
func CheckLeftRecursion() (ret bool) {
	return checkLeftRecursion(make(map[string]bool), RootSymbol)
}

// 使用 DFS 检测左递归，rec 用于记录已经检测过的非终结符
func checkLeftRecursion(rec map[string]bool, curr *NonTerminalSymbol) bool {
	// 如果已经检测过，返回
	if rec[curr.Name] {
		return true
	}
	// 否则加入 map
	rec[curr.Name] = true

	for _, production := range curr.ProductionTable {
		// 长度为 0，返回
		if len(production.BodySymbol) == 0 {
			continue
		}
		// 是非终结符
		if symbol, ok := production.BodySymbol[0].(*NonTerminalSymbol); ok {
			if checkLeftRecursion(rec, symbol) {
				return true
			}
		}
	}
	return false
}

// 提取左因子
func ExtractLeftFactor() {
}

// 检测左因子
func (s *NonTerminalSymbol) LeftFactored() bool {
	m := map[string]bool{}
	name := ""
	for _, production := range s.ProductionTable {
		for _, symbol := range production.BodySymbol {
			switch symbol := symbol.(type) {
			case *TerminalSymbol:
				name = symbol.Name
			case *NonTerminalSymbol:
				name = symbol.Name
			case *GrammarSymbol:
				name = symbol.Name
			}
			if m[name] {
				return true
			}
			m[name] = true
		}
	}
	return false
}

// 判断是否是 LL(1) 文法
func CheckLL1() bool {
	// 检测左递归
	if CheckLeftRecursion() {
		return false
	}

	// 检测左因子
	for _, symbol := range GrammarSymbolTable {
		if nt, ok := symbol.(*NonTerminalSymbol); ok {
			if nt.LeftFactored() {
				return false
			}
		}
	}

	// 检测 FIRST 集是否有交集
	for _, symbol1 := range GrammarSymbolTable {
		for _, symbol2 := range GrammarSymbolTable {
			if symbol2 == symbol1 {
				continue
			}

			nt1, ok1 := symbol1.(*NonTerminalSymbol)
			nt2, ok2 := symbol2.(*NonTerminalSymbol)
			if !ok1 || !ok2 {
				continue
			}

			result := intersectMaps(nt1.First(), nt2.First())
			if len(result) > 0 {
				return false
			}
		}
	}
	return true
}

func intersectMaps[T comparable, U any](map1 map[T]U, map2 map[T]U) map[T]U {
	result := make(map[T]U)

	for key1 := range map1 {
		if _, ok := map2[key1]; ok {
			result[key1] = map1[key1]
		}
	}

	return result
}
