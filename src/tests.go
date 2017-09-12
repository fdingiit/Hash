package main

import (
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		panic(os.Args)
	}

	funcs := []func(){nil, testcase1, testcase2, testcase3, testcase4}
	tc, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if tc > len(funcs)-1 || tc < 1 {
		panic(tc)
	}

	funcs[tc]()
}
