package main

import (
	"fmt"
)

func main() {

	lp := NewParser("log_files")
	err := lp.ParseLines("log.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// lp.PrintErrorLog()
	// lp.PrintWarnLog()
	// lp.PrintInfoLog()
	lp.MostFrequentIDs()
}
