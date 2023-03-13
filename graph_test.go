package main

import "testing"

func TestGenerateBasicNFA(t *testing.T) {
	_ = generateBasicNFA(DriverChar, 1)
}

func TestUnionNFA(t *testing.T) {
	_ = unionNFA(&Graph{}, &Graph{})
}

func TestProductNFA(t *testing.T) {
	_ = productNFA(&Graph{}, &Graph{})
}

func TestPlusClosureNFA(t *testing.T) {
	_ = plusClosureNFA(&Graph{})
}
