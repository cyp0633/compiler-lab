package main

import (
	"fmt"
	"testing"
)

func TestRangeChars(t *testing.T) {
	base := len(CharsetTable)
	id := rangeChars('a', 'z')
	printCharset(id)
	// 结果：1 个段，a-z
	if len(CharsetTable) != 1+base || CharsetTable[base].FromChar != 'a' || CharsetTable[base].ToChar != 'z' {
		t.Fail()
	}
}

func TestUnionChars(t *testing.T) {
	base := len(CharsetTable)
	id := unionChars('a', 'z')
	printCharset(id)
	// 结果：2 个段，a-a、z-z
	if len(CharsetTable) != base+2 || CharsetTable[base].FromChar != 'a' || CharsetTable[base].ToChar != 'a' || CharsetTable[base+1].FromChar != 'z' || CharsetTable[base+1].ToChar != 'z' {
		t.Fail()
	}

	id = unionChars('a', 'a')
	printCharset(id)
	// 结果：1 个段，a-a
	if len(CharsetTable) != base+3 || CharsetTable[base+2].FromChar != 'a' || CharsetTable[base+2].ToChar != 'a' {
		t.Fail()
	}
}

func TestUnionCharsetAndChar(t *testing.T) {
	base := len(CharsetTable)

	// c 在 charset 中
	baseID := rangeChars('a', 'c')
	id := unionCharsetAndChar(baseID, 'c')
	printCharset(id)
	// 结果：1 个段，a-c
	if len(CharsetTable) != base+2 || CharsetTable[base+1].FromChar != 'a' || CharsetTable[base+1].ToChar != 'c' {
		t.Fail()
	}

	// c 在 charset 外的边缘
	id = unionCharsetAndChar(baseID, 'd')
	printCharset(id)
	// 结果：1 个段，a-d
	if len(CharsetTable) != base+3 || CharsetTable[base+2].FromChar != 'a' || CharsetTable[base+2].ToChar != 'd' {
		t.Fail()
	}

	// c 在 charset 外，不搭边
	id = unionCharsetAndChar(baseID, 'e')
	printCharset(id)
	// 结果：2 个段，a-c, e-e
	if len(CharsetTable) != base+5 || CharsetTable[base+3].FromChar != 'a' || CharsetTable[base+3].ToChar != 'c' || CharsetTable[base+4].FromChar != 'e' || CharsetTable[base+4].ToChar != 'e' {
		t.Fail()
	}
}

func TestUnionTwoCharsets(t *testing.T) {
	base := len(CharsetTable)

	// 不搭边
	id1 := rangeChars('b', 'e')
	id2 := rangeChars('h', 'l')
	id3 := unionTwoCharsets(id1, id2)
	printCharset(id3)
	// 结果：2 个段，b-e， h-l
	if len(CharsetTable) != base+4 || CharsetTable[base+2].FromChar != 'b' || CharsetTable[base+2].ToChar != 'e' || CharsetTable[base+3].FromChar != 'h' || CharsetTable[base+3].ToChar != 'l' {
		t.Errorf("Case 1 failed")
	}

	// 两个长一点的
	id4 := rangeChars('a', 'c')
	id5 := rangeChars('g', 'z')
	id6 := unionTwoCharsets(id4, id5)
	printCharset(id6)
	// 结果：2 个段，a-c， g-z
	if len(CharsetTable) != base+8 || CharsetTable[base+6].FromChar != 'a' || CharsetTable[base+6].ToChar != 'c' || CharsetTable[base+7].FromChar != 'g' || CharsetTable[base+7].ToChar != 'z' {
		t.Errorf("Case 2 failed")
	}

	// 有多段的
	id7 := unionTwoCharsets(id3, id6)
	printCharset(id7)
	// 结果：2 个段，a-e， g-z
	if len(CharsetTable) != base+10 || CharsetTable[base+8].FromChar != 'a' || CharsetTable[base+8].ToChar != 'e' || CharsetTable[base+9].FromChar != 'g' || CharsetTable[base+9].ToChar != 'z' {
		t.Errorf("Case 3 failed")
	}

	// 后面跟前面重合的
	id8 := rangeChars('e', 'g')
	id9 := unionTwoCharsets(id4, id8)
	id10 := unionTwoCharsets(id3, id9)
	printCharset(id10)
	// 结果：1 个段，a-l
	if len(CharsetTable) != base+14 || CharsetTable[base+13].FromChar != 'a' || CharsetTable[base+13].ToChar != 'l' {
		t.Errorf("Case 4 failed")
	}
}

func TestDifference(t *testing.T) {
	base := len(CharsetTable)

	// 不在里面的
	id1 := rangeChars('a', 'e')
	id2 := difference(id1, 'f')
	printCharset(id2)
	// 结果：1 个段，a-e
	if len(CharsetTable) != base+2 || CharsetTable[base+1].FromChar != 'a' || CharsetTable[base+1].ToChar != 'e' {
		t.Errorf("Case 1 failed")
	}

	// 在里面的
	id3 := difference(id1, 'c')
	printCharset(id3)
	// 结果：2 个段，a-b, d-e
	if len(CharsetTable)!=base+4 || CharsetTable[base+2].FromChar != 'a' || CharsetTable[base+2].ToChar != 'b' || CharsetTable[base+3].FromChar != 'd' || CharsetTable[base+3].ToChar != 'e' {
		t.Errorf("Case 2 failed")
	}

	// 在边上的
	id4 := difference(id1, 'a')
	printCharset(id4)
	// 结果：1 个段，b-e
	if len(CharsetTable)!=base+5 || CharsetTable[base+4].FromChar != 'b' || CharsetTable[base+4].ToChar != 'e' {
		t.Errorf("Case 3 failed")
	}
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
