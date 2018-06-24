// dcat implements the W3C Data Catalog Vocabulary (DCAT) specification: https://www.w3.org/TR/vocab-dcat/
// Copious notes within the package are pulled directly from the spec, and are
// intented to help get acquainted with the spec.
//
// JSON is the principle encoding/decoding format.
//
// From the spec:
// DCAT is an RDF vocabulary designed to facilitate interoperability between data catalogs published on the Web.
// This document defines the schema and provides examples for its use.
// By using DCAT to describe datasets in data catalogs, publishers increase discoverability and enable applications easily to consume metadata from multiple catalogs. It further enables decentralized publishing of catalogs and facilitates federated dataset search across sites. Aggregated DCAT metadata can serve as a manifest file to facilitate digital preservation.
package dcat

import (
	"github.com/datatogether/linked_data"
	"time"
)

// Namespaces is a map of prefixes & descriptive urls.
// The namespace for DCAT is http://www.w3.org/ns/dcat#. However, it should be noted that DCAT makes
// extensive use of terms from other vocabularies, in particular Dublin Core.
// DCAT itself defines a minimal set of classes and properties of its own.
var Namespaces = map[string]string{
	"dcat":   "http://www.w3.org/ns/dcat#",
	"dct":    "http://purl.org/dc/terms/",
	"dctype": "http://purl.org/dc/dcmitype/",
	"foaf":   "http://xmlns.com/foaf/0.1/",
	"rdf":    "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"rdfs":   "http://www.w3.org/2000/01/rdf-schema#",
	"skos":   "http://www.w3.org/2004/02/skos/core#",
	"vcard":  "http://www.w3.org/2006/vcard/ns#",
	"xsd":    "http://www.w3.org/2001/XMLSchema#",
}

// Catalog is A data catalog is a curated collection of metadata about datasets.
// Typically, a web-based data catalog is represented as a single instance of this class.
type Catalog struct {
	// A name given to the catalog.
	Title string `json:"dct:title,omitempty"`
	// A free-text account of the catalog.
	Description string `json:"dct:description,omitempty"`
	// Date of formal issuance (e.g., publication) of the catalog.
	Issued time.Time `json:"dct:issued,omitempty"`
	// Most recent date on which the catalog was changed, updated or modified.
	Modified time.Time `json:"dct:modified,omitempty"`
	// The language of the catalog. This refers to the language used in the textual
	// metadata describing titles, descriptions, etc. of the datasets in the catalog.
	// Usage Note: Multiple values can be used. The publisher might also choose to describe the
	//             language on the dataset level (see dataset language).
	// TODO - make a dct:LinguisticSystem instance
	Language []string `json:"dct:language,omitempty"`
	// The homepage of the catalog.
	// Usage Note: foaf:homepage is an inverse functional property (IFP) which means that it
	//             should be unique and precisely identify the catalog. This allows smushing
	//             various descriptions of the catalog when different URIs are used.
	Homepage string `json:"foaf:homepage,omitempty"`
	// The entity responsible for making the catalog online.
	// Usage Note: Resources of type foaf:Agent are recommended as values for this property.
	// TODO - make an instance of Organization / Person
	Publisher string `json:"dct:publisher,omitempty"`
	// The geographical area covered by the catalog.
	// TODO - should be a dct:Location instance
	Spatial string `json:"dct:spatial,omitempty"`
	// The knowledge organization system (KOS) used to classify catalog's datasets.
	// TODO - should be a dcat:Catalog
	Themes []string `json:"dcat:themeTaxonomy,omitempty"`
	// This links to the license document under which the catalog is made available and *not the datasets*.
	// Even if the license of the catalog applies to all of its datasets and distributions, it should be replicated on each distribution.
	License string `json:"dct:license,omitempty"`
	// This describes the rights under which the catalog can be used/reused and *not the datasets*.
	// Even if theses rights apply to all the catalog datasets and distributions, it should be replicated on each distribution.
	Rights string `json:"dct:rights,omitempty"`
	// A dataset that is part of the catalog.
	Dataset *Dataset `json:"dcat:dataset,omitempty"`
	// A catalog record that is part of the catalog.
	Record *CatalogRecord `json:"dcat:record,omitempty"`
}

