package main

import (
	"fmt"
	"net/http"
)

type pingEndpoint struct{}

func (ep pingEndpoint) Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	}
}

func (ep pingEndpoint) Path() string {
	return "/ping"
}

var PluggableEndpoint pingEndpoint
