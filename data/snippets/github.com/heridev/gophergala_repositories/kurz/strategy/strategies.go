package strategy

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/FGM/kurz/storage"
	"log"
)

// StrategyMap adds helper methods to a map of AliasingStrategy
type StrategyMap map[string]AliasingStrategy

// IsValid() checks whether a strategy is registered for a given name.
func (m StrategyMap) IsValid(name string) bool {
	var ok bool
	_, ok = m[name]
	return ok
}

// Strategies is a global instance of the registered strategies.
var Strategies StrategyMap

func MakeStrategies(list []AliasingStrategy) StrategyMap {
	ret := make(map[string]AliasingStrategy, len(list))

	for _, s := range list {
		ret[s.Name()] = s
	}

	return ret
}

// StatisticsMap adds helper methods to a map of AliasingStrategy use counts.
type StatisticsMap map[string]int64

// Statistics is a global instance of the AliasingStrategy use counts.
var Statistics StatisticsMap

// Get() fetches the AliasingStrategy use counts from the database and returns the updated Statistics.
func (ss StatisticsMap) Refresh(s storage.Storage) StatisticsMap {
	var err error
	var strategyResult sql.NullString
	var countResult sql.NullInt64

	sql := `
SELECT strategy, COUNT(*)
FROM shorturl
GROUP BY strategy
	`

	rows, err := s.DB.Query(sql)
	if err != nil {
		log.Printf("Failed querying database for strategy statistics: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&strategyResult, &countResult); err != nil {
			log.Fatal(err)
		}

		if !Strategies.IsValid(strategyResult.String) {
			log.Fatalf("'%#v' is not a valid strategy\n", strategyResult)
		}
		ss[strategyResult.String] = countResult.Int64
	}

	for name, _ := range Strategies {
		_, ok := ss[name]
		if !ok {
			ss[name] = 0
		}
	}

	return ss
}

// String() implements the fmt.Stringer interface.
func (ss StatisticsMap) String() string {
	var buf bytes.Buffer

	for name, count := range ss {
		buf.WriteString(fmt.Sprintf("%-8s: %d\n", name, count))
	}

	return fmt.Sprint(buf.String())
}

// init() initializes the global variables provided by the package.
func init() {
	Strategies = MakeStrategies([]AliasingStrategy{
		baseStrategy{},
		HexCrc32Strategy{},
		ManualStrategy{},
	})

	Statistics = make(map[string]int64, len(Strategies))
}
