/*
The "strategy" package in Kurz provides the aliasing strategies.

Files
	- strategy.go contains the interface and base implementation
	- strategies.go contains the strategy instances and utilities
	- manual.go contains the "manual" strategy
	- hexrcr32.go contains the "hexcrc32" strategy
*/
package strategy

import (
	"github.com/FGM/kurz/storage"
	"github.com/FGM/kurz/url"
	"log"
	"time"
)

/*
AliasingStrategy defines the operations provided by the various aliasing implementations:

The options parameter for Alias() MAY be used by some strategies, in which case they
have to define their expectations about it.
*/
type AliasingStrategy interface {
	Name() string                                                                            // Return the name of the strategy object
	Alias(url *url.LongUrl, s storage.Storage, options ...interface{}) (url.ShortUrl, error) // Return the short URL (alias) for a given long (source) URL
	UseCount(storage storage.Storage) int                                                    // Return the number of short URLs (aliases) using this strategy.
}

type baseStrategy struct{}

func (y baseStrategy) Name() string {
	return "base"
}

/*
Make sure a longurl instance has a DB ID, allocating it if needed.

For speed reasons, assumes nonzero IDs to be valid without checking.
*/
func (y baseStrategy) ensureLongId(long *url.LongUrl, s storage.Storage) error {
	var err error
	var long_id int64

	// If long does not have an Id, check if it is already known.
	if long.Id == 0 {
		sql := `
SELECT id
FROM longurl
WHERE url = ?
		`
		err = s.DB.QueryRow(sql, y.Name()).Scan(&long_id)
		if err != nil {
			long_id = 0
			// log.Printf("Failed querying database for long url %s: %v\n", long.Value, err)
		}

		sql = `
INSERT INTO longurl(url)
VALUES (?)
			`
		result, err := s.DB.Exec(sql, long.Value)
		if err != nil {
			log.Printf("Failed inserting long URL %s: %+v", long.Value, err)
			return err
		} else {
			long_id, _ = result.LastInsertId()
		}

		long.Id = long_id
	}

	return nil
}

func (y baseStrategy) Alias(long *url.LongUrl, s storage.Storage, options ...interface{}) (url.ShortUrl, error) {
	var short url.ShortUrl
	var sql string
	var err error

	err = y.ensureLongId(long, s)
	if err != nil {
		return short, err
	}

	/** TODO
	 * - validate alias is available
	 */
	short = url.ShortUrl{
		Value:       long.Value,
		ShortFor:    *long,
		Domain:      long.Domain(),
		Strategy:    y.Name(),
		SubmittedBy: storage.CurrentUser(),
		SubmittedOn: time.Now().UTC().Unix(),
		IsEnabled:   true,
	}

	sql = `
INSERT INTO shorturl(url, longurl, domain, strategy, submittedBy, submittedInfo, isEnabled)
VALUES (?, ?, ?, ?, ?, ?, ?)
		`
	result, err := s.DB.Exec(sql, short.Value, short.ShortFor.Id, short.Domain, short.Strategy, short.SubmittedBy.Id, short.SubmittedOn, short.IsEnabled)
	if err != nil {
		log.Printf("Failed inserting short %s: %#v", short.Value, err)
	} else {
		short.Id, _ = result.LastInsertId()
	}

	return short, err
}

/**
Any nonzero result is likely an error.
*/
func (y baseStrategy) UseCount(s storage.Storage) int {
	sql := `
SELECT COUNT(*)
FROM shorturl
WHERE strategy = ?
	`
	var count int
	err := s.DB.QueryRow(sql, y.Name()).Scan(&count)
	if err != nil {
		count = 0
		log.Printf("Failed querying database for base strategy use count: %v\n", err)
	}

	return count
}
