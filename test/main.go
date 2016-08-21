// Command test tests a model on a testing corpus.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"github.com/unixpickle/sentigraph"
	"github.com/unixpickle/serializer"
)

const (
	ModelArg  = 1
	CorpusArg = 2
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "model_file corpus.csv")
		os.Exit(1)
	}
	model := readModel()
	sampleChan := readSamples()
	statusChan := make(chan bool)

	var wg sync.WaitGroup
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			runSamples(model, sampleChan, statusChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(statusChan)
	}()

	printStatuses(statusChan)
}

func runSamples(model sentigraph.Model, samples <-chan *sentigraph.Sample, statuses chan<- bool) {
	for sample := range samples {
		output := model.Classify(sample.Contents)
		statuses <- output == sample.Sentiment
	}
}

func printStatuses(statusChan <-chan bool) {
	var total int
	var correct int
	for status := range statusChan {
		total++
		if status {
			correct++
		}
		fmt.Printf("\rGot %d/%d (%.2f%%)     ", correct, total,
			float64(correct)/float64(total)*100)
	}
	fmt.Println("")
}

func readModel() sentigraph.Model {
	modelData, err := ioutil.ReadFile(os.Args[ModelArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read model:", err)
		os.Exit(1)
	}
	modelObj, err := serializer.DeserializeWithType(modelData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to deserialize model:", err)
		os.Exit(1)
	}
	model, ok := modelObj.(sentigraph.Model)
	if !ok {
		fmt.Fprintf(os.Stderr, "Invalid model type: %T\n", modelObj)
		os.Exit(1)
	}
	return model
}

func readSamples() <-chan *sentigraph.Sample {
	corpusFile, err := os.Open(os.Args[CorpusArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open corpus:", err)
		os.Exit(1)
	}
	defer corpusFile.Close()
	samples, err := sentigraph.ReadSamples(corpusFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse corpus:", err)
		os.Exit(1)
	}

	sampleChan := make(chan *sentigraph.Sample, len(samples))
	for _, sample := range samples {
		sampleChan <- sample
	}
	close(sampleChan)

	return sampleChan
}
