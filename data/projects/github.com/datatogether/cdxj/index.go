package cdxj

import (
	"github.com/datatogether/warc"
)

// Index is a list of cdxj records. above a standard slice of record pointers,
// it includes methods for common patterns for adding and removing records
type Index []*Record

// AddWARCRecords adds a list of warc records to the index, creating
// cdxj records for each WARC record
func (index Index) AddWARCRecords(recs warc.Records) (Index, error) {
	var err error
	for _, rec := range recs {
		index, err = index.AddWARCRecord(rec)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// AddWARCRecord adds creates a cdxj record from a WARC record and adds it to the index
func (index Index) AddWARCRecord(rec *warc.Record) (Index, error) {
	cdxjRec, err := CreateRecord(rec)
	if err != nil {
		return index, err
	}
	return index.AddRecord(cdxjRec), nil
}

// AddRecord adds a single record to the list, unless it's exact
// URI is already present
func (index Index) AddRecord(rec *Record) Index {
	for _, r := range index {
		if r.URI == rec.URI {
			return index
		}
	}
	return append(index, rec)
}
