package lab2

// LR(0) 项
type LR0Item struct {
	NonTerminalSymbol *NonTerminalSymbol // 非终结符
	Production        *Production        // 产生式
	DotPosition       int                // 点的位置
	Type              itemCategory
}

// LR(0) 项的类型
type itemCategory int

const (
	//核心项
	CoreItem itemCategory = iota
	//非核心项
	NonCoreItem
)

// LR(0) 项集
type ItemSet struct {
	// 状态序号
	ID int
	// LR(0) 项目表
	ItemTable []*LR0Item
}

// 变迁边
type TransitionEdge struct {
	DriverSymbol *GrammarSymbol // 转换符
	FromItemSet  *ItemSet       // 起始状态
	ToItemSet    *ItemSet       // 终止状态
}

// LR(0) 自动机
type DFA struct {
	// 开始项集
	StartItemSet *ItemSet
	// 变迁边表
	EdgeTable []*TransitionEdge
}

// LR(0) DFA 项集表
var ItemSetTable []*ItemSet

type ActionCell struct {
	// 状态序号
	StateID int
	// 终结符名称
	TerminalSymbolName string
	// 动作
	Type ActionCategory
	// 动作编号
	ActionID int
}

type ActionCategory int

const (
	// 规约
	Reduce ActionCategory = iota
	// 移入
	Shift
	// 接受
	Accept
)

type GotoCell struct {
	// 状态序号
	StateID int
	// 非终结符名称
	NonTerminalSymbolName string
	// 转换状态序号
	NextStateID int
}

// 分析表 Action 部分（如何改 map？）
var ActionTable []*ActionCell

// 分析表 Goto 部分
var GotoTable []*GotoCell

type ProductionInfo struct {
	// 产生式编号
	ID int
	// 头部非终结符
	HeadName string
	// 文法符个数
	BodySize int
}

// 产生式概述表
var ProductionInfoTable []*ProductionInfo

// 项目集的 Closure 函数
func (itemSet *ItemSet) Closure() (closureSet *ItemSet) {
	// 加入 I 中已有的项目
	closureSet = copyItemSet(itemSet)
	// 循环直到 closureSet 不再增大
	lastSize := 0
	for len(closureSet.ItemTable) != lastSize {
		lastSize = len(closureSet.ItemTable)
		// 遍历 closureSet 中的每个项目
		for _, item := range closureSet.ItemTable {
			// 寻找 \alpha \cdot B \beta
			// 即 dotPosition 处之后为非终结符的项目
			B := item.Production.BodySymbol[item.DotPosition]
			B, ok := B.(*NonTerminalSymbol)
			if !ok {
				continue
			}
			// 重新遍历 I
			for _, item1 := range itemSet.ItemTable {
				// 寻找 B -> \cdot \gamma
				if item1.NonTerminalSymbol != B || item1.DotPosition != 0 {
					continue
				}
				// 将 B -> \cdot \gamma 加入 closureSet
				item2 := &LR0Item{
					NonTerminalSymbol: item1.NonTerminalSymbol,
					Production:        item1.Production,
					DotPosition:       item1.DotPosition,
					Type:              item1.Type,
				}
				closureSet.ItemTable = append(closureSet.ItemTable, item2)
			}
		}
	}
	return
}

// 深拷贝项目集
func copyItemSet(itemSet *ItemSet) (newItemSet *ItemSet) {
	newItemSet = &ItemSet{
		ID: maxItemSetID() + 1,
	}
	for _, item := range itemSet.ItemTable {
		// 深拷贝时不需要继续拷贝到 NonTerminal 和 Production
		// 因为语法分析时没必要改变这些东西
		item1 := &LR0Item{
			NonTerminalSymbol: item.NonTerminalSymbol,
			Production:        item.Production,
			DotPosition:       item.DotPosition,
			Type:              item.Type,
		}
		newItemSet.ItemTable = append(newItemSet.ItemTable, item1)
	}
	// 加入项目集表
	ItemSetTable = append(ItemSetTable, newItemSet)
	return
}

// 项目集表的最大 ID
func maxItemSetID() (maxID int) {
	if len(ItemSetTable) == 0 {
		return -1
	}
	maxID = ItemSetTable[len(ItemSetTable)-1].ID
	return
}
