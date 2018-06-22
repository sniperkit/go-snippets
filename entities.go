package entities

const (
	ObjectTypeTreeMimeType string = "clawio/tree"
)
const (
	ObjectTypeTree ObjectType = iota
	ObjectTypeBLOB
)

type (
	// ObjectType indicates if the object is either a Tree or a BLOB.
	ObjectType int

	// ObjectInfo represents the metadata information retrieved
	// from an object, either tree or blob.
	ObjectInfo struct {
		PathSpec string     `json:"pathspec"`
		Size     int64      `json:"size"`
		Type     ObjectType `json:"type"`
		MimeType string     `json:"mime_type"`
		Checksum string     `json:"checksum"`
	}

	// User represents an user of the system.
	// They are created by the authentication service.
	User struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
)
