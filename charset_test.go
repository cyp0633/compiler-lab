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
	id := unionCharsetAndChar(baseID, 'b')
	printCharset(id)
}

func TestUnionTwoCharsets(t *testing.T) {
	unionTwoCharsets(1, 2)
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
