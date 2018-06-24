# sciencebase
--
    import "github.com/datatogether/linked_data/sciencebase"

sciencebase is a package for working with USGS science data.

Details from the sciencebase site:

    ScienceBase provides a digital data repository for scientific data assets
    of many different types.
    This is an evolving capability of ScienceBase, driven by the needs of individual
    communities who use ScienceBase as their data management platform.
    Certain types of data such as shapefiles, GeoTIFF images, and a few others are
    able to be served from the Repository using appropriate types of web service
    technology (e.g., OGC-WMS, OGC-WCS, etc.) for streaming-type uses in addition to
    original file downloads.

more info:
https://my.usgs.gov/confluence/display/sciencebase/ScienceBase+Services+and+Models

## Usage

#### type Catalog

```go
type Catalog struct {
	Total    int     `json:"total"`
	Took     string  `json:"took"`
	Selflink *Link   `json:"selflink"`
	Items    []*Item `json:"items"`
}
```


#### type Date

```go
type Date struct {
	Type       string `json:"type"`
	DateString string `json:"dateString"`
	Label      string `json:"label"`
}
```


#### type DistLink

```go
type DistLink struct {
	Uri       string      `json:"uri"`
	Title     string      `json:"title"`
	Type      string      `json:"type"`
	TypeLabel string      `json:"typeLabel"`
	Rel       string      `json:"rel"`
	Name      string      `json:"name"`
	Files     interface{} `json:"files"`
}
```


#### type File

```go
type File struct {
	Name             string `json:"name"`
	Title            string `json:"title"`
	ContentType      string `json:"contentType"`
	ContentEncoding  string `json:"contentEncoding"`
	PathOnDisk       string `json:"pathOnDisk"`
	Processed        bool   `json:"processed"`
	ProcessToken     string `json:"processToken"`
	ImageWidth       int    `json:"imageWidth"`
	ImageHeight      int    `json:"imageHeight"`
	Size             int64  `json:"size"`
	DateUploaded     string `json:"dateUploaded"`
	OriginalMetadata bool   `json:"originalMetadata"`
	UseForPreview    bool   `json:"useForPreview"`
	S3Object         string `json:"s3Object"`
	// screw you poorly specified values
	Checksum interface{} `json:"checksum"`
	Url      string      `json:"url"`
}
```


#### type Item

```go
type Item struct {
	Link              *Link                    `json:"link,omitempty"`
	RelatedItems      map[string]*Link         `json:"relatedItems,omitempty"`
	Id                string                   `json:"id,omitempty"`
	Title             string                   `json:"title,omitempty"`
	AlternateTitles   []string                 `json:"alternateTitles,omitempty"`
	Summary           string                   `json:"summary,omitempty"`
	Body              string                   `json:"body,omitempty"`
	Citation          string                   `json:"citation,omitempty"`
	Purpose           string                   `json:"purpose,omitempty"`
	Provenance        map[string]interface{}   `json:"provenance,omitempty"`
	HasChildren       bool                     `json:"hasChildren,omitempty"`
	ParentId          string                   `json:"parentId,omitempty"`
	Contacts          []map[string]interface{} `json:"contacts,omitempty"`
	WebLinks          []interface{}            `json:"webLinks,omitempty"`
	BrowseCategories  []interface{}            `json:"browseCategories,omitempty"`
	BrowseTypes       []interface{}            `json:"browseTypes,omitempty"`
	SystemTypes       []interface{}            `json:"systemTypes,omitempty"`
	Tags              []*Tag                   `json:"tags,omitempty"`
	Dates             []*Date                  `json:"dates,omitempty"`
	Spatial           map[string]interface{}   `json:"spatial,omitempty"`
	Facets            []interface{}            `json:"facets,omitempty"`
	Files             []*File                  `json:"files,omitempty"`
	DistributionLinks []*DistLink              `json:"distributionLinks,omitempty"`
	PreviewImage      map[string]interface{}   `json:"previewImage,omitempty"`
}
```

A ScienceBase Item is the basic unit of data within the ScienceBase Catalog. An
Item can be anything from a metadata description of a data service to a dataset
in the Digital Repository to a description of a project. Items share a common
set of simple metadata elements based on the Dublin Core Metadata Element Set
and may be extended with additional properties.

#### func (*Item) ChildrenJsonUrl

```go
func (i *Item) ChildrenJsonUrl() string
```

#### type Link

```go
type Link struct {
	Rel string `json:"rel"`
	Url string `json:"url"`
}
```


#### func (*Link) JsonUrl

```go
func (l *Link) JsonUrl() string
```

#### type Provenance

```go
type Provenance struct {
	Annotation  string `json:"annotation,omitempty"`
	DataSource  string `json:"dataSource,omitempty"`
	DateCreated string `json:"dateCreated,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}
```


#### type Tag

```go
type Tag struct {
	Type   string `json:"type"`
	Scheme string `json:"scheme"`
	Name   string `json:"name"`
}
```
