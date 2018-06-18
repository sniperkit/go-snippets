package network

type MediaType int

const (
	MEDIATYPE_TEXT_PLAIN                  MediaType = 0
	MEDIATYPE_TEXT_XML                    MediaType = 1
	MEDIATYPE_TEXT_CSV                    MediaType = 2
	MEDIATYPE_TEXT_HTML                   MediaType = 3
	MEDIATYPE_IMAGE_GIF                   MediaType = 21
	MEDIATYPE_IMAGE_JPEG                  MediaType = 22
	MEDIATYPE_IMAGE_PNG                   MediaType = 23
	MEDIATYPE_IMAGE_TIFF                  MediaType = 24
	MEDIATYPE_AUDIO_RAW                   MediaType = 25
	MEDIATYPE_VIDEO_RAW                   MediaType = 26
	MEDIATYPE_APPLICATION_LINK_FORMAT     MediaType = 40
	MEDIATYPE_APPLICATION_XML             MediaType = 41
	MEDIATYPE_APPLICATION_OCTET_STREAM    MediaType = 42
	MEDIATYPE_APPLICATION_RDFXML          MediaType = 43
	MEDIATYPE_APPLICATION_SOAPXML         MediaType = 44
	MEDIATYPE_APPLICATION_ATOMXML         MediaType = 45
	MEDIATYPE_APPLICATION_XMPPXML         MediaType = 46
	MEDIATYPE_APPLICATION_EXI             MediaType = 47
	MEDIATYPE_APPLICATION_FASTINFOSET     MediaType = 48
	MEDIATYPE_APPLICATION_SOAPFASTINFOSET MediaType = 49
	MEDIATYPE_APPLICATION_JSON            MediaType = 50
	MEDIATYPE_APPLICATION_X_OBIT_BINARY   MediaType = 51
	MEDIATYPE_TEXT_PLAIN_VND_OMA_LWM2M    MediaType = 1541
	MEDIATYPE_TLV_VND_OMA_LWM2M           MediaType = 1542
	MEDIATYPE_JSON_VND_OMA_LWM2M          MediaType = 1543
	MEDIATYPE_OPAQUE_VND_OMA_LWM2M        MediaType = 1544
)

const (
	METHOD_GET    = "GET"
	METHOD_PUT    = "PUT"
	METHOD_POST   = "POST"
	METHOD_DELETE = "DELETE"
	METHOD_PATCH  = "PATCH"
)
