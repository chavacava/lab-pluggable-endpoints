package main

import (
	"flag"
	"log"
	"net/http"
	"plugin"
	"strings"

	"github.com/gorilla/mux"
)

// PluggableEndpoint is the interface to be implemented by endpoints willing to plug-into this server
type PluggableEndpoint interface {
	Path() string
	Handler() func(w http.ResponseWriter, r *http.Request)
}

func main() {
	endpoints := flag.String("e", "", "(comma separated list) of available endpoints")
	flag.Parse()

	r := mux.NewRouter()
	epNames := strings.Split(*endpoints, ",")

	for _, epn := range epNames {
		// open the endpoint's dynamic library
		pluginPath := "./endpoints/" + epn + "/" + epn + ".so"
		plug, err := plugin.Open(pluginPath)
		if err != nil {
			log.Fatalf("Unable to load endpoint %s: %v", epn, err)
		}

		// look in the library for the symbol (named) PluggableEndpoint
		endpointSymbol, err := plug.Lookup("PluggableEndpoint")
		if err != nil {
			log.Fatalf("Unable to lookup for a PluggableEndpoint symbol in the plugin: %v", err)
		}

		// check if the symbol implements the interface PluggableEndpoint
		var endpoint PluggableEndpoint
		endpoint, ok := endpointSymbol.(PluggableEndpoint)
		if !ok {
			log.Fatalf("The symbol is not of type %T but of type %T", endpoint, endpointSymbol)
		}

		// add the enpoint to the router
		r.HandleFunc(endpoint.Path(), endpoint.Handler())
	}

	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(":5000", r))
}
