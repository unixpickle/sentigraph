package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/unixpickle/sentigraph"
)

func writeCSV(w io.Writer, points []*DataPoint) {
	sort.Sort(PointSorter(points))

	writer := csv.NewWriter(w)
	for _, point := range points {
		var sentiment string
		switch point.Sentiment {
		case sentigraph.Negative:
			sentiment = "-1"
		case sentigraph.Positive:
			sentiment = "1"
		case sentigraph.Neutral:
			sentiment = "0"
		}
		record := []string{fmt.Sprintf("%.06f", point.Position), sentiment}
		if err := writer.Write(record); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write output:", err)
			os.Exit(1)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write output:", err)
		os.Exit(1)
	}
}
