package lab34

import (
	"bufio"
	"fmt"
	"os"
)

func OpenFile() {
	// 打开源代码 code.tny
	file, err := os.Open("demo.tny")
	if err != nil {
		_ = fmt.Errorf("failed to open file: %s", err.Error())
		panic(err)
	}
	sourceScanner = bufio.NewScanner(file)
	// 打开输出文件 code.tm
	file, err = os.Create("demo.tm")
	if err != nil {
		_ = fmt.Errorf("failed to create file: %s", err.Error())
		panic(err)
	}

	outputFile = bufio.NewWriter(file)
}
