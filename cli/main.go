package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aquilax/crossword"
)

func getWord(w string) string {
	if strings.Contains(w, "/") {
		w = strings.Split(w, "/")[0]
	}
	return strings.ToLower(strings.TrimSpace(w))
}

func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var words []crossword.Word

	cols := 30
	rows := 30

	limit := max(cols, rows)

	file, err := os.Open(os.Args[1])
	defer file.Close()

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		word := getWord(line)
		if len([]rune(word)) <= limit {
			words = append(words, crossword.Word{word, word})
		}
	}
	cw := crossword.New(cols, rows, words)
	cw.Generate(0, 2)
	fmt.Print(cw)
}
