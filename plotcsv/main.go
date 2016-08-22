// Command plotcsv generates a CSV file containing the
// sentiments throughout the course of a body of text.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/unixpickle/sentigraph"
)

const (
	ModelArg  = 1
	TextArg   = 2
	OutputArg = 3
)

type SentenceInfo struct {
	Text     string
	Position float64
}

type DataPoint struct {
	Sentiment sentigraph.Sentiment
	Position  float64
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0],
			"model_file text_file ouput.csv")
		os.Exit(1)
	}

	sentences := readSentences()
	dataPoints := classifySentences(sentences)

	var points []*DataPoint
	for point := range dataPoints {
		points = append(points, point)
		fmt.Printf("\rGot mood %v at position %.04f      ",
			point.Sentiment, point.Position)
	}

	fmt.Println()

	outFile, err := os.Create(os.Args[OutputArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create output:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	writeCSV(outFile, points)
}

func readSentences() <-chan *SentenceInfo {
	text, err := ioutil.ReadFile(os.Args[TextArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read text file:", err)
		os.Exit(1)
	}

	res := make(chan *SentenceInfo)

	go func() {
		var curSentence []string
		fields := strings.Fields(string(text))
		for i, word := range fields {
			curSentence = append(curSentence, word)
			if sentenceEnded(word) {
				sentence := &SentenceInfo{
					Text:     strings.Join(curSentence, " "),
					Position: float64(i) / float64(len(fields)),
				}
				res <- sentence
				curSentence = nil
			}
		}
		close(res)
	}()

	return res
}

func classifySentences(sentences <-chan *SentenceInfo) <-chan *DataPoint {
	model, err := sentigraph.ReadModel(os.Args[ModelArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read model:", err)
		os.Exit(1)
	}
	resChan := make(chan *DataPoint)
	var wg sync.WaitGroup
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sentence := range sentences {
				sent := model.Classify(sentence.Text)
				resChan <- &DataPoint{
					Sentiment: sent,
					Position:  sentence.Position,
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()
	return resChan
}

func sentenceEnded(s string) bool {
	if s == "Dr." || s == "Mr." || s == "Mrs." || s == "Ms." {
		return false
	}
	return strings.HasSuffix(s, ".") || strings.HasSuffix(s, "?") ||
		strings.HasSuffix(s, "!")
}
