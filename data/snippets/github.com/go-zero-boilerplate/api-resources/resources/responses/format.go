package responses

//Format is the enum for the representation type, currently we have JSON and XML
type Format int

const (
	FormatNone Format = iota

	//FormatJSON is the JSON format
	FormatJSON

	//FormatXML is the XML format
	FormatXML
)
