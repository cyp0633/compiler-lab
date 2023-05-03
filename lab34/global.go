package lab34

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

func (t *tokenType) String() string {
	switch *t {
	case ifToken:
		return "if"
	case thenToken:
		return "then"
	case elseToken:
		return "else"
	case endToken:
		return "end"
	case repeatToken:
		return "repeat"
	case untilToken:
		return "until"
	case readToken:
		return "read"
	case writeToken:
		return "write"
	case idToken:
		return "id"
	case numToken:
		return "num"
	case assignToken:
		return "assign"
	case eqToken:
		return "equal"
	case ltToken:
		return "less than"
	case plusToken:
		return "plus"
	case minusToken:
		return "minus"
	case timesToken:
		return "times"
	case overToken:
		return "over"
	case lparenToken:
		return "left paren"
	case rparenToken:
		return "right paren"
	case semicolonToken:
		return "semicolon"
	case eofToken:
		return "end of file"
	case errorToken:
		return "error"
	default:
		return "unknown token"
	}
}
