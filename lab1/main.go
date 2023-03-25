package lab1

func main() {
	testRegex()
}

// 测试正则表达式 (a|b)*abb
func testRegex() {
	// r1 -> a
	r1 := lab1.generateBasicNFA(DriverChar, 'a')
	lab1.printGraph(r1)
	// r2 -> b
	r2 := lab1.generateBasicNFA(DriverChar, 'b')
	lab1.printGraph(r2)
	// r3 -> r1 | r2
	r3 := lab1.unionNFA(r1, r2)
	printGraph(r3)
	// r4 -> r3*
	r4 := kleeneClosureNFA(r3)
	printGraph(r4)
	// r5 -> r4 r1
	r5 := productNFA(r4, r1)
	printGraph(r5)
	// r6 -> r5 r2
	r6 := productNFA(r5, r2)
	printGraph(r6)
	// r7 -> r6 r2
	r7 := productNFA(r6, r2)
	printGraph(r7)
}
