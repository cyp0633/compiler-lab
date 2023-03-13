package main

import "testing"

func TestGenerateBasicNFA(t *testing.T) {
	_ = generateBasicNFA(DriverChar, 1)
}

func TestUnionNFA(t *testing.T) {
	_ = unionNFA(&Graph{}, &Graph{})
}
