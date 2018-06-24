package GhostWriter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"compress/bzip2"
	"compress/gzip"
	"encoding/gob"
	"encoding/xml"
	"github.com/mvryan/fasttag"
	"math"
	"flag"
)

type SEBadge struct {
	AttrClass    string `xml:" Class,attr"  json:",omitempty"`
	AttrDate     string `xml:" Date,attr"  json:",omitempty"`
	AttrId       string `xml:" Id,attr"  json:",omitempty"`
	AttrName     string `xml:" Name,attr"  json:",omitempty"`
	AttrTagBased string `xml:" TagBased,attr"  json:",omitempty"`
	AttrUserId   string `xml:" UserId,attr"  json:",omitempty"`
}
type SEBadges struct {
	SEBadgeRow []*SEBadge `xml:" row,omitempty" json:"row,omitempty"`
}

type SEComment struct {
	AttrCreationDate    string `xml:" CreationDate,attr"  json:",omitempty"`
	AttrId              string `xml:" Id,attr"  json:",omitempty"`
	AttrPostId          string `xml:" PostId,attr"  json:",omitempty"`
	AttrScore           string `xml:" Score,attr"  json:",omitempty"`
	AttrText            string `xml:" Text,attr"  json:",omitempty"`
	AttrUserDisplayName string `xml:" UserDisplayName,attr"  json:",omitempty"`
	AttrUserId          string `xml:" UserId,attr"  json:",omitempty"`
}
type SEComments struct {
	ChiRow []*SEComment `xml:" row,omitempty" json:"row,omitempty"`
}

type SEPostHistory struct {
	AttrComment           string `xml:" Comment,attr"  json:",omitempty"`
	AttrCreationDate      string `xml:" CreationDate,attr"  json:",omitempty"`
	AttrId                string `xml:" Id,attr"  json:",omitempty"`
	AttrPostHistoryTypeId string `xml:" PostHistoryTypeId,attr"  json:",omitempty"`
	AttrPostId            string `xml:" PostId,attr"  json:",omitempty"`
	AttrRevisionGUID      string `xml:" RevisionGUID,attr"  json:",omitempty"`
	AttrText              string `xml:" Text,attr"  json:",omitempty"`
	AttrUserDisplayName   string `xml:" UserDisplayName,attr"  json:",omitempty"`
	AttrUserId            string `xml:" UserId,attr"  json:",omitempty"`
}
type SEPostHistories struct {
	ChiRow []*SEPostHistory `xml:" row,omitempty" json:"row,omitempty"`
}

type SEPostLink struct {
	AttrCreationDate  string `xml:" CreationDate,attr"  json:",omitempty"`
	AttrId            string `xml:" Id,attr"  json:",omitempty"`
	AttrLinkTypeId    string `xml:" LinkTypeId,attr"  json:",omitempty"`
	AttrPostId        string `xml:" PostId,attr"  json:",omitempty"`
	AttrRelatedPostId string `xml:" RelatedPostId,attr"  json:",omitempty"`
}
type SEPostLinks struct {
	ChiRow []*SEPostLink `xml:" row,omitempty" json:"row,omitempty"`
}

type SEPost struct {
	AttrAcceptedAnswerId      string `xml:" AcceptedAnswerId,attr"  json:",omitempty"`
	AttrAnswerCount           string `xml:" AnswerCount,attr"  json:",omitempty"`
	AttrBody                  string `xml:" Body,attr"  json:",omitempty"`
	AttrClosedDate            string `xml:" ClosedDate,attr"  json:",omitempty"`
	AttrCommentCount          string `xml:" CommentCount,attr"  json:",omitempty"`
	AttrCreationDate          string `xml:" CreationDate,attr"  json:",omitempty"`
	AttrFavoriteCount         string `xml:" FavoriteCount,attr"  json:",omitempty"`
	AttrId                    string `xml:" Id,attr"  json:",omitempty"`
	AttrLastActivityDate      string `xml:" LastActivityDate,attr"  json:",omitempty"`
	AttrLastEditDate          string `xml:" LastEditDate,attr"  json:",omitempty"`
	AttrLastEditorDisplayName string `xml:" LastEditorDisplayName,attr"  json:",omitempty"`
	AttrLastEditorUserId      string `xml:" LastEditorUserId,attr"  json:",omitempty"`
	AttrOwnerDisplayName      string `xml:" OwnerDisplayName,attr"  json:",omitempty"`
	AttrOwnerUserId           string `xml:" OwnerUserId,attr"  json:",omitempty"`
	AttrParentId              string `xml:" ParentId,attr"  json:",omitempty"`
	AttrPostTypeId            string `xml:" PostTypeId,attr"  json:",omitempty"`
	AttrScore                 string `xml:" Score,attr"  json:",omitempty"`
	AttrTags                  string `xml:" Tags,attr"  json:",omitempty"`
	AttrTitle                 string `xml:" Title,attr"  json:",omitempty"`
	AttrViewCount             string `xml:" ViewCount,attr"  json:",omitempty"`
}
type SEPosts struct {
	ChiRow []*SEPost `xml:" row,omitempty" json:"row,omitempty"`
}

