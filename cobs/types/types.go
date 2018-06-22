package types

type ImageInfo struct {
	OriginalName string `json:"original_name"`
	NewName      string `json:"new_name"`
	Architecture string `json:"architecture"`
	Tag          string `json:"tag"`
}
