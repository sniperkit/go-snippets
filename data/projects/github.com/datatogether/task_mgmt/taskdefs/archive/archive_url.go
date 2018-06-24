package archive

type ArchiveUrl struct {
	Url              string `json:"url"`              // url to resource to be added
	Checksum         string `json:"checksum"`         // optional checksum to check resp against
	ipfsApiServerUrl string `json:"ipfsApiServerUrl"` // url of IPFS api server
}
