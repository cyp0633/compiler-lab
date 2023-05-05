package main

import "compiler-lab/lab34"

func main() {
	lab34.OpenFile()
	tree := lab34.Parse()
	lab34.BuildSymtab(tree)
	lab34.TypeCheck(tree)
	lab34.CodeGen(tree)
}
