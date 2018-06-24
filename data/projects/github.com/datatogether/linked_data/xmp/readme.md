# xmp
--
    import "github.com/archivers-space/xmp"

XMP is a package for parsing Extensible Metadata Platform documents. This
package includes lots of comments to help make sense of XMP for the purposes of
metadata extraction & conversion to other metadata formats.

from: https://en.wikipedia.org/wiki/Extensible_Metadata_Platform The Extensible
Metadata Platform (XMP) is an ISO standard (ISO 16684-1), originally created by
Adobe Systems Inc., for the creation, processing and interchange of standardized
and custom metadata for digital documents and data sets. XMP standardizes a data
model, a serialization format and core properties for the definition and
processing of extensible metadata. It also provides guidelines for embedding XMP
information into popular image, video and document file formats, such as JPEG
and PDF, without breaking their readability by applications that do not support
XMP. Therefore, the non-XMP metadata have to be reconciled with the XMP
properties.

## Usage

#### type DublinCore

```go
type DublinCore struct {
	// An entity responsible for making contributions to the resource
	// Examples of a contributor include a person, an organization, or a service.
	// Typically, the name of a contributor should be used to indicate the entity.
	// XMP usage is a list of contributors. These contributors should not include those listed in dc:creator.
	Contributor RDFSeq `xml:"contributor"`
	// The spatial or temporal topic of the resource, the spatial applicability of the resource,
	// or the jurisdiction under which the resource is relevant.
	// XMP usage is the extent or scope of the resource.
	Coverage string `xml:"coverage"`
	// An entity primarily responsible for making the resource.
	// Examples of a creator include a person, an organization, or a
	// service. Typically, the name of a creator should be used to indicate the entity.
	// XMP usage is a list of creators. Entities should be listed in order of decreasing precedence,
	// if such order is significant.
	Creator RDFSeq `xml:"creator"`
	// A point or period of time associated with an event in the life cycle of the resource.
	Date *time.Time `xml:"date"`
	// An account of the resource.
	// XMP usage is a list of textual descriptions of the content of the resource, given in various languages.
	Description RDFAlt `xml:"description"`
	// XMP usage is a MIME type.
	Format string `xml:"format"`
	// An unambiguous reference to the resource within a given context.
	Identifier string `xml:"identifier"`
	// A language of the resource.
	// XMP usage is a list of languages used in the content of the resource.
	// TODO - RDFSeq is a guess
	Language RDFSeq `xml:"language"`
	// An entity responsible for making the resource available
	// Examples of a publisher include a person, an organization, or a
	// service. Typically, the name of a publisher should be used to indicate the entity.
	//  XMP usage is a list of publishers.
	// TODO - RDFSeq is a guess
	Publisher RDFSeq `xml:"publisher"`
	// A related resource.
	// Recommended best practice is to identify the related resource
	// by means of a string conforming to a formal identification system.
	// XMP usage is a list of related resources.
	// TODO - RDFSeq is a guess
	Relation RDFSeq `xml:"relation"`
	// Information about rights held in and over the resource.
	// ypically, rights information includes a statement about various property
	// rights associated with the resource, including intellectual property rights.
	// XMP usage is a list of informal rights statements, given in various languages.
	// TODO - RDFAlt is a guess
	Rights RDFAlt `xml:"rights"`
	// A related resource from which the described resource is derived.
	// The described resource may be derived from the related resource in whole or in part.
	// Recommended best practice is to identify the related resource by means of a string
	// conforming to a formal identification system.
	Source string `xml:"source"`
	// The topic of the resource.
	// Typically, the subject will be represented using keywords, key phrases, or
	// classification codes. Recommended best practice is to use a controlled vocabulary.
	// To describe the spatial or temporal topic of the resource, use the dc:coverage element.
	// XMP usage is a list of descriptive phrases or keywords that specify the content of the resource.
	Subject RDFBag `xml:"subject"`
	// A name given to the resource.
	// Typically, a title will be a name by which the resource is formally known.
	// XMP usage is a title or name, given in various languages.
	Title RDFAlt `xml:"title"`
	// The nature or genre of the resource.
	// Recommended best practice is to use a controlled vocabulary such as the DCMI Type
	// Vocabulary [DCMITYPE]. To describe the file format, physical medium, or dimensions of the
	// resource, use the dc:format element.
	// See the dc:format entry for clarification of the XMP usage of that element.
	Type RDFBag `xml:"type"`
}
```

