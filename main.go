package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type ByLength struct {
	totalQuestions int64
	numWords       int64
}

type Measures struct {
	byLength map[int]ByLength
}

type Strategy struct {
	name string
	f    func(word string) int
	m    *Measures
}

func newMeasures() *Measures {
	return &Measures{
		byLength: make(map[int]ByLength),
	}
}

func (m *Measures) record(wordLength int, numQuestions int) {
	byLength := m.byLength[wordLength]
	byLength.numWords += 1
	byLength.totalQuestions += int64(numQuestions)
	m.byLength[wordLength] = byLength
}

func readWords(fileName string) []string {
	payload, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	words := strings.Split(string(payload), "\n")
	for i := range words {
		words[i] = strings.TrimSpace(words[i])
	}
	return words
}

func strategy0(word string) int {
	numQuestions := 0
	for _, c := range word {
		index := int(c - 'a')
		numQuestions += index + 1
	}
	numQuestions += 26
	return numQuestions
}

func strategy1(word string) int {
	numQuestions := 0
	for _, c := range word {
		index := int(c - 'a')
		numQuestions += index + 1
		numQuestions += 1
	}
	return numQuestions
}

func getStrategies() []Strategy {
	return []Strategy{
		{name: "strategy0", f: strategy0, m: newMeasures()},
		{name: "strategy1", f: strategy1, m: newMeasures()},
	}
}

func printHeaderCsv(strategies []Strategy) {
	cols := []string{"word_length"}
	for _, strategy := range strategies {
		cols = append(cols, strategy.name)
	}
	fmt.Printf("%s\n", strings.Join(cols, ","))
}

func printWordFrequencies(strategies []Strategy) {
	// for words of length 1 to 20 (inclusive)
	for i := 1; i <= 20; i++ {
		// print word length
		fmt.Printf("%d", i)
		for _, strat := range strategies {
			// ...and average number of questions for strategy
			if bl, ok := strat.m.byLength[i]; ok {
				fmt.Printf(",%.2f", float64(bl.totalQuestions)/float64(bl.numWords))
			} else {
				// except if we have no data, in which case, just print 0
				fmt.Printf(",0")
			}
		}
		fmt.Printf("\n")
	}
}

func printStrategiesCsv(strategies []Strategy) {
	printHeaderCsv(strategies)
	printWordFrequencies(strategies)
}

func updateStrategiesMeasures(words []string, strategies []Strategy) {
	for _, w := range words {
		for _, strategy := range strategies {
			numQuestions := strategy.f(w)
			strategy.m.record(len(w), numQuestions)
		}
	}
}

func measureLetterFrequency(words []string) {
	frequencies := make([]int64, 26)
	var total int64
	for _, w := range words {
		for _, c := range w {
			frequencies[int(c-'a')] += 1
			total += 1
		}
	}

	fmt.Printf("letter,frequency\n")
	for i := 0; i < 26; i++ {
		letter := 'a' + rune(i)
		frequency := float64(frequencies[i]) / float64(total)
		fmt.Printf("%v,%.5f\n", string([]rune{letter}), frequency)
	}
}

func main() {
	words := readWords("words.txt")
	// strategies := getStrategies()
	// updateStrategiesMeasures(words, strategies)
	// printStrategiesCsv(strategies)
	measureLetterFrequency(words)
}