func (c Catalog) Class() string {
	return "dcat:Catalog"
}

func (c *Catalog) GetDatasets() []ld.Dataset {
	return []ld.Dataset{c.Dataset}
}

// A record in a data catalog, describing a single dataset.
// Usage Note: This class is optional and not all catalogs will use it.
// It exists for catalogs where a distinction is made between metadata about a dataset and metadata
// about the dataset's entry in the catalog.
// For example, the publication date property of the dataset reflects the date when the information was
// originally made available by the publishing agency, while the publication date of the catalog record
// is the date when the dataset was added to the catalog. In cases where both dates differ, or where only
// the latter is known, the publication date should only be specified for the catalog record.
// Notice that the W3C PROV Ontology [prov-o] allows describing further provenance information such as the
// details of the process and the agent involved in a particular change to a dataset.
type CatalogRecord struct {
	// A name given to the record.
	Title string `json:"dct:title"`
	// free-text account of the record.
	Description string `json:"dct:description,omitempty"`
	// The date of listing the corresponding dataset in the catalog.
	Issued time.Time `json:"dct:issued,omitempty"`
	// Most recent date on which the catalog entry was changed, updated or modified.
	Modified time.Time `json:"dct:modified,omitempty"`
	// Links the catalog record to the dcat:Dataset resource described in the record.
	PrimaryTopic string `json:"foaf:primaryTopic,omitempty"`
}

func (c CatalogRecord) Class() string {
	return "dcat:CatalogRecord"
}

// A collection of data, published or curated by a single agent,
// and available for access or download in one or more formats.
// Usage Note: This class represents the actual dataset as published by the dataset publisher.
// In cases where a distinction between the actual dataset and its entry in the catalog is necessary
// (because metadata such as modification date and maintainer might differ),
// the catalog record class can be used for the latter.
type Dataset struct {
	// A name given to the dataset.
	Title string `json:"dct:title"`
	// free-text account of the dataset.
	Description string `json:"dct:description,omitempty"`
	// Date of formal issuance (e.g., publication) of the dataset.
	Issued time.Time `json:"dct:issued,omitempty"`
	// Most recent date on which the dataset was changed, updated or modified.
	Modified time.Time `json:"dct:modified,omitempty"`
	// The language of the dataset.
	// Usage Notes:
	//   * This overrides the value of the catalog language in case of conflict.
	//   * If the dataset is available in multiple languages, use multiple values for this property.
	//     If each language is available separately, define an instance of dcat:Distribution for each language and
	//     describe the specific language of each distribution using dct:language (i.e. the dataset will have multiple
	//     dct:language values and each distribution will have one of these languages as value of its dct:language property).
	// TODO - make a formal Language instance
	Language []string `json:"dct:language,omitempty"`
	// An entity responsible for making the dataset available.
	// Usage Note: Resources of type foaf:Agent are recommended as values for this property.
	Publisher string `json:"dct:publisher,omitempty"`
	// The frequency at which dataset is published.
	// TODO - make this a dct:Frequency instance
	Frequency string `json:"dct:accrualPeriodicity,omitempty"`
	// A unique identifier of the dataset.
	// Usage Note: The identifier might be used as part of the URI of the dataset, but still having it represented explicitly is useful.
	Identifier string `json:"dct:identifier,omitempty"`
	// Spatial coverage of the dataset.
	// TODO - dct:Location (A spatial region or named place)
	Spatial string `json:"dct:spatial,omitempty"`
	// The temporal period that the dataset covers.
	// TODO - make this a dct:PeriodOfTime (An interval of time that is named or defined by its start and end dates)
	Temporal string `json:"dct:temporal,omitempty"`
	// The main category of the dataset. A dataset can have multiple themes.
	// Usage Note: The set of skos:Concepts used to categorize the datasets are
	// organized in a skos:ConceptScheme describing all the categories and their relations in the catalog.
	Theme []string `json:"dcat:theme,omitempty"`
	// A keyword or tag describing the dataset.
	Keyword []string `json:"dcat:keyword,omitempty"`
	// Link a dataset to relevant contact information which is provided using VCard [vcard-rdf].
	// TODO - make this a vcard contact instance
	ContactPoint string `json:"dcat:contactPoint,omitempty"`
	// A Web page that can be navigated to in a Web browser to gain access to the dataset, its distributions and/or additional information.
	// Usage Note: If the distribution(s) are accessible only through a landing page (i.e. direct download URLs are not known),
	//             then the landing page link should be duplicated as accessURL on a distribution.
	LandingPage string `json:"dcat:landingPage,omitempty"`
}

