package main

// 字符集
type Charset struct {
	// 字符集 ID
	IndexID int
	// 段 ID。对每个字符集，段 ID 从 0 递增
	SegmentID int
	// 起始字符
	FromChar rune
	// 终止字符
	ToChar rune
}

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

// 字符集表
var CharsetTable []*Charset

// 字符的范围运算
//
// 生成一个字符集，包含从 fromChar 到 toChar 的所有字符
func rangeChars(fromChar rune, toChar rune) (indexID int) {
	var c = Charset{FromChar: fromChar, ToChar: toChar}
	c.IndexID = CharsetTable[len(CharsetTable)-1].SegmentID + 1
	c.SegmentID = 0
	CharsetTable = append(CharsetTable, &c)
	return c.IndexID
}

// 字符集的并运算
//
// 生成一个字符集，包含 c1 和 c2
func unionChars(c1 rune, c2 rune) (indexID int) {
	var cs1 = Charset{FromChar: c1, ToChar: c1}
	cs1.IndexID = CharsetTable[len(CharsetTable)-1].SegmentID + 1
	cs1.SegmentID = 0
	var cs2 = Charset{FromChar: c2, ToChar: c2}
	cs2.IndexID = CharsetTable[len(CharsetTable)-1].SegmentID + 1
	cs2.SegmentID = 0
	CharsetTable = append(CharsetTable, &cs1, &cs2)
	return cs1.IndexID
}

// 字符集与字符之间的并运算
//
// 将原有的字符集与新字符合并
func unionCharsetAndChar(indexID int, c rune) (newIndexID int) {
	var cs = CharsetTable[indexID]
	switch {
	// c 在 cs 中，直接拷贝
	case c >= cs.FromChar && c <= cs.ToChar:
	// c 在 cs 正前面，合并
	case c == cs.FromChar-1:
		cs.FromChar = c
	// c 在 cs 正后面，合并
	case c == cs.ToChar+1:
		cs.ToChar = c
	// 俩不挨着
	default:
	}
	return
}
