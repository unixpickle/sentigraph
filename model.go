package sentigraph

import (
	"fmt"
	"io/ioutil"

	"github.com/unixpickle/serializer"
)

// A Model learns to classify the sentiment of text.
type Model interface {
	serializer.Serializer

	// Classify uses the current model to classify
	// the given text.
	// It is not valid to call this before the model
	// has been trained using the Train routine.
	Classify(text string) Sentiment

	// Train trains the model on the set of samples.
	// Depending on the model, this may be interactive
	// with the command-line user.
	Train(samples []*Sample)
}

// ReadModel reads a model from a file.
func ReadModel(path string) (Model, error) {
	modelData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	modelObj, err := serializer.DeserializeWithType(modelData)
	if err != nil {
		return nil, err
	}
	model, ok := modelObj.(Model)
	if !ok {
		return nil, fmt.Errorf("invalid model type: %T", modelObj)
	}
	return model, nil
}

// Models maps model names to functions which construct
// new instances of those models.
var Models = map[string]func() Model{
	"forest": func() Model {
		return &Forest{}
	},
	"forestBigraph": func() Model {
		return &Forest{Bigraph: true}
	},
	"bayes": func() Model {
		return &Bayes{}
	},
	"bayesBigraph": func() Model {
		return &Bayes{Bigraph: true}
	},
}
