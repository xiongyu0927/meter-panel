package main

import "fmt"

func main() {
	test("A")
}

func test(ss string) {
	if ss == "A" {
		var x int
	}
	if ss == "B" {
		var x string
	}
	fmt.Println(x)
}
