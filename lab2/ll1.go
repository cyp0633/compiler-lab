package lab2

import (
	"github.com/google/go-cmp/cmp"
)

// LL1 分析表
var LL1AnalysisTable = map[struct {
	*NonTerminalSymbol
	string
}]*Production{}

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
// 这个不用先检测
func ExtractLeftFactor() {
	// 对每个非终结符提取左因子
	for _, symbol := range GrammarSymbolTable {
		symbol, ok := symbol.(*NonTerminalSymbol)
		if !ok {
			continue
		}

		// 循环直到没有左因子
		for symbol.LeftFactored() {
			// 统计每个左因子的产生式数量
			// map 的 key 是文法符指针
			m := map[interface{}][]*Production{}
			for _, production := range symbol.ProductionTable {
				symbol := production.BodySymbol[0]
				m[symbol] = append(m[symbol], production)
			}

			// 找出数量最多的左因子
			var maxLeftFactor interface{}
			maxCount := 0
			for leftFactor, productions := range m {
				if len(productions) > maxCount {
					maxLeftFactor = leftFactor
					maxCount = len(productions)
				}
			}
			// 要扩展的左因子的产生式 slice
			prods := m[maxLeftFactor]

			// 生成新的非终结符
			newSymbol := &NonTerminalSymbol{
				GrammarSymbol: GrammarSymbol{
					Name: symbol.Name + "'",
					Type: NonTerminal,
				},
				NumOfProduction: maxCount,
				ProductionTable: make([]*Production, maxCount),
			}
			GrammarSymbolTable = append(GrammarSymbolTable, newSymbol)

			// 在有公因子的产生式数量不变的情况下，试着扩展左因子长度
			extendFactor := []interface{}{maxLeftFactor}
		extend:
			for {
				// 本轮扩展使用的左因子
				if len(prods[0].BodySymbol) < len(extendFactor)+1 {
					break extend
				}
				candidate := prods[0].BodySymbol[len(extendFactor)]

				// 检查每个产生式是否也能扩展
				for _, production := range prods {
					if len(production.BodySymbol) < len(extendFactor)+1 || !cmp.Equal(production.BodySymbol[len(extendFactor)], candidate) {
						break extend
					}
				}

				// 一直没 break，那就是可以扩展
				extendFactor = append(extendFactor, candidate)
			}

			// 生成新的产生式
			for index, production := range prods {
				newProduction := &Production{
					BodySymbol: production.BodySymbol[len(extendFactor):],
					BodySize:   production.BodySize - len(extendFactor),
					ID:         production.ID, // 可以循环利用，反正这个后面会删掉
				}
				if newProduction.BodySize == 0 {
					newProduction.BodySymbol = []interface{}{&epsilonSymbol}
				}
				newSymbol.ProductionTable[index] = newProduction
			}

			// 调整旧的文法符
			index := 0
			// 因为 prods 是按照 ProductionTable 的顺序排列的
			// 所以可以直接用一个 index 一起
			for originalIndex := 0; originalIndex < len(symbol.ProductionTable); originalIndex++ {
				if index >= len(prods) {
					break
				}
				production := symbol.ProductionTable[originalIndex]
				if !cmp.Equal(production, prods[index]) {
					continue
				}
				// 直接删掉
				symbol.ProductionTable = append(symbol.ProductionTable[:originalIndex], symbol.ProductionTable[originalIndex+1:]...)
				originalIndex--
				index++
			}
			// 添加一个 A \to \alpha A'
			symbol.ProductionTable = append(symbol.ProductionTable, &Production{
				BodySymbol: append(extendFactor, newSymbol),
				BodySize:   len(extendFactor) + 1,
				ID:         maxProdID() + 1,
			})
			symbol.NumOfProduction = len(symbol.ProductionTable)
		}
	}
}

// 检测左因子
func (s *NonTerminalSymbol) LeftFactored() bool {
	m := map[string]bool{}
	name := ""
	for _, production := range s.ProductionTable {
		symbol := production.BodySymbol[0]
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

	// 检查 FOLLOW 是否有交集
	for _, symbol := range GrammarSymbolTable {
		symbol, ok := symbol.(*NonTerminalSymbol)
		if !ok {
			continue
		}
		set := make(map[interface{}]bool)
		for _, production := range symbol.ProductionTable {
			for k, v := range production.Select(symbol) {
				if set[k] && v {
					return false
				}
				set[k] = v
			}
		}
	}
	return true
}

// SELECT 集
func (p *Production) Select(s *NonTerminalSymbol) map[interface{}]bool {
	// \alpha 为空
	if p.BodySize == 1 && p.BodySymbol[0] == &epsilonSymbol {
		if s.FollowSet == nil {
			Follow()
		}
		// FOLLOW(A) \cup FIRST(\alpha)
		set := make(map[interface{}]bool)
		for k, v := range s.FollowSet {
			set[k] = v
		}
		for k, v := range p.First() {
			set[k] = v
		}
		return set
	} else {
		return p.First()
	}
}

// 构造 LL(1) 分析表
func BuildLL1AnalysisTable() {
	// 遍历非终结符 A
	for _, symbol := range GrammarSymbolTable {
		symbol, ok := symbol.(*NonTerminalSymbol)
		if !ok {
			continue
		}

		// 遍历产生式 A -> \alpha
		for _, production := range symbol.ProductionTable {
			// 对 a \in FIRST(\alpha)
			for a := range production.First() {
				switch a := a.(type) {
				case *GrammarSymbol:
					// 如果 a 是 epsilon
					if a.Type == Null {
						// 对 b \in FOLLOW(A)
						for b := range symbol.FollowSet {
							var name string
							switch b := b.(type) {
							case *TerminalSymbol:
								name = b.Name
							case *GrammarSymbol:
								// 不考虑空输入
								if b == &epsilonSymbol {
									continue
								}
								name = b.Name
							default:
								panic("unknown type")
							}
							// M[A,b] = A -> \alpha
							LL1AnalysisTable[struct {
								*NonTerminalSymbol
								string
							}{symbol, name}] = production
						}
					}
				// 如果 a 是终结符，M[A,a] = A -> \alpha
				case *TerminalSymbol:
					LL1AnalysisTable[struct {
						*NonTerminalSymbol
						string
					}{symbol, a.Name}] = production
				// FIRST 里应该不会有其他东西吧？
				default:
					panic("unknown type")
				}
			}
		}
	}
	// 其他情况，直接检测 key 是否存在即可
}
