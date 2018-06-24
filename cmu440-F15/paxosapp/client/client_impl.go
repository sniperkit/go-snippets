package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/cmu440-F15/paxosapp/collectlinks"
	"github.com/cmu440-F15/paxosapp/common"
	"github.com/cmu440-F15/paxosapp/pagerank"
	"github.com/cmu440-F15/paxosapp/rpc/paxosrpc"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"time"
)

type Status int

const (
	OK   Status = iota + 1 // Status OK
	Fail                   // Status Failure
)

const (
	tolerance         = 0.0001 //the default tolerance when running pagerank algorithm
	followProbability = 0.85   //the defaul probabilitic assumption of following a link
)

type clientNode struct {
	//the connection with the master node
	conn       common.Conn
	myHostPort string
	//a channel to receive all url's relationship before writing to ScrapeStore
	linkRelationChan chan linkRelation
	//a channel to receive all url's
	allLinksChan chan string
	//a channel to receive all url's to crawl next time
	nextLinkChan chan string
	//a map to indicate whether the link has been crawled or not
	visited map[string]bool
	//the http client who fetches the page
	httpClient http.Client
	//a mapping between the url and its id for pagerank calculation
	idMap  map[string]int
	urlMap map[int]string
	//the next id for next url
	nextId int
}

type linkRelation struct {
	follower string   //The root url
	followee []string //The urls followed by root url
}

type CrawlArgs struct {
	RootUrl  string //The root url where crawling started
	NumPages int    //The number of web pages to crawl
}

type CrawlReply struct {
	Status Status
}

type GetLinksArgs struct {
	Url string //The Url requested
}

type GetLinksReply struct {
	List []string // The list of all urls
}

type GetRankArgs struct {
	Url string //The Url requested
}

type GetRankReply struct {
	Value float64 //The rank
}

type PageRankArgs struct {
	//Nothing to provide here
}

type PageRankReply struct {
	Status Status
}

// NewClientNode creates a new ClientNode. This function should return only when
// it successfully dials to a master node, and should return a non-nil error if the
// master node could not be reached.
// masterHostPort is the list of hostnames and port numbers to master nodes
func NewClientNode(myHostPort string, masterHostPort []string) (ClientNode, error) {
	log.SetOutput(ioutil.Discard)
	defer log.Println("Leaving NewClientNode")
	log.Println("myhostport is", myHostPort, ", hostPort of masterNode to connect is", masterHostPort)
	var a ClientNode
	//Pick a master node and dial on it
	index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(masterHostPort))
	numRetries := 3
	dialer, err := rpc.DialHTTP("tcp", masterHostPort[index])
	for err != nil {
		numRetries -= 1
		dialer, err = rpc.DialHTTP("tcp", masterHostPort[index])
		if numRetries <= 0 {
			return nil, errors.New("Fail to dial to master:" + masterHostPort[index])
		}
	}
	//Save the hostport and dialer for master node, initialize basic member variables
	var conn common.Conn
	conn.HostPort = masterHostPort
	conn.Dialer = dialer
	node := clientNode{}
	node.conn = conn
	node.myHostPort = myHostPort
	node.linkRelationChan = make(chan linkRelation)
	node.allLinksChan = make(chan string)
	node.nextLinkChan = make(chan string)
	node.visited = make(map[string]bool)
	node.idMap = make(map[string]int)
	node.urlMap = make(map[int]string)
	// Initialize the http client for the crawler
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	node.httpClient = http.Client{Transport: transport}
	a = &node
	return a, nil
}

// Crawls the internet with given root url until it reaches the number of pages sepecified.
func (cn *clientNode) Crawl(args *CrawlArgs, reply *CrawlReply) error {
	defer log.Println("Leavin Crawl on client", cn.myHostPort)
	log.Println("Crawl invoked on ", cn.myHostPort)
	//set root url as the starting point of crawling
	rootUrl := args.RootUrl
	go func() { cn.allLinksChan <- rootUrl }()
	go cn.checkVisited()
	//start to crawl and count the number of pages crawled!
	count := 0
	for {
		select {
		case uri := <-cn.nextLinkChan:
			//get the next link and crawl all pages starting from there
			if count < args.NumPages {
				rel, err := cn.doCrawl(uri)
				if err == nil {
					count += 1
					//Writing the crawl result to ScrapeStore
					var appendArgs paxosrpc.AppendArgs
					appendArgs.Key = rel.follower
					appendArgs.Value = rel.followee
					var appendReply paxosrpc.AppendReply
					log.Println("Calling PaxosNode.Append")
					for _, v := range appendArgs.Value {
						log.Println("Writing", v, "with key", appendArgs.Key)
					}
					cn.conn.Dialer.Call("PaxosNode.Append", &appendArgs, &appendReply)
				}
			} else {
				return nil
			}
		}
	}
	return nil
}

