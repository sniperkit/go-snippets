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
package pod

import (
	"github.com/datatogether/linked_data"
)

// Map of values for Catalog.ConformsTo
var SchemaVersions = map[string]string{
	"https://project-open-data.cio.gov/v1.1/schema": "1.1",
	"https://project-open-data.cio.gov/v1.0/schema": "1.0",
}

// These fields describe the entire Public Data Listing catalog file.
// Publishers can also use the describedBy field to reference the default JSON Schema file used to define the schema
// (https://project-open-data.cio.gov/v1.1/schema/catalog.json) or they may refer to their own JSON Schema file if
// they have extended the schema with additional schema definitions.
// Similarly, @context can be used to reference the default JSON-LD Context used to define the schema
// (https://project-open-data.cio.gov/v1.1/schema/catalog.jsonld) or publishers can refer to their own if they have
// extended the schema with additional linked data vocabularies. See the Catalog section under
// Further Metadata Field Guidance for more details.
type Catalog struct {
	// Metadata Context  URL or JSON object for the JSON-LD Context that defines the schema used.
	Context string `json:"@context,omitempty" required:"no"`
	// Metadata Catalog ID IRI for the JSON-LD Node Identifier of the Catalog. This should be the URL of the data.json file itself.
	Id string `json:"@id,omitempty" required:"no"`
	// Metadata Type IRI for the JSON-LD data type. This should be dcat:Catalog for the Catalog.
	Type string `json:"@type,omitempty" required:"no"`
	// Schema Version URI that identifies the version of the Project Open Data schema being used.
	ConformsTo string `json:"conformsTo" required:"yes"`
	// Data Dictionary URL for the JSON Schema file that defines the schema used.
	DescribedBy string `json:"describedBy,omitempty" required:"no"`
	// Dataset A container for the array of Dataset objects. See Dataset Fields below for details.
	Dataset []*Dataset `json:"dataset" required:"yes"`
}

func (c *Catalog) GetDatasets() []ld.Dataset {
	ds := make([]ld.Dataset, len(c.Dataset))
	for i, d := range c.Dataset {
		ds[i] = d
	}
	return ds
}

// @TODO - part of a planned marshalling / @type upgrade
// func (c Catalog) Class() string {
// 	return "dcat:Catalog"
// }

// See the Further Metadata Field Guidance (https://project-open-data.cio.gov/v1.1/schema/#further-metadata-field-guidance)
// section to learn more about the use of each element, including the range of valid entries where appropriate.
// Consult the field mappings to find the equivalent v1.0, DCAT, Schema.org, and CKAN fields.
type Dataset struct {
	// Metadata Type IRI for the JSON-LD data type. This should be dcat:Dataset for each Dataset.
	Type string `json:"@type,omitempty" required:"no"`
	// Title Human-readable name of the asset. Should be in plain English and include sufficient detail to facilitate search and discovery.
	Title string `json:"title" required:"yes"`
	// Description Human-readable description (e.g., an abstract) with sufficient detail to enable a user to quickly understand whether the asset is of interest.
	Description string `json:"description" required:"yes"`
	// Tags  Tags (or keywords) help users discover your dataset; please include terms that would be used by technical and non-technical users.
	// Keyword []string `json:"keyword" required:"yes"`
	// TODO - this cause doe data.json so stupid
	Keyword interface{} `json:"keyword" required:"yes"`
	// Last Update Most recent date on which the dataset was changed, updated or modified.
	// @TODO - this is currently a string b/c timestamps aren't stored in proper JSON timestamp format
	// will need custom JSON marshal/unmarshal
	Modified string `json:"modified" required:"yes"`
	// Publisher The publishing entity and optionally their parent organization(s).
	Publisher *Organization `json:"publisher" required:"yes"`
	// Contact Name and Email  Contact person’s name and email for the asset.
	ContactPoint *Contact `json:"contactPoint" required:"yes"`
	// Unique Identifier A unique identifier for the dataset or API as maintained within an Agency catalog or database.
	Identifier string `json:"identifier" required:"yes"`
	// Public Access Level The degree to which this dataset could be made publicly-available, regardless of whether it has been made available. Choices: public (Data asset is or could be made publicly available to all without restrictions), restricted public (Data asset is available under certain use restrictions), or non-public (Data asset is not available to members of the public).
	AccessLevel string `json:"accessLevel" required:"yes"`
	// Bureau Code Federal agencies, combined agency and bureau code from OMB Circular A-11, Appendix C (PDF, CSV) in the format of 015:11.
	BureauCode []string `json:"bureauCode" required:"yest"`
	// Program Code  Federal agencies, list the primary program related to this data asset, from the Federal Program Inventory. Use the format of 015:001.
	ProgramCode []string `json:"programCode" required:"yes"`
	// License The license or non-license (i.e. Public Domain) status with which the dataset or API has been published. See Open Licenses for more information.
	// Required If-Applicable
	License string `json:"license,omitempty" required:"no"`
	// Rights  This may include information regarding access or restrictions based on privacy, security, or other policies. This should also serve as an explanation for the selected “accessLevel” including instructions for how to access a restricted file, if applicable, or explanation for why a “non-public” or “restricted public” data asset is not “public,” if applicable. Text, 255 characters.
	// Required If-Applicable
	Rights string `json:"rights,omitempty" requird:"no"`
	// Spatial The range of spatial applicability of a dataset. Could include a spatial region like a bounding box or a named place.
	// Required If-Applicable
	Spatial string `json:"spatial,omitempty" required:"no"`
	// Temporal  The range of temporal applicability of a dataset (i.e., a start and end date of applicability for the data).
	// Required  If-Applicable
	Temporal string `json:"temporal,omitempty" required:"no"`
	// Distribution  A container for the array of Distribution objects. See Dataset Distribution Fields below for details.
	// Required If-Applicable
	Distribution []*Distribution `json:"distribution,omitempty" required:"no"`
	// Frequency The frequency with which dataset is published.
	AccrualPeriodicity string `json:"accrualPeriodicity,omitempty" required:"no"`
	// Data Standard URI used to identify a standardized specification the dataset conforms to.
	ConformsTo string `json:"conformsTo,omitempty" required:"no"`
	// Data Quality  Whether the dataset meets the agency’s Information Quality Guidelines (true/false)
	DataQuality bool `json:"dataQuality" required:"no"`
	// Data Dictionary URL to the data dictionary for the dataset. Note that documentation other than a data dictionary can be referenced using Related Documents (references).
	DescribedBy string `json:"describedBy,omitempty" required:"no"`
	// Data Dictionary Type  The machine-readable file format (IANA Media Type also known as MIME Type) of the dataset’s Data Dictionary (describedBy).
	DescribedByType string `json:"describedByType,omitempty" required:"no"`
	// Collection  The collection of which the dataset is a subset.
	IsPartOf string `json:"isPartOf,omitempty" required:"no"`
	// Release Date  Date of formal issuance.
	// @TODO - this is currently a string b/c timestamps aren't stored in proper JSON timestamp format
	// will need custom JSON marshal/unmarshal
	Issued string `json:"issued,omitempty" required:"no"`
	// Language  The language of the dataset.
	Language []string `json:"language,omitempty" required:"no"`
	// Homepage URL  This field is not intended for an agency’s homepage (e.g. www.agency.gov), but rather if a dataset has a human-friendly hub or landing page that users can be directed to for all resources tied to the dataset.
	LandingPage string `json:"landingPage,omitempty" required:"no"`
	// Primary IT Investment UII For linking a dataset with an IT Unique Investment Identifier (UII).
	PrimaryITInvestmentUII string `json:"primaryITInvestmentUII,omitempty" required:"no"`
	// Related Documents Related documents such as technical information about a dataset, developer documentation, etc.
	References []string `json:"references,omitempty" required:"no"`
	// System of Records If the system is designated as a system of records under the Privacy Act of 1974, provide the URL to the System of Records Notice related to this dataset.
	SystemOfRecords string `json:"systemOfRecords,omitempty" required:"no"`
	// Category  Main thematic category of the dataset.
	Theme []string `json:"theme,omitempty" required:"no"`
}

