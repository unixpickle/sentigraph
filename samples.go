package sentigraph

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

// A Sentiment indicates the positivity/negativity of a
// statement or idea.
type Sentiment int

const (
	Positive Sentiment = iota
	Neutral
	Negative
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

// Normalize applies some basic rewrite rules to ensure
// that text from social media posts don't contain
// specific digital information and misspellings.
func (s *Sample) Normalize() string {
	words := strings.Fields(s.Contents)
	var newFields []string
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			continue
		}
		word = strings.ToLower(word)
		word = removeRepeatedLetters(word)
		newFields = append(newFields, word)
	}
	return strings.Join(newFields, " ")
}

// removeRepeatedLetters removes occurrences of letters so
// that no letter is repeated more than twice.
// This was suggested in
// http://cs.stanford.edu/people/alecmgo/papers/TwitterDistantSupervision09.pdf.
func removeRepeatedLetters(s string) string {
	var res string
	var last rune
	var count int
	for _, ch := range s {
		if ch == last {
			count++
		} else {
			last = ch
			count = 1
		}
		if count <= 2 {
			res += string(ch)
		}
	}
	return res
}
