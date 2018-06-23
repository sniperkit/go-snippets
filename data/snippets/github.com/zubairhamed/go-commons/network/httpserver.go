package network

import (
	"log"
	"net/http"
	"strconv"
)

func NewDefaultHttpServer(port string) *HttpServer {
	return &HttpServer{
		port: port,
		routes: []*Route{},
	}
}

type HttpServer struct {
	port   string
	routes []*Route
}

func (h *HttpServer) Start() {
	h.serveServer()
}

func (h *HttpServer) serveServer() {
	wh := &WrappedHandler{
		routes: h.routes,
	}
	log.Println("Serving HTTP on port ", h.port)
	err := http.ListenAndServe(h.port, wh)
	if err != nil {
		log.Fatal(err)
	}

}

func (h *HttpServer) NewRoute(path string, method string, fn RouteHandler) *Route {
	route := CreateNewRoute(path, method, fn)
	h.routes = append(h.routes, route)

	return route
}

type WrappedHandler struct {
	routes []*Route
}

func (wh *WrappedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, attrs, err := MatchingRoute(r.URL.Path, r.Method, nil, wh.routes)

	if err != nil {
		if err == ERR_NO_MATCHING_ROUTE {
			http.NotFound(w, r)
			return
		}

		if err == ERR_NO_MATCHING_METHOD {
			// TODO
			return
		}

		if err == ERR_UNSUPPORTED_CONTENT_FORMAT {
			// TODO
			return
		}
	} else {
		req := NewRequestFromHttp(attrs)
		resp := route.Handler(req).(*HttpResponse)

		SendHttpResponse(resp, w, r)
		return
	}
}

func SendHttpResponse(response *HttpResponse, w http.ResponseWriter, r *http.Request) {
	w.Write(response.Payload.GetBytes())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func NewRequestFromHttp(attrs map[string]string) *HttpRequest {
	return &HttpRequest{
		attrs: attrs,
	}
}

type HttpRequest struct {
	attrs map[string]string
}

func (c *HttpRequest) GetAttributes() map[string]string {
	return c.attrs
}

func (c *HttpRequest) GetAttribute(o string) string {
	return c.attrs[o]
}

func (c *HttpRequest) GetAttributeAsInt(o string) int {
	attr := c.GetAttribute(o)
	i, _ := strconv.Atoi(attr)

	return i
}

type HttpResponse struct {
	contentType   string
	err           error
	TemplateModel interface{}
	Payload       MessagePayload
}

func (r *HttpResponse) GetError() error {
	return r.err
}

func (r *HttpResponse) GetPayload() []byte {
	return r.Payload.GetBytes()
}

func (r *HttpResponse) GetContentType() string {
	return r.contentType
}

func (r *HttpResponse) GetTemplateModel() interface{} {
	return r.TemplateModel
}
