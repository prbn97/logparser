package main

import (
	"fmt"
	"paulo/parser"
)

func main() {

	lp := parser.New("log_files")
	err := lp.ParseLines("log.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	lp.PrintErrorLog()
}