func (d *Dataset) GetId() string {
	return d.Identifier
}

func (d *Dataset) GetDistributions() []ld.Distribution {
	dst := make([]ld.Distribution, len(d.Distribution))
	for i, ds := range d.Distribution {
		dst[i] = ds
	}
	return dst
}

// @TODO - part of a planned marshalling / @type upgrade
// func (d Dataset) Class() string {
// 	return "dcat:Dataset"
// }

// Within a dataset, distribution is used to aggregate the metadata specific to a dataset’s resources
// (accessURL and downloadURL), which may be described using the following fields.
// Each distribution should contain one accessURL or downloadURL.
// A downloadURL should always be accompanied by mediaType.
type Distribution struct {
	// Metadata Type IRI for the JSON-LD data type. This should be dcat:Distribution for each Distribution.
	Type string `json:"@type,omitempty" required:"no"`
	// Access URL  URL providing indirect access to a dataset, for example via API or a graphical interface.
	AccessURL string `json:"accessURL,omitempty" required:"no"`
	// Data Standard URI used to identify a standardized specification the distribution conforms to.
	ConformsTo string `json:"conformsTo,omitempty" required:"no"`
	// Data Dictionary URL to the data dictionary for the distribution found at the downloadURL. Note that documentation
	// other than a data dictionary can be referenced using Related Documents as shown in the expanded fields.
	DescribedBy string `json:"describedBy,omitempty" required:"no"`
	// Data Dictionary Type  The machine-readable file format (IANA Media Type or MIME Type) of the distribution’s describedBy URL.
	DescribedByType string `json:"describedByType,omitempty" required:"no"`
	// Description Human-readable description of the distribution.
	Description string `json:"description,omitempty" required:"no"`
	// Download URL  URL providing direct access to a downloadable file of a dataset.
	DownloadURL string `json:"downloadURL,omitempty" required:"no"`
	// Format  A human-readable description of the file format of a distribution.
	Format string `json:"format,omitempty" required:"no"`
	// Media Type  The machine-readable file format (IANA Media Type or MIME Type) of the distribution’s downloadURL.
	MediaType string `json:"mediaType,omitempty" required:"no"`
	// Title Human-readable name of the distribution.
	Title string `json:"title,omitempty" required:"no"`
}

func (d *Distribution) GetDownloadUrl() string {
	return d.DownloadURL
}

// @TODO - part of a planned marshalling / @type upgrade
// func (d Distribution) Class() string {
// 	return "dcat:Distribution"
// }

type Organization struct {
	Type              string        `json:"@type,omitempty" required:"no"`
	Name              string        `json:"name" required:"yes"`
	SubOrganizationOf *Organization `json:"subOrganizationOf,omitempty"`
}

type Contact struct {
	Type     string `json:"@type,omitempty" required:"no"`
	Fn       string `json:"fn" required:"yes"`
	HasEmail string `json:"hasEmail" required:"yes"`
}
