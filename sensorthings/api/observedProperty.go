package api

import (
	"errors"

	entities "github.com/FObersteiner/gosta-core"
	gostErrors "github.com/FObersteiner/gosta-server/errors"
	"github.com/FObersteiner/gosta-server/sensorthings/odata"
)

// GetObservedProperty todo
func (a *APIv1) GetObservedProperty(id interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error) {
	op, err := a.db.GetObservedProperty(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(op, qo)

	return op, nil
}

// GetObservedPropertyByDatastream todo
func (a *APIv1) GetObservedPropertyByDatastream(datastreamID interface{}, qo *odata.QueryOptions, path string) (*entities.ObservedProperty, error) {
	op, err := a.db.GetObservedPropertyByDatastream(datastreamID, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(op, qo)

	return op, nil
}

// GetObservedProperties todo
func (a *APIv1) GetObservedProperties(qo *odata.QueryOptions, path string) (*entities.ArrayResponse, error) {
	ops, count, hasNext, err := a.db.GetObservedProperties(qo)
	if err != nil {
		return nil, err
	}

	for idx, item := range ops {
		i := *item
		a.SetLinks(&i, qo)
		ops[idx] = &i
	}

	var data interface{} = ops

	return a.createArrayResponse(count, hasNext, path, qo, data), nil
}

// PostObservedProperty todo
func (a *APIv1) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	_, err := containsMandatoryParams(op)
	if err != nil {
		return nil, err
	}

	nop, err2 := a.db.PostObservedProperty(op)
	if err2 != nil {
		return nil, []error{err2}
	}

	nop.SetAllLinks(a.config.GetExternalServerURI())

	return nop, nil
}

// PatchObservedProperty patches a given ObservedProperty
func (a *APIv1) PatchObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	if op.Datastreams != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch ObservedProperty"))
	}

	return a.db.PatchObservedProperty(id, op)
}

// PutObservedProperty patches a given ObservedProperty
func (a *APIv1) PutObservedProperty(id interface{}, op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	nop, err2 := a.db.PutObservedProperty(id, op)
	if err2 != nil {
		return nil, []error{err2}
	}

	nop.SetAllLinks(a.config.GetExternalServerURI())

	return nop, nil
}

// DeleteObservedProperty deletes a given ObservedProperty from the database
func (a *APIv1) DeleteObservedProperty(id interface{}) error {
	return a.db.DeleteObservedProperty(id)
}
