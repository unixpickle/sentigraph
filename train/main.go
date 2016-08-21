// Command train trains a model on a CSV file.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/unixpickle/sentigraph"
	"github.com/unixpickle/serializer"
)

const (
	ModelArg     = 1
	ModelPathArg = 2
	DataPathArg  = 3
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "model_name model_path data.csv")
		fmt.Fprintln(os.Stderr, "\nAvailable models:")
		for _, model := range modelNames() {
			fmt.Fprintln(os.Stderr, " -", model)
		}
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	var model sentigraph.Model

	modelData, err := ioutil.ReadFile(os.Args[ModelPathArg])
	if err == nil {
		modelObj, err := serializer.DeserializeWithType(modelData)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to deserialize existing model:", err)
			os.Exit(1)
		}
		var ok bool
		model, ok = modelObj.(sentigraph.Model)
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid model type: %T\n", modelObj)
			os.Exit(1)
		}
		log.Println("Loaded existing model from file.")
	} else {
		constructor, ok := sentigraph.Models[os.Args[ModelArg]]
		if !ok {
			fmt.Fprintln(os.Stderr, "Unknown model:", os.Args[1])
			os.Exit(1)
		}
		model = constructor()
	}

	dataFile, err := os.Open(os.Args[DataPathArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open data:", err)
		os.Exit(1)
	}
	defer dataFile.Close()
	samples, err := sentigraph.ReadSamples(dataFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse data:", err)
		os.Exit(1)
	}

	model.Train(samples)

	data, err := serializer.SerializeWithType(model)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to serialize model:", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(os.Args[ModelPathArg], data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write model file:", err)
		os.Exit(1)
	}
}

func modelNames() []string {
	var res []string
	for model := range sentigraph.Models {
		res = append(res, model)
	}
	sort.Strings(res)
	return res
}
