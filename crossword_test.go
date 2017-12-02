package crossword

import (
	"math/rand"
	"testing"
)

func TestNew(t *testing.T) {
	cols := 10
	rows := 20
	cr := New(cols, rows, []Word{
		Word{"a", "a"},
		Word{"aa", "aa"},
	})
	if rows != len(cr.grid) {
		t.Error("Wrong rows")
	}
	col := (cr.grid)[0]
	if cols != len(col) {
		t.Error("Wrong columns")
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func generateWord(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateWords(n, max int) []Word {
	var words []Word
	var word Word
	for i := 0; i < n; i++ {
		word = Word{
			generateWord((i%max - 3) + 3),
			"",
		}
		words = append(words, word)
	}
	return words
}

func BenchmarkGenerate(b *testing.B) {
	rand.Seed(6)
	w := generateWords(1000, 10)
	cols := 15
	rows := 17
	for n := 0; n < b.N; n++ {
		cr := New(cols, rows, w)
		cr.Generate(0, 5)
	}
}
