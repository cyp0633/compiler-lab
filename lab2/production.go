package lab2

import "compiler-lab/lab1"

type Symbol interface{}

// 文法符
type GrammarSymbol struct {
	Name string     // 符号名
	Type symbolType // 符号类型
}

type symbolType int

const (
	NonTerminal symbolType = iota // 非终结符
	Terminal                      // 终结符
	Null                          // 空，epsilon
)

// 终结符
type TerminalSymbol struct {
	GrammarSymbol
	Category lab1.LexemeCategory // 终结符类别
}

// 非终结符
type NonTerminalSymbol struct {
	GrammarSymbol
	ProductionTable      []*Production               // 非终结符的产生式表
	NumOfProduction      int                         // 产生式数量
	FirstSet             map[TerminalSymbol]bool     // 该非终结符的 First 函数值
	FollowSet            map[TerminalSymbol]bool     // 该非终结符的 Follow 函数值
	DependentSetInFollow map[*NonTerminalSymbol]bool // 该非终结符的 Follow 函数中依赖的非终结符
}

// 产生式
type Production struct {
	ID         int                      // 产生式编号
	BodySize   int                      // 该产生式的文法符个数
	BodySymbol []interface{}            // 该产生式的文法符表
	FirstSet   map[*TerminalSymbol]bool // 该产生式的 First 函数值
}

// 文法符表
var grammarSymbolTable []*GrammarSymbol

// 开始符
var rootSymbol *NonTerminalSymbol

// 语法分析表项（格子）
type Cell struct {
	NonTerminalSymbol *NonTerminalSymbol // 非终结符
	TerminalSymbol    *TerminalSymbol    // 终结符
	Production        *Production        // 产生式
}

// 语法分析表
var parseTable []*Cell

// 终结符的 FIRST 函数
//
// 其实就是它自己
func (t *TerminalSymbol) First() (m map[TerminalSymbol]bool) {
	m = make(map[TerminalSymbol]bool)
	m[*t] = true
	return
}
func (t symbolType) String() string {
	switch t {
	case NonTerminal:
		return "Non-terminal"
	case Terminal:
		return "Terminal"
	case Null:
		return "Null"
	default:
		return "Unknown"
	}
}

func (s *GrammarSymbol) Self() *GrammarSymbol {
	return s
}
