package core_http_server

import "net/http"

type HTTPServer struct {
	mux *http.ServeMux
}
