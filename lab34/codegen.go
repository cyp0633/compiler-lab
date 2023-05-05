package lab34

// codegen.go 生成虚拟机代码

const PROGRAM_COUNTER int = 7 // 程序计数器寄存器
const MEMORY_POINTER int = 6
const GLOBAL_POINTER int = 5
const ACCUMULATOR1 int = 0
const ACCUMULATOR2 int = 1

var tmpOffset int = 0

func CodeGen(tree *treeNode) {
	emitRM("LD", MEMORY_POINTER, 0, ACCUMULATOR1)
	emitComment("Load maxaddress from location 0")
	emitRM("ST", ACCUMULATOR1, 0, ACCUMULATOR1)
	emitComment("Clear location 0")
	cGen(tree)
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
}

// 对一个语句节点，生成代码
func genStmt(t *treeNode) {
	switch t.stmt {
	case ifStmt:
		emitComment("-> if")
		p1 := t.child[0]
		p2 := t.child[1]
		p3 := t.child[2]
		// 生成条件表达式的代码
		cGen(p1)
		savedLoc1 := emitSkip(1)
		emitComment("if: jump to else belongs here")
		cGen(p2)
		savedLoc2 := emitSkip(1)
		emitComment("if: jump to end belongs here")
		currLoc := emitSkip(0)
		emitBackup(savedLoc1)
		emitAbsRM("JEQ", ACCUMULATOR1, currLoc)
		emitComment("if: jmp to else")
		emitRestore()
		cGen(p3)
		currLoc = emitSkip(0)
		emitBackup(savedLoc2)
		emitAbsRM("LDA", PROGRAM_COUNTER, currLoc)
		emitComment("if: jmp to end")
		emitRestore()
		emitComment("<- if")
	case repeatStmt:
		emitComment("-> repeat")
		p1 := t.child[0]
		p2 := t.child[1]
		savedLoc1 := emitSkip(0)
		emitComment("repeat: jump after body comes back here")
		cGen(p1)
		cGen(p2)
		emitAbsRM("JEQ", ACCUMULATOR1, savedLoc1)
		emitComment("<- repeat")
	case assignStmt:
		emitComment("-> assign")
		cGen(t.child[0])
		loc, _ := lookupSymtab(t.attr)
		emitRM("ST", ACCUMULATOR1, loc, GLOBAL_POINTER)
		emitComment("<- assign")
	case readStmt:
		emitComment("-> read")
		emitRO("IN", ACCUMULATOR1, 0, 0)
		loc, _ := lookupSymtab(t.attr)
		emitRM("ST", ACCUMULATOR1, loc, GLOBAL_POINTER)
		emitComment("<- read")
	case writeStmt:
		emitComment("-> write")
		cGen(t.child[0])
		emitRO("OUT", ACCUMULATOR1, 0, 0)
		emitComment("<- write")
	}
}

func genExp(t *treeNode) {
	switch t.expr {
	case constExpr:
		emitComment("-> Const")
		emitRM("LDC", ACCUMULATOR1, t.val, 0)
		emitComment("load const")
		emitComment("<- Const")
	case idExpr:
		emitComment("-> Id")
		loc, _ := lookupSymtab(t.attr)
		emitRM("LD", ACCUMULATOR1, loc, GLOBAL_POINTER)
		emitComment("load id value")
		emitComment("<- Id")
	case opExpr:
		emitComment("-> Op")
		cGen(t.child[0])
		emitRM("ST", ACCUMULATOR1, tmpOffset, MEMORY_POINTER)
		tmpOffset--
		emitComment("op: push left")
		cGen(t.child[1])
		tmpOffset++
		emitRM("LD", ACCUMULATOR2, 0, MEMORY_POINTER)
		emitComment("op: load left")
		switch t.op {
		case plusToken:
			emitRO("ADD", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op +")
		case minusToken:
			emitRO("SUB", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op -")
		case timesToken:
			emitRO("MUL", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op *")
		case overToken:
			emitRO("DIV", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op /")
		case ltToken:
			emitRO("SUB", ACCUMULATOR1, ACCUMULATOR2, ACCUMULATOR1)
			emitComment("op <")
			emitRM("JLT", ACCUMULATOR1, 2, PROGRAM_COUNTER)
			emitComment("br if true")
			emitRM("LDC", ACCUMULATOR1, 0, ACCUMULATOR1)
			emitComment("false case")
			emitRM("LDA", PROGRAM_COUNTER, 1, PROGRAM_COUNTER)
			emitComment("unconditional jmp")
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
