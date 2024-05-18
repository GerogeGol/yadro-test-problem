package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GerogeGol/yadro-test-problem/domain/scan"
)

func main() {
	if len(os.Args) == 1 {
		panic("no specified file")
	}
	filepath := os.Args[1]

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := &strings.Builder{}
	line, err := scan.ScanInputData(file, buf)
	if err != nil {
		fmt.Println(line)
		return
	}
	fmt.Print(buf)

}
