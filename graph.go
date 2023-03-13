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
	state0 := State{StateID: 0, StateType: StateUnmatch, Category: LexemeNull}
	state1 := State{StateID: 1, StateType: StateMatch, Category: LexemeNull}
	edge := Edge{DriverType: driver, DriverID: driverID, FromState: 0, NextState: 1}
	g.EdgeTable = append(g.EdgeTable, &edge)
	g.StateTable = append(g.StateTable, &state0, &state1)
	g.NumOfStates = 2
	return
}

// NFA 的并运算
func unionNFA(g1, g2 *Graph) (g *Graph) {
	// 拷贝两个 NFA
	g1 = copyNFA(g1)
	g2 = copyNFA(g2)
	// 最简并之前的预处理
	unionNFAPreprocess(g1)
	unionNFAPreprocess2(g1)
	unionNFAPreprocess(g2)
	unionNFAPreprocess2(g2)
	// 合并
	g = new(Graph)
	// 添加新的起始和接受
	g.NumOfStates = g1.NumOfStates + g2.NumOfStates
	startState := State{StateID: 0, StateType: StateUnmatch}
	endState := State{StateID: g.NumOfStates - 1, StateType: StateMatch}
	g.StateTable = append(g.StateTable, &startState, &endState)
	// 将 g1 的状态拷贝进来，除了 0 和最后一个状态
	for i := 1; i < g1.NumOfStates-1; i++ {
		g.StateTable = append(g.StateTable, g1.StateTable[i])
	}
	// 将 g1 的边拷贝进来，对指向最后一个状态的边的序号进行修改
	for _, edge := range g1.EdgeTable {
		if edge.NextState == g1.NumOfStates-1 {
			edge.NextState = g.NumOfStates - 1
		}
		g.EdgeTable = append(g.EdgeTable, edge)
	}
	// 将 g2 的状态拷贝进来，除了 0 和最后一个状态，并对序号进行修改
	for i := 1; i < g2.NumOfStates-1; i++ {
		g.StateTable = append(g.StateTable, g2.StateTable[i])
		g.StateTable[i+g1.NumOfStates-1].StateID += g1.NumOfStates - 1
	}
	// 将 g2 的边拷贝进来，对指向每一个边的序号进行修改
	for _, edge := range g2.EdgeTable {
		if edge.FromState != 0 {
			edge.FromState += g1.NumOfStates - 1
		}
		edge.NextState += g1.NumOfStates - 1
	}
	return
}

// NFA 并运算的预处理
//
// 也就是根据有无出边在前后加状态
func unionNFAPreprocess(g *Graph) {
	hasInEdge, hasOutEdge := inOutEdge(g)
	// 若 0 有入边，则新建一个 0 状态，用 epsilon 转换连接到原来的 0，并修改序号
	if hasInEdge {
		g.NumOfStates++
		// 将原来边的序号加 1
		for _, edge := range g.EdgeTable {
			edge.FromState++
			edge.NextState++
		}
		// 将原来状态的序号加 1
		for _, state := range g.StateTable {
			state.StateID++
		}
		// 新建一个 0 状态
		state0 := State{StateID: 0, StateType: StateUnmatch, Category: LexemeNull}
		// 新建一个 epsilon 转换
		edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: 0, NextState: 1}
		g.EdgeTable = append(g.EdgeTable, &edge)
		g.StateTable = append(g.StateTable, &state0)
	}
	// 若最后一个状态有出边，则新建一个状态，从原来的最后一个状态用 epsilon 转换连接到新的状态
	if hasOutEdge {
		g.NumOfStates++
		// 修改接受状态
		g.StateTable[g.NumOfStates-2].StateType = StateUnmatch
		// 新建一个状态
		state := State{StateID: g.NumOfStates - 1, StateType: StateMatch, Category: LexemeNull}
		// 新建一个 epsilon 转换
		edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 2, NextState: g.NumOfStates - 1}
		g.EdgeTable = append(g.EdgeTable, &edge)
		g.StateTable = append(g.StateTable, &state)
	}
}

// NFA 并运算的预处理第二步
//
// 独立带 category 的状态
func unionNFAPreprocess2(g *Graph) {
	// 若最后一个状态带 category，则新建一个状态，从原来的最后一个状态用 epsilon 转换连接到新的状态
	if g.StateTable[g.NumOfStates-1].Category != LexemeNull {
		g.NumOfStates++
		// 修改接受状态
		g.StateTable[g.NumOfStates-2].StateType = StateUnmatch
		// 新建一个状态
		state := State{StateID: g.NumOfStates - 1, StateType: StateMatch, Category: LexemeNull}
		// 新建一个 epsilon 转换
		edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 2, NextState: g.NumOfStates - 1}
		g.EdgeTable = append(g.EdgeTable, &edge)
		g.StateTable = append(g.StateTable, &state)
	}
}

// （深）拷贝 NFA
//
// 防止编辑同一个指针时修改了原来的 NFA
func copyNFA(g *Graph) (gCopy *Graph) {
	gCopy = new(Graph)
	gCopy.NumOfStates = g.NumOfStates
	for _, edge := range g.EdgeTable {
		edgeCopy := Edge{DriverType: edge.DriverType, DriverID: edge.DriverID, FromState: edge.FromState, NextState: edge.NextState}
		gCopy.EdgeTable = append(gCopy.EdgeTable, &edgeCopy)
	}
	for _, state := range g.StateTable {
		stateCopy := State{StateID: state.StateID, StateType: state.StateType, Category: state.Category}
		gCopy.StateTable = append(gCopy.StateTable, &stateCopy)
	}
	return
}