func (d Dataset) Class() string {
	return "dcat:Dataset"
}

func (d *Dataset) GetId() string {
	return d.Identifier
}

func (d *Dataset) GetDistributions() []ld.Distribution {
	// TODO
	return nil
}

// Represents a specific available form of a dataset. Each dataset might be available in different forms,
// these forms might represent different formats of the dataset or different endpoints.
// Examples of distributions include a downloadable CSV file, an API or an RSS feed
// Usage Note: This represents a general availability of a dataset it implies no information about the actual
//             access method of the data, i.e. whether it is a direct download, API, or some through Web page.
//             The use of dcat:downloadURL property indicates directly downloadable distributions.
type Distribution struct {
	// A name given to the distribution.
	Title string `json:"dct:title"`
	// free-text account of the distribution.
	Description string `json:"dct:description,omitempty"`
	// Date of formal issuance (e.g., publication) of the distribution.
	Issued time.Time `json:"dct:issued,omitempty"`
	// Most recent date on which the distribution was changed, updated or modified.
	Modified time.Time `json:"dct:modified,omitempty"`
	// This links to the license document under which the distribution is made available.
	License string `json:"dct:license,omitempty"`
	// Information about rights held in and over the distribution.
	// Usage Note: dct:license, which is a sub-property of dct:rights, can be used to link a distribution
	//             to a license document. However, dct:rights allows linking to a rights statement that can
	//             include licensing information as well as other information that supplements the licence
	//             such as attribution.
	Rights string `json:"dct:rights,omitempty"`
	// A landing page, feed, SPARQL endpoint or other type of resource that gives access to the distribution of the dataset
	// Usage Notes:
	//    * Use accessURL, and not downloadURL, when it is definitely not a download or when you are not sure whether it is.
	//    * If the distribution(s) are accessible only through a landing page (i.e. direct download URLs are not known),
	//      then the landing page link should be duplicated as accessURL on a distribution.
	AccessUrl string `json:"dcat:accessURL,omitempty"`
	// A file that contains the distribution of the dataset in a given format
	// Usage Notes: dcat:downloadURL is a specific form of dcat:accessURL. Nevertheless, DCAT does not define dcat:downloadURL
	//              as a subproperty of dcat:accessURL not to enforce this entailment as DCAT profiles may wish to impose a
	//              stronger separation where they only use accessURL for non-download locations.
	DownloadUrl string `json:"dcat:downloadURL,omitempty"`
	// The size of a distribution in bytes.
	// Usage Note: The size in bytes can be approximated when the precise size is not known.
	ByteSize int64 `json:"dcat:byteSize,omitempty"`
	// The media type of the distribution as defined by IANA.
	// Usage Note: This property should be used when the media type of the distribution is defined in IANA,
	//             otherwise dct:format may be used with different values.
	MediaType string `json:"dcat:mediaType,omitempty"`
	// The file format of the distribution.
	// Usage Note: dcat:mediaType should be used if the type of the distribution is defined by IANA.
	Format string `json:"dct:format,omitempty"`
}

func (d Distribution) Class() string {
	return "dcat:Distribution"
}

// The knowledge organization system (KOS) used to represent themes/categories of datasets in the catalog.
type ConceptScheme struct {
	// TODO - implement skos package & come back here
}

func (c ConceptScheme) Class() string {
	return "skos:ConceptScheme"
}

// A category or a theme used to describe datasets in the catalog.
// Usage Note: It is recommended to use either skos:inScheme or skos:topConceptOf on every skos:Concept
//             used to classify datasets to link it to the concept scheme it belongs to.
//             This concept scheme is typically associated with the catalog using dcat:themeTaxonomy
type Concept struct {
	// TODO - implement skos package & come back here
}

func (c Concept) Class() string {
	return "skos:Concept"
}

// Usage Note: FOAF provides sufficient properties to describe these entities.
type OrganizationPerson struct {
	// TODO - implement foaf package & come back here
}
