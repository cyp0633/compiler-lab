package lab3

type tokenType int

const (
	// 保留词

	ifToken tokenType = iota
	thenToken
	elseToken
	endToken
	repeatToken
	untilToken
	readToken
	writeToken

	// 多字符 token

	idToken
	numToken

	// 特殊符号

	assignToken    // :=
	eqToken        // =
	ltToken        // <
	plusToken      // +
	minusToken     // -
	timesToken     // *
	overToken      // /
	lparenToken    // (
	rparenToken    // )
	semicolonToken // ;

	// 编译器情况
	eofToken   // 文件末尾
	errorToken // 出错
)

// 当前行数
var lineNo int = 0
