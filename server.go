package main

import (
	"fmt"
	"os"

	"github.com/echojc/lernkartei/dict/glosbe"
)

func main() {
	for _, arg := range os.Args[1:] {
		words, _ := glosbe.NewWord(arg)
		for _, word := range words {
			fmt.Printf("%+v\n", word)
		}
	}
}
