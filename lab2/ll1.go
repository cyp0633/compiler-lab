package lab2

// 消除左递归
func LeftRecursionElimination() {

}

// 对所有非终结符检测左递归
func CheckLeftRecursion() (ret bool) {
	return checkLeftRecursion(make(map[string]bool), RootSymbol)
}

// 使用 DFS 检测左递归，rec 用于记录已经检测过的非终结符
func checkLeftRecursion(rec map[string]bool, curr *NonTerminalSymbol) (ret bool) {
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
			ret = ret || checkLeftRecursion(rec, symbol)
		}
	}
	return
}

// 提取左因子
func ExtractLeftFactor() {
}

// 检测左因子
func CheckLeftFactor() (ret bool) {
	return
}
