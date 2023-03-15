package main

import (
	"github.com/google/go-cmp/cmp"
)

// 子集构造法的 move 函数，state 经过 input 可达的状态
func (g *Graph) Move(state int, input int, inputType driverType) (ac map[int]bool) {
	ac = make(map[int]bool)
	queue := make(chan int, g.NumOfStates)
	queue <- state
	for {
		if len(queue) == 0 {
			break
		}
		temp := <-queue
		for _, edge := range g.EdgeTable {
			// 检查 state 的 input 出边
			if edge.FromState == temp && edge.DriverType == DriverChar && edge.DriverID == input {
				// 如果没有添加过，就添加到 map 中
				if _, ok := ac[edge.NextState]; !ok {
					ac[edge.NextState] = true
				}
			}
		}
	}
	return
}

// 子集构造法的 move 函数，stateSet 经过 input 可达的状态
func (g *Graph) MoveSet(stateSet map[int]bool, input int, inputType driverType) (ac map[int]bool) {
	ac = make(map[int]bool)
	queue := make(chan int, g.NumOfStates)
	for state := range stateSet {
		queue <- state
	}
	for {
		if len(queue) == 0 {
			break
		}
		temp := <-queue
		for _, edge := range g.EdgeTable {
			// 检查 state 的 input 出边
			if edge.FromState == temp && edge.DriverType == DriverChar && edge.DriverID == input {
				// 如果没有添加过，就添加到 map 中
				if _, ok := ac[edge.NextState]; !ok {
					ac[edge.NextState] = true
				}
			}
		}
	}
	return
}

// 子集构造法的 epsilonClosure 函数，state 经过 epsilon 可达的状态（包含自身）
func (g *Graph) EpsilonClosure(state int) (ac map[int]bool) {
	ac = make(map[int]bool)
	// 先添加自身
	ac[state] = true
	// 使用 buffered channel 作为队列，先添加 state
	queue := make(chan int, g.NumOfStates)
	queue <- state
	for {
		if len(queue) == 0 {
			break
		}
		temp := <-queue
		for _, edge := range g.EdgeTable {
			// 检查 state 的出边中的 epsilon 转换
			if edge.FromState == temp && edge.DriverType == DriverNull {
				// 如果没有添加过，就添加到队列和 map 中
				if _, ok := ac[edge.NextState]; !ok {
					ac[edge.NextState] = true
					queue <- edge.NextState
				}
			}
		}
	}
	return
}

// 子集构造法的 EpsilonClosure 函数，但输入的是一个状态集合
func (g *Graph) EpsilonClosureSet(stateSet map[int]bool) (ac map[int]bool) {
	ac = make(map[int]bool)
	// 使用 buffered channel 作为队列
	queue := make(chan int, g.NumOfStates)
	for state := range stateSet {
		// 先添加自身
		ac[state] = true
		// 先添加 state
		queue <- state
	}
	for {
		if len(queue) == 0 {
			break
		}
		temp := <-queue
		for _, edge := range g.EdgeTable {
			// 检查 state 的出边中的 epsilon 转换
			if edge.FromState == temp && edge.DriverType == DriverNull {
				// 如果没有添加过，就添加到队列和 map 中
				if _, ok := ac[edge.NextState]; !ok {
					ac[edge.NextState] = true
					queue <- edge.NextState
				}
			}
		}
	}
	return
}

// 子集构造法的 DTran 函数
func (g *Graph) DTran(state int) (ac []int) {
	return
}

// 使用子集构造法从 NFA 转换为 DFA
func (g *Graph) SubsetConstruction() (gNew *Graph) {
	// 新节点表，key 为构造的子集，value 为新节点的编号
	newNodes := make(map[*map[int]bool]int)
	gNew = &Graph{}
	// 第一个节点，为初始状态的 epsilonClosure
	firstNode := g.EpsilonClosure(0)
	newNodes[&firstNode] = 0
	gNew.StateTable = append(gNew.StateTable, &State{StateID: 0, StateType: StateUnmatch, Category: LexemeNull})
	gNew.NumOfStates++
	// 构建输入集，key 为 {driverID, driverType}，value 为 true
	inputSet := make(map[struct {
		int
		driverType
	}]bool)
	for _, edge := range g.EdgeTable {
		if edge.DriverType != DriverNull {
			inputSet[struct {
				int
				driverType
			}{edge.DriverID, edge.DriverType}] = true
		}
	}
	// 待遍历的新节点队列
	queue := make(chan *map[int]bool, g.NumOfStates*10)
	queue <- &firstNode
	// 遍历每个新节点
	for {
		if len(queue) == 0 {
			break
		}
		currNode := <-queue
		// 遍历每种输入
		for key := range inputSet {
			node := g.EpsilonClosureSet(g.MoveSet(*currNode, key.int, key.driverType))
			// 检查 DTran 在新节点表中的新分配序号，不在则 -1
			nodeInTable := -1
			for key := range newNodes {
				if cmp.Equal(*key, node) {
					nodeInTable = newNodes[key]
					break
				}
			}
			// 如果新节点不在节点表中，就添加到节点表和队列中
			if nodeInTable == -1 {
				newNodes[&node] = gNew.NumOfStates
				nodeInTable = gNew.NumOfStates
				state := State{StateID: gNew.NumOfStates, StateType: StateUnmatch, Category: LexemeNull}
				gNew.StateTable = append(gNew.StateTable, &state)
				gNew.NumOfStates++
				queue <- &node
			}
			// 添加边
			edge := Edge{FromState: newNodes[currNode], NextState: nodeInTable, DriverID: key.int, DriverType: key.driverType}
			gNew.EdgeTable = append(gNew.EdgeTable, &edge)
		}
	}
	// 将最后一个状态设为接受
	gNew.StateTable[gNew.NumOfStates-1].StateType = StateMatch
	return
}