// RunRageRank retrieves all the crawled link relationships from ScrapeStore and runs
// pagerank algorithm on them. It also saves the page rank of each url back to ScrapeStore
// for further references
func (cn *clientNode) RunPageRank(args *PageRankArgs, reply *PageRankReply) error {
	log.Println("RunPageRank invoked on ", cn.myHostPort)
	//Get all page relationships from ScrapeStore
	var getAllLinksArgs paxosrpc.GetAllLinksArgs
	var getAllLinksReply paxosrpc.GetAllLinksReply
	log.Println("Invoking PaxosNode.GetAllLinks on", cn.conn.HostPort)
	err := cn.conn.Dialer.Call("PaxosNode.GetAllLinks", &getAllLinksArgs, &getAllLinksReply)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Calculating page rank on all links...")
	//calculate page rank!
	pageRankEngine := pagerank.New()
	for k, list := range getAllLinksReply.LinksMap {
		followerId := cn.getId(k)
		for _, v := range list {
			followeeId := cn.getId(v)
			pageRankEngine.Link(followerId, followeeId)
		}
	}
	//Save all page rank results to ScrapeStore
	pageRankEngine.Rank(followProbability, tolerance, func(label int, rank float64) {
		fmt.Println(cn.urlMap[label], rank*100)
		var putRankArgs paxosrpc.PutRankArgs
		putRankArgs.Key = cn.urlMap[label]
		putRankArgs.Value = rank * 100
		var putRankReply paxosrpc.PutRankReply
		err := cn.conn.Dialer.Call("PaxosNode.PutRank", &putRankArgs, &putRankReply)
		if err != nil {
			fmt.Println(err)
		}
	})
	return nil
}

// getId returns the id mapped from given url
func (cn *clientNode) getId(url string) int {
	id, ok := cn.idMap[url]
	//If url doesn't exist in the current url - id map, create a new mapping for it
	if !ok {
		cn.idMap[url] = cn.nextId
		cn.urlMap[cn.nextId] = url
		cn.nextId += 1
	}
	return id
}

// GetRank fetches the page rank for the requested url from ScrapeStore
func (cn *clientNode) GetRank(args *GetRankArgs, reply *GetRankReply) error {
	log.Println("GetRank invoked on ", cn.myHostPort)
	var rankArgs paxosrpc.GetRankArgs
	rankArgs.Key = args.Url
	var rankReply paxosrpc.GetRankReply
	err := cn.conn.Dialer.Call("PaxosNode.GetRank", &rankArgs, &rankReply)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("GetRank returns value:", rankReply.Value)
	reply.Value = rankReply.Value
	return nil
}

// GetLinks fetches all links contained in the given link (if any) from ScrapeStore
func (cn *clientNode) GetLinks(args *GetLinksArgs, reply *GetLinksReply) error {
	log.Println("GetLink invoked on ", cn.myHostPort)
	var linksArgs paxosrpc.GetLinksArgs
	linksArgs.Key = args.Url
	var linksReply paxosrpc.GetLinksReply
	log.Println("Calling PaxosNode.GetLinks with key", linksArgs.Key)
	err := cn.conn.Dialer.Call("PaxosNode.GetLinks", &linksArgs, &linksReply)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("GetLinks returns list:")
	for _, v := range linksReply.Value {
		fmt.Println(v)
	}
	reply.List = linksReply.Value
	return nil
}

// checkVisited checks if all links crawled are visited. If not, push it to the
// nextLinksChan to start crawling from there
func (cn *clientNode) checkVisited() {
	for val := range cn.allLinksChan {
		if !cn.visited[val] {
			cn.visited[val] = true
			cn.nextLinkChan <- val
		}
	}
}

// doCrawl fetches the webpage and collects the links contained in this webpage.
// It returns a linkRelation which contains the following relationship of a root url
// and all links this root url points to.
func (cn *clientNode) doCrawl(uri string) (linkRelation, error) {
	var rel linkRelation
	rel.follower = common.RemoveHttp(uri)
	rel.followee = make([]string, 0)
	resp, err := cn.httpClient.Get(uri)
	if err != nil {
		log.Println(err)
		return rel, errors.New("Failed to fetch page" + uri)
	}
	defer resp.Body.Close()
	fmt.Println(uri)
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		absolute := common.GetAbsoluteUrl(link, uri)
		if absolute != "" {
			go func() { cn.allLinksChan <- absolute }()
			removeHttps := common.RemoveHttp(absolute)
			rel.followee = append(rel.followee, removeHttps)
		}
	}
	return rel, nil
}
