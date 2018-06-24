// linked_data (imported as "ld") is a packged for working with a number of different linked-data structures
// This package provides general metadata inference functions for working with any of the supported
// metadata formats
// Each metadata specification is broken out into it's own subpackage, users that only need a single
// linked-data spec are encouraged to import needed packages directly.
package ld

// Catalog represents a collection of Datasets
type Catalog interface {
	GetDatasets() []Dataset
}

// Dataset is is a single Dataset
type Dataset interface {
	GetId() string
	GetDistributions() []Distribution
}

// Distribution is a method for downloading a dataset
type Distribution interface {
	GetDownloadUrl() string
}
