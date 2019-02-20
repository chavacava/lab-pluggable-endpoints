package main

import (
	"flag"
	"fmt"
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
	r.HandleFunc("/update/{epn}", update(r))
	epNames := strings.Split(*endpoints, ",")

	for _, epn := range epNames {
		endpoint, err := retrievePluggableEndpoint(epn)
		if err != nil {
			log.Fatalf("Unable to create a PluggableEndpoint for %s: %v", epn, err)
		}

		// add the endpoint to the router
		r.HandleFunc(endpoint.Path(), endpoint.Handler())
	}

	log.Print("Server is running")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func update(router *mux.Router) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		epn := vars["epn"]

		endpoint, err := retrievePluggableEndpoint(epn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unable to create a PluggableEndpoint for %s: %v", epn, err)))
			return
		}

		router.HandleFunc(endpoint.Path(), endpoint.Handler())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Added new endpoint %s", endpoint.Path())))
	}
}

func retrievePluggableEndpoint(epn string) (PluggableEndpoint, error) {
	// open the endpoint's dynamic library
	pluginPath := "./endpoints/" + epn + "/" + epn + ".so"
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load endpoint %s: %v", epn, err)
	}

	// look in the library for the symbol (named) PluggableEndpoint
	endpointSymbol, err := plug.Lookup("PluggableEndpoint")
	if err != nil {
		return nil, fmt.Errorf("Unable to lookup for a PluggableEndpoint symbol in the plugin: %v", err)
	}

	// check if the symbol implements the interface PluggableEndpoint
	var endpoint PluggableEndpoint
	endpoint, ok := endpointSymbol.(PluggableEndpoint)
	if !ok {
		return nil, fmt.Errorf("The symbol is not of type %T but of type %T", endpoint, endpointSymbol)
	}

	return endpoint, nil
}
