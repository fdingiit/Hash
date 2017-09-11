package main

import "fmt"

func testFind(h *Hash, str string) {
	var ok bool
	var err error

	if ok, err = h.Find(str); err != nil {
		panic(err)
	} else if ok {
		fmt.Println("Found", str)
	} else {
		fmt.Println("Did not find", str)
	}
}

func testRemove(h *Hash, str string) {
	var rStr string
	var err error

	rStr, err = h.Remove(str)
	if err != nil {
		panic(err)
	}
	if rStr == "" {
		fmt.Println("String", str, "not found, not removed")
	} else {
		fmt.Println("Removed string =", str)
	}
}
