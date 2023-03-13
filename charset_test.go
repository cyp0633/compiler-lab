package main

import "testing"

func TestRangeChars(t *testing.T) {
	rangeChars('a', 'z')
}

func TestUnionChars(t *testing.T) {
	unionChars('a', 'z')
}

func TestUnionCharsetAndChar(t *testing.T) {
	unionCharsetAndChar(1, 'a')
}

func TestUnionTwoCharsets(t *testing.T) {
	unionTwoCharsets(1, 2)
}

func TestDifference(t *testing.T) {
	difference(1, 'a')
}
