package crawler

import (
	"encoding/gob"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"bufio"
)

type visit struct {
	rawUrl       string
	requestUrl   string
	foundUrls    []string
	indexInArray int
}

type job struct {
	refererUrl string
	rawUrl     string
}

type Crawler struct {
	options      *Options
	VisitedUrls  []string
	LinksOnPages []int

	Matrix map[string]map[string]int
	visits map[string]visit

	queue []job
}

func NewCrawler(opts *Options) *Crawler {
	ret := new(Crawler)
	ret.options = opts

	ret.visits = make(map[string]visit)
	ret.Matrix = make(map[string]map[string]int)
	ret.VisitedUrls = []string{}
	ret.LinksOnPages = []int{}

	return ret
}

func (c *Crawler) addToMatrix(refererUrl string, requestUrl string) {
	if c.Matrix[refererUrl] == nil {
		c.Matrix[refererUrl] = make(map[string]int)
	}
	c.Matrix[refererUrl][requestUrl] += 1
}

func (c *Crawler) Run(entryUrl string, outputf *os.File) {
	v, err := c.visitUrl(entryUrl)

	if err != nil {
		return
	}

	c.visits[v.rawUrl] = v
	c.visits[v.requestUrl] = v
	v.indexInArray = len(c.VisitedUrls)
	c.VisitedUrls = append(c.VisitedUrls, v.requestUrl)
	c.LinksOnPages = append(c.LinksOnPages, len(v.foundUrls))
	fmt.Println(len(c.VisitedUrls), "—", "New link:", v.requestUrl)

	for _, u := range v.foundUrls {
		c.queue = append(c.queue, job{v.requestUrl, u})
	}

	for len(c.queue) != 0 {
		q := c.queue[0]
		c.queue = c.queue[1:]

		if prevVisit, ok := c.visits[q.rawUrl]; ok {
			c.addToMatrix(q.refererUrl, prevVisit.requestUrl)
			continue
		}

		if c.options.MaxVisits > 0 && len(c.VisitedUrls) >= c.options.MaxVisits {
			continue
		}

		v, err := c.visitUrl(q.rawUrl)

		if err != nil {
			fmt.Printf("%v (Referer: %v)\n", err, q.refererUrl)
			continue
		}

		if prevVisit, ok := c.visits[v.requestUrl]; !ok {
			v.indexInArray = len(c.VisitedUrls)
			c.visits[v.rawUrl] = v
			c.visits[v.requestUrl] = v
			c.VisitedUrls = append(c.VisitedUrls, v.requestUrl)
			c.LinksOnPages = append(c.LinksOnPages, len(v.foundUrls))
			fmt.Println(len(c.VisitedUrls), "—", "New link:", v.requestUrl)
			for _, u := range v.foundUrls {
				c.queue = append(c.queue, job{v.requestUrl, u})
			}
		} else {
			c.visits[v.rawUrl] = prevVisit
		}

		c.addToMatrix(q.refererUrl, v.requestUrl)
	}

	if outputf != nil {
		fmt.Println("Saving to a file...")
		c.saveMatrix(outputf)
	} else {
		fmt.Println("Calculating PageRank...")
		p := c.evaluatePagerank()

		sum := 0.0
		for i, pagerank := range p {
			sum += pagerank
			fmt.Println(c.VisitedUrls[i], pagerank)
		}

		fmt.Println(sum)
	}
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return
}

func (c *Crawler) isSameHost(u *url.URL, aUrl string) bool {
	au, err := url.Parse(aUrl)

	if err != nil {
		return false
	}

	return au.Host == u.Host
}

func isHtml(res *http.Response) bool {
	contentType := res.Header.Get("Content-type")

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == "text/html" {
			return true
		}
	}
	return false
}

func (c *Crawler) visitUrl(urlString string) (visit, error) {
	res, err := http.Head(urlString)

	v := visit{
		requestUrl: res.Request.URL.String(),
		rawUrl:     urlString,
		foundUrls:  []string{},
	}

	if !isHtml(res) {
		fmt.Println(v.requestUrl, "is not HTML page")
		return v, nil
	}

	if err != nil {
		return v, errors.New("ERROR: Failed to crawl \"" + urlString + "\"")
	}

	res, err = http.Get(urlString)

	if err != nil {
		return v, errors.New("ERROR: Failed to crawl \"" + urlString + "\"")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return v, errors.New("ERROR: Failed to crawl \"" + urlString + "\" (code " + strconv.Itoa(res.StatusCode) + ")")
	}

	z := html.NewTokenizer(res.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return v, nil
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, aUrl := getHref(t)
			if !ok {
				continue
			}

			au, err := url.Parse(strings.Trim(aUrl, " "))

			if err != nil {
				continue
			}

			if au.Scheme == "" {
				au.Scheme = res.Request.URL.Scheme
			}

			if au.Host == "" {
				au.Host = res.Request.URL.Host
			}

			if !strings.HasPrefix(au.Path, "/") {
				au.Path = res.Request.URL.Path + au.Path
			}

			au.Fragment = ""

			aUrl = au.String()

			if (au.Scheme != "http" && au.Scheme != "https") || (c.options.SameHostOnly && !c.isSameHost(res.Request.URL, aUrl)) {
				continue
			}

			v.foundUrls = append(v.foundUrls, aUrl)
		}
	}
}

