package core_http_server

import "net/http"

type Route struct {
	Methode string
	Path    string
	Handler http.HandlerFunc
}

func NewRoute(methode, path string, handler http.HandlerFunc) Route {
	return Route{
		Methode: methode,
		Path:    path,
		Handler: handler,
	}
}
