# sentigraph

The goal of this project is to use Machine Learning to graph the mood (positive/negative) over a long piece of text (e.g. a book). In other words, it will graph sentiment over time.

# Example

Using this tool, you can turn a book (in this case Under The Dome by Steven King) into a map of emotion like this:

![Under the Dome heatmap](examples/UnderTheDome.png)

The red bars indicate negative mood, the green indicate positive mood, and white indicates neutral.

# Usage

This project has three different components, so there are three main steps to using it.

## Train a classifier

The first step is to train a Machine Learning algorithm to determine the sentiment of a piece of text. You must download a training corpus for this (I recommend the one at [http://help.sentiment140.com/for-students/](http://help.sentiment140.com/for-students/)). You will need to pick a location to save the trained classifier (I'll use `/path/to/classifier`):

```
$ go run train/*.go bayes /path/to/classifier /path/to/training.csv
```

This will take several minutes to run, and once it's done you will have a classifier.

## Create a CSV for some text

The next step is to generate a CSV file with the sentiment of each sentence in the body of text you would like to graph. To do this, do the following:

```
$ go run plotcsv/*.go /path/to/classifier /path/to/text.txt /path/to/sentiments.csv
```

This will generate a file at `/path/to/sentiments.csv` containing sentiments for each sentence in the text file `/path/to/text.txt`.

## Graph the sentiments

Finally, to create a graphical image of the previously generated CSV file, you can do the following:

```
$ go run graph/*.go /path/to/sentiments.csv /path/to/graph.png heat
```

That will create a sentiment heat map out of the CSV file.
