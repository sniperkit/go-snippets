package xmp

import (
	"encoding/json"
	"time"
)

// POD is a high-level interpretation of the Project Open Data
// data.json spec:
// https://project-open-data.cio.gov/v1.1/schema/
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

func (p *POD) AsObject() (map[string]interface{}, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	obj := map[string]interface{}{}
	err = json.Unmarshal(data, &obj)
	return obj, err
}

// AsPOD turns an XMPPacket into a Project Open Data struct
func (p *XMPPacket) AsPOD() *POD {
	return &POD{
		Title:       p.RDF.Description.Title.Default(),
		Description: p.RDF.Description.Description.Default(),
		Keyword:     p.RDF.Description.Subject.Default(),
		Modified:    p.RDF.Description.ModifyDate,
		Publisher:   p.RDF.Description.Creator.String(),
		// TODO
	}
}

// MarshalPODJSON renders XMP metadata as Project Open Data Metadata
func (p *XMPPacket) MarshalPODJSON() ([]byte, error) {
	return json.Marshal(p.AsPOD())
}
