package main

import (
	"crawler"
	"errors"
	"log"
	"market"
	"regexp"
	"stock"
	"stockdb"
)

var shanghai = "http://quote.eastmoney.com/stocklist.html#sh"
var shenzhen = "http://quote.eastmoney.com/stocklist.html#sz"

var stockRe = regexp.MustCompile(`(?P<name>[[\p{Han}[:word:]\*]+)\((?P<code>[[:word:]]+)\)`)

func shenZhenDBBuidler() ([]*stock.Stock, error) {
	s, err := dBBuilder(shenzhen)
	if err != nil {
		return nil, errors.New("Build Db failed for: Shenzhen")
	}

	return s, nil
}

func shangHaiDBBuidler() ([]*stock.Stock, error) {
	s, err := dBBuilder(shanghai)
	if err != nil {
		return nil, errors.New("Build Db failed for: Shenzhen")
	}

	return s, nil
}

func dBBuilder(page string) ([]*stock.Stock, error) {
	stocks, err := crawler.Crawl(shenzhen, "div.quotebody>div#quotesearch>ul>li")
	if err != nil {
		log.Println("Cannot get DB")
		return nil, errors.New("Cannot get DB")
	}
	db := make([]*stock.Stock, 0, len(stocks))

	for _, s := range stocks {
		match := stockRe.FindStringSubmatch(s)
		if len(match) == 3 {
			var s stock.Stock
			s.Name = match[1]
			s.Code = match[2]
			db = append(db, &s)
		}
	}

	if len(db) == 0 {
		log.Println("No stock find for market shenzhen")
		return nil, errors.New("No stock find for market shenzhen")
	}

	return db, nil
}

func main() {
	db := stockdb.StockDB{}
	db.Store = make(map[string]*market.Market, 2)
	db.RegisterMarket("ShenZhen", shenZhenDBBuidler)
	db.RegisterMarket("ShangHai", shangHaiDBBuidler)
	db.Dump()
}
