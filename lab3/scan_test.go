package lab3

import (
	"bufio"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	file, err := os.Open("demo.tny")
	if err != nil {
		return
	}
	sourceScanner = bufio.NewScanner(file)
	code := m.Run()
	os.Exit(code)
}

func TestGetNextChar(t *testing.T) {
	var c byte
	var err error
	for err == nil {
		c, err = getNextChar()
		t.Log(string(c))
	}

}

func TestGetToken(t *testing.T) {
	var token tokenType
	for token != eofToken {
		token = GetToken()
		t.Logf("%s", token.String())
	}
}
