package main

import (
	"fmt"
	"testing"
)

func TestRangeChars(t *testing.T) {
	id := rangeChars('a', 'z')
	printCharset(id)
}

func TestUnionChars(t *testing.T) {
	id := unionChars('a', 'z')
	printCharset(id)
}

func TestUnionCharsetAndChar(t *testing.T) {
	// c 在 charset 中
	baseID := rangeChars('a', 'c')
	id := unionCharsetAndChar(baseID, 'c')
	printCharset(id)

	// c 在 charset 外的边缘
	id = unionCharsetAndChar(baseID, 'd')
	printCharset(id)

	// c 在 charset 外，不搭边
	id = unionCharsetAndChar(baseID, 'e')
	printCharset(id)
}

func TestUnionTwoCharsets(t *testing.T) {
	// 不搭边
	id1 := rangeChars('b', 'e')
	id2 := rangeChars('h', 'l')
	id3 := unionTwoCharsets(id1, id2)
	printCharset(id3)

	// 两个长一点的
	id4 := rangeChars('a', 'c')
	id5 := rangeChars('g', 'z')
	id6 := unionTwoCharsets(id4, id5)
	printCharset(id6)

	// 有多段的
	id7 := unionTwoCharsets(id3, id6)
	printCharset(id7)
}

func TestDifference(t *testing.T) {
	difference(1, 'a')
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
