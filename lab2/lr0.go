package lab2

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

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
	// LR(0) 项目表（其实是个集合）
	ItemTable map[LR0Item]struct{}
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
		for item := range closureSet.ItemTable {
			// 寻找 \alpha \cdot B \beta
			// 即 dotPosition 处之后为非终结符的项目
			B, ok := item.Production.BodySymbol[item.DotPosition].(*NonTerminalSymbol)
			if !ok {
				continue
			}
			// 遍历 B 的产生式
			for _, production := range B.ProductionTable {
				// 对 B \to \gamma，添加项目 B \to \cdot \gamma
				item1 := LR0Item{
					NonTerminalSymbol: B,
					Production:        production,
					DotPosition:       0,
					Type:              NonCoreItem,
				}
				if _, ok := closureSet.ItemTable[item1]; !ok {
					closureSet.ItemTable[item1] = struct{}{}
				}
			}
		}
	}
	return
}

// 深拷贝项目集
func copyItemSet(itemSet *ItemSet) (newItemSet *ItemSet) {
	newItemSet = &ItemSet{
		ID:        maxItemSetID() + 1,
		ItemTable: map[LR0Item]struct{}{},
	}
	for item := range itemSet.ItemTable {
		// 深拷贝时可以直接复制指针
		// 毕竟 item 不会变
		newItemSet.ItemTable[item] = struct{}{}
	}
	// 加入项目集表
	ItemSetTable = append(ItemSetTable, newItemSet)
	return
}

// 穷举项目集变迁
func (itemSet *ItemSet) ExhaustTransition() {
	// 驱动符，就是点之后的符号，可能是终结符或非终结符（均为指针）
	drivers := map[interface{}]struct{}{}
	// 遍历项目集中的每个项目
	for item := range itemSet.ItemTable {
		// 将项目 item 的点后的符号加入驱动符集
		if item.DotPosition < len(item.Production.BodySymbol) {
			drivers[item.Production.BodySymbol[item.DotPosition]] = struct{}{}
		}
	}

	// 对每一种驱动符，新建一个项目集
	// key 为驱动符，一个指针
	newItemSets := map[interface{}]*ItemSet{}
	for driver := range drivers {
		newItemSets[driver] = &ItemSet{ItemTable: map[LR0Item]struct{}{}}
	}

	// 遍历项目集中的每个项目
	// 此处并不需要驱动符一层项目再一层，省点时间
	for item := range itemSet.ItemTable {
		// 如果已经是归约/接受项目（A \to \cdot \alpha），则不需要变迁
		if item.DotPosition == len(item.Production.BodySymbol) {
			continue
		}
		// 取出该驱动符对应的项目集
		itemSet := newItemSets[item.Production.BodySymbol[item.DotPosition]].ItemTable
		// 将项目 item 的点后移一位
		item1 := LR0Item{
			NonTerminalSymbol: item.NonTerminalSymbol,
			Production:        item.Production,
			DotPosition:       item.DotPosition + 1,
			Type:              item.Type,
		}
		// 将新项目加入项目集
		itemSet[item1] = struct{}{}
	}

	// 对每一个新项目集，求 Closure
	for driver, itemSet := range newItemSets {
		// 求 Closure
		closureSet := itemSet.Closure()
		// 检查是否已存在
		for _, set := range ItemSetTable {
			if cmp.Equal(set.ItemTable, closureSet.ItemTable) {
				// 已存在，将新项目集指向已存在的项目集
				newItemSets[driver] = set
				break
			}
		}
		
	}
}

// 项目集表的最大 ID
func maxItemSetID() (maxID int) {
	if len(ItemSetTable) == 0 {
		return -1
	}
	maxID = ItemSetTable[len(ItemSetTable)-1].ID
	return
}

func (item *LR0Item) String() (str string) {
	str = item.NonTerminalSymbol.Name + " -> "
	for index, symbol := range item.Production.BodySymbol {
		if index == item.DotPosition {
			str += "• "
		}
		switch symbol := symbol.(type) {
		case *NonTerminalSymbol:
			str += symbol.Name + " "
		case *TerminalSymbol:
			str += symbol.Name + " "
		default:
			panic("unknown symbol type")
		}
	}
	return
}

func (set *ItemSet) String() (str string) {
	str = fmt.Sprintf("Item set #%d:\n", set.ID)
	for item := range set.ItemTable {
		str += item.String() + "\n"
	}
	return
}