type SETag struct {
	AttrCount         string `xml:" Count,attr"  json:",omitempty"`
	AttrExcerptPostId string `xml:" ExcerptPostId,attr"  json:",omitempty"`
	AttrId            string `xml:" Id,attr"  json:",omitempty"`
	AttrTagName       string `xml:" TagName,attr"  json:",omitempty"`
	AttrWikiPostId    string `xml:" WikiPostId,attr"  json:",omitempty"`
}
type SETags struct {
	ChiRow []*SETag `xml:" row,omitempty" json:"row,omitempty"`
}

type SEUser struct {
	AttrAboutMe         string `xml:" AboutMe,attr"  json:",omitempty"`
	AttrAccountId       string `xml:" AccountId,attr"  json:",omitempty"`
	AttrAge             string `xml:" Age,attr"  json:",omitempty"`
	AttrCreationDate    string `xml:" CreationDate,attr"  json:",omitempty"`
	AttrDisplayName     string `xml:" DisplayName,attr"  json:",omitempty"`
	AttrDownVotes       string `xml:" DownVotes,attr"  json:",omitempty"`
	AttrId              string `xml:" Id,attr"  json:",omitempty"`
	AttrLastAccessDate  string `xml:" LastAccessDate,attr"  json:",omitempty"`
	AttrLocation        string `xml:" Location,attr"  json:",omitempty"`
	AttrProfileImageUrl string `xml:" ProfileImageUrl,attr"  json:",omitempty"`
	AttrReputation      string `xml:" Reputation,attr"  json:",omitempty"`
	AttrUpVotes         string `xml:" UpVotes,attr"  json:",omitempty"`
	AttrViews           string `xml:" Views,attr"  json:",omitempty"`
	AttrWebsiteUrl      string `xml:" WebsiteUrl,attr"  json:",omitempty"`
}
type SEUsers struct {
	ChiRow []*SEUser `xml:" row,omitempty" json:"row,omitempty"`
}

// Graciously open sourced by github.com/dustin/go-wikiparse under an MIT styled license (10TH of July 2017).
type MWSiteInfo struct {
	SiteName   string `xml:"sitename"`
	Base       string `xml:"base"`
	Generator  string `xml:"generator"`
Case       string `xml:"case"`
	Namespaces []struct {
		Key   string `xml:"key,attr"`
	Case  string `xml:"case,attr"`
		Value string `xml:",chardata"`
	} `xml:"namespaces>namespace"`
}

// A Contributor is a user who contributed a revision.
type MWContributor struct {
	ID       uint64 `xml:"id"`
	Username string `xml:"username"`
}

// A Redirect to another Page.
type MWRedirect struct {
	Title string `xml:"title,attr"`
}

// A Revision to a page.
type MWRevision struct {
	ID          uint64        `xml:"id"`
	Timestamp   string        `xml:"timestamp"`
	Contributor MWContributor `xml:"contributor"`
	Comment     string        `xml:"comment"`
	Text        string        `xml:"text"`
}

// A Page in the wiki.
type MWPage struct {
	Title     string       `xml:"title"`
	ID        uint64       `xml:"id"`
	Redir     MWRedirect     `xml:"redirect"`
	Revisions []MWRevision `xml:"revision"`
	Ns        uint64       `xml:"ns"`
}

type MWAbstractRoot struct {
	MWAbstractFeed *MWAbstractFeed `xml:" feed,omitempty" json:"feed,omitempty"`
}

type MWAbstractAbstract struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type MWAbstractAnchor struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type MWAbstractDoc struct {
	MWAbstractAbstract *MWAbstractAbstract `xml:" abstract,omitempty" json:"abstract,omitempty"`
	MWAbstractLinks    *MWAbstractLinks    `xml:" links,omitempty" json:"links,omitempty"`
	MWAbstractTitle    *MWAbstractTitle    `xml:" title,omitempty" json:"title,omitempty"`
	MWAbstractUrl      *MWAbstractUrl      `xml:" url,omitempty" json:"url,omitempty"`
}

type MWAbstractFeed struct {
	MWAbstractDoc []*MWAbstractDoc `xml:" doc,omitempty" json:"doc,omitempty"`
}

type MWAbstractLink struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type MWAbstractLinks struct {
	MWAbstractSublink []*MWAbstractSublink `xml:" sublink,omitempty" json:"sublink,omitempty"`
}

type MWAbstractSublink struct {
	AttrLinktype     string            `xml:" linktype,attr"  json:",omitempty"`
	MWAbstractAnchor *MWAbstractAnchor `xml:" anchor,omitempty" json:"anchor,omitempty"`
	MWAbstractLink   *MWAbstractLink   `xml:" link,omitempty" json:"link,omitempty"`
}

type MWAbstractTitle struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type MWAbstractUrl struct {
	Text string `xml:",chardata" json:",omitempty"`
}

