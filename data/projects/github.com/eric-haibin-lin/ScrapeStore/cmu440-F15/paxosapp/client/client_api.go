package client

type ClientNode interface {
	// Crawls the internet with given root url until it reaches the number of pages sepecified.
	Crawl(args *CrawlArgs, reply *CrawlReply) error

	// GetLinks fetches all links contained in the given link (if any) from ScrapeStore
	GetLinks(args *GetLinksArgs, reply *GetLinksReply) error

	// RunRageRank retrieves all the crawled link relationships from ScrapeStore and runs
	// pagerank algorithm on them. It also saves the page rank of each url back to ScrapeStore
	// for further references
	RunPageRank(args *PageRankArgs, reply *PageRankReply) error

	// GetRank fetches the page rank for the requested url from ScrapeStore
	GetRank(args *GetRankArgs, reply *GetRankReply) error
}
