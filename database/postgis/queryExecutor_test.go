package postgis

import (
	entities "github.com/FObersteiner/gosta-core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddRelationToEntity(t *testing.T) {
	// parent entities.Entity, subEntities []entities.Entity
	// arrange
	thing := &entities.Thing{}
	location := &entities.Location{}
	historicalLocation := &entities.HistoricalLocation{}
	datastream := &entities.Datastream{}
	observation := &entities.Observation{}
	sensor := &entities.Sensor{}
	observedProperty := &entities.ObservedProperty{}
	foi := &entities.FeatureOfInterest{}

	// act
	addRelationToEntity(thing, []entities.Entity{location})
	addRelationToEntity(thing, []entities.Entity{historicalLocation})
	addRelationToEntity(thing, []entities.Entity{datastream})
	addRelationToEntity(location, []entities.Entity{historicalLocation})
	addRelationToEntity(location, []entities.Entity{thing})
	addRelationToEntity(historicalLocation, []entities.Entity{thing})
	addRelationToEntity(historicalLocation, []entities.Entity{location})
	addRelationToEntity(datastream, []entities.Entity{observation})
	addRelationToEntity(datastream, []entities.Entity{thing})
	addRelationToEntity(datastream, []entities.Entity{sensor})
	addRelationToEntity(datastream, []entities.Entity{observedProperty})
	addRelationToEntity(sensor, []entities.Entity{datastream})
	addRelationToEntity(observedProperty, []entities.Entity{datastream})
	addRelationToEntity(observation, []entities.Entity{datastream})
	addRelationToEntity(observation, []entities.Entity{foi})
	addRelationToEntity(foi, []entities.Entity{observation})

	// assert
	assert.Equal(t, len(thing.Locations), 1)
	assert.Equal(t, len(thing.HistoricalLocations), 1)
	assert.Equal(t, len(thing.Datastreams), 1)
	assert.Equal(t, len(location.HistoricalLocations), 1)
	assert.Equal(t, len(location.Things), 1)
	assert.NotNil(t, historicalLocation.Thing)
	assert.Equal(t, len(historicalLocation.Locations), 1)
	assert.Equal(t, len(datastream.Observations), 1)
	assert.NotNil(t, datastream.Thing)
	assert.NotNil(t, datastream.Sensor)
	assert.NotNil(t, datastream.ObservedProperty)
	assert.Equal(t, len(sensor.Datastreams), 1)
	assert.Equal(t, len(observedProperty.Datastreams), 1)
	assert.NotNil(t, observation.Datastream)
	assert.NotNil(t, observation.FeatureOfInterest)
	assert.Equal(t, len(foi.Observations), 1)
}