// stop word list from SMART (Salton,1971).  Available at ftp://ftp.cs.cornell.edu/pub/smart/english.stop
var StopWordsSlice = []string{
	"a", "a's", "able", "about", "above", "according", "accordingly", "across", "actually", "after", "afterwards", "again", "against", "ain't", "all", "allow", "allows", "almost", "alone", "along", "already", "also", "although", "always", "am", "among", "amongst", "an", "and", "another", "any", "anybody", "anyhow", "anyone", "anything", "anyway", "anyways", "anywhere", "apart", "appear", "appreciate", "appropriate", "are", "aren't", "around", "as", "aside", "ask", "asking", "associated", "at", "available", "away", "awfully", "b", "be", "became", "because", "become", "becomes", "becoming", "been", "before", "beforehand", "behind", "being", "believe", "below", "beside", "besides", "best", "better", "between", "beyond", "both", "brief", "but", "by", "c", "c'mon", "c's", "came", "can", "can't", "cannot", "cant", "cause", "causes", "certain", "certainly", "changes", "clearly", "co", "com", "come", "comes", "concerning", "consequently", "consider", "considering", "contain", "containing", "contains", "corresponding", "could", "couldn't", "course", "currently", "d", "definitely", "described", "despite", "did", "didn't", "different", "do", "does", "doesn't", "doing", "don't", "done", "down", "downwards", "during", "e", "each", "edu", "eg", "eight", "either", "else", "elsewhere", "enough", "entirely", "especially", "et", "etc", "even", "ever", "every", "everybody", "everyone", "everything", "everywhere", "ex", "exactly", "example", "except", "f", "far", "few", "fifth", "first", "five", "followed", "following", "follows", "for", "former", "formerly", "forth", "four", "from", "further", "furthermore", "g", "get", "gets", "getting", "given", "gives", "go", "goes", "going", "gone", "got", "gotten", "greetings", "h", "had", "hadn't", "happens", "hardly", "has", "hasn't", "have", "haven't", "having", "he", "he's", "hello", "help", "hence", "her", "here", "here's", "hereafter", "hereby", "herein", "hereupon", "hers", "herself", "hi", "him", "himself", "his", "hither", "hopefully", "how", "howbeit", "however", "i", "i'd", "i'll", "i'm", "i've", "ie", "if", "ignored", "immediate", "in", "inasmuch", "inc", "indeed", "indicate", "indicated", "indicates", "inner", "insofar", "instead", "into", "inward", "is", "isn't", "it", "it'd", "it'll", "it's", "its", "itself", "j", "just", "k", "keep", "keeps", "kept", "know", "knows", "known", "l", "last", "lately", "later", "latter", "latterly", "least", "less", "lest", "let", "let's", "like", "liked", "likely", "little", "look", "looking", "looks", "ltd", "m", "mainly", "many", "may", "maybe", "me", "mean", "meanwhile", "merely", "might", "more", "moreover", "most", "mostly", "much", "must", "my", "myself", "n", "name", "namely", "nd", "near", "nearly", "necessary", "need", "needs", "neither", "never", "nevertheless", "new", "next", "nine", "no", "nobody", "non", "none", "noone", "nor", "normally", "not", "nothing", "novel", "now", "nowhere", "o", "obviously", "of", "off", "often", "oh", "ok", "okay", "old", "on", "once", "one", "ones", "only", "onto", "or", "other", "others", "otherwise", "ought", "our", "ours", "ourselves", "out", "outside", "over", "overall", "own", "p", "particular", "particularly", "per", "perhaps", "placed", "please", "plus", "possible", "presumably", "probably", "provides", "q", "que", "quite", "qv", "r", "rather", "rd", "re", "really", "reasonably", "regarding", "regardless", "regards", "relatively", "respectively", "right", "s", "said", "same", "saw", "say", "saying", "says", "second", "secondly", "see", "seeing", "seem", "seemed", "seeming", "seems", "seen", "self", "selves", "sensible", "sent", "serious", "seriously", "seven", "several", "shall", "she", "should", "shouldn't", "since", "six", "so", "some", "somebody", "somehow", "someone", "something", "sometime", "sometimes", "somewhat", "somewhere", "soon", "sorry", "specified", "specify", "specifying", "still", "sub", "such", "sup", "sure", "t", "t's", "take", "taken", "tell", "tends", "th", "than", "thank", "thanks", "thanx", "that", "that's", "thats", "the", "their", "theirs", "them", "themselves", "then", "thence", "there", "there's", "thereafter", "thereby", "therefore", "therein", "theres", "thereupon", "these", "they", "they'd", "they'll", "they're", "they've", "think", "third", "this", "thorough", "thoroughly", "those", "though", "three", "through", "throughout", "thru", "thus", "to", "together", "too", "took", "toward", "towards", "tried", "tries", "truly", "try", "trying", "twice", "two", "u", "un", "under", "unfortunately", "unless", "unlikely", "until", "unto", "up", "upon", "us", "use", "used", "useful", "uses", "using", "usually", "uucp", "v", "value", "various", "very", "via", "viz", "vs", "w", "want", "wants", "was", "wasn't", "way", "we", "we'd", "we'll", "we're", "we've", "welcome", "well", "went", "were", "weren't", "what", "what's", "whatever", "when", "whence", "whenever", "where", "where's", "whereafter", "whereas", "whereby", "wherein", "whereupon", "wherever", "whether", "which", "while", "whither", "who", "who's", "whoever", "whole", "whom", "whose", "why", "will", "willing", "wish", "with", "within", "without", "won't", "wonder", "would", "would", "wouldn't", "x", "y", "yes", "yet", "you", "you'd", "you'll", "you're", "you've", "your", "yours", "yourself", "yourselves", "z", "zero",
}

