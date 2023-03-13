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
	// 原来的字符集可能不只有一段
	var oldCharset []*Charset
	maxID := CharsetTable[len(CharsetTable)-1].IndexID
	// 将老的字符集拷贝一份（不懂为什么非要创建新的）
	newCharset := copyCharset(oldCharset, maxID+1)
	for _, csTemp := range CharsetTable {
		if csTemp.IndexID == indexID {
			oldCharset = append(oldCharset, csTemp)
		}
	}
	if len(oldCharset) == 0 {
		return -1
	}
	// 遍历老的字符集各段，看看能不能合并进去
	for _, csTemp := range oldCharset {
		// 在中间，无需其他操作即可合并
		if c >= csTemp.FromChar && c <= csTemp.ToChar {
			CharsetTable = append(CharsetTable, newCharset...)
			return maxID + 1
		}
	}
	// 看看是不是在边缘，合进来
	// 不能与上个放在同一个循环里，防止两个重合
	for index, csTemp := range oldCharset {
		if c == csTemp.FromChar-1 || c == csTemp.ToChar+1 {
			// 将新字符集对应段的边缘合并
			if c == csTemp.FromChar-1 {
				newCharset[index].FromChar = c
			} else {
				newCharset[index].ToChar = c
			}
			CharsetTable = append(CharsetTable, newCharset...)
			return maxID + 1
		}
	}
	// 完全没法合并了
	// 按顺序插入新的字符集
	for index, csTemp := range newCharset {
		if c < csTemp.FromChar {
			var csNew = Charset{IndexID: maxID + 1, SegmentID: index, FromChar: c, ToChar: c}
			newCharset = append(newCharset[:index], append([]*Charset{&csNew}, newCharset[index:]...)...)
			// 重新设置段 ID
			for i := index + 1; i < len(newCharset); i++ {
				newCharset[i].SegmentID = i
			}
			CharsetTable = append(CharsetTable, newCharset...)
			return maxID + 1
		}
	}
	// 也许是最后一个
	var csNew = Charset{IndexID: maxID + 1, SegmentID: len(newCharset), FromChar: c, ToChar: c}
	newCharset = append(newCharset, &csNew)
	CharsetTable = append(CharsetTable, newCharset...)
	return maxID + 1
}

// 将一个字符集复制一份
func copyCharset(oldCharset []*Charset, newIndex int) (newCharset []*Charset) {
	for _, csTemp := range oldCharset {
		var csNew = Charset{IndexID: newIndex, SegmentID: csTemp.SegmentID, FromChar: csTemp.FromChar, ToChar: csTemp.ToChar}
		newCharset = append(newCharset, &csNew)
	}
	return
}
