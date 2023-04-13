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
	// 核心项
	// S' \to \cdot S 和所有点不在左端的
	CoreItem itemCategory = iota
	// 非核心项
	// 除那个以外所有点在左端的
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
type TransitionKey struct {
	DriverSymbol interface{} // 转换符，某个语法符的指针
	FromItemSet  *ItemSet    // 起始状态
}

// LR(0) 自动机
var DFA struct {
	// 开始项集
	StartItemSet *ItemSet
	// map 优化的变迁边表
	// 通常查询更快，极端情况下也不会更慢
	EdgeSet map[TransitionKey]*ItemSet
}

// LR(0) DFA 项集表
var ItemSetTable []*ItemSet

// Action 动作类型，没有对应表项代表错误
type ActionCategory int

const (
	// 规约
	Reduce ActionCategory = iota
	// 移入
	Shift
	// 接受
	Accept
)

// LR Action 表
// 由状态 ID 和终结符名称，到动作类型和编号的映射
var ActionTable map[struct {
	// 当前栈顶状态序号
	StateID int
	// 待读入的终结符名称
	TerminalSymbolName string
}]struct {
	// 动作类型
	Type ActionCategory
	// 动作编号，如归约的产生式编号和移进的下个状态
	ActionID int
}

// LR Goto 表
// 由状态 ID 和非终结符名称，到转换状态序号的映射
var GotoTable map[struct {
	// 当前栈顶状态序号
	StateID int
	// 非终结符名称
	NonTerminalSymbolName string
}]int

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
			// 如果是 \alpha \cdot 即归约/接受项目，跳过
			if item.DotPosition == len(item.Production.BodySymbol) {
				continue
			}
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
// 不再直接加入项目集规范族（即 ItemSetTable）
func copyItemSet(itemSet *ItemSet) (newItemSet *ItemSet) {
	newItemSet = &ItemSet{
		ItemTable: map[LR0Item]struct{}{},
	}
	for item := range itemSet.ItemTable {
		// 深拷贝时可以直接复制指针
		// 毕竟 item 不会变
		newItemSet.ItemTable[item] = struct{}{}
	}
	return
}

// 项目集的 Goto 函数
//
// 需要自行进行 Closure！！！
func (itemSet *ItemSet) Goto(X interface{}) (gotoSet *ItemSet) {
	gotoSet = &ItemSet{
		ID:        itemSet.ID,
		ItemTable: map[LR0Item]struct{}{},
	}

	for item := range itemSet.ItemTable {
		// 如果是 \alpha \cdot 即归约/接受项目，跳过
		if item.DotPosition == len(item.Production.BodySymbol) {
			continue
		}
		// 寻找 \alpha \cdot X \beta
		// 即 dotPosition 处之后为 X 的项目
		if item.Production.BodySymbol[item.DotPosition] == X {
			// 对 \alpha X \cdot \beta，添加项目 \alpha X \to \cdot \beta
			item1 := LR0Item{
				NonTerminalSymbol: item.NonTerminalSymbol,
				Production:        item.Production,
				DotPosition:       item.DotPosition + 1,
				Type:              NonCoreItem,
			}
			if _, ok := gotoSet.ItemTable[item1]; !ok {
				gotoSet.ItemTable[item1] = struct{}{}
			}
		}
	}
	return
}

// 穷举（某个）项目集的变迁
func (itemSet *ItemSet) ExhaustTransition() {
	drivers := itemSet.driver()

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
		// 新的项目！
		item1 := LR0Item{
			NonTerminalSymbol: item.NonTerminalSymbol,
			Production:        item.Production,
			DotPosition:       item.DotPosition + 1,
			Type:              CoreItem, // 一定是核心项啦
		}
		// 将新项目加入项目集
		itemSet[item1] = struct{}{}
	}

	// 对每一个新项目集，求 Closure
	for driver, itemSet := range newItemSets {
		// 求 Closure
		closureSet := itemSet.Closure()
		// 检查是否已存在
		exist := false
		for _, set := range ItemSetTable {
			if cmp.Equal(set.ItemTable, closureSet.ItemTable) {
				// 已存在，将新项目集指向已存在的项目集
				newItemSets[driver] = set
				exist = true
				break
			}
		}
		// 不存在，加入项目集规范族
		if !exist {
			newItemSets[driver] = closureSet
			closureSet.ID = maxItemSetID() + 1
			ItemSetTable = append(ItemSetTable, closureSet)
		}
	}
}

