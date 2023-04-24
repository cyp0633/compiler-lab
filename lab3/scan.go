// 实验三：词法分析器构造工具的实现
package lab3

import "fmt"

// 词法分析器 DFA 的状态
type stateType int

// 词法分析器 DFA 的状态
const (
	start     stateType = iota
	inAssign            // :=
	inComment           // { ... }
	inNum               // 123
	inId                // abc
	inDone              // 结束
)

// 当前行
var lineBuf string = ""

// lineBuf 中当前字符的位置
var linePos int = 0

// 当前行数
var lineNo int = 0

// 是否已经读取到文件末尾
var eof bool = false

// 保留词表
var reservedWords = map[string]tokenType{
	"if":     ifToken,
	"then":   thenToken,
	"else":   elseToken,
	"end":    endToken,
	"repeat": repeatToken,
	"until":  untilToken,
	"read":   readToken,
	"write":  writeToken,
}
// map 类型又不需要查找函数了

// getNextChar 从 lineBuf 中读取一个字符，
// 如果 lineBuf 为空则从输入流中读取一行
func getNextChar() int {
	// 本行已经读取完毕
	if linePos >= len(lineBuf) {
		lineNo++
		// TODO: 替换为从文件中读取一行
		n, err := fmt.Scanf("%s", &lineBuf)
		if err != nil {
			fmt.Println(err)
		}
		if n == 0 { // 读取到文件尾
			eof = true
			return -1
		}
		linePos = 1
		return int(lineBuf[0])
	} else {
		linePos++
		return int(lineBuf[linePos-1])
	}
}

// ungetNextChar 将 lineBuf 中的一个字符退回
func ungetNextChar() {
	linePos--
}
