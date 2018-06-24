package jsonld

import (
	ld "github.com/datatogether/linked_data"
)

// UnmarshalCatalog returns a linked data catalog from a byte slice
func UnmarshalCatalog(data []byte) (ld.Catalog, error) {
	detect := struct {
		ConformsTo string
	}{}

	err := json.Unmarshal(data, &detect)
	if err != nil {
		return nil, err
	}

	if pod.SchemaVersions[detect.ConformsTo] != "" {
		catalog := &pod.Catalog{}
		err = json.Unmarshal(data, catalog)
		return catalog, err
	}

	catalog := &dcat.Catalog{}
	err = json.Unmarshal(data, catalog)
	return catalog, err
}