// 得到一个项目集的所有驱动符
func (itemSet *ItemSet) driver() (drivers map[interface{}]struct{}) {
	// 驱动符，就是点之后的符号，可能是终结符或非终结符（均为指针）
	// 由于 ItemTable 是个 map，无法指定驱动符出现的顺序
	drivers = map[interface{}]struct{}{}

	// 遍历项目集中的每个项目
	for item := range itemSet.ItemTable {
		// 将项目 item 的点后的符号加入驱动符集
		if item.DotPosition < len(item.Production.BodySymbol) {
			drivers[item.Production.BodySymbol[item.DotPosition]] = struct{}{}
		}
	}
	return
}

// 得到一个项目集对应的核心项集
func (itemSet *ItemSet) core() (coreSet *ItemSet) {
	coreSet = &ItemSet{
		ID: itemSet.ID,
	}
	for item := range itemSet.ItemTable {
		if item.Type == CoreItem {
			coreSet.ItemTable[item] = struct{}{}
		}
	}
	return
}

// 构造 LR(0) DFA
func BuildDFA() {
	// 初始化变迁边表
	DFA.EdgeSet = map[TransitionKey]*ItemSet{}

	// 找到初始项集
	// 其应该有且仅有初始符对应的一个核心项
	initialItem := LR0Item{
		NonTerminalSymbol: RootSymbol,
		Production:        RootSymbol.ProductionTable[0],
		DotPosition:       0,
		Type:              CoreItem,
	}
	for _, itemSet := range ItemSetTable {
		if _, ok := itemSet.ItemTable[initialItem]; ok {
			DFA.StartItemSet = itemSet
		}
	}

	// 遍历项目集规范族
	for _, currItemSet := range ItemSetTable {
		drivers := currItemSet.driver()

		// 对每一种驱动符，新建一个项目集
		// key 为驱动符，一个指针
		newItemSets := map[interface{}]*ItemSet{}
		for driver := range drivers {
			newItemSets[driver] = &ItemSet{ItemTable: map[LR0Item]struct{}{}}
		}

		// 遍历项目集中的每个项目，求出 GOTO 的核心项
		for item := range currItemSet.ItemTable {
			// 如果已经是归约/接受项目（A \to \cdot \alpha），则不需要变迁
			if item.DotPosition == len(item.Production.BodySymbol) {
				continue
			}
			// 取出该驱动符对应的项目集
			itemSet := newItemSets[item.Production.BodySymbol[item.DotPosition]].ItemTable
			// 将项目 item 的点后移一位
			// 新的项目！
			item1 := LR0Item{
				NonTerminalSymbol: item.NonTerminalSymbol,
				Production:        item.Production,
				DotPosition:       item.DotPosition + 1,
				Type:              CoreItem, // 一定是核心项啦
			}
			// 将新项目加入项目集
			itemSet[item1] = struct{}{}
		}

		// 对每个新项目集，找出对应的 DFA 状态，关联驱动符
		for driver, newItemSet := range newItemSets {
			// 求 Closure
			closureSet := newItemSet.Closure()
			// 寻找已知的项集（状态）
			for _, set := range ItemSetTable {
				if cmp.Equal(set.ItemTable, closureSet.ItemTable) {
					// 已存在，将新项目集指向已存在的项目集
					DFA.EdgeSet[struct {
						DriverSymbol interface{}
						FromItemSet  *ItemSet
					}{driver, currItemSet}] = set
					break
				}
			}
		}
	}
}

// 检查是否为 SLR(1) 文法
// 需要先构造项目集规范族
func CheckSLR1() bool {
	// 遍历项目集规范族
	for _, itemSet := range ItemSetTable {
		// 归约项集和移进项对应终结符集
		reduceSet := map[LR0Item]struct{}{}
		shiftSymbols := map[TerminalSymbol]struct{}{}
		// 遍历项目集中的每个项目
		for item := range itemSet.ItemTable {
			// 归约项
			if item.DotPosition == len(item.Production.BodySymbol) {
				// 跳过接受项
				if item.NonTerminalSymbol == RootSymbol {
					continue
				}
				reduceSet[item] = struct{}{}
			} else if symbol, ok := item.Production.BodySymbol[item.DotPosition].(*TerminalSymbol); ok {
				// 移进项
				// 点后面是终结符
				shiftSymbols[*symbol] = struct{}{}
			}
		}

		// 如果归约和移进二者至多存在一种，则是LR(0)，更是 SLR(1)
		if (len(reduceSet) == 0) != (len(shiftSymbols) == 0) {
			continue
		}

		// 否则，求交集
		tempSet := map[TerminalSymbol]struct{}{}
		// 先把移进项的终结符加入
		for symbol := range shiftSymbols {
			tempSet[symbol] = struct{}{}
		}
		// 然后是归约项左端的 follow 集
		Follow()
		for item := range reduceSet {
			// 对每个项的 follow 先看看是不是存在，再加入
			for symbol := range item.NonTerminalSymbol.FollowSet {
				terminal, ok := symbol.(*TerminalSymbol)
				if !ok {
					// 是个 epsilon
					continue
				}
				if _, ok := tempSet[*terminal]; ok {
					// 有交集，不是 SLR(1)
					return false
				}
				tempSet[*terminal] = struct{}{}
			}
		}
	}
	// 都没问题，是 SLR(1)
	return true
}

