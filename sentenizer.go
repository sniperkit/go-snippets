// So... I became aware of map keys could have multiple address pointers - put in another way: There can be multiple pointers to a map key, each with a different address. - So some of the variables in the Graph can be updates, but as of now I do not need to, or want to, because it works, and it tooks a while to get to this point.
package main
import (
	"fmt"
	"os"
	"github.com/mvryan/fasttag"
	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/data"
	"sort"
	"io/ioutil"
	"math"
	"strings"
)
type Sentence struct {
	Sentence string
	PosTags []string
}
type Noun struct {
	Noun string
	Occurence uint
}
type nounSorter struct {
	nouns []*Noun
	by func(n1, n2 *Noun) bool
}
type ByNoun func(n1, n2 *Noun) bool
func (by ByNoun) Sort(nouns []*Noun) {
	ns := &nounSorter{
		nouns: nouns,
		by: by,
	}
	sort.Sort(ns)
}
func (ns *nounSorter) Len() int {
	return len(ns.nouns)
}
func (ns *nounSorter) Swap(i, j int) {
	ns.nouns[i], ns.nouns[j] = ns.nouns[j], ns.nouns[i]
}
func (ns *nounSorter) Less(i, j int) bool {
	return ns.by(ns.nouns[i], ns.nouns[j])
}
type TriGram struct {
	Nouns [3]*Noun
	Occurence uint
}
type trigramSorter struct {
	trigrams []*TriGram
	by func(n1, n2 *TriGram) bool
}
type ByTriGram func(t1, t2 *TriGram) bool
func (by ByTriGram) Sort(trigrams []*TriGram) {
	ts := &trigramSorter{
		trigrams: trigrams,
		by: by,
	}
	sort.Sort(ts)
}
func (ts *trigramSorter) Len() int {
	return len(ts.trigrams)
}
func (ts *trigramSorter) Swap(i, j int) {
	ts.trigrams[i], ts.trigrams[j] = ts.trigrams[j], ts.trigrams[i]
}
func (ts *trigramSorter) Less(i, j int) bool {
	return ts.by(ts.trigrams[i], ts.trigrams[j])
}
type Pointers struct {
	Nouns []*Noun
	TriGrams []*TriGram
}
type Graph struct {
	Sentences []Sentence
	Nouns []Noun
	TriGrams []TriGram

	qaNouns map[string]uint
	qaNounToSentences map[string][]uint

	qaTriGrams map[[3]string]uint
	qaNounToTriGram map[string][]uint
	qaTriGramToSentences map[*TriGram][]uint

	NounMaxOccurence uint
	NounMinOccurence uint
	NounSumOccurence uint
	NounMeanOccurence float64
	NounStdDev float64
	NounZScore float64

	TriGramMaxOccurence uint
	TriGramMinOccurence uint
	TriGramSumOccurence uint
	TriGramMeanOccurence float64
	TriGramStdDev float64
	TriGramZScore float64

	Final Pointers
	Print []*Sentence
}
func (graph *Graph) LinearRegression() {
	graph.NounMinOccurence = 10000 // Default init value.
	// var sumWords uint
	for noun, _ := range graph.qaNounToSentences {
		if graph.NounMaxOccurence < uint(len(graph.qaNounToSentences[noun])) {
			graph.NounMaxOccurence = uint(len(graph.qaNounToSentences[noun]))
		}
		if graph.NounMinOccurence > uint(len(graph.qaNounToSentences[noun])) {
			graph.NounMinOccurence = uint(len(graph.qaNounToSentences[noun]))
		}
		graph.NounSumOccurence += uint(len(graph.qaNounToSentences[noun]))
	}
	graph.NounMeanOccurence = float64(float64(graph.NounSumOccurence) / float64(len(graph.Nouns)))
	graph.NounStdDev = math.Sqrt(float64(graph.NounSumOccurence)/graph.NounMeanOccurence)
	graph.NounZScore = graph.NounStdDev / float64(graph.NounMaxOccurence-graph.NounMinOccurence)
	graph.TriGramMinOccurence = 10000
	for trigramPtr, _ := range graph.qaTriGramToSentences {
		if len(graph.qaTriGramToSentences[trigramPtr]) == 0 {
			delete(graph.qaTriGramToSentences, trigramPtr)
			*trigramPtr = TriGram{}
			continue
		}
		if graph.TriGramMaxOccurence < uint(len(graph.qaTriGramToSentences[trigramPtr])) {
			graph.TriGramMaxOccurence = uint(len(graph.qaTriGramToSentences[trigramPtr]))
		}
		if graph.TriGramMinOccurence > uint(len(graph.qaTriGramToSentences[trigramPtr])) {
			graph.TriGramMinOccurence = uint(len(graph.qaTriGramToSentences[trigramPtr]))
		}
		graph.TriGramSumOccurence += uint(len(graph.qaTriGramToSentences[trigramPtr]))
	}
	graph.TriGramMeanOccurence = float64(float64(graph.TriGramSumOccurence) / float64(len(graph.qaTriGramToSentences)))
		graph.TriGramStdDev = math.Sqrt(float64(graph.TriGramSumOccurence) / graph.TriGramMeanOccurence)
	graph.TriGramZScore = graph.TriGramStdDev / float64(graph.TriGramMaxOccurence-graph.TriGramMinOccurence)
}
func main() {
	// var dict map[string]bool // TODO: read in words from the british and american dictionary and check whether or not a noun from the texts are part of it, if not continue to next noun. (or do a levenshtein spell-corrector, but that requires are loop for the or whole or part of the dictionary). - NOTE: If you want, the most common noun so far is a newline carriage, which is stripped away by the linux/gnu command "cat txts/* | tr '\n' ' ' > txts/"stripped newlines.txt""
	// If you want you can replace newlinens with following
	/*re = regexp.MustCompile(`\r?\n`)
     input = re.ReplaceAllString(input, " ") */

	var mygraph Graph
	var trigramBuffer [3]string
	var counter uint8 = 0
	var previousOffset uint = 0
	files, err := GetFilesFromDir("/home/naamik/go/src/legacy/pressure679/GhostWriter/txts/")
	// var buffer []*string
	if err != nil { panic(err) }
	dataAsset, err := data.Asset("data/english.json")
	if err != nil { panic(err) }
	trainingData, err := sentences.LoadTraining(dataAsset)
	if err != nil { panic(err) }
	sentenceTokenizer := sentences.NewSentenceTokenizer(trainingData)
	mygraph.qaTriGrams = make(map[[3]string]uint)
	mygraph.qaNouns = make(map[string]uint)
	mygraph.qaNounToSentences = make(map[string][]uint)
	mygraph.qaNounToTriGram = make(map[string][]uint)
	mygraph.qaTriGramToSentences = make(map[*TriGram][]uint)

	fmt.Println("* Metadata")
	var started bool
	for fileNum, _ := range files {
		osFile, err := os.Open("/home/naamik/go/src/legacy/pressure679/GhostWriter/txts/" + files[fileNum])
		if err != nil { panic(err) }
		buffer, err := ioutil.ReadAll(osFile)
		if err != nil { panic(err) }
		if len(mygraph.Sentences) == 0 { previousOffset = 0 } else {
			previousOffset = uint(len(mygraph.Sentences))
		}
		absSentences := sentenceTokenizer.Tokenize(string(buffer))
		for nssnum, _ := range absSentences {
			mygraph.Sentences = append(mygraph.Sentences, Sentence{Sentence: absSentences[nssnum].Text})
			words := fasttag.WordsToSlice(absSentences[nssnum].Text)
			posTags := fasttag.BrillTagger(words)
			mygraph.Sentences[int(previousOffset) + nssnum].PosTags = posTags
		}
		for sentenceNum, _ := range mygraph.Sentences[previousOffset + 1:] {
			words := fasttag.WordsToSlice(mygraph.Sentences[sentenceNum].Sentence)
			for wordNum, _ := range words {
				if string(mygraph.Sentences[sentenceNum].PosTags[wordNum][0]) != "N" { continue }
				lowercaseWord := strings.ToLower(words[wordNum])
				mygraph.qaNounToSentences[lowercaseWord] = append(mygraph.qaNounToSentences[lowercaseWord], previousOffset + uint(sentenceNum))
				// trigram queue appendage
				if !started {
					if counter == 1 {
						trigramBuffer[2] = lowercaseWord
					} else if counter == 2 {
						trigramBuffer[1] = trigramBuffer[2]
						trigramBuffer[2] = lowercaseWord
					} else if counter == 3 {
						trigramBuffer[0] = trigramBuffer[1]
						trigramBuffer[1] = trigramBuffer[2]
						trigramBuffer[2] = lowercaseWord
						counter = 0
						started = true
					}
				} else {
					trigramBuffer[0] = trigramBuffer[1]
					trigramBuffer[1] = trigramBuffer[2]
					trigramBuffer[2] = lowercaseWord
				}
				// init
				if len(mygraph.Nouns) == 0 {
					mygraph.qaNouns[lowercaseWord] = 0
					mygraph.Nouns = append(mygraph.Nouns, Noun{Noun:lowercaseWord})
					mygraph.Final.Nouns = append(mygraph.Final.Nouns, &mygraph.Nouns[0])
				}
				// append
				if _, exists := mygraph.qaNouns[lowercaseWord]; !exists {
					mygraph.Nouns = append(mygraph.Nouns, Noun{Noun:lowercaseWord})
					mygraph.Final.Nouns = append(mygraph.Final.Nouns, &mygraph.Nouns[len(mygraph.Nouns)-1])
					// mygraph.Nouns[len(mygraph.Nouns)-1].Occurence = 1
					mygraph.qaNouns[lowercaseWord] = uint(len(mygraph.Nouns) - 1)
				} else { // increment occurence
					mygraph.Nouns[mygraph.qaNouns[lowercaseWord]].Occurence++
				}
				// wait for trigrambuffer
				if len(mygraph.Nouns) < 2 { continue }
				// append
				if trigramBuffer[0] != "" {
					if _, exists := mygraph.qaTriGrams[trigramBuffer]; !exists {
						mygraph.TriGrams = append(mygraph.TriGrams, TriGram{Nouns: [3]*Noun{&mygraph.Nouns[mygraph.qaNouns[trigramBuffer[2]]], &mygraph.Nouns[mygraph.qaNouns[trigramBuffer[1]]], &mygraph.Nouns[mygraph.qaNouns[trigramBuffer[0]]]}})
						// mygraph.TriGrams[len(mygraph.TriGrams)-1].Occurence = 1
						mygraph.Final.TriGrams = append(mygraph.Final.TriGrams, &mygraph.TriGrams[len(mygraph.TriGrams)-1])
						mygraph.qaTriGrams[trigramBuffer] = uint(len(mygraph.TriGrams) - 1)
						mygraph.qaNounToTriGram[lowercaseWord] = append(mygraph.qaNounToTriGram[lowercaseWord], uint(len(mygraph.TriGrams) - 1))
						mygraph.qaTriGramToSentences[&mygraph.TriGrams[mygraph.qaTriGrams[trigramBuffer]]] = append(mygraph.qaTriGramToSentences[&mygraph.TriGrams[mygraph.qaTriGrams[trigramBuffer]]], previousOffset + uint(sentenceNum))
					} else { // increment
						mygraph.TriGrams[mygraph.qaTriGrams[trigramBuffer]].Occurence++
					}
				}
				counter++
			}
			trigramBuffer = [3]string{"", "", ""}
		}
		fmt.Println("  Read \"" + files[fileNum] + "\"")
	}
	fmt.Println("\nnumber of sentences:", len(mygraph.Sentences))
	fmt.Println("number of nouns:", len(mygraph.Nouns))
	fmt.Println("number of trigrams:", len(mygraph.TriGrams))

	mygraph.LinearRegression()
	nounByOccurence := func(n1, n2 *Noun) bool {
		return n1.Occurence < n2.Occurence
	}
	trigramByOccurence := func(t1, t2 *TriGram) bool {
		return t1.Nouns[0].Occurence + t1.Nouns[1].Occurence + t1.Nouns[2].Occurence < t2.Nouns[0].Occurence + t2.Nouns[1].Occurence + t2.Nouns[2].Occurence
	}
	ByNoun(nounByOccurence).Sort(mygraph.Final.Nouns)
	ByTriGram(trigramByOccurence).Sort(mygraph.Final.TriGrams)
	
	fmt.Println("\nSum of trigram:", len(mygraph.Final.TriGrams))
	fmt.Println("Sum of nouns:", len(mygraph.Final.Nouns))
	fmt.Println("\nTrigram max occurence:", mygraph.TriGramMaxOccurence)
	fmt.Println("Trigram min occurence:", mygraph.TriGramMinOccurence)
	fmt.Println("Trigram mean occurence:", mygraph.TriGramMeanOccurence)
	fmt.Println("Trigram std dev:", mygraph.TriGramStdDev)
	fmt.Println("\nNoun max occurence:", mygraph.NounMaxOccurence)
	fmt.Println("Noun min occurence:", mygraph.NounMinOccurence)
	fmt.Println("Noun mean occurence:", mygraph.NounMeanOccurence)
	fmt.Println("Noun std dev:", mygraph.NounStdDev)

	// fmt.Println(mygraph.qaNounToSentences["introduction"])
	fmt.Println("\n" + mygraph.Final.Nouns[0].Noun, mygraph.Final.Nouns[0].Occurence)
	fmt.Println(mygraph.Final.Nouns[len(mygraph.Final.Nouns)-1].Noun, mygraph.Final.Nouns[len(mygraph.Final.Nouns)-1].Occurence)
	fmt.Println(mygraph.Nouns[len(mygraph.Nouns)-1].Noun, mygraph.Nouns[len(mygraph.Nouns)-1].Occurence)
	fmt.Println(mygraph.Nouns[0].Noun, mygraph.Nouns[0].Occurence)
	var nounNum int = len(mygraph.Final.Nouns) - 1
	var uniqueSentences map[*string]bool = make(map[*string]bool)
	var uniqueTrigrams map[*TriGram]bool = make(map[*TriGram]bool)
	var nounStarted, sentencerStarted, trigrammerStarted bool = false, false, false

	final, err := os.Create("/home/naamik/go/src/legacy/pressure679/GhostWriter/sentenizer_final.org")
	if err != nil { panic(err) }
	for nounNum = len(mygraph.Final.Nouns) - 1; nounNum > - 1; nounNum-- {
		// fmt.Fprintln(final, "*", mygraph.Final.Nouns[nounNum].Noun)
		for trigramNum := len(mygraph.qaNounToTriGram[mygraph.Nouns[nounNum].Noun]); trigramNum > -1; trigramNum-- {
			if uniqueTrigrams[&mygraph.TriGrams[trigramNum]] { continue }
			if !nounStarted {
				fmt.Fprintln(final, "*", mygraph.Final.Nouns[nounNum].Noun)
				nounStarted = true
			}
			if !trigrammerStarted {
				fmt.Fprintln(final, "** Trigrams")
				trigrammerStarted = true
			}
			uniqueTrigrams[&mygraph.TriGrams[trigramNum]] = true
			fmt.Fprint(final, "   ")
			for trigramNounNum, _ := range mygraph.TriGrams[trigramNum].Nouns {
				fmt.Fprint(final, mygraph.TriGrams[trigramNum].Nouns[trigramNounNum].Noun)
				if trigramNounNum < 3 { fmt.Fprint(final, ", ") } else { fmt.Println() }
			}
			if trigramNum < len(mygraph.qaNounToTriGram[mygraph.Nouns[nounNum].Noun]) - 1 { fmt.Fprint(final, ", ") } else { fmt.Fprintln(final) }
		}
		for sentenceNum, _ := range mygraph.qaNounToSentences[mygraph.Final.Nouns[nounNum].Noun] {
			if uniqueSentences[&mygraph.Sentences[sentenceNum].Sentence] { continue }
			if !nounStarted {
				fmt.Fprintln(final, "*", mygraph.Final.Nouns[nounNum].Noun)
				nounStarted = true
			}
			if !sentencerStarted { 
				fmt.Fprintln(final, "** Sentences")
				sentencerStarted = true
			}
			uniqueSentences[&mygraph.Sentences[sentenceNum].Sentence] = true
			fmt.Fprintln(final, "*** " + mygraph.Sentences[sentenceNum].Sentence)
		}
		nounStarted = false
		trigrammerStarted = false
		sentencerStarted = false
	}
	err = final.Close()
	if err != nil { panic(err) }
	fmt.Println("Done.")
}
func GetFilesFromDir(directory string) (files []string, err error) {
osFileInfo, err := ioutil.ReadDir(directory); if err != nil { return nil, err }
for _, fileInfo := range osFileInfo { if !fileInfo.IsDir() { files = append(files, fileInfo.Name()) } }; return files, err
}
