package main

import (
	"compiler-lab/lab34"
	"fmt"
)

func main() {
	lab34.OpenFile()
	fmt.Println("------Parsing------")
	tree := lab34.Parse()
	fmt.Println("------Syntax tree------")
	lab34.PrintTree(tree, 0)
	fmt.Println("------Semantic analysis------")
	lab34.BuildSymtab(tree)
	lab34.TypeCheck(tree)
	fmt.Println("------Code generation------")
	lab34.CodeGen(tree)
}