DublinCore is the "dc:" namespace often seen in xmp meta. It's pretty sweet.
https://en.wikipedia.org/wiki/Dublin_Core http://dublincore.org For the XMP
flavour, see XMP section 8.3

#### type POD

```go
type POD struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Keyword      []string   `json:"keyword"`
	Created      *time.Time `json:"created"`
	Modified     *time.Time `json:"modified"`
	Publisher    string     `json:"publisher"`
	ContactPoint string     `json:"contactPoint"`
	Identifier   string     `json:"identifier"`
	AccessLevel  string     `json:"accessLevel"`
	BureauCode   string     `json:"bureauCode"`
	ProgramCode  string     `json:"programCode"`
	License      string     `json:"license"`
	Rights       string     `json:"rights"`
}
```

POD is a high-level interpretation of the Project Open Data data.json spec:
https://project-open-data.cio.gov/v1.1/schema/

#### type RDFAlt

```go
type RDFAlt struct {
	RDFList RDFList `xml:"Alt"`
}
```

RDF Langauge Alternative

#### func (RDFAlt) Default

```go
func (d RDFAlt) Default() string
```

#### type RDFBag

```go
type RDFBag struct {
	Items RDFList `xml:"Bag"`
}
```


#### func (RDFBag) Default

```go
func (b RDFBag) Default() []string
```

#### type RDFList

```go
type RDFList struct {
	Items []string `xml:"li"`
}
```


#### func (RDFList) Default

```go
func (a RDFList) Default() string
```

#### type RDFSeq

```go
type RDFSeq struct {
	Items RDFList `xml:"Seq"`
}
```


#### func (RDFSeq) Default

```go
func (seq RDFSeq) Default() []string
```

#### func (RDFSeq) String

```go
func (seq RDFSeq) String() string
```

#### type XMPBasic

```go
type XMPBasic struct {
	// The date and time the resource was created. For a digital file, this need not match a
	// file-system  creation time. For a freshly created resource, it should be close to that time,
	// modulo the time taken to write the file. Later file transfer, copying, and so on, can make the
	// file-system time arbitrarily different.
	CreateDate *time.Time `xml:"CreateDate"`
	// The name of the first known tool used to create the resource.
	// TODO - the spec says this should be an "AgentName"
	CreatorTool string `xml:"CreatorTool"`
	// An unordered array of text strings that unambiguously identify the resource within a
	// given context. An array item may be qualified with xmpidq:Scheme (see 8.7, “xmpidq namespace”)
	// to denote the formal identification system to which that identifier conforms.
	// NOTE: The xmp:Identifier property was added because dc:identifier has been defined in
	// the original XMP specification as a single identifier instead of as an array, and changing
	// dc:identifier to an array would break compatibility with existing XMP processors.
	Identifier RDFBag `xml:"Identifier"`
	// A word or short phrase that identifies a resource as a member of a user- defined collection.
	// NOTE: One anticipated usage is to organize resources in a file browser.
	Label string `xml:"Label"`
	// The date and time that any metadata for this resource was last changed.
	// It should be the same as or more recent than xmp:ModifyDate.
	MetadataDate *time.Time `xml:"MetadataDate"`
	// The date and time the resource was last modified.
	// NOTE: The value of this property is not necessarily the same as the file’s
	// system modification date because it is typically set before the file is saved.
	ModifyDate *time.Time `xml:"ModifyDate"`
	// A user-assigned rating for this file. The value shall be -1 or in the range [0..5],
	// where -1 indicates “rejected” and 0 indicates “unrated”. If xmp:Rating is not present,
	// a value of 0 should be assumed.
	// NOTE: Anticipated usage is for a typical “star rating” UI, with the addition of a notion of rejection.
	Rating float32 `xml:"Rating"`
}
```

The XMP basic namespace contains properties that provide basic descriptive
information. XMP spec Section 8.4

