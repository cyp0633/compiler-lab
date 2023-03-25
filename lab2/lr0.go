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
