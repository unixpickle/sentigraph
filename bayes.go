package sentigraph

import (
	"encoding/json"
	"log"
	"math"
	"strings"

	"github.com/unixpickle/serializer"
)

func init() {
	var b Bayes
	serializer.RegisterTypedDeserializer(b.SerializerType(), DeserializeBayes)
}

type Bayes struct {
	// Bigraph should be set to true if bigraphs are to
	// be used in addition to unigraphs.
	Bigraph bool

	// MinimumProb is the minimum conditional probability
	// that a feature can be given (i.e. the probability
	// of a feature which does not occur).
	MinimumProb float64

	// Sentiments stores the unconditional probabilities of
	// each possible sentiment.
	Sentiments map[Sentiment]float64

	// Conditional stores the probabilities of each text
	// feature conditioned on a Sentiment.
	Conditional map[Sentiment]map[string]float64

	// Features stores the unconditional probability of
	// each feature.
	Features map[string]float64
}

// DeserializeBayes deserializes a Bayes model.
func DeserializeBayes(d []byte) (*Bayes, error) {
	var res Bayes
	if err := json.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Classify returns the most likely classification for
// the given piece of text.
func (b *Bayes) Classify(text string) Sentiment {
	features := b.features(text)

	var bestLogProb float64
	var bestSentiment Sentiment

	for i, sentiment := range AllSentiments {
		var logProb float64
		sentProb := b.Sentiments[sentiment]
		if sentProb == 0 {
			continue
		}
		conditional := b.Conditional[sentiment]
		for _, feature := range features {
			if b.Features[feature] == 0 {
				continue
			}
			prob := conditional[feature] * sentProb / b.Features[feature]
			if prob == 0 {
				prob = b.MinimumProb
			}
			logProb += math.Log(prob)
		}
		if logProb > bestLogProb || i == 0 {
			bestLogProb = logProb
			bestSentiment = sentiment
		}
	}

	return bestSentiment
}

// Train regenerates the Bayes classifier using the
// given list of samples.
func (b *Bayes) Train(s []*Sample) {
	b.MinimumProb = 1 / float64(len(s))
	b.Sentiments = map[Sentiment]float64{}
	b.Features = map[string]float64{}
	b.Conditional = map[Sentiment]map[string]float64{}
	for _, sent := range AllSentiments {
		b.Conditional[sent] = map[string]float64{}
	}

	log.Println("Counting features...")
	for _, sample := range s {
		b.Sentiments[sample.Sentiment]++
		for _, feature := range b.features(sample.Contents) {
			b.Features[feature]++
			b.Conditional[sample.Sentiment][feature]++
		}
	}

	log.Println("Normalizing features...")
	for _, sent := range AllSentiments {
		conditional := b.Conditional[sent]
		for feature, count := range b.Features {
			conditional[feature] /= float64(count)
		}
	}
	for sent, count := range b.Sentiments {
		b.Sentiments[sent] = count / float64(len(s))
	}
	for feature, count := range b.Features {
		b.Features[feature] = count / float64(len(s))
	}
}

// SerializerType gives the unique ID used to serialize
// Bayes instances with the serializer package.
func (b *Bayes) SerializerType() string {
	return "github.com/unixpickle/sentigraph.Bayes"
}

// Serialize serializes the bayes classifier.
func (b *Bayes) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Bayes) features(text string) []string {
	fields := strings.Fields(Normalize(text))
	featureSet := map[string]bool{}
	for i, f := range fields {
		featureSet[f] = true
		if i > 0 && b.Bigraph {
			featureSet[fields[i-1]+" "+f] = true
		}
	}
	slice := make([]string, 0, len(featureSet))
	for f := range featureSet {
		slice = append(slice, f)
	}
	return slice
}
