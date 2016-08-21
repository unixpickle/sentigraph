package sentigraph

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

// A Sentiment indicates the positivity/negativity of a
// statement or idea.
type Sentiment int

var AllSentiments = []Sentiment{Neutral, Negative, Positive}

const (
	Neutral Sentiment = iota
	Negative
	Positive
)

// Sample is a single textual training or testing sample.
type Sample struct {
	Contents  string
	Sentiment Sentiment
}

// ReadSamples reads samples from a CSV stream.
//
// The format of the CSV data is inferred automatically
// based on some popular corpora:
//
// * Corpus: http://help.sentiment140.com/for-students/
//   * Format: "0"/"2"/"4",ignored,ignored,ignored,ignored,tweet_body
//
// Other corpora may be supported in the future.
func ReadSamples(r io.Reader) ([]*Sample, error) {
	source := csv.NewReader(r)
	first, err := source.Read()
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if len(first) < 2 {
		return nil, errors.New("unknown data format")
	}
	if len(first) == 6 && (first[0] == "0" || first[0] == "2" || first[0] == "4") {
		return read024Samples(source, first)
	}
	return nil, errors.New("unknown data format")
}

func read024Samples(r *csv.Reader, first []string) ([]*Sample, error) {
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	records = append([][]string{first}, records...)
	samples := make([]*Sample, len(records))
	for i, record := range records {
		samples[i] = &Sample{Contents: record[len(record)-1]}
		switch record[0] {
		case "0":
			samples[i].Sentiment = Negative
		case "2":
			samples[i].Sentiment = Neutral
		case "4":
			samples[i].Sentiment = Positive
		default:
			return nil, fmt.Errorf("record %d: invalid sentiment %s",
				i, record[0])
		}
	}
	return samples, nil
}
