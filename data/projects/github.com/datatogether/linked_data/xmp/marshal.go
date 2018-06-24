package xmp

import (
	"encoding/json"
)

// POD is a high-level interpretation of the Project Open Data
// data.json spec:
// https://project-open-data.cio.gov/v1.1/schema/
// type POD struct {
// 	Title        string    `json:"title"`
// 	Description  string    `json:"description"`
// 	Keyword      []string  `json:"keyword"`
// 	Created      time.Time `json:"created"`
// 	Modified     time.Time `json:"modified"`
// 	Publisher    string    `json:"publisher"`
// 	ContactPoint string    `json:"contactPoint"`
// 	Identifier   string    `json:"identifier"`
// 	AccessLevel  string    `json:"accessLevel"`
// 	BureauCode   string    `json:"bureauCode"`
// 	ProgramCode  string    `json:"programCode"`
// 	License      string    `json:"license"`
// 	Rights       string    `json:"rights"`
// }

// // AsPOD turns an XMPPacket into a Project Open Data struct
// func (p *XMPPacket) AsPOD() *POD {
// 	return &POD{
// 		Title:       p.RDF.Description.Title.Default(),
// 		Description: p.RDF.Description.Description.Default(),
// 		Keyword:     p.RDF.Description.Subject.Default(),
// 		Modified:    p.RDF.Description.ModifyDate,
// 		Publisher:   p.RDF.Description.Creator.String(),
// 		// TODO
// 	}
// }

// AsPODObject returns XMPPacket as a map[string]interface that
func (p *XMPPacket) AsPODObject() map[string]interface{} {
	d := p.RDF.Description
	data := map[string]interface{}{}

	if d.Title.Default() != "" {
		data["title"] = d.Title.Default()
	}
	if d.Description.Default() != "" {
		data["description"] = d.Description.Default()
	}
	if len(d.Subject.Default()) > 0 {
		data["keyword"] = d.Subject.Default()
	}
	if d.ModifyDate != "" {
		if date, err := unmarshalXmpDate(d.ModifyDate); err == nil {
			data["modified"] = date
		}
	}
	if d.CreateDate != "" {
		if date, err := unmarshalXmpDate(d.CreateDate); err == nil {
			data["created"] = date
		}
	}
	if len(d.Publisher.Default()) > 0 {
		data["publisher"] = d.Publisher.Default()
	}
	if len(d.Language.Default()) > 0 {
		data["language"] = d.Language.DefaultString()
	}

	// data["contactPoint"] = p.RDF.Description
	// data["identifier"] = p.RDF.Description.DublinCore.Identifier
	// data["accessLevel"] = p.RDF.Description
	// data["bureauCode"] = p.RDF.Description
	// data["programCode"] = p.RDF.Description
	// data["license"] = p.RDF.Description.License.Default()

	if d.Rights.Default() != "" {
		data["rights"] = p.RDF.Description.Rights.Default()
	}

	return data
}

// MarshalPODJSON renders XMP metadata as Project Open Data Metadata
func (p *XMPPacket) MarshalPODJSON() ([]byte, error) {
	return json.Marshal(p.AsPODObject())
}