// TODO: Fix error, "Hash tables typically grow automatically when the number of elements increases above a certain threshold. When it does, it reallocates a new underlying array and copies all of the old elements into the new one." - dsnet, https://github.com/golang/go/issues/11865, 4/18-18 (original answer to issue from 25/7-15).
type Links []*string
func (links Links) appendLinkIfNotExists(link *string) (bool, uint) { for num, _ := links { if strings.EqualFold(*links[num], &link) { return true, uint(num) } } links = append(links, link); return false, uint(0) }
type Articles []string
func (articles Articles) appendArticleIfNotExists(article string) (bool, uint) { for num, _ := range articles { if strings.EqualFold(articles[num], article) { return true, uint(num) } } articles = append(articles, article); return false, uint(0) }
type Words []string
func (words Words) appendWordIfNotExists(word string) (bool, uint) { for num, _ := range words { if strings.EqualFold(words[num], word) { return true, uint(num) } } words = append(words, word); return false, uint(0) }
// This is just a search which starts 4 goroutines each searching it's quarter of the Words list (I have no idea whether or not if a golang compiler optimizes such).
func (words Words) search(word string) (bool, uint) {
	length := len(words)
	// Possible read/write queue here for each goroutine.
	go func() {
		for i := 0; i < length / 4; length++ {
			if strings.EqualFold(words[i], word) { return false, uint(i) }
		}
	}()
	go func() {
		for i := length / 4; i < length / 2; length++ {
			if strings.EqualFold(words[i], word) { return false, uint(i) }
		}
	}()
	go func() {
		for i := length / 2; i < length / 2 + length / 4; length++ {
			if strings.EqualFold(words[i], word) { return false, uint(i) }
		}
	}()
	go func() {
		for i := length / 2 + length / 4; i < length / 4 * 3; length++ {
			if strings.EqualFold(words[i], word) { return false, uint(i) }
		}
	}()
	return true, uint(0)
}

type HashTable struct { Words Words; Links Links; Articles Articles }

type Amount []byte
func (amount Amount) increment() { amount = []byte(strconv.Itoa(strconv.ParseInt(string(amount), 10, 0)+1)) }

type WordMetaData struct { ZScore ZScore; Extremum Extremum }
type Sentence struct { Start, End []byte }
type Sentences []Sentence
type WordData struct { Sentences Sentences; MetaData WordMetaData }

type Article struct { Sections map[string]string; References []string; Links map[*string]bool; Nouns map[*string]WordData; WordMaxOccurence map[*string]Amount; WordMinOccurence map[*string]Amount; IndexOffset int64; ZScore ZScore; MinLinkedOccurence, MaxLinkedOccurence Amount }

// ****, I just realized int's are optimal for binary methods in a large system (for returning a left or right bit shift in a super object/struct/hashtable). - 4/19/2018 - nevertheless, uints can hold larger numeric values, and given hardware of optimized quality multiple memory retrievals can be easy, although that is for large motherboard systems with many n amount of transistors.
type ZScore []byte
func (lrValue ZScore) declare(StdDev float64, Max, Min uint) { lrValue = ZScore(strconv.FormatFloat(StdDev / float64(Max - Min), 'f', -1, 64)) }
type Extremum []byte
func (lrValue Extremum) declare(Num uint, StdDev float64, Max, Min uint) { lrValue = Extremum(strconv.FormatFloat(float64(Num) * StdDev / float64(Max - Min), 'f', -1, 64)) }

type DataGraphArticles map[*string]Article
func (dataGraphArticles DataGraphArticles) New(article *string) { dataGraphArticles[article] = make(Article) }
type DataGraphWords map[*string][]*string
func (words DataGraphWords) appendIfNotExists(article, word *string) (bool, uint) {
	for num, _ := range words[article] {
		if strings.EqualFolds(*words[article][num], *word) { return true, uint(num) }
	}
	return false, uint(0)
}
type DataGraph struct { Articles DataGraphArticles; Words DataGraphWords /* map key is an article, items from the slices are pointers to words in the HashTable */; AlphabeticIndex [26]uint; SumLinks uint; HashTable HashTable; Final Final }

type Final struct { Articles []Articles; Words []Words }

