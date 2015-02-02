package handlers

import (
	"net/http"

	"github.com/pivotal-golang/envoy/domain"
)

type Cataloger interface {
	Catalog() domain.Catalog
}

type CatalogHandler struct {
	cataloger Cataloger
}

func NewCatalogHandler(cataloger Cataloger) CatalogHandler {
	return CatalogHandler{
		cataloger: cataloger,
	}
}

func (handler CatalogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	respond(w, http.StatusOK, handler.cataloger.Catalog())
}
