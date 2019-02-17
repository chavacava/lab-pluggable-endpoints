package main

// This file is a plugin.
// Build it with:
//		go build -buildmode=plugin -o hello.so hello.go

import (
	"fmt"
	"net/http"
)

// The type implementing the PluggableEndpoint interface (defined elsewhere)
type helloEndpoint struct{}

// Handler handles a HTTP request
func (ep helloEndpoint) Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello!")
	}
}

// Path yields the path (or route) for this endpoint
func (ep helloEndpoint) Path() string {
	return "/hello"
}

// PluggableEndpoint is a PUBLIC variable of the type of this pluggable endpoint
// By convention, its name is that of the interface it implements
var PluggableEndpoint helloEndpoint