// The amount of articles and words sorted will be determined by the amount of memory available. As of now an estimation of amount is articles with a depth of 3 link-levels from root-/base-article. (~5GB RAM depending on the amount of links and reoccuring links/words - calculated with a chosen median of 100 links per article)

// TODO: Add keywords  for Bloom's taxonomy
/* type BloomsTaxonomy struct {
	L1 []string{}
	L2 []string{}
	L3 []string{}
	L4 []string{}
	L5 []string{}
	L6 []string{}
	L7 []string{}
} */
const Levels [5]string = [5]string{"remember", "understand", "apply", "analyze", "evaluate", "create"}

func GetFilesFromArticlesDir(directory *string) (files []string, err error) {
	osFileInfo, err := ioutil.ReadDir(*directory); if err != nil { return nil, err }; for _, fileInfo := range osFileInfo { if !fileInfo.IsDir() { files = append(files, fileInfo.Name()) } }; return
}
func DecompressBZip(directory, file *string) (ioReader io.Reader, err error) {
	osFile, err := os.Open(*directory + "/" + *file); if err != nil { return nil, err }; if err != nil { return nil, err }; ioReader = bzip2.NewReader(osFile); return ioReader, nil
}

// TODO (4/19-15): Update to the new data structure. - Edit (4/19/2018): and add the added methods for each struct into the below methods' loops.

// TODO: To ease on memory consumption do the WriteIndexAnContentData method after a certain amount of memory. - The size of the file can be retrieved within the WriteIndexAnContentData method.
func (graph Graph) ReadWikiXML(directory *string) (err error) {
	files, err := GetFilesFromArticlesDir(directory); if err != nil { return err }; referenceRE := regexp.MustCompile("<ref>(.+)</ref>"); linkRE := regexp.MustCompile("[[(.+)]]"); sectionRE := regexp.MustCompile("[=]{2,5}(.+)[=]{2,5}"); var sentenceIndex, previousSentenceIndex, curNounIndex int = -1, -1, -1; var mwPage MWPage; var sections  map[string]string; var sectionTitle string; var references []string; sentenceStop := ";:,.!?";
	for _, file := range files {
		ioReader, err := DecompressBZip(directory, &file); if err != nil { return err }; xmlDecoder := xml.NewDecoder(ioReader); err = xmlDecoder.Decode(&mwPage); if err != nil { if strings.EqualFold(err.Error(), io.EOF.Error()) { break } else { return err } }
		if mwPage.Revisions[0].Text != "" {
			var article Article
			graph.Articles[mwPage.Title] = Article{}
			graph.HashTable.Articles = make(Articles, 0)
			mwPage.Revisions[0].Text = strings.Replace(mwPage.Revisions[0].Text, "&lt;", "<", -1); mwPage.Revisions[0].Text = strings.Replace(mwPage.Revisions[0].Text, "&gt;", ">", -1); mwPage.Revisions[0].Text = strings.Replace(mwPage.Revisions[0].Text, "&quot;", "\"", -1); mwPage.Revisions[0].Text = strings.Replace(mwPage.Revisions[0].Text, "&amp;", "&", -1);
			sectionIndex := sectionRE.FindAllStringIndex(mwPage.Revisions[0].Text, -1); for sectionNum, _ := range sectionIndex { sectionTitle = strings.Trim(mwPage.Revisions[0].Text[sectionIndex[sectionNum][0]:sectionIndex[sectionNum][1]], "="); if sectionNum == 0 { sections["Abstract"] = mwPage.Revisions[0].Text[:sectionIndex[sectionNum][1]-1] } else if sectionNum < len(sectionIndex) - 1 { sections[sectionTitle] = mwPage.Revisions[0].Text[sectionIndex[sectionNum][1]:sectionIndex[sectionNum+1][0]] } else { sections[sectionTitle] = mwPage.Revisions[0].Text[sectionIndex[sectionNum][1]:len(mwPage.Revisions[0].Text)] }
				refIndex := referenceRE.FindAllStringIndex(mwPage.Revisions[0].Text, -1); for refNum, _ := range refIndex { references = append(references, mwPage.Revisions[0].Text[refIndex[refNum][0]:refIndex[refNum][1]]); mwPage.Revisions[0].Text = strings.Replace(mwPage.Revisions[0].Text, mwPage.Revisions[0].Text[refIndex[refNum][0]:refIndex[refNum][1]], "[r:" + strconv.Itoa(refNum) + "]", -1) }
				linkIndex := linkRE.FindAllStringIndex(mwPage.Revisions[0].Text, -1); for linkNum, _ := range linkIndex { link := strings.Split(mwPage.Revisions[0].Text[linkIndex[linkNum][0]:linkIndex[linkNum][1]], "|")[2:][0]; /* TODO */ graph.HashTable.Articles = append(graph.HashTable.Articles, link); if graph.Articles[mwPage.Title].Links[*] == false { graph.Articles[link] = Article{}; parentVertex := make([]*string, 0); parentVertex = append(parentVertex, &mwPage.Title) } else { /* graph.Articles[link].ParentVertice.appendChildParentPointer(graph, mwPage.Title, link); */ ; article.Links.appendChild(graph.Articles[link], mwPage.Title) }; mwPage.Revisions[0].Text = strings.Replace(mwPage.Revisions[0].Text, mwPage.Revisions[0].Text[linkIndex[linkNum][0]:linkIndex[linkNum][1]], "[l:" + strconv.Itoa(linkNum) + "]", -1) }
				article.Sections = sections
			};
			words := fasttag.WordsToSlice(mwPage.Revisions[0].Text); posTags := fasttag.BrillTagger(words); for posNum, _ := range posTags { if string(posTags[posNum]) != "N" { continue }; for curNounIndex = strings.Index(mwPage.Revisions[0].Text[curNounIndex:], words[posNum]); curNounIndex != -1; {
				for i := 0; i < 3; i++ { previousSentenceIndex, sentenceIndex = strings.LastIndexAny(mwPage.Revisions[0].Text[previousSentenceIndex:curNounIndex], sentenceStop), strings.IndexAny(mwPage.Revisions[0].Text[curNounIndex+len(words[posNum])+sentenceIndex+1:], sentenceStop) };
				article.Nouns[words[posNum]].Sentences.appendSentence(Sentence{Start: []byte(strconv.Itoa(previousSentenceIndex)), End: []byte(strconv.Itoa(sentenceIndex))})
			};
				graph.Articles[mwPage.Title] = article
			};
		};
	}; return err
}

