package handlers

import (
	"net/http"

	"github.com/FObersteiner/gosta-server/sensorthings/models"
	"github.com/FObersteiner/gosta-server/sensorthings/odata"
	"github.com/FObersteiner/gosta-server/sensorthings/rest/writer"
)

// handleGetRequest is the default function to handle incoming GET requests
func handleGetRequest(
	w http.ResponseWriter, e *models.Endpoint, r *http.Request, h *func(q *odata.QueryOptions, path string,
	) (interface{}, error), indentJSON bool, maxEntities int, externalURI string,
) {
	// Parse query options from request
	queryOptions, err := odata.GetQueryOptions(r, maxEntities)
	if err != nil && len(err) > 0 {
		writer.SendError(w, err, indentJSON)

		return
	}

	// Run the handler func such as Api.GetThingById
	handler := *h

	data, err2 := handler(queryOptions, externalURI+r.URL.RawPath)
	if err2 != nil {
		writer.SendError(w, []error{err2}, indentJSON)

		return
	}

	writer.SendJSONResponse(w, http.StatusOK, data, queryOptions, indentJSON)
}
