package lab34

import (
	"fmt"
	"strconv"
)

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

func Parse() (t *treeNode) {
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

// if 语句的匹配
//
// if <expression> then <stmtSequence> [else <stmtSequence>] end
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

// repeat 语句的匹配
//
// repeat <stmtSequence> until <expression>
func repeatStatement() *treeNode {
	t := newStatementNode(repeatStmt)
	match(repeatToken)
	t.child[0] = stmtSequence()
	match(untilToken)
	t.child[1] = expression()
	return t
}

// 赋值语句的匹配
//
// <id> := <expression>
func assignStatement() *treeNode {
	t := newStatementNode(assignStmt)
	if currToken == idToken {
		t.attr = tokenString
	}
	match(idToken)
	match(assignToken)
	t.child[0] = expression()
	return t
}

// read 语句的匹配
//
// read <id>
func readStatement() *treeNode {
	t := newStatementNode(readStmt)
	match(readToken)
	if currToken == idToken {
		t.attr = tokenString
	}
	match(idToken)
	return t
}

// write 语句的匹配
//
// write <expression>
func writeStatement() *treeNode {
	t := newStatementNode(writeStmt)
	match(writeToken)
	t.child[0] = expression()
	return t
}

// 表达式的匹配：处理 < 和 = 运算符
// 优先级最低
func expression() *treeNode {
	t := simpleExpression()
	if currToken == ltToken || currToken == eqToken {
		p := newExpressionNode(opExpr)
		p.child[0] = t
		p.op = currToken
		t = p
		match(currToken)
		t.child[1] = simpleExpression()
	}
	return t
}

// 表达式的匹配：处理 + 和 - 运算符
// 优先级中等
func simpleExpression() *treeNode {
	t := term()
	for currToken == plusToken || currToken == minusToken {
		p := newExpressionNode(opExpr)
		p.child[0] = t
		p.op = currToken
		t = p
		match(currToken)
		t.child[1] = term()
	}
	return t
}

// 表达式的匹配：处理 * 和 / 运算符
// 优先级最高
func term() *treeNode {
	t := factor()
	for currToken == timesToken || currToken == overToken {
		p := newExpressionNode(opExpr)
		p.child[0] = t
		p.op = currToken
		t = p
		match(currToken)
		t.child[1] = factor()
	}
	return t
}

// 表达式的匹配：处理数字、标识符和括号
func factor() (t *treeNode) {
	switch currToken {
	case numToken:
		t = newExpressionNode(constExpr)
		if currToken == numToken {
			var err error
			t.val, err = strconv.Atoi(tokenString)
			if err != nil {
				syntaxError("unexpected token -> %s" + currToken.String())
			}
			match(numToken)
		}
	case idToken:
		t = newExpressionNode(idExpr)
		if currToken == idToken {
			t.attr = tokenString
		}
		match(idToken)
	case lparenToken:
		match(lparenToken)
		t = expression()
		match(rparenToken)
	default:
		syntaxError("unexpected token -> %s" + currToken.String())
		currToken = GetToken()
	}
	return
}

// 新建一个语法树语句节点
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