#### type XMPDescription

```go
type XMPDescription struct {
	DublinCore
	XMPBasic
	XMPRights
	XMPMM
	XMPIdq
}
```


#### type XMPIdq

```go
type XMPIdq struct {
	// A qualifier providing the name of the formal identification scheme
	// used for an item in the xmp:Identifier array.
	Scheme string
}
```

The xmpidq namespace contains a single qualifier that defines the scheme used in
the xmp:Identifier array. XMP Section 8.7

#### type XMPMM

```go
type XMPMM struct {
	// A reference to the resource from which this one is derived.
	// This should be a minimal reference, in which missing components can be
	// assumed to be unchanged.
	// NOTE A rendition might need to specify only the
	// xmpMM:InstanceID and xmpMM:RenditionClass of the original.
	// TODO - this is actually a "ResourceRef" type
	DerivedFrom string
	// The common identifier for all versions and renditions of a resource.
	// TODO - this is actually a GUID type
	DocumentId string
	// An identifier for a specific incarnation of a resource,
	// updated each time a file is saved.
	// TODO - this is actually a GUID type
	InstanceId string
	// The common identifier for the original resource from which the current
	// resource is derived. For example, if you save a resource to a different format,
	// then save that one to another format, each save operation should generate a new
	// xmpMM:DocumentID that uniquely identifies the resource in that format,
	// but should retain the ID of the source file here.
	// TODO - this is actually a GUID type
	OriginalDocumentId string
	// The rendition class name for this resource. This property should be absent or
	// set to 'default' for a resource that is not a derived rendition.
	// See definitions of rendition (3.7) and version (3.9).
	// TODO - this is actually a RenditionClass type
	RenditionClass string
	// Can be used to provide additional rendition parameters that are
	// too complex or verbose to encode in xmpMM:RenditionClass.
	RenditionParams string
}
```

The XMP Media Management namespace contains properties that provide information
regarding the identification, composition, and history of a resource. XMP spec
Section 8.6

#### type XMPPacket

```go
type XMPPacket struct {
	RDF XMPRDFMeta `xml:"RDF"`
}
```

XMPPacket is an Instance of the XMP Data Model. Also known as an XMP document.

#### func  Unmarshal

```go
func Unmarshal(data []byte) (*XMPPacket, error)
```
Unmarshal parses a stream of bytes into an XMPPacket

#### func (*XMPPacket) AsPOD

```go
func (p *XMPPacket) AsPOD() *POD
```
AsPOD turns an XMPPacket into a Project Open Data struct

#### func (*XMPPacket) MarshalPODJSON

```go
func (p *XMPPacket) MarshalPODJSON() ([]byte, error)
```
MarshalPODJSON renders XMP metadata as Project Open Data Metadata

#### type XMPRDFMeta

```go
type XMPRDFMeta struct {
	Description XMPDescription `xml:"Description"`
}
```

XMP is itself a subset of the Resource Description Framework. RDF is an OG
(Original Gangter) web (W3) standard for data interchange.
https://www.w3.org/TR/rdf11-concepts/ All XMPPackets are a single instance of an
RDF document, with limitations defined by the XMP spec TODO - consider giving
RDF it's own package

#### type XMPRights

```go
type XMPRights struct {
	// A Web URL for a rights management certificate.
	// NOTE: This is a normal (non-URI) simple value because of historical usage.
	Certificate string `xml:"Certificate"`
	// When true, indicates that this is a rights-managed resource.
	// When false, indicates that this is a public-domain resource.
	// Omit if the state is unknown.
	Marked bool `xml:"Marked"`
	// A list of legal owners of the resource.
	Owner RDFBag `xml:"Owner"`
	// A collection of text instructions on how a resource can be legally used,
	// given in a variety of languages.
	UsageTerms RDFAlt `xml:"UsageTerms"`
	// A Web URL for a statement of the ownership and usage rights for this resource.
	// NOTE: This is a normal (non-URI) simple value because of historical usage.
	WebStatement string
}
```

These XMP properties are intended to provide a means of rights expression. They
are not intended to provide digital rights management (DRM) controls. XMP spec
Section 8.5
