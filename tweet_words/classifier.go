package tweet_words


import (
	"github.com/jbrukh/bayesian"
	"strings"
)

var classifier *bayesian.Classifier

func init() {

	myClassifier, err := bayesian.NewClassifierFromFile("classifier.gob");
	if err != nil {
		panic(err)
	}
	classifier = myClassifier
}

func ClassifyTweet(tweetText string) (string) {
	tweetWordList :=   strings.FieldsFunc(tweetText, func(r rune) bool {
            switch r {
            case ' ':
                return true
            }
            return false
    })
    _, likely, _ := classifier.LogScores(tweetWordList)
    return string(classifier.Classes[likely])
}