func (articles *Articles) appendChild(childArticle Article, parentArticleTitle string) {
	childArticle.ParentVertice = append(childArticle.ParentVertice, &parentArticleTitle)
	*articles = append(*articles, childArticle)
}
func (sentences Sentences) appendSentence(sentence Sentence) {
	sentences = append(sentences, sentence)
}

func (graph Graph) WriteIndexAndContentData(directory *string) (err error) {
	indexFile, err := os.OpenFile(*directory+"/"+"index.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770); if err != nil { return err }; defer indexFile.Close(); contentFile, err := os.OpenFile(*directory+"/"+"content.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770); if err != nil { return err }; defer contentFile.Close()
	var gzipBuffer bytes.Buffer; compressor := gzip.NewWriter(&gzipBuffer); var sectionNum int = -1; var cFSize int
	for articleTitle, _ := range graph.Articles { /* format: (sectionTitle\{gzipped(sectionContent)&(gzipped(sectionReferenceContent)_.)&(gzipped(sectionLinkContent)_.)/.) */
		// For the sake of 1 vs multiple HDD Operations one buffer is made for the compressed data and that buffer is used for 1 write operation (the individual types could be asserted to the compressed form and then individually written with no effect on overall PC performance, but effect on program performance happens. - Extra memory is consumed by the RAM by below loops though.
		for sectionTitle, _ := range graph.Articles[articleTitle].Sections {
			gzipBuffer.WriteString(sectionTitle + "{");
			compressor.Write([]byte(graph.Articles[articleTitle].Sections[sectionTitle]))
			for refNum, _ := range graph.Articles[articleTitle].References { compressor.Write([]byte(graph.Articles[articleTitle].References[refNum])); if refNum < len(graph.Articles[articleTitle].References) { gzipBuffer.WriteString("|") } }; gzipBuffer.WriteString("&")
			sectionNum++; if sectionNum < len(graph.Articles[articleTitle].Sections) { gzipBuffer.WriteString("/") }
		}
		fmt.Fprintln(contentFile, gzipBuffer); cFSize = cFSize + len(gzipBuffer.Bytes()) + 1; fmt.Fprintln(indexFile, articleTitle + ":" + string(strconv.Itoa(cFSize))); gzipBuffer.Reset(); compressor.Reset(&gzipBuffer) }; return
}

func (article *Article) compress(compressor gzip.Writer, gzipBuffer bytes.Buffer) {
	for linkNum, childArticle := range article.ParentVertice { compressor.Write([]byte(*childArticle)); if linkNum < len(article.ParentVertice) { gzipBuffer.WriteString("|") }; linkNum++ }
}

func (graph Graph) ReadIndex_StageOne(directory *string) (err error) {
	indexFile, err := os.Open(*directory + "/index.txt"); bufioReader := bufio.NewReader(indexFile); var alphabet string = "abcdefghijklmnopqrstuvwxyz"; var alphabeticCounter uint8 = 0
	for { line, _, err := bufioReader.ReadLine(); if err != nil { if err == io.EOF { break } else { return err } }; parts := strings.Split(string(line), ":"); offset, err := strconv.Atoi(parts[1]); if alphabeticCounter == 0 { graph.AlphabeticIndex[alphabeticCounter] = int64(offset); alphabeticCounter++ } else { if !strings.EqualFold(string(parts[1][0]), string(alphabet[alphabeticCounter])) { graph.AlphabeticIndex[alphabeticCounter] = int64(offset) } } }; return
}

