package rest

const (
	charsetUTF8 = "charset=UTF-8"
)

// HTTP methods
const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

// MIME type definitions
const (
	MIMEApplicationJSON                  MIME = "application/json"
	MIMEApplicationJSONCharsetUTF8       MIME = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            MIME = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 MIME = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   MIME = "application/xml"
	MIMEApplicationXMLCharsetUTF8        MIME = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          MIME = "text/xml"
	MIMETextXMLCharsetUTF8               MIME = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  MIME = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              MIME = "application/protobuf"
	MIMEApplicationMsgpack               MIME = "application/msgpack"
	MIMETextHTML                         MIME = "text/html"
	MIMETextHTMLCharsetUTF8              MIME = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        MIME = "text/plain"
	MIMETextPlainCharsetUTF8             MIME = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    MIME = "multipart/form-data"
	MIMEOctetStream                      MIME = "application/octet-stream"
)

// MIME is a string which implements the ContentType interface
type MIME string

// ContentType returns the content-type header for the mime type
func (m MIME) ContentType() string {
	return string(m)
}

func (m MIME) String() string {
	return string(m)
}
