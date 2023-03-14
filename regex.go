package main

type RegularExpression struct {
	RegularId int
	Name      string
	// 正则运算符，共有 7 种：‘=’, ‘~’，‘-’, ‘|’，‘.’, ‘*’, ‘+’, ‘?’
	OperatorSymbol rune
	// 左操作数
	OperandID1 int
	// 右操作数
	OperandID2 int
	// 左操作数类型
	Type1 OperandType
	// 右操作数类型
	Type2 OperandType
	// 运算结果类别
	ResultType OperandType
	// 词的 category 属性值
	Category LexemeCategory
	// 对应的 NFA
	PNFA *Graph
}

// OperandType 表示操作数的类型
//
// 不一定是那些 iota，也可能是 ASCII 数值（type1 = CHAR）
type OperandType int

const (
	// 一元运算，无右操作数
	OperandNull OperandType = iota
	// 字符
	OperandChar
	// 字符集
	OperandCharset
	// 正则表达式
	OperandRegular
)

type LexemeCategory int

const (
	// 无 category 值
	LexemeNull LexemeCategory = iota
	// 整数常量
	LexemeIntegerConst
	// 实数常量
	LexemeFloatConst
	// 科学计数法常量
	LexemeScientificConst
	// 数值运算词
	LexemeNumericOperator
	// 注释
	LexemeNote
	// 字符串常量
	LexemeStringConst
	// 空格常量
	LexemeSpaceConst
	// 比较运算词
	LexemeCompareOperator
	// 变量词
	LexemeID
	// 逻辑运算词
	LexemeLogicOperator
)

func (l LexemeCategory) String() string {
	switch l {
	case LexemeNull:
		return "null"
	case LexemeIntegerConst:
		return "integer"
	case LexemeFloatConst:
		return "float"
	case LexemeScientificConst:
		return "scientific"
	case LexemeNumericOperator:
		return "numeric operator"
	case LexemeNote:
		return "note"
	case LexemeStringConst:
		return "string"
	case LexemeSpaceConst:
		return "space"
	case LexemeCompareOperator:
		return "compare operator"
	case LexemeID:
		return "id"
	case LexemeLogicOperator:
		return "logic operator"
	default:
		return "unknown"
	}
}
