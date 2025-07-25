package handlers

import (
	"net/http"

	entities "github.com/FObersteiner/gosta-core"
	"github.com/FObersteiner/gosta-server/sensorthings/models"
	"github.com/FObersteiner/gosta-server/sensorthings/rest/reader"
	"github.com/FObersteiner/gosta-server/sensorthings/rest/writer"
)

// handlePutRequest todo: currently almost same as handlePostRequest, merge if it stays like this
func handlePutRequest(w http.ResponseWriter, e *models.Endpoint, r *http.Request, entity entities.Entity, h *func() (interface{}, []error), indentJSON bool) {
	if !reader.CheckContentType(w, r, indentJSON) {
		return
	}

	byteData := reader.CheckAndGetBody(w, r, indentJSON)
	if byteData == nil {
		return
	}

	err := reader.ParseEntity(entity, byteData)
	if err != nil {
		writer.SendError(w, []error{err}, indentJSON)

		return
	}

	handle := *h

	data, err2 := handle()
	if err2 != nil {
		writer.SendError(w, err2, indentJSON)

		return
	}

	selfLink := entity.GetSelfLink()
	w.Header().Add("Location", selfLink)
	writer.SendJSONResponse(w, http.StatusOK, data, nil, indentJSON)
}