func (graph Graph) TfIdf() (err error) {
	var uint_wordOccurenceSum, uint_minOccurence, uint_maxOccurence, uint_wordSum uint = 0, 0, 0, 0; var f64_mean, f64_stdDev float64 = 0.0, 0.0

	// MeanLinks has nothing to do with TF-IDF, but with weighing articles based on links, so instead of looping again in the WeighArticles method the counting is done here.
	graph.SumLinks = 0
	for articleTitle, _ := range graph.Articles {
		min, err := strconv.ParseUint(string(graph.MinLinksOccurence[&Article]), 10, 0); if err != nil { return err }
		max, err := strconv.ParseUint(string(graph.MinLinksOccurence[&Article]), 10, 0); if err != nil { return err }
		switch {
		case max < uint(len(graph.Articles[articleTitle].Links)):
			graph.MinLinksOccurence[&] = []byte(strconv.FormatUint(uint(len(graph.Articles[articleTitle].Links)), 10));
		case min > uint(len(graph.Articles[articleTitle].Links)):
			graph.MaxLinksOccurence[*articleTitle] = []byte(strconv.FormatUint(uint(len(graph.Articles[articleTitle].Links)), 10));
		}
		graph.SumLinks += uint(len(graph.Articles[articleTitle].Links)); for wordTitle, _ := range graph.Articles[articleTitle].Nouns { graph.Articles[articleTitle] += uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)) switch { case uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)) > uint_maxOccurence: uint_maxOccurence = uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)); case uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)) < uint_minOccurence: uint_minOccurence = uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)) } uint_wordSum += uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)) } }
	f64_mean = float64(uint_wordOccurenceSum) / float64(uint_wordSum)

	for articleTitle, _ := range graph.Articles { for wordTitle, _ := range graph.Articles[articleTitle].Nouns { f64_stdDev += math.Pow(float64(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences))-f64_mean, float64(2)) } }
	// f64_stdDev = math.Sqrt(float64(uint_wordSum)*f64_stdDev)
	f64_stdDev = math.Sqrt(float64(uint_wordSum)*f64_mean)

	// TODO: "Bug", uint_maxOccurence and uint_minOccurence has to be for each and every word (fix ZScore.declare method)
	for articleTitle, _ := range graph.Articles { for wordTitle, _ := range graph.Articles[articleTitle].Nouns { graph.Articles[articleTitle].Nouns[wordTitle].MetaData.ZScore.declare(f64_stdDev, uint_maxOccurence, uint_minOccurence); graph.Articles[articleTitle].Nouns[wordTitle].MetaData.Extremum.declare(uint(len(graph.Articles[articleTitle].Nouns[wordTitle].Sentences)), f64_stdDev, uint_maxOccurence, uint_minOccurence) } }
	return err
}