func (c *Crawler) pagerankIterate(probabilityOfTransitionToRandomPage float64, p []float64, danglingNodes []int) []float64 {
	size := len(p)
	innerProduct := 0.0

	for _, danglingNode := range danglingNodes {
		innerProduct += p[danglingNode]
	}

	innerProductOverSize := innerProduct / float64(size)

	norm := 0.0
	new_p := make([]float64, size)

	for i := 0; i < size; i++ {
		sum := 0.0

		for j := 0; j < size; j++ {
			numberOfLinkFromJToI := c.Matrix[c.VisitedUrls[j]][c.VisitedUrls[i]]
			if numberOfLinkFromJToI > 0 {
				sum += p[j] * float64(numberOfLinkFromJToI) / float64(c.LinksOnPages[j])
			}
		}
		new_p[i] = c.options.FollowingProb*(sum+innerProductOverSize) + probabilityOfTransitionToRandomPage

		norm += new_p[i]
	}

	antinorm := 1.0 / norm

	for j := 0; j < size; j++ {
		new_p[j] *= antinorm
	}

	return new_p
}

func (c *Crawler) pagerankIterateParallel(channel chan struct {
	int
	float64
}, i int, p []float64) {
	size := len(p)

	result := 0.0

	for j := 0; j < size; j++ {
		numberOfLinkFromJToI := c.Matrix[c.VisitedUrls[j]][c.VisitedUrls[i]]
		if numberOfLinkFromJToI > 0 {
			result += p[j] * float64(numberOfLinkFromJToI) / float64(c.LinksOnPages[j])
		}
	}

	channel <- struct {
		int
		float64
	}{i, result}
}

func calculateChange(p, new_p []float64) float64 {
	acc := 0.0

	for i, pForI := range p {
		acc += math.Abs(pForI - new_p[i])
	}

	return acc
}

func (c *Crawler) calculateDanglingNodes() []int {
	danglingNodes := make([]int, 0, len(c.LinksOnPages))

	for i, numberOutLinksForI := range c.LinksOnPages {
		if numberOutLinksForI == 0 {
			danglingNodes = append(danglingNodes, i)
		}
	}

	return danglingNodes
}

func (c *Crawler) evaluatePagerank() []float64 {
	size := len(c.VisitedUrls)
	inverseOfSize := 1.0 / float64(size)
	probabilityOfTransitionToRandomPage := (1.0 - c.options.FollowingProb) / float64(size)
	danglingNodes := c.calculateDanglingNodes()

	p := make([]float64, size)
	for i := range p {
		p[i] = inverseOfSize
	}

	change := 2.0

	if c.options.Parallel {
		channel := make(chan struct {
			int
			float64
		})
		for change > c.options.Tolerance {
			for i := 0; i < size; i++ {
				go c.pagerankIterateParallel(channel, i, p)
			}

			innerProduct := 0.0

			for _, danglingNode := range danglingNodes {
				innerProduct += p[danglingNode]
			}

			innerProductOverSize := innerProduct / float64(size)

			new_p := make([]float64, size)
			norm := 0.0
			for i := 0; i < size; i++ {
				pair := <-channel
				new_p[pair.int] = c.options.FollowingProb*(pair.float64+innerProductOverSize) + probabilityOfTransitionToRandomPage
				norm += pair.float64
			}
			antinorm := 1.0 / norm

			for j := 0; j < size; j++ {
				new_p[j] *= antinorm
			}
			change = calculateChange(p, new_p)
			p = new_p
		}
	} else {
		for change > c.options.Tolerance {
			new_p := c.pagerankIterate(probabilityOfTransitionToRandomPage, p, danglingNodes)
			change = calculateChange(p, new_p)
			p = new_p
		}
	}

	return p
}

func (c *Crawler) ParseMatrix(file *os.File) {
	fmt.Println("Parsing file...")
	if c.options.FileType == "bin" {
		dec := gob.NewDecoder(file)
		err := dec.Decode(c)
		if err != nil {
			panic(err)
		}
	} else {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		size, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		c.VisitedUrls = make([]string, size)
		c.Matrix = make(map[string]map[string]int)
		c.LinksOnPages = make([]int, size)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			from, fromErr := strconv.Atoi(fields[0])
			to, toErr := strconv.Atoi(fields[1])
			if fromErr != nil || toErr != nil {
				panic(err)
			}
			c.VisitedUrls[from] = fields[0]
			c.VisitedUrls[to] = fields[1]
			c.LinksOnPages[from] += 1
			c.addToMatrix(fields[0], fields[1])
		}
	}
	fmt.Println("Calculating PageRank...")

	p := c.evaluatePagerank()

	sum := 0.0
	for i, pagerank := range p {
		sum += pagerank
		fmt.Println(c.VisitedUrls[i], pagerank)
	}

	fmt.Println(sum)
}

func (c *Crawler) saveMatrix(file *os.File) {
	if c.options.FileType == "bin" {
		enc := gob.NewEncoder(file)
		err := enc.Encode(&c)
		if err != nil {
			panic(err)
		}
	} else {
		size := len(c.VisitedUrls)
		file.WriteString(strconv.Itoa(size) + "\n")
		for fromIndex, fromUrl := range c.VisitedUrls {
			for toIndex := 0; toIndex < size; toIndex++ {
				for i := 0; i < c.Matrix[fromUrl][c.VisitedUrls[toIndex]]; i++  {
					file.WriteString(strconv.Itoa(fromIndex) + "\t" + strconv.Itoa(toIndex) + "\n")
				}
			}
		}
	}
}
