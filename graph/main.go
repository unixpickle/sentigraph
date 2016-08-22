// Command graph generates a sentiment graph for a
// sentiment CSV file (which can be generated using
// the plotcsv command).
package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/unixpickle/sentigraph"
)

const (
	InputArg  = 1
	OutputArg = 2
	StyleArg  = 3
)

type DataPoint struct {
	Sentiment sentigraph.Sentiment
	Position  float64
}

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0],
			"input.csv output.png [style]")
		fmt.Fprintln(os.Stderr, "Available styles:")
		fmt.Fprintln(os.Stderr, " - line (default)")
		fmt.Fprintln(os.Stderr, " - heat")
		os.Exit(1)
	}

	style := "line"
	if len(os.Args) > StyleArg {
		style = os.Args[StyleArg]
	}

	data := readData()

	outFile, err := os.Create(os.Args[OutputArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create output:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	var img image.Image
	switch style {
	case "line":
		img = lineGraph(data)
	case "heat":
		img = heatGraph(data)
	default:
		fmt.Fprintln(os.Stderr, "Unknown style:", style)
		os.Exit(1)
	}

	if err := png.Encode(outFile, img); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode output:", err)
		os.Exit(1)
	}
}

func readData() []*DataPoint {
	inFile, err := os.Open(os.Args[InputArg])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open input:", err)
		os.Exit(1)
	}
	defer inFile.Close()

	r := csv.NewReader(inFile)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read input:", err)
		os.Exit(1)
	}

	output := make([]*DataPoint, len(records))
	for i, record := range records {
		if len(record) != 2 {
			fmt.Fprintln(os.Stderr, "Invalid number of columns:", len(record))
			os.Exit(1)
		}
		pos, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Row %d: invalid position: %s\n",
				i, record[0])
			os.Exit(1)
		}
		var sent sentigraph.Sentiment
		switch record[1] {
		case "0":
			sent = sentigraph.Neutral
		case "-1":
			sent = sentigraph.Negative
		case "1":
			sent = sentigraph.Positive
		default:
			fmt.Fprintf(os.Stderr, "Row %d: invalid sentiment: %s\n",
				i, record[1])
			os.Exit(1)
		}
		output[i] = &DataPoint{
			Sentiment: sent,
			Position:  pos,
		}
	}

	return output
}
