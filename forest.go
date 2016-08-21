package sentigraph

import (
	"runtime"
	"strings"

	"github.com/unixpickle/weakai/idtrees"
)

// ForestSize is the size of the random forests build by
// Forest.Train().
const ForestSize = 100

// A Forest classifies text documents using a random
// forest of decision trees.
type Forest struct {
	// Bigraph is true if bigraphs should be used in
	// addition to unigraphs.
	Bigraph bool

	// Forest is the learned model.
	// It is nil if no model has been trained.
	Forest idtrees.Forest
}

// Classify classifies the text using the forest.
// It is only valid to call this is f.Forest is
// non-nil (i.e. if the Forest has been trained).
func (f *Forest) Classify(text string) Sentiment {
	classes := f.Forest.Classify(newForestSampleText(f.Bigraph, text))

	var maxClass Sentiment
	var maxVal float64
	for class, prob := range classes {
		if prob >= maxVal {
			maxVal = prob
			maxClass = class.(Sentiment)
		}
	}
	return maxClass
}

// Train generates a forest for the training data.
func (f *Forest) Train(data []*Sample) {
	samples := make([]idtrees.Sample, len(data))
	features := map[string]bool{}
	for i, d := range data {
		fs := newForestSample(f.Bigraph, d)
		samples[i] = fs
		for feature := range fs.features {
			features[feature] = true
		}
	}

	attrs := make([]idtrees.Attr, 0, len(features))
	for feature := range features {
		attrs = append(attrs, feature)
	}

	f.Forest = idtrees.BuildForest(ForestSize, samples, attrs, len(samples)/2, 0,
		func(s []idtrees.Sample, attrs []idtrees.Attr) *idtrees.Tree {
			return idtrees.ID3(s, attrs, runtime.GOMAXPROCS(0))
		})
}

type forestSample struct {
	features map[string]bool
	class    Sentiment
}

func newForestSample(bigraphs bool, s *Sample) *forestSample {
	f := newForestSampleText(bigraphs, s.Contents)
	f.class = s.Sentiment
	return f
}

func newForestSampleText(bigraphs bool, t string) *forestSample {
	words := strings.Fields(Normalize(t))
	res := &forestSample{features: map[string]bool{}}
	for _, w := range words {
		res.features[w] = true
	}
	if bigraphs {
		for i := 1; i < len(words); i++ {
			res.features[words[i-1]+" "+words[i]] = true
		}
	}
	return res
}

func (f *forestSample) Attr(attr idtrees.Attr) idtrees.Val {
	return f.features[attr.(string)]
}

func (f *forestSample) Class() idtrees.Class {
	return f.class
}
