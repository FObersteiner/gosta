package api

import (
	entities "github.com/FObersteiner/gosta-core"
	"github.com/FObersteiner/gosta-server/sensorthings/odata"

	"errors"

	gostErrors "github.com/FObersteiner/gosta-server/errors"
)

// GetHistoricalLocation retrieves a single HistoricalLocation by id
func (a *APIv1) GetHistoricalLocation(id interface{}, qo *odata.QueryOptions, path string) (*entities.HistoricalLocation, error) {
	hl, err := a.db.GetHistoricalLocation(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(hl, qo)

	return hl, nil
}

// GetHistoricalLocations retrieves all HistoricalLocations
func (a *APIv1) GetHistoricalLocations(qo *odata.QueryOptions, path string) (*entities.ArrayResponse, error) {
	hl, count, hasNext, err := a.db.GetHistoricalLocations(qo)

	return processHistoricalLocations(a, hl, qo, path, count, hasNext, err)
}

// GetHistoricalLocationsByLocation retrieves all HistoricalLocations linked to a given location
func (a *APIv1) GetHistoricalLocationsByLocation(locationID interface{}, qo *odata.QueryOptions, path string) (*entities.ArrayResponse, error) {
	hl, count, hasNext, err := a.db.GetHistoricalLocationsByLocation(locationID, qo)

	return processHistoricalLocations(a, hl, qo, path, count, hasNext, err)
}

// GetHistoricalLocationsByThing retrieves all HistoricalLocations linked to a given thing
func (a *APIv1) GetHistoricalLocationsByThing(thingID interface{}, qo *odata.QueryOptions, path string) (*entities.ArrayResponse, error) {
	hl, count, hasNext, err := a.db.GetHistoricalLocationsByThing(thingID, qo)

	return processHistoricalLocations(a, hl, qo, path, count, hasNext, err)
}

func processHistoricalLocations(a *APIv1, historicalLocations []*entities.HistoricalLocation, qo *odata.QueryOptions, path string, count int, hasNext bool, err error) (*entities.ArrayResponse, error) {
	for idx, item := range historicalLocations {
		i := *item
		a.SetLinks(&i, qo)
		historicalLocations[idx] = &i
	}

	var data interface{} = historicalLocations

	return a.createArrayResponse(count, hasNext, path, qo, data), nil
}

// PostHistoricalLocation adds a new HistoricalLocation to the database
func (a *APIv1) PostHistoricalLocation(hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error) {
	_, err := containsMandatoryParams(hl)
	if err != nil {
		return nil, err
	}

	l, err2 := a.db.PostHistoricalLocation(hl)
	if err2 != nil {
		return nil, []error{err2}
	}

	l.SetAllLinks(a.config.GetExternalServerURI())

	return l, nil
}

// PutHistoricalLocation adds a new HistoricalLocation to the database
func (a *APIv1) PutHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, []error) {
	l, err2 := a.db.PutHistoricalLocation(id, hl)
	if err2 != nil {
		return nil, []error{err2}
	}

	l.SetAllLinks(a.config.GetExternalServerURI())

	return l, nil
}

// PatchHistoricalLocation updates the given HistoricalLocation in the database
func (a *APIv1) PatchHistoricalLocation(id interface{}, hl *entities.HistoricalLocation) (*entities.HistoricalLocation, error) {
	if hl.Locations != nil || hl.Thing != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch HistoricalLocation"))
	}

	return a.db.PatchHistoricalLocation(id, hl)
}

// DeleteHistoricalLocation deletes a given HistoricalLocation from the database
func (a *APIv1) DeleteHistoricalLocation(id interface{}) error {
	return a.db.DeleteHistoricalLocation(id)
}
