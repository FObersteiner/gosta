package config

import (
	"fmt"

	"github.com/FObersteiner/gosta-server/sensorthings/models"
	"github.com/FObersteiner/gosta-server/sensorthings/rest/endpoint"
	"github.com/FObersteiner/gosta-server/sensorthings/rest/handlers"
)

// CreateCreateObservationsEndpoint constructs the CreateObservations endpoint configuration
func CreateCreateObservationsEndpoint(externalURL string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Name:       "CreateObservations",
		OutputInfo: false,
		URL:        fmt.Sprintf("%s/%s/%s", externalURL, models.APIPrefix, "CreateObservations"),
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationPost, Path: "/v1.0/createobservations", Handler: handlers.HandlePostCreateObservations},
		},
	}
}
