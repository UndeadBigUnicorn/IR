package my_automaton

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestSparseAutomaton(t *testing.T) {
	words := []string{"banana", "bananas"}

	for n := 2; n < 3; n++ {

		for _, word := range words {
			a := NewSparseAutomaton(word, n)

			for _, query := range words {

				fmt.Printf("Testing query %s vs word %s, max distance %d\n\n", query, word, n)

				state := a.Start()
				for i, b := range query {

					state = a.Step(state, byte(b))
					canMatch, isMatch := a.CanMatch(state), a.IsMatch(state)

					fmt.Printf(" Query: %s, Match? %v, CanMatch? %v\n", query[:i+1], isMatch, canMatch)

					if isMatch && !canMatch {
						t.Errorf("IsMatch is true, canMatch must be true too")
					}
					if !canMatch {
						break
					}

				}
				fmt.Println("----")
			}
		}
	}

}


func TestTree(t *testing.T) {
	tree := NewTree()
	words := []string{"banana", "bananas", "bnaana", "world"}
	nonwords := []string{"sdfsdfsd", "hellos", "jeolls", "ello", "wrlds"}
	for _, word := range words {
		tree.Insert(word)
	}

	for _, word := range words {
		if !tree.Exists(word) {
			t.Error("Not found", word)
		}
	}

	for _, word := range nonwords {
		if tree.Exists(word) {
			t.Error("found", word)
		}
	}

	matchTest := map[string][]string{
		"banana": {"banana", "bananas", "bnaana"},
		"world":  {"world"},
		"wordl":  {"world"},
		"fordl":  {},
		"bnarna": {"banana", "bnaana"},
		"bananr": {"banana", "bananas"},
	}

	for k, expected := range matchTest {
		matches := tree.FuzzyMatches(k, 2)

		assert.Equal(t, len(matches), len(expected))
		for _, m := range matches {
			assert.Contains(t, expected, m)
		}

	}

}

var (
	tree      *Tree
	testwords []string
)

func TestExampleTree(t *testing.T) {

	tree := NewTree()
	tree.Insert("banana")
	tree.Insert("bananas")
	tree.Insert("cabana")
	tree.Insert("cabasa")

	fmt.Println(tree.FuzzyMatches("banana", 2))
	// XOutput:
	// [banana bananas cabana]
}

func TestMain(m *testing.M) {

	tree = NewTree()
	testwords = SampleEnglish()
	for _, word := range testwords {
		tree.Insert(word)
	}

	rc := m.Run()

	os.Exit(rc)

}

func BenchmarkTrie(b *testing.B) {

	for i := 0; i < b.N; i++ {
		tree.FuzzyMatches("holocaust", 2)
	}
}

func SampleEnglish() []string {
	file, err := os.Open("./big.txt")
	if err != nil {
		fmt.Println(err)
		return testwords
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	// Count the words.
	count := 0
	exp, _ := regexp.Compile("[a-zA-Z]+")

	for scanner.Scan() {
		words := exp.FindAll([]byte(scanner.Text()), -1)
		for _, word := range words {
			if len(word) > 1 {
				testwords = append(testwords, strings.ToLower(string(word)))
				count++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	fmt.Println("Read", len(testwords), "words")
	return testwords
}
