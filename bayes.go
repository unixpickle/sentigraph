package sentigraph

import (
	"encoding/json"
	"log"
	"math"
	"strings"

	"github.com/unixpickle/serializer"
)

// BayesSmoothing is the number of each feature to add to
// all sentiments in order to "smooth" zero probabilities.
// A value of 1 is specifically called Laplace smoothing.
const BayesSmoothing = 1

// BayesMinFeatureCount is the minimum number of times a
// feature must appear in order to be used.
const BayesMinFeatureCount = 2

func init() {
	var b Bayes
	serializer.RegisterTypedDeserializer(b.SerializerType(), DeserializeBayes)
}

type Bayes struct {
	// Bigraph should be set to true if bigraphs are to
	// be used in addition to unigraphs.
	Bigraph bool

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

	bestLogProb := math.Inf(-1)
	var bestSentiment Sentiment

	for _, sentiment := range AllSentiments {
		var logProb float64
		sentProb := b.Sentiments[sentiment]
		if sentProb == 0 {
			continue
		}
		for feature, condProb := range b.Conditional[sentiment] {
			var prob float64
			if features[feature] {
				prob = condProb / b.Features[feature]
			} else {
				prob = (1 - condProb) / (1 - b.Features[feature])
			}
			prob *= sentProb
			logProb += math.Log(prob)
		}
		if logProb > bestLogProb {
			bestLogProb = logProb
			bestSentiment = sentiment
		}
	}

	return bestSentiment
}

// Train regenerates the Bayes classifier using the
// given list of samples.
func (b *Bayes) Train(s []*Sample) {
	b.Sentiments = map[Sentiment]float64{}
	b.Features = map[string]float64{}
	b.Conditional = map[Sentiment]map[string]float64{}
	for _, sent := range AllSentiments {
		b.Conditional[sent] = map[string]float64{}
	}

	log.Println("Counting features...")
	for _, sample := range s {
		b.Sentiments[sample.Sentiment]++
		for feature := range b.features(sample.Contents) {
			if _, ok := b.Features[feature]; !ok {
				b.Features[feature] = BayesSmoothing
				for _, m := range b.Conditional {
					m[feature] = BayesSmoothing
				}
			}
			b.Features[feature]++
			b.Conditional[sample.Sentiment][feature]++
		}
	}

	log.Println("Pruning features...")
	for feature, count := range b.Features {
		if int(count-BayesSmoothing+0.5) < BayesMinFeatureCount {
			delete(b.Features, feature)
			for _, m := range b.Conditional {
				delete(m, feature)
			}
		}
	}

	log.Println("Normalizing", len(b.Features), "features...")
	for sent, count := range b.Sentiments {
		conditional := b.Conditional[sent]
		for feature := range conditional {
			conditional[feature] /= count
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

func (b *Bayes) features(text string) map[string]bool {
	fields := strings.Fields(SeparatePunctuation(Normalize(text)))
	res := map[string]bool{}
	for i, f := range fields {
		res[f] = true
		if i > 0 && b.Bigraph {
			res[fields[i-1]+" "+f] = true
		}
	}
	return res
}
