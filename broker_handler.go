package envoy

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pivotal-cf-experimental/envoy/internal/handlers"
	"github.com/pivotal-cf-experimental/envoy/internal/middleware"
)

func NewBrokerHandler(broker Broker) http.Handler {
	catalogHandler := handlers.NewCatalogHandler(broker)
	provisionHandler := handlers.NewProvisionHandler(broker)
	bindHandler := handlers.NewBindHandler(broker)
	unbindHandler := handlers.NewUnbindHandler(broker)
	deprovisionHandler := handlers.NewDeprovisionHandler(broker)

	routes := map[string]http.Handler{
		"GET /v2/catalog":                                                          middleware.NewAuthenticator(catalogHandler, broker),
		"PUT /v2/service_instances/{instance_id}":                                  middleware.NewAuthenticator(provisionHandler, broker),
		"PUT /v2/service_instances/{instance_id}/service_bindings/{binding_id}":    middleware.NewAuthenticator(bindHandler, broker),
		"DELETE /v2/service_instances/{instance_id}/service_bindings/{binding_id}": middleware.NewAuthenticator(unbindHandler, broker),
		"DELETE /v2/service_instances/{instance_id}":                               middleware.NewAuthenticator(deprovisionHandler, broker),
	}

	router := mux.NewRouter()
	for endpoint, handler := range routes {
		parts := strings.Split(endpoint, " ")
		router.Handle(parts[1], handler).Methods(parts[0])
	}

	return router
}
