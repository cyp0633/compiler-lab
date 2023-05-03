package lab34

import "fmt"

// parse.go 语法分析器

// 表达式数值类型
type exprType int

const (
	intExpr exprType = iota
	boolExpr
	voidExpr
)

// 语句类型，如赋值、循环等
type stmtKind int

const (
	ifStmt stmtKind = iota
	repeatStmt
	assignStmt
	readStmt
	writeStmt
)

// 表达式类型，如运算、常量等
type exprKind int

const (
	opExpr exprKind = iota
	constExpr
	idExpr
)

// 节点类型
type nodeKind int

const (
	stmtNode nodeKind = iota
	exprNode
)

// 语法树节点
type treeNode struct {
	child   [3]*treeNode // 子节点
	sibling *treeNode    // 兄弟节点
	lineNo  int          // 行号
	op      tokenType    // 操作符
	val     int          // 值
	attr    string       // 属性
	stmt    stmtKind     // 语句类型
	expr    exprKind     // 表达式类型
	typ     exprType     // 表达式数值类型
	node    nodeKind     // 节点类型
}

// 当前 token
var currToken tokenType

func parse() (t *treeNode) {
	token := GetToken()
	t = stmtSequence()
	if token != eofToken {
		syntaxError("Code ends before file\n")
	}
	return
}

// 用于显示语法错误
func syntaxError(msg string) {
	fmt.Printf("line %d: %s\n", lineNo, msg)
}

// 试图匹配一个 token，如果匹配，则往下读
func match(expected tokenType) {
	if currToken == expected {
		currToken = GetToken()
	} else {
		syntaxError("unexpected token -> %s" + currToken.String())
	}
}

func stmtSequence() *treeNode {
	t := statement()
	p := t
	for currToken != eofToken && currToken != endToken && currToken != elseToken && currToken != untilToken {
		match(semicolonToken)
		q := statement()
		if q != nil {
			if t == nil {
				p = q
				t = p
			} else {
				p.sibling = q
				p = q
			}
		}
	}
	return t
}

// 根据不同的语句返回对应的语法树节点
func statement() *treeNode {
	var t *treeNode
	switch currToken {
	case ifToken:
		t = ifStatement()
	case repeatToken:
		t = repeatStatement()
	case idToken:
		t = assignStatement()
	case readToken:
		t = readStatement()
	case writeToken:
		t = writeStatement()
	default:
		syntaxError("unexpected token -> %s" + currToken.String())
		currToken = GetToken()
	}
	return t
}

func ifStatement() *treeNode {
	t := newStatementNode(ifStmt)
	match(ifToken)
	t.child[0] = expression()
	match(thenToken)
	t.child[1] = stmtSequence()
	if currToken == elseToken {
		match(elseToken)
		t.child[2] = stmtSequence()
	}
	match(endToken)
	return t
}

func repeatStatement() *treeNode {
	return nil
}

func assignStatement() *treeNode {
	return nil
}

func readStatement() *treeNode {
	return nil
}

func writeStatement() *treeNode {
	return nil
}

func expression() *treeNode {
	return nil
}

func simpleExpression() *treeNode {
	return nil
}

func term() *treeNode {
	return nil
}

func factor() *treeNode {
	return nil
}

// 新建一个语法树声明节点
func newStatementNode(kind stmtKind) *treeNode {
	return &treeNode{
		child:   [3]*treeNode{},
		sibling: nil,
		node:    stmtNode,
		stmt:    kind,
		lineNo:  lineNo,
	}
}

// 新建一个语法树表达式节点
func newExpressionNode(kind exprKind) *treeNode {
	return &treeNode{
		child:   [3]*treeNode{},
		sibling: nil,
		node:    exprNode,
		expr:    kind,
		lineNo:  lineNo,
		typ:     voidExpr,
	}
}
