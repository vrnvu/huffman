package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/disiqueira/gotree"
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

type Node struct {
	Letters    []rune  `json:"letters"`
	Probabilty float64 `json:"probabilty"`
	Children   []*Node `json:"children"`
}

type ByAscendingProbability []*Node

func (a ByAscendingProbability) Len() int           { return len(a) }
func (a ByAscendingProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAscendingProbability) Less(i, j int) bool { return a[i].Probabilty < a[j].Probabilty }

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

func (n *Node) string() string {
	return fmt.Sprintf("%q (%.2f%%)", string(n.Letters), n.Probabilty*100)
}

func makeTree(tree gotree.Tree, node *Node) {
	for _, n := range node.Children {
		t := tree.Add(n.string())
		makeTree(t, n)
	}
}

var rootNode *Node

func init() {
	bs, err := ioutil.ReadFile("huffman.json")
	if err != nil {
		panic(err)
	}
	var node Node
	json.Unmarshal(bs, &node)
	rootNode = &node
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

func strategy2(word string) int {
	frequencyLetter := "eiaonsrtlcupdmhgybfvkwxzqj"
	numQuestions := 0
	for _, c := range word {
		for _, l := range frequencyLetter {
			numQuestions += 1
			if c == l {
				break
			}
		}
		numQuestions += 1
	}
	return numQuestions
}

func strategy3(word string) int {
	return (5 + 1) * len(word)
}

func strategy4(word string) int {
	// dummy tree
	numQuestions := 0
	for _, c := range word {
		if c == 'i' || c == 'e' {
			numQuestions += 2
		} else {
			numQuestions += 6
		}
		numQuestions += 1
	}
	return numQuestions
}

func strategy5(word string) int {
	numQuestions := 0
	for _, wordLetter := range word {
		node := rootNode
		for len(node.Letters) > 1 {
		searchChild:
			for _, child := range node.Children {
				for _, nodeLetter := range child.Letters {
					if nodeLetter == wordLetter {
						node = child
						break searchChild
					}
				}
			}
			// every time we travel to a new node, that's a question!
			numQuestions++
		}

		// ask if the word is over
		numQuestions++
	}
	return numQuestions
}

func getStrategies() []Strategy {
	return []Strategy{
		{name: "strategy0", f: strategy0, m: newMeasures()},
		{name: "strategy1", f: strategy1, m: newMeasures()},
		{name: "strategy2", f: strategy2, m: newMeasures()},
		{name: "strategy3", f: strategy3, m: newMeasures()},
		{name: "strategy4", f: strategy4, m: newMeasures()},
		{name: "strategy5", f: strategy5, m: newMeasures()},
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

func measureLetterFrequency(words []string) []float64 {
	frequenciesCount := make([]int64, 26)
	frequencies := make([]float64, 26)
	var total int64

	for _, w := range words {
		for _, c := range w {
			frequenciesCount[int(c-'a')] += 1
			total += 1
		}
	}

	for i := 0; i < 26; i++ {
		frequencies[i] = float64(frequenciesCount[i]) / float64(total)

	}

	return frequencies
}

func printLetterFrequency(frequencies []float64) {
	fmt.Printf("letter,frequency\n")
	for i := 0; i < 26; i++ {
		letter := 'a' + rune(i)
		fmt.Printf("%v,%.5f\n", string([]rune{letter}), frequencies[i])
	}
}

func makeLeafsFromFrequencies(frequencies []float64) []*Node {
	var leafs []*Node
	for i := 0; i < 26; i++ {
		letter := 'a' + rune(i)
		frequency := frequencies[i]
		leafs = append(leafs, &Node{
			Letters:    []rune{letter},
			Probabilty: frequency,
			Children:   nil,
		})
	}
	return leafs
}

func buildRecursivelySortedComposedNodes(leafs []*Node) []*Node {
	for len(leafs) > 1 {
		sort.Sort(ByAscendingProbability(leafs))
		l := leafs[0]
		r := leafs[1]
		newNode := &Node{
			Letters:    append(l.Letters, r.Letters...),
			Probabilty: l.Probabilty + r.Probabilty,
			Children:   []*Node{l, r},
		}
		leafs = append([]*Node{newNode}, leafs[2:]...)
	}
	return leafs
}

func buildFrequencyDecissionTree(leafs []*Node) gotree.Tree {
	tree := gotree.New(leafs[0].string())
	makeTree(tree, leafs[0])
	return tree
}

func serializeFrequencyDecissionTree(fileName string, leafs []*Node) {
	bs, _ := json.MarshalIndent(leafs[0], "", "")
	_ = ioutil.WriteFile(fileName, bs, 0644)
	fmt.Println("serialized to " + fileName)
}

func main() {
	words := readWords("words.txt")
	strategies := getStrategies()
	updateStrategiesMeasures(words, strategies)
	printStrategiesCsv(strategies)

	// frequencies := measureLetterFrequency(words)
	// printLetterFrequency(frequencies)

	// leafs := makeLeafsFromFrequencies(frequencies)
	// leafs = buildRecursivelySortedComposedNodes(leafs)

	// tree := buildFrequencyDecissionTree(leafs)
	// fmt.Println(tree.Print())

	// serializeFrequencyDecissionTree("huffman.json", leafs)
}
