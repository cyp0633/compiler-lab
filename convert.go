package main

// 子集构造法的 move 函数，state 经过 input 可达的状态
func (g *Graph) Move(state int, input int) (ac map[int]bool) {
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
	for temp := range queue {
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
	return
}
