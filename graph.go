package main

// 图，描述 NFA 或 DFA
type Graph struct {
	// 图 ID
	GraphId int
	// 状态数
	NumOfStates int
	// 边的列表
	EdgeTable []*Edge
	// 状态列表
	StateTable []*State
}

// 图的边
type Edge struct {
	// 从哪个状态出发
	FromState int
	// 到哪个状态
	NextState int
	// driver
	DriverID int
	// driver 类型
	DriverType driverType
}

// 边的转换类型
type driverType int

const (
	// epsilon 转换
	DriverNull driverType = iota
	// 字符
	DriverChar
	// 字符集
	DriverCharset
)

// 图的状态（节点）
type State struct {
	// 状态 ID
	StateID int
	// 状态类型
	StateType stateType
	// 类型
	Category LexemeCategory
}

type stateType int

const (
	// 接受状态
	StateMatch stateType = iota
	// 非接受状态
	StateUnmatch
)

// 生成基本 NFA，只有 0 和 1
func generateBasicNFA(driver driverType, driverID int) (g *Graph) {
	g = new(Graph)
	state0 := State{StateID: 0, StateType: StateUnmatch, Category: LexemeIntegerConst}
	state1 := State{StateID: 1, StateType: StateMatch, Category: LexemeIntegerConst}
	edge := Edge{DriverType: driver, DriverID: driverID, FromState: 0, NextState: 1}
	g.EdgeTable = append(g.EdgeTable, &edge)
	g.StateTable = append(g.StateTable, &state0, &state1)
	g.NumOfStates = 2
	return
}

// NFA 的并运算
func unionNFA(g1, g2 *Graph) (g *Graph) {
	g = new(Graph)
	g.EdgeTable = append(g.EdgeTable, g1.EdgeTable...)
	g.EdgeTable = append(g.EdgeTable, g2.EdgeTable...)
	g.StateTable = append(g.StateTable, g1.StateTable...)
	g.StateTable = append(g.StateTable, g2.StateTable...)
	g.NumOfStates = g1.NumOfStates + g2.NumOfStates
	return
}
