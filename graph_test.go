package main

import "testing"

func TestGenerateBasicNFA(t *testing.T) {
	_ = generateBasicNFA(DriverChar, 1)
}
