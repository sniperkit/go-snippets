package stockdb

import (
	"errors"
	"log"
	"market"
)

type StockDB struct {
	Store map[string]*market.Market
}

func (sdb *StockDB) RegisterMarket(name string, dbbuilder market.MarketDBBuilder) error {
	if _, ok := sdb.Store[name]; !ok {
		sdb.Store[name] = &market.Market{
			Name:      name,
			DBBuilder: dbbuilder,
		}

		err := sdb.Store[name].BuildDB()
		if err != nil {
			return errors.New("Cannot build stock DB for market " + name)
		}

		return nil
	}

	return errors.New("Same Market already exist")
}

func (sdb *StockDB) Dump() {
	for k, v := range sdb.Store {
		log.Println("All stocks for: ", k)
		for c, n := range v.DB {
			log.Println("Code: ", c, " Name: ", n)
		}
		log.Println("All stocks for: ", k, ", Total count: ", len(v.DB))
	}
}
