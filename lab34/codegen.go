package lab34

// codegen.go 生成虚拟机代码

const PROGRAM_COUNTER int = 7 // 程序计数器寄存器
const MEMORY_POINTER int = 6
const GLOBAL_POINTER int = 5
const ACCUMULATOR1 int = 0
const ACCUMULATOR2 int = 1

var tmpOffset int = 0

func CodeGen(tree *treeNode) {
	// 初始化
	emitRM("LD", MEMORY_POINTER, 0, ACCUMULATOR1)
	emitComment("Load maxaddress from location 0")
	emitRM("ST", ACCUMULATOR1, 0, ACCUMULATOR1)
	emitComment("Clear location 0")
	// 对根节点开始遍历，生成代码
	cGen(tree)
	// 结束
	emitRO("HALT", 0, 0, 0)
}

func cGen(t *treeNode) {
	if t == nil {
		return
	}
	switch t.node {
	case stmtNode:
		genStmt(t)
	case exprNode:
		genExp(t)
	}
	cGen(t.sibling)
}

// 对一个语句节点，生成代码
func genStmt(t *treeNode) {
	switch t.stmt {
	case ifStmt:
		// if <expr> then <stmt1> else <stmt2> end
		// 或 if <expr> then <stmt1> end
		emitComment("-> if")
		p1 := t.child[0]
		p2 := t.child[1]
		p3 := t.child[2]
		// 生成条件表达式的代码
		cGen(p1)
		// 空一个跳转 stmt1 末尾的位置
		savedLoc1 := emitSkip(1)
		emitComment("if: jump to else belongs here")
		// 生成 stmt1 的代码
		cGen(p2)
		// 空一个指令，供跳转到 stmt2 结束
		savedLoc2 := emitSkip(1)
		emitComment("if: jump to end belongs here")
		// 获得 stmt2 的开始位置
		currLoc := emitSkip(0)
		// 回到 expr 判断结束处，写入条件跳转语句
		emitBackup(savedLoc1)
		emitAbsRM("JEQ", ACCUMULATOR1, currLoc)
		emitComment("if: jmp to else")
		// 回到 stmt2 开始处
		emitRestore()
		// 生成 stmt2 语句的代码
		cGen(p3)
		// 获得 stmt2 结束后的位置
		currLoc = emitSkip(0)
		// 回到 stmt1 结束处，无条件跳转到 stmt2 结束处
		emitBackup(savedLoc2)
		emitAbsRM("LDA", PROGRAM_COUNTER, currLoc)
		emitComment("if: jmp to end")
		// 回到 stmt2 结束处
		emitRestore()
		emitComment("<- if")
	case repeatStmt:
		// repeat <stmt1> until <expr>
		emitComment("-> repeat")
		p1 := t.child[0]
		p2 := t.child[1]
		// 保存循环前的位置
		savedLoc1 := emitSkip(0)
		emitComment("repeat: jump after body comes back here")
		// 执行语句
		cGen(p1)
		// 判断条件
		cGen(p2)
		// 如果符合，跳转到 savedLoc1，即循环前
		emitAbsRM("JEQ", ACCUMULATOR1, savedLoc1)
		emitComment("<- repeat")
	case assignStmt:
		// <id> := <expr>
		emitComment("-> assign")
		// 计算表达式的值，存入 ACCUMULATOR1
		cGen(t.child[0])
		// 找到 id 的内存偏移
		loc, _ := lookupSymtab(t.attr)
		// 将表达式的值存入 id 的内存偏移
		emitRM("ST", ACCUMULATOR1, loc, GLOBAL_POINTER)
		emitComment("<- assign")
	case readStmt:
		// read <id>
		emitComment("-> read")
		// 读取输入，存入 ACCUMULATOR1
		emitRO("IN", ACCUMULATOR1, 0, 0)
		loc, _ := lookupSymtab(t.attr)
		// 将输入的值存入 id 的内存偏移
		emitRM("ST", ACCUMULATOR1, loc, GLOBAL_POINTER)
		emitComment("<- read")
	case writeStmt:
		emitComment("-> write")
		// 计算表达式的值，存入 ACCUMULATOR1
		cGen(t.child[0])
		// 将 ACCUMULATOR1 的值输出
		emitRO("OUT", ACCUMULATOR1, 0, 0)
		emitComment("<- write")
	}
}

func genExp(t *treeNode) {
	switch t.expr {
	case constExpr:
		emitComment("-> Const")
		// 将常量存入 ACCUMULATOR1
		emitRM("LDC", ACCUMULATOR1, t.val, 0)
		emitComment("load const")
		emitComment("<- Const")
	case idExpr:
		emitComment("-> Id")
		// 找到 id 的内存偏移
		loc, _ := lookupSymtab(t.attr)
		// 将该地址对应的值存入 ACCUMULATOR1
		emitRM("LD", ACCUMULATOR1, loc, GLOBAL_POINTER)
		emitComment("load id value")
		emitComment("<- Id")
	case opExpr:
		emitComment("-> Op")
		// 求左子表达式的值，存入 ACCUMULATOR1
		cGen(t.child[0])
		// 将 ACCUMULATOR1 的值存入内存偏移 tmpOffset
		emitRM("ST", ACCUMULATOR1, tmpOffset, MEMORY_POINTER)
		// 将 tmpOffset 加 1，即指向下一个内存偏移（当然后面存的不能重合了）
		tmpOffset--
		emitComment("op: push left")
		// 求右子表达式的值，存入 ACCUMULATOR1
		cGen(t.child[1])
		tmpOffset++
		// 将左子表达式的值存入 ACCUMULATOR2
		emitRM("LD", ACCUMULATOR2, tmpOffset, MEMORY_POINTER)
		emitComment("op: load left")
		switch t.op {
		case plusToken:
			// AC1 = AC2 + AC1
			emitRO("ADD", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op +")
		case minusToken:
			// AC1 = AC2 - AC1
			emitRO("SUB", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op -")
		case timesToken:
			// AC1 = AC2 * AC1
			emitRO("MUL", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op *")
		case overToken:
			// AC1 = AC2 / AC1
			emitRO("DIV", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op /")
		case ltToken:
			// AC1 = AC2 - AC1
			emitRO("SUB", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op <")
			// 若 AC1<AC2（新 AC1<0），则往后跳两个指令
			emitRM("JLT", ACCUMULATOR1, 2, PROGRAM_COUNTER)
			emitComment("br if true")
			// 若 AC1>=AC2，则 AC1=0
			emitRM("LDC", ACCUMULATOR1, 0, ACCUMULATOR1)
			emitComment("false case")
			// 无条件跳一条指令
			emitRM("LDA", PROGRAM_COUNTER, 1, PROGRAM_COUNTER)
			emitComment("unconditional jmp")
			// 若 AC1<AC2，则 AC1=1
			emitRM("LDC", ACCUMULATOR1, 1, ACCUMULATOR1)
			emitComment("true case")
		case eqToken:
			emitRO("SUB", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op ==")
			emitRM("JEQ", ACCUMULATOR1, 2, PROGRAM_COUNTER)
			emitComment("br if true")
			emitRM("LDC", ACCUMULATOR1, 0, ACCUMULATOR1)
			emitComment("false case")
			emitRM("LDA", PROGRAM_COUNTER, 1, PROGRAM_COUNTER)
			emitComment("unconditional jmp")
			emitRM("LDC", ACCUMULATOR1, 1, ACCUMULATOR1)
			emitComment("true case")
		default:
			emitComment("BUG: Unknown operator")
		}
		emitComment("<- Op")
	}
}
