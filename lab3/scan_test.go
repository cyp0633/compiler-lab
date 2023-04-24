package lab3

import (
	"bufio"
	"os"
	"testing"
)

func TestGetNextChar(t *testing.T) {
	file, err := os.Open("demo.tny")
	if err != nil {
		t.Error("Open file failed", err)
	}
	sourceScanner = bufio.NewScanner(file)
	var c byte
	for err == nil {
		c, err = getNextChar()
		t.Log(string(c))
	}
	
}
