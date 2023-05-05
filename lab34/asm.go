package lab34

import (
	"fmt"
	"io"
)

// asm.go 产生 TM 虚拟机可运行的代码，对应原 code.c

var emitLoc = 0
var highEmitLoc = 0

// 汇编代码目标文件
var outputFile io.Writer

// 输出注释
//
// * <comment>
func emitComment(comment string) {
	fmt.Fprintf(outputFile, "* %s\n", comment)
}

// 寄存器间的操作(register only)，不含注释的输出
//
// <emitLoc>: <op> <target>,<reg1>,<reg2>
func emitRO(op string, target, reg1, reg2 int) {
	fmt.Fprintf(outputFile, "%3d:  %5s  %d,%d,%d\n", emitLoc, op, target, reg1, reg2)
	emitLoc++
	if highEmitLoc < emitLoc {
		highEmitLoc = emitLoc
	}
}

// 寄存器和内存的操作(register & memory)，不含注释的输出
//
// <emitLoc>: <op> <target>,<offset>(<base>)
func emitRM(op string, target, offset, base int) {
	fmt.Fprintf(outputFile, "%3d:  %5s  %d,%d(%d)\n", emitLoc, op, target, offset, base)
	emitLoc++
	if highEmitLoc < emitLoc {
		highEmitLoc = emitLoc
	}
}

// 寄存器和内存的操作，但引用内存中的绝对地址
func emitAbsRM(op string, target, loc int) {
	fmt.Fprintf(outputFile, "%3d:  %5s  %d,%d(%d)\n", emitLoc, op, target, loc, 0)
	emitLoc++
	if highEmitLoc < emitLoc {
		highEmitLoc = emitLoc
	}
}

// 留出特定长度的汇编代码空间，并返回跳过前的输出位置
func emitSkip(len int) int {
	i := emitLoc
	emitLoc += len
	if highEmitLoc < emitLoc {
		highEmitLoc = emitLoc
	}
	return i
}

// 回填
func emitBackup(loc int) {
	if loc > highEmitLoc {
		emitComment("BUG in emitBackup") // 不能回填到没达到的我i之
	}
	emitLoc = loc
}

// 恢复
func emitRestore() {
	emitLoc = highEmitLoc
}
