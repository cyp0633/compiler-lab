package lab34

import "fmt"

// symtab.go 符号表

// 符号表
var symtab map[string]*struct {
	refLines []int // 引用行号列表
	memLoc   int   // 内存位置
} = make(map[string]*struct {
	refLines []int
	memLoc   int
})

// 当前使用的内存位置
var memLoc int = 0

// 查找符号表
// 返回内存位置和是否找到
func lookupSymtab(name string) (location int, ok bool) {
	if item, ok := symtab[name]; ok {
		return item.memLoc, true
	}
	return 0, false
}

// 写入符号表
// 符号名，当前出现行数
func insertSymtab(name string, lineNo int) {
	if item, ok := symtab[name]; ok {
		item.refLines = append(item.refLines, lineNo)
	} else {
		symtab[name] = &struct {
			refLines []int
			memLoc   int
		}{[]int{lineNo}, memLoc}
		memLoc++
	}
}

func printSymtab() {
	fmt.Printf("Variable name\tLocation\tLine Number\n")
	for name, item := range symtab {
		fmt.Printf("%-14s\t%-8d\t", name, item.memLoc)
		for _, lineNo := range item.refLines {
			fmt.Printf("%-4d ", lineNo)
		}
		fmt.Printf("\n")
	}
}
