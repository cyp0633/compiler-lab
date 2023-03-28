package lab2

import (
	"compiler-lab/lab1"
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

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
	ID         int                     // 产生式编号
	BodySize   int                     // 该产生式的文法符个数
	BodySymbol []interface{}           // 该产生式的文法符表
	FirstSet   map[TerminalSymbol]bool // 该产生式的 First 函数值
}

// 文法符表
var GrammarSymbolTable []interface{}

// 开始符
var RootSymbol *NonTerminalSymbol

// 语法分析表项（格子）
type Cell struct {
	NonTerminalSymbol *NonTerminalSymbol // 非终结符
	TerminalSymbol    *TerminalSymbol    // 终结符
	Production        *Production        // 产生式
}

// 语法分析表
var ParseTable []*Cell

// epsilon
// 看起来要用好多次，就先定义好了
var epsilonSymbol = TerminalSymbol{
	GrammarSymbol: GrammarSymbol{
		Name: "epsilon",
		Type: Null,
	},
}

// 产生式的 FIRST 函数
func (p *Production) First() map[TerminalSymbol]bool {
	var symbol interface{}
	var index int

	// 只有 epsilon，就直接返回
	if p.BodySize == 1 && cmp.Equal(p.BodySymbol[0], &epsilonSymbol) {
		return map[TerminalSymbol]bool{epsilonSymbol: true}
	}

	// 遍历整个产生式的文法符，找到第一个非终结符
	for _, symbol = range p.BodySymbol {
		// s 是非终结符
		if st, ok := symbol.(*NonTerminalSymbol); ok {
			sf := st.First()
			// 将非终结符的 First 函数值加入该非终结符的 First 函数值
			for k, v := range sf {
				p.FirstSet[k] = v
			}
			break
		}
	}

	// 再次遍历产生式的文法符，找到第一个不可以推导出 epsilon 的非终结符
	for index, symbol = range p.BodySymbol {
		sf := First(symbol)
		// 如果该非终结符的 First 函数值中包含 epsilon，则继续遍历
		if _, ok := sf[epsilonSymbol]; ok {
			continue
		}
		// 否则将该非终结符的 First 函数值加入该非终结符的 First 函数值
		for k, v := range sf {
			p.FirstSet[k] = v
		}
		break
	}

	// 如果上次遍历完发现全是空，就加入 epsilon
	if index == len(p.BodySymbol)-1 {
		p.FirstSet[epsilonSymbol] = true
	}
	return p.FirstSet
}

// 非终结符的 FIRST 函数
func (nt *NonTerminalSymbol) First() map[TerminalSymbol]bool {
	// 如果已经生成过了，就不要再生成了
	if nt.FirstSet != nil && len(nt.FirstSet) != 0 {
		return nt.FirstSet
	}

	// 寻找 epsilon 的产生式（仅含 epsilon）
	for _, p := range nt.ProductionTable {
		if p.BodySize == 1 && p.BodySymbol[0].(*GrammarSymbol).Type == Null {
			// 如果存在将 epsilon 加入该非终结符的 First 函数值
			nt.FirstSet[epsilonSymbol] = true
		}
	}

	// 对每个产生式，调用产生式的 FIRST 函数并合并
	for _, production := range nt.ProductionTable {
		pf := production.First()
		for k, v := range pf {
			nt.FirstSet[k] = v
		}
	}
	return nt.FirstSet
}

// 终结符的 FIRST 函数
//
// 其实就是它自己
func (t *TerminalSymbol) First() (m map[TerminalSymbol]bool) {
	m = make(map[TerminalSymbol]bool)
	m[*t] = true
	return
}

// 所有语法符的 FIRST 函数
//
// 根据类型推导调用不同的 FIRST
func First(s interface{}) (m map[TerminalSymbol]bool) {
	switch s := s.(type) {
	case *Production:
		return s.First()
	case *NonTerminalSymbol:
		return s.First()
	case *TerminalSymbol:
		return s.First()
	case *GrammarSymbol:
		m = make(map[TerminalSymbol]bool)
		m[epsilonSymbol] = true
		return
	default:
		panic("Unknown type")
	}
}

// 对所有非终结符求 FOLLOW 集合
func Follow() {
	// 找到初始符，加入 #
	RootSymbol.FollowSet[TerminalSymbol{
		GrammarSymbol: GrammarSymbol{
			Name: "#",
			Type: Terminal,
		},
	}] = true

	// 循环，直到 FOLLOW 都不变
	changed := true
	for changed {
		// 遍历每个非终结符，寻找 A
		for _, A := range GrammarSymbolTable {
			// 如果不是非终结符，跳过
			if reflect.TypeOf(A) != reflect.TypeOf(RootSymbol) {
				continue
			}

			// 对非终结符 A，遍历每个产生式
			for _, production := range A.(*NonTerminalSymbol).ProductionTable {
				// 对每个产生式，找出文法符中的非终结符 B
				for index, B := range production.BodySymbol {
					// 如果不是非终结符，跳过
					if reflect.TypeOf(B) != reflect.TypeOf(RootSymbol) {
						continue
					}

					// 先看成 A \Rightarrow \alpha B \beta，求 FIRST(\beta)
					// 将 \beta 部分组合成一个产生式
					tempProduction := Production{
						BodySymbol: production.BodySymbol[index+1:],
						BodySize:   len(production.BodySymbol[index+1:]),
					}
					betaFirst := tempProduction.First()

					// 如果 FIRST(\beta) 中包含 epsilon，或者 \beta 为空
					// 为 \alpha B 的形式
					if _, ok := betaFirst[epsilonSymbol]; ok || index == len(production.BodySymbol)-1 {
						// 则将 FOLLOW(A) 加入 FOLLOW(B)
						l := len(B.(*NonTerminalSymbol).FollowSet)
						for k, v := range A.(*NonTerminalSymbol).FollowSet {
							B.(*NonTerminalSymbol).FollowSet[k] = v
						}
						// 监测长度变化，以判断添加
						if !changed && l != len(B.(*NonTerminalSymbol).FollowSet) {
							changed = true
						}
					} else {
						// 否则为 \alpha B \beta 的形式
						// 将 FIRST(\beta)-\epsilon 加入 FOLLOW(B)
						for k, v := range betaFirst {
							l := len(B.(*NonTerminalSymbol).FollowSet)
							if !cmp.Equal(k, epsilonSymbol) {
								B.(*NonTerminalSymbol).FollowSet[k] = v
							}
							// 监测长度变化，以判断添加
							if !changed && l != len(B.(*NonTerminalSymbol).FollowSet) {
								changed = true
							}
						}
					}
				}
			}
		}
	}
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

func (s *NonTerminalSymbol) String() (str string) {
	str += fmt.Sprintf("Nonterminal %v", s.GrammarSymbol.Name)
	for _, p := range s.ProductionTable {
		str += "\n" + p.String()
	}
	str += fmt.Sprintf("\nFirst: %v\nFollow: %v", s.FirstSet, s.FollowSet)
	return
}

func (p *Production) String() (str string) {
	str += "->"
	for _, v := range p.BodySymbol {
		v := v.(*GrammarSymbol)
		str += v.Name + " "
	}
	str += fmt.Sprintf("(ID: %v)", p.ID)
	return
}