func (graph Graph) WriteMetaData(directory, baseArticle string) (err error) {
	mdatFile, err := os.OpenFile(directory+"/"+baseArticle+".mdat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770);
	if err != nil { return err };
	defer mdatFile.Close();
	var buffer bytes.Buffer;
	gobEnc := gob.NewEncoder(&buffer);
	for articleTitle, _ := range graph.Articles { /* format: articleTitle:word:gob_encoded(WordData) */
		for word, _ := range graph.Articles[articleTitle].Nouns {
			_, err = buffer.Write([]byte(word + ":"))
			err = gobEnc.Encode(graph.Articles[articleTitle].Nouns[word]);
			if err != nil { return err };
		}
		fmt.Fprintln(mdatFile, articleTitle + ":" + buffer.String());
		buffer.Reset() };

	return
}
/* func (graph Graph) ReadIndex_StageTwo(directory *string, graphOffset int, letterOffset *uint8) (err error) {
	indexFile, err := os.Open(*directory + "/index.txt");
	bufioReader := bufio.NewReader(indexFile);
	var alphabet string = "abcdefghijklmnopqrstuvwxyz" // TODO: sequence the readline loop to divide the index into 250 parts for every letter in the alphabet to read one part 
	_, err = indexFile.Seek(graph.AlphabeticIndex[int(*letterOffset)], -1);
	// if err != nil { return err 
	// if err != nil { return func (graph Graph) WriteMetaData(directory, baseArticle string) (err error) {
	if err != nil { return err }
	return
} */
func (graph Graph) ReadMetaData(directory, baseArticle string) (err error) {
	file, err := os.OpenFile(directory+"/"+baseArticle+".mdat", os.O_RDONLY, 0700);
	if err != nil { return err };
	defer file.Close();
	bufioReader := bufio.NewReader(file);
	var line []byte;
	var buffer bytes.Buffer;
	gobDecoder := gob.NewDecoder(&buffer);
	var table WordMetaData
	for {
		line, _, err = bufioReader.ReadLine();
		if err != nil {
			if err == io.EOF {
				break
			};
			return err
		}
		unformattedLine := strings.Split(string(line), ":")
		_, err = buffer.Read([]byte(unformattedLine[2]))
		if err != nil {
			return err
		};
		err = gobDecoder.Decode(&table)
		if err != nil {
			return err
		}
		graph.Articles[unformattedLine[0]].Nouns[unformattedLine[1]].MetaData.declare(table)
	}
	return
}
func (wordMetaData WordMetaData) declare(metadata WordMetaData) {
	wordMetaData = metadata
}

func (graph Graph) WeighArticles() {
	for articleTitle, _ := range graph.Articles {
		for _, linkArticle1 := graph.Articles[articleTitle1].Links {
			if graph.Articles[articleTitle] == *linkArticle1 { linkArticle1.ParentVertice.appendLink(&linkArticle1); graph.Articles[articleTitle].LinkedAmount.increment() }
			for _, linkArticle2 := range *linkArticle1 {
				if graph.Articles[articleTitle] == *linkArticle2 { linkArticle2.ParentVertice.appendLink(&linkArticle2); graph.Articles[articleTitle].LinkedAmount.increment() }
			}
		}
	}
	stdDev := math.Sqrt(float64(graph.SumLinks)*(uint(len(graph.Articles)) / graph.SumLinks))
	for articleTitle, _ := range graph.Articles {
		// TODO: "Bug", graph.MaxLinksOccurence and graph.MinLinksOccurence has to be for each and every article, and grapgh.*LinksOccurence has to be variables inside graph.Articles[<article title>]. - also fix TfIdf method's uint_maxOccurence and uint_minOccurence counter, they have to be variables declared by an assignment from graph.Articles[<article title>].Sentences.
		graph.Articles[articleTitle].ZScore.declare(stdDev, graph.MaxLinksOccurence, Min uint)
	}
}

// TODO: This is a method which will produce a final product from the wikipedia articles, e.g this will take up extra space (duplicate text from content.dat and <article>.mdat) so it should preferrably be used on a disk drive with enough space or on a USB.
func (graph Graph) WriteTxt() (err error) {
	var depthCounter uint8 = 0
	// fWriter := bufio.NewWriter(ioWriter)
	indexFile, err := os.Create("index-" + articleName + ".org"); if err != nil { return err }; defer indexFile.Close()
	for articleDepth, _ := range graph.Final.Articles {
		for articleRankAndTitle, _ := range graph.Final.Articles[articleDepth] {
			file, err := os.Create(*graph.Final.Articles[articleDepth] + ".org"); if err != nil { return err }; defer file.Close()
			file.WriteString("* " + *graph.Final.Articles[articleDepth]); file.WriteString("** Sections"); for sectionName, sectionText := range graph.Articles[articleRankAndTitle].Sections { file.WriteString("*** " + sectionName); file.WriteString("    " + sectionText) } }
		/* TODO: Update the articles.Pageranks to fit with Graph.Final.Articles */
		for num, item := range graph.Articles[articleRankAndTitle] {
			depthCounter++
			indexFile.WriteString("* " + strconv.Itoa(int(depth)), "-", articleRankAndTitle)
			if depthCounter == 7 { depthCounter = 0; break }
			/* pagerankIndex := make(map[float64]string)
			for article, pagerank := range articles[articleName].Pageranks[depth] {
				float64Pr, err := strconv.ParseFloat(string(pagerank), 64)
				if err != nil { return err }
				pageranks[articleName] = append(pageranks[articleName], Pagerank(float64Pr))
				pagerankIndex[float64Pr] = article
			}
			sort.Sort(sort.Reverse(SortedPageranks(pageranks[articleName])))
			for pagerank, article := range pagerankIndex {
				indexFile.WriteString("** " + article + " - " + strconv.FormatFloat(pagerank, 'f', 6, 64))
			} */
		}
	}
	return nil
}

func main() {
	readDir := flag.String("readDir", "", "The directory contatining the bzipped wikipedia files (full articles, preferrably the ones separated into several files depending on your available RAM).")
	writeDir := flag.String("writeDir", "", "The directory the chosen graph of the chosen base article will be written to. If no value is given a directory with the name of the base article is made.")
	flag.Parse()
	

	// TODO (latest update as of 4/15-18): 1) decide ReadIndexStageTwo for disk read performance optimization, 2) calculate ZScore for Article with the LR functions, 3) implement a dijkstra function for Articles (see legacy/dijkstra in GOPATH/src/legacy), 4) A SortBy method for graph.Articles[...].Nouns[...].MetaData.* and graph.Articles[...].ZScore (the structs require pointers and these will have the actual SortBy methods and Sort requirements (len, less, swap(...))).
	//    TODO: add the WriteTxt method from the sample-mysql/boltdb.go file - in GOPATH/src/legacy/pressure679/WikiPagerankDB.
	//       TODO: Actually utilize the existing code to the endpoint (WriteTxt method).
	//    TODO: Add a chatbot (markov/viterbi chain generator with a keyword adder from Bloom's taxonomical model) functionality which reads existing user content for active learning. - this will include the dijkstra method/package. - Edit (4/17-18): make a BST/decision tree (actually a C4.5 decision tree by nature) which depends it's decision by parsing which child node to choose from a matrix set (Final struct) which parses amount of occurences of keywords from the Final struct in an input text (text given from user).
}
