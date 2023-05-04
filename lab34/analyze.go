package lab34

import "fmt"

// analyze.go 语义分析器

// 遍历语法树
// 自定义前序和后序处理函数
func traverse(t *treeNode, preProc func(*treeNode), postProc func(*treeNode)) {
	if t == nil {
		return
	}
	if preProc != nil {
		preProc(t)
	}
	for i := 0; i < 3; i++ {
		traverse(t.child[i], preProc, postProc)
	}
	if postProc != nil {
		postProc(t)
	}
	traverse(t.sibling, preProc, postProc)
}

// 将某个语法树节点的信息插入符号表
func insertNode(t *treeNode) {
	switch t.node {
	case stmtNode: // 对于声明，只有赋值和读语句需要插入符号表
		if t.stmt == assignStmt || t.stmt == readStmt {
			insertSymtab(t.attr, t.lineNo)
		}
	case exprNode: // 对于表达式，其中的标识符都要记在符号表中
		if t.expr == idExpr {
			insertSymtab(t.attr, t.lineNo)
		}
	}
}

// 遍历，构建符号表
// 其实前序后序不会有什么差别
func BuildSymtab(t *treeNode) {
	traverse(t, insertNode, nil)
}

// 产生一个类型错误
func typeError(t *treeNode, msg string) {
	fmt.Printf("Type error at line %d: %s\n", t.lineNo, msg)
}

// 对 t 做类型检查，反正就 bool 和 integer 两种类型
func checkNode(t *treeNode) {
	switch t.node {
	case exprNode:
		switch t.expr {
		case opExpr: // 运算符表达式，两个 child 类型均为 integer
			if t.child[0].typ != intExpr || t.child[1].typ != intExpr {
				typeError(t, "Op applied to non-integer")
			}
			if t.op == eqToken || t.op == ltToken { // 返回的是比较结果
				t.typ = boolExpr
			} else { // 返回的是计算结果
				t.typ = intExpr
			}
		default: // const 或 id，直接返回类型
			t.typ = intExpr
		}
	case stmtNode:
		switch t.stmt {
		case ifStmt: // if 语句，条件表达式类型为 bool
			if t.child[0].typ != boolExpr {
				typeError(t.child[0], "If test is not Boolean")
			}
		case assignStmt: // 赋值语句，表达式类型为 integer
			if t.child[0].typ != intExpr {
				typeError(t.child[0], "Assignment of non-integer value")
			}
		case writeStmt: // 写语句，表达式类型为 integer
			if t.child[0].typ != intExpr {
				typeError(t.child[0], "Write of non-integer value")
			}
		case repeatStmt: // repeat 语句，条件表达式类型为 bool
			if t.child[1].typ != boolExpr {
				typeError(t.child[1], "Repeat test is not Boolean")
			}
		}
	}
}

// 对整个语法树做类型检查
// 在进行完毕类型推导后，再做节点类型检查
func TypeCheck(t *treeNode) {
	traverse(t, nil, checkNode)
}
