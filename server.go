package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/echojc/lernkartei/dict"
	"github.com/echojc/lernkartei/dict/glosbe"
)

func main() {
	var out []dict.Word
	for _, arg := range os.Args[1:] {
		words, _ := glosbe.NewWord(arg)
		out = append(out, words...)
	}

	bs, _ := json.Marshal(out)
	fmt.Println(string(bs))
}
