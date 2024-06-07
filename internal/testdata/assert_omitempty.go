package main

import (
	"encoding/json"
	"fmt"
	"testdata/a"
)

// This binary serves to assert that all of the fields that are marked omitempty
// in the testdata are actually omitted when encoding json.

func main() {
	g := a.Good{
		SliceCap0: make([]int, 0, 10),
	}
	jb, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(jsonOrder(jb)))
}

// Force consistent field ordering for string comparison in test.
func jsonOrder(bytes []byte) []byte {
	var a any
	err := json.Unmarshal(bytes, &a)
	if err != nil {
		panic(err.Error())
	}
	b, err := json.Marshal(a)
	if err != nil {
		panic(err.Error())
	}
	return b
}
