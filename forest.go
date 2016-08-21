package sentigraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/unixpickle/serializer"
	"github.com/unixpickle/weakai/idtrees"
)

// ForestSampleCountEnvVar is an environment variable
// which specifies the number of samples to use for
// generating each tree.
const ForestSampleCountEnvVar = "FOREST_SAMPLE_COUNT"

func init() {
	var f Forest
	var t treeSerializer
	serializer.RegisterTypedDeserializer(f.SerializerType(), DeserializeForest)
	serializer.RegisterTypedDeserializer(t.SerializerType(), deserializeTreeSerializer)
}

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

// DeserializeForest deserializes a forest.
func DeserializeForest(d []byte) (*Forest, error) {
	slice, err := serializer.DeserializeSlice(d)
	if err != nil {
		return nil, err
	}
	if len(slice) < 1 {
		return nil, errors.New("invalid Forest slice")
	}
	intVal, ok := slice[0].(serializer.Int)
	if !ok {
		return nil, errors.New("invalid Forest slice")
	}
	var res Forest
	res.Bigraph = intVal == 1
	for _, t := range slice[1:] {
		tree, ok := t.(*treeSerializer)
		if !ok {
			return nil, errors.New("invalid Forest slice")
		}
		res.Forest = append(res.Forest, tree.Tree())
	}
	return &res, nil
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
	log.Println("Creating samples...")
	samples := make([]idtrees.Sample, len(data))
	features := map[string]bool{}
	for i, d := range data {
		fs := newForestSample(f.Bigraph, d)
		samples[i] = fs
		for feature := range fs.features {
			features[feature] = true
		}
	}

	log.Println("Created", len(samples), "samples with", len(features), "features")

	attrs := make([]idtrees.Attr, 0, len(features))
	for feature := range features {
		attrs = append(attrs, feature)
	}

	subsampleCount := len(samples) / 2

	if countStr := os.Getenv(ForestSampleCountEnvVar); countStr != "" {
		var err error
		subsampleCount, err = strconv.Atoi(countStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid %s: %s", ForestSampleCountEnvVar, countStr)
		}
	}

	f.Forest = idtrees.BuildForest(ForestSize, samples, attrs, subsampleCount, 0,
		func(s []idtrees.Sample, attrs []idtrees.Attr) *idtrees.Tree {
			return idtrees.ID3(s, attrs, runtime.GOMAXPROCS(0))
		})
}

// SerializerType gives the unique ID used to serialize
// Forests with the serializer package.
func (f *Forest) SerializerType() string {
	return "github.com/unixpickle/sentigraph.Forest"
}

// Serialize serializes the random forest.
func (f *Forest) Serialize() ([]byte, error) {
	var bigraph serializer.Int
	if f.Bigraph {
		bigraph = 1
	}
	serializers := make([]serializer.Serializer, len(f.Forest)+1)
	serializers[0] = bigraph
	for i, t := range f.Forest {
		serializers[i+1] = newTreeSerializer(t)
	}
	return serializer.SerializeSlice(serializers)
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

type treeSerializer struct {
	Classification map[int]float64 `json:"c"`
	Attr           string          `json:"w"`
	TrueBranch     *treeSerializer `json:"t"`
	FalseBranch    *treeSerializer `json:"f"`
}

func deserializeTreeSerializer(d []byte) (*treeSerializer, error) {
	var t treeSerializer
	if err := json.Unmarshal(d, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func newTreeSerializer(t *idtrees.Tree) *treeSerializer {
	if t.Classification != nil {
		c := map[int]float64{}
		for sentiment, prob := range t.Classification {
			c[int(sentiment.(Sentiment))] = prob
		}
		return &treeSerializer{Classification: c}
	}
	return &treeSerializer{
		Attr:        t.Attr.(string),
		TrueBranch:  newTreeSerializer(t.ValSplit[true]),
		FalseBranch: newTreeSerializer(t.ValSplit[false]),
	}
}

func (t *treeSerializer) Tree() *idtrees.Tree {
	if t.Classification != nil {
		c := map[idtrees.Class]float64{}
		for class, val := range t.Classification {
			c[Sentiment(class)] = val
		}
		return &idtrees.Tree{Classification: c}
	}
	return &idtrees.Tree{
		Attr: t.Attr,
		ValSplit: map[idtrees.Val]*idtrees.Tree{
			true:  t.TrueBranch.Tree(),
			false: t.FalseBranch.Tree(),
		},
	}
}

func (t *treeSerializer) SerializerType() string {
	return "github.com/unixpickle/sentigraph.treeSerializer"
}

func (t *treeSerializer) Serialize() ([]byte, error) {
	return json.Marshal(t)
}