// 填写 LR(0) 分析表
// 需要构造项目集规范族和 DFA
func FillLR0ParsingTable() {
	// 非终结符，GOTO 列
	nonTerminals := map[string]struct{}{}
	// 终结符，ACTION 列
	terminals := map[string]struct{}{}
	// 遍历语法符表，填入所有的终结符和非终结符
	for _, symbol := range GrammarSymbolTable {
		switch symbol := symbol.(type) {
		case NonTerminalSymbol:
			nonTerminals[symbol.Name] = struct{}{}
		case TerminalSymbol:
			terminals[symbol.Name] = struct{}{}
		default:
			panic("unknown symbol type")
		}
	}
	// 加入 # 终结符
	terminals["#"] = struct{}{}

	// 遍历项目集规范族
	for _, itemSet := range ItemSetTable {
		// 对其中的每个项目
		for item := range itemSet.ItemTable {
			// 判断项目类型
			// 为什么不用一个 switch？因为特么的 go 不支持 switch 里做声明
			if item.NonTerminalSymbol == RootSymbol && item.DotPosition == len(item.Production.BodySymbol) {
				// 接受项目 S' \to S \cdot
				// ACTION[i,#] = accept
				ActionTable[struct {
					StateID            int
					TerminalSymbolName string
				}{itemSet.ID, "#"}] = struct {
					Type     ActionCategory
					ActionID int
				}{Accept, -1}
			} else if item.DotPosition == len(item.Production.BodySymbol) {
				// 归约项目 A \to \alpha \cdot
				// 对所有非终结符 a 或 #，ACTION[i,a] = reduce j
				for symbol := range terminals {
					ActionTable[struct {
						StateID            int
						TerminalSymbolName string
					}{itemSet.ID, symbol}] = struct {
						Type     ActionCategory
						ActionID int
					}{Reduce, item.Production.ID}
				}
			} else if symbol, ok := item.Production.BodySymbol[item.DotPosition].(*TerminalSymbol); ok {
				// 移进项目 A \to \alpha \cdot a \beta
				// action[i,a] = shift j
				targetState := DFA.EdgeSet[TransitionKey{symbol, itemSet}]
				ActionTable[struct {
					StateID            int
					TerminalSymbolName string
				}{itemSet.ID, symbol.Name}] = struct {
					Type     ActionCategory
					ActionID int
				}{Shift, targetState.ID}
			} else if symbol, ok := item.Production.BodySymbol[item.DotPosition].(*NonTerminalSymbol); ok {
				// 待约项目 A \to \alpha \cdot B
				// goto[i,B] = j
				targetState := DFA.EdgeSet[TransitionKey{symbol, itemSet}]
				GotoTable[struct {
					StateID               int
					NonTerminalSymbolName string
				}{itemSet.ID, symbol.Name}] = targetState.ID
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
	if item.DotPosition == len(item.Production.BodySymbol) {
		str += "•"
	}
	if item.Type == CoreItem {
		str += "\t(Core)"
	} else {
		str += "\t(NonCore)"
	}
	return
}

func (set *ItemSet) String() (str string) {
	str = fmt.Sprintf("Item set #%d, %d items:\n", set.ID, len(set.ItemTable))
	for item := range set.ItemTable {
		str += item.String() + "\n"
	}
	return
}

func (edge *TransitionKey) String() (str string) {
	var name string
	switch driver := edge.DriverSymbol.(type) {
	case *NonTerminalSymbol:
		name = "nonterminal " + driver.Name
	case *TerminalSymbol:
		name = "terminal " + driver.Name
	default:
		panic("unknown driver type")
	}
	str = fmt.Sprintf("State from #%v, driver: %v", edge.FromItemSet.ID, name)
	return
}