// 连接运算
func productNFA(g1, g2 *Graph) (g *Graph) {
	// g1 的出边和 g2 的入边
	hasInEdge, hasOutEdge := false, false
	for _, edge := range g1.EdgeTable {
		if edge.NextState == g1.NumOfStates-1 {
			hasOutEdge = true
		}
	}
	for _, edge := range g2.EdgeTable {
		if edge.FromState == 0 {
			hasInEdge = true
		}
	}
	// 将 g1 的状态拷贝进来
	g = copyNFA(g1)
	// 如果 g1 的出边和 g2 的入边都存在，则再加一个状态
	// 原最后一个状态用 epsilon 转换连接到新的状态
	// 新的状态用 epsilon 转换连接到 g2 的第一个状态
	if hasInEdge && hasOutEdge {
		g.NumOfStates++
		// 新建一个状态
		state := State{StateID: g.NumOfStates - 1, StateType: StateUnmatch, Category: LexemeNull}
		// 新建一个 epsilon 转换
		edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 2, NextState: g.NumOfStates - 1}
		g.EdgeTable = append(g.EdgeTable, &edge)
		g.StateTable = append(g.StateTable, &state)
	}
	// 建立到 g2 第一个状态的转换
	edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 1, NextState: g.NumOfStates}
	g.EdgeTable = append(g.EdgeTable, &edge)
	// 将 g2 的状态拷贝进来，并修改序号
	g2Copy := copyNFA(g2)
	for _, edge := range g2Copy.EdgeTable {
		edge.FromState += g.NumOfStates - 1
		edge.NextState += g.NumOfStates - 1
	}
	for _, state := range g2Copy.StateTable {
		state.StateID += g.NumOfStates - 1
	}
	g.NumOfStates += g2Copy.NumOfStates - 1
	g.EdgeTable = append(g.EdgeTable, g2Copy.EdgeTable...)
	g.StateTable = append(g.StateTable, g2Copy.StateTable...)
	return
}

// 正闭包运算
//
// 1 或更多次
func plusClosureNFA(g *Graph) *Graph {
	g = copyNFA(g)
	hasInEdge, hasOutEdge := inOutEdge(g)
	// 不管怎样，先从最后一个到第一个加一个 epsilon 转换
	edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 1, NextState: 0}
	g.EdgeTable = append(g.EdgeTable, &edge)
	// 由于加入了一个 epsilon 转换，所以不能直接使用 unionNFAPreprocess1 加状态
	// 如果有入边，就加一个 0，用 epsilon 连接原来的 0
	if hasInEdge {
		g.NumOfStates++
		// 将原来边的序号加 1
		for _, edge := range g.EdgeTable {
			edge.FromState++
			edge.NextState++
		}
		// 将原来状态的序号加 1
		for _, state := range g.StateTable {
			state.StateID++
		}
		// 新建一个 0 状态
		state0 := State{StateID: 0, StateType: StateUnmatch, Category: LexemeNull}
		// 新建一个 epsilon 转换
		edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: 0, NextState: 1}
		g.EdgeTable = append(g.EdgeTable, &edge)
		g.StateTable = append(g.StateTable, &state0)
	}
	// 若最后一个状态有出边，则新建一个状态，从原来的最后一个状态用 epsilon 转换连接到新的状态
	if hasOutEdge {
		g.NumOfStates++
		// 修改接受状态
		g.StateTable[g.NumOfStates-2].StateType = StateUnmatch
		// 新建一个状态
		state := State{StateID: g.NumOfStates - 1, StateType: StateMatch, Category: LexemeNull}
		// 新建一个 epsilon 转换
		edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 2, NextState: g.NumOfStates - 1}
		g.EdgeTable = append(g.EdgeTable, &edge)
		g.StateTable = append(g.StateTable, &state)
	}
	return g
}

// Kleene 闭包运算
//
// 0 或更多次
func kleeneClosureNFA(g *Graph) *Graph {
	g = plusClosureNFA(g)
	// 再从最后一个到第一个加一个 epsilon 转换就行了
	edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: g.NumOfStates - 1, NextState: 0}
	g.EdgeTable = append(g.EdgeTable, &edge)
	return g
}

// 0 或 1 次
func zeroOrOneNFA(g *Graph) *Graph {
	g = copyNFA(g)
	// 加状态
	unionNFAPreprocess(g)
	// 从第一个状态到最后一个状态加一个 epsilon 转换
	edge := Edge{DriverType: DriverNull, DriverID: 0, FromState: 0, NextState: g.NumOfStates - 1}
	g.EdgeTable = append(g.EdgeTable, &edge)
	return g
}

func inOutEdge(g *Graph) (hasInEdge, hasOutEdge bool) {
	for _, edge := range g.EdgeTable {
		if edge.NextState == g.NumOfStates-1 {
			hasOutEdge = true
		}
		if edge.FromState == 0 {
			hasInEdge = true
		}
	}
	return
}
