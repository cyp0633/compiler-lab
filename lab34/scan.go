// 实验三：词法分析器构造工具的实现
//
// 实验四：TINY 语言编译器的实现
package lab34

// scan.go 词法分析器

import (
	"bufio"
	"fmt"
	"io"
)

// 词法分析器 DFA 的状态
type stateType int

// 词法分析器 DFA 的状态
const (
	start     stateType = iota
	inAssign            // :=
	inComment           // { ... }
	inNum               // 123
	inId                // abc
	done                // 结束
)

// 当前行
var lineBuf string = ""

// lineBuf 中当前字符的位置
var linePos int = 0

// 是否已经读取到文件末尾
var eof bool = false

// 源代码 scanner
var sourceScanner *bufio.Scanner

// 当前 token 内容
var tokenString = ""

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
func getNextChar() (byte, error) {
	// 本行已经读取完毕
	if linePos >= len(lineBuf) {
		lineNo++
		ok := sourceScanner.Scan()
		if !ok { // 文件末尾
			err := sourceScanner.Err()
			if err == nil {
				// EOF 的话不会返回 err
				eof = true
				return 0, io.EOF
			} else {
				// 真的出现错误了
				lineBuf = ""
				linePos = 0
				return 0, err
			}
		}
		lineBuf = sourceScanner.Text()
		linePos = 1
		return lineBuf[0], nil
	} else {
		linePos++
		return lineBuf[linePos-1], nil
	}
}

// ungetNextChar 将 lineBuf 中的一个字符退回
func ungetNextChar() {
	if !eof {
		linePos--
	}
}

// getToken 从源文件中读取一个 token
func GetToken() tokenType {
	// 当前状态
	state := start
	// 当前 token 类型
	var currToken tokenType
	// 清空当前 token 字符串
	tokenString = ""

	// 是否将当前字符保存到 tokenString 中
	// 如一串数字要保存，但空格啊回车之类的不用
	save := false
	for state != done {
		save = true
		c, err := getNextChar()
		// 根据当前状态和读取到的字符决定下一步的操作
		switch state {
		case start:
			if c >= '0' && c <= '9' {
				// 数字
				state = inNum
			} else if isAlpha(c) {
				// 字母
				// 先输进来再判断是不是保留字
				state = inId
			} else if c == '{' {
				// 注释
				save = false
				state = inComment
			} else if c == ':' {
				// 只有赋值才用到冒号
				state = inAssign
			} else if c == ' ' || c == '\n' || c == '\t' {
				// 空格，回车，tab
				save = false
			} else {
				state = done
				if err == io.EOF {
					save = false
					currToken = eofToken
					break
				}
				switch c {
				case '=':
					currToken = eqToken
				case '<':
					currToken = ltToken
				case '+':
					currToken = plusToken
				case '-':
					currToken = minusToken
				case '*':
					currToken = timesToken
				case '/':
					currToken = overToken
				case '(':
					currToken = lparenToken
				case ')':
					currToken = rparenToken
				case ';':
					currToken = semicolonToken
				default:
					currToken = errorToken
				}
			}
		case inComment:
			save = false
			if err == io.EOF {
				state = done
				currToken = eofToken
			} else if c == '}' {
				state = start
			}
		case inAssign:
			state = done
			if c == '=' {
				currToken = assignToken
			} else {
				// 并不是一个赋值语句
				ungetNextChar()
				save = false
				currToken = errorToken
			}
		case inNum:
			if c < '0' || c > '9' {
				// 数字结束
				ungetNextChar()
				save = false
				state = done
				currToken = numToken
			}
		case inId:
			if !isAlpha(c) {
				// 字母结束
				ungetNextChar()
				save = false
				state = done
				currToken = idToken
			}
		default:
			// 包含 done
			// 发生错误
			fmt.Println("Scanner Bug: state = ", state)
			state = done
			currToken = errorToken
		}
		if save {
			tokenString += string(c)
		}
		if state == done {
			if currToken == idToken {
				// 判断是否为保留字
				if tokenType, ok := reservedWords[tokenString]; ok {
					currToken = tokenType
				}
			}
		}
	}
	fmt.Printf("\t%d: %v\t%v\n", lineNo, currToken.String(), tokenString)
	return currToken
}

// 判断是否为字母
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
