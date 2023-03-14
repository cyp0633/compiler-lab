package main

import (
	"fmt"
	"testing"
)

func TestRangeChars(t *testing.T) {
	id := rangeChars('a', 'z')
	printCharset(id)
	// 结果：1 个段，a-z
	if len(CharsetTable) < 1 || CharsetTable[0].FromChar != 'a' || CharsetTable[0].ToChar != 'z' {
		t.Fail()
	}
}

func TestUnionChars(t *testing.T) {
	id := unionChars('a', 'z')
	printCharset(id)
	// 结果：2 个段，a-a、z-z
	if len(CharsetTable) < 2 || CharsetTable[0].FromChar != 'a' || CharsetTable[0].ToChar != 'a' || CharsetTable[1].FromChar != 'z' || CharsetTable[1].ToChar != 'z' {
		t.Fail()
	}
}

func TestUnionCharsetAndChar(t *testing.T) {
	// c 在 charset 中
	baseID := rangeChars('a', 'c')
	id := unionCharsetAndChar(baseID, 'c')
	printCharset(id)
	// 结果：1 个段，a-c
	if len(CharsetTable) != 2 || CharsetTable[1].FromChar != 'a' || CharsetTable[1].ToChar != 'c' {
		t.Fail()
	}

	// c 在 charset 外的边缘
	id = unionCharsetAndChar(baseID, 'd')
	printCharset(id)
	// 结果：1 个段，a-d
	if len(CharsetTable) != 3 || CharsetTable[2].FromChar != 'a' || CharsetTable[2].ToChar != 'd' {
		t.Fail()
	}

	// c 在 charset 外，不搭边
	id = unionCharsetAndChar(baseID, 'e')
	printCharset(id)
	// 结果：2 个段，a-c, e-e
	if len(CharsetTable) != 5 || CharsetTable[3].FromChar != 'a' || CharsetTable[3].ToChar != 'c' || CharsetTable[4].FromChar != 'e' || CharsetTable[4].ToChar != 'e' {
		t.Fail()
	}
}

func TestUnionTwoCharsets(t *testing.T) {
	// 不搭边
	id1 := rangeChars('b', 'e')
	id2 := rangeChars('h', 'l')
	id3 := unionTwoCharsets(id1, id2)
	printCharset(id3)
	// 结果：2 个段，b-e， h-l

	// 两个长一点的
	id4 := rangeChars('a', 'c')
	id5 := rangeChars('g', 'z')
	id6 := unionTwoCharsets(id4, id5)
	printCharset(id6)
	// 结果：2 个段，a-c， g-z

	// 有多段的
	id7 := unionTwoCharsets(id3, id6)
	printCharset(id7)
	// 结果：2 个段，a-e， g-z

	// 后面跟前面重合的
	id8 := rangeChars('e', 'g')
	id9 := unionTwoCharsets(id4, id8)
	id10 := unionTwoCharsets(id3, id9)
	printCharset(id10)
	// 结果：1 个段，a-l
}

func TestDifference(t *testing.T) {
	// 不在里面的
	id1 := rangeChars('a', 'e')
	id2 := difference(id1, 'f')
	printCharset(id2)
	// 结果：1 个段，a-e

	// 在里面的
	id3 := difference(id1, 'c')
	printCharset(id3)
	// 结果：2 个段，a-b, d-e

	// 在边上的
	id4 := difference(id1, 'a')
	printCharset(id4)
	// 结果：1 个段，b-e
}

// 打印一个 charset
func printCharset(indexID int) {
	fmt.Printf("Printing indexID #%v:\n", indexID)
	for _, csTemp := range CharsetTable {
		if csTemp.IndexID == indexID {
			fmt.Printf("segment #%v: %c-%c\n", csTemp.SegmentID, csTemp.FromChar, csTemp.ToChar)
		}
	}
}
