//The runner for ClientNode

package main

import (
	"flag"
	"fmt"
	"github.com/cmu440-F15/paxosapp/client"
	"github.com/cmu440-F15/paxosapp/common"
	"log"
	"strings"
	"io/ioutil"
)

var (
	port       = flag.String("port", "", "port for client node")
	opt        = flag.String("opt", "", "option for client node, either crawl or pagerank")
	num        = flag.Int("num", 10, "number of webpages to crawl")
	url        = flag.String("url", "http://news.google.com", "root url to start crawling with")
	masterPort = flag.String("masterPort", "", "port for master node")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

// main routine creates a client node, connects with a master node and invokes 
// corresponding method on the client 
func main() {
	flag.Parse()
	log.SetOutput(ioutil.Discard)
	myHostPort := "localhost:" + *port
	masterHostPort := strings.Split(*masterPort, ",")
	log.Println("crunner creates client node on", myHostPort, "with master on ", *masterPort)

	// Create and start the client Node.
	cli, err := client.NewClientNode(myHostPort, masterHostPort)
	if err != nil {
		log.Fatalln("Failed to create client node:", err)
	}
	
	//Crawl
	if *opt == "crawl" {
		fmt.Println("crunner invokes Crawl on client")
		var args client.CrawlArgs
		var reply client.CrawlReply
		absoluteUrl := *url
		//Attach http:// prefix for raw domain name 
		if !strings.Contains(*url, "http://") && !strings.Contains(*url, "https://") {
			absoluteUrl = "http://" + *url
		}
		args.RootUrl = absoluteUrl
		args.NumPages = *num
		cli.Crawl(&args, &reply)
	//The the links contained in the root url
	} else if *opt == "getlink" {
		fmt.Println("crunner invokes GetLink on client")
		var getLinkArgs client.GetLinksArgs
		getLinkArgs.Url = common.RemoveHttp(*url)
		var getLinkReply client.GetLinksReply
		cli.GetLinks(&getLinkArgs, &getLinkReply)
	//Run the page rank algorithm on all url dataset we have crawled so far
	} else if *opt == "pagerank" {
		fmt.Println("crunner invokes RunPageRank on client")
		var pageRankArgs client.PageRankArgs
		var pageRankReply client.PageRankReply
		cli.RunPageRank(&pageRankArgs, &pageRankReply)
	//Get the page rank we calculated for this root url
	} else if *opt == "getrank" {
		fmt.Println("crunner invokes GetRank on client")
		var getRankArgs client.GetRankArgs
		getRankArgs.Url = common.RemoveHttp(*url)
		var getRankReply client.GetRankReply
		cli.GetRank(&getRankArgs, &getRankReply)
	}
}
