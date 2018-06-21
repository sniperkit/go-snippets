package market

import (
	"errors"
	"log"
	"stock"
)

type Market struct {
	DB        map[string]string
	Name      string
	DBBuilder MarketDBBuilder
}

type MarketDBBuilder func() ([]*stock.Stock, error)

func (m *Market) BuildDB() error {
	if m.DBBuilder == nil {
		log.Println("Market ", m.Name, " is not initialized")
		return errors.New("Builder is not set")
	}
	stocks, err := m.DBBuilder()
	if err != nil {
		log.Println(err)
		return errors.New("DB build failed")
	}

	m.DB = make(map[string]string, len(stocks))
	for _, stock := range stocks {
		m.DB[stock.Code] = stock.Name
	}

	return nil
}
