package postgis

import (
	entities "github.com/gost/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObservationParamFactory(t *testing.T) {
	// arrange
	phenomenonTime := "2015-03-06T00:00:00.000Z"
	resultTime := "2015-03-06T00:00:00.000Z"
	validTime := "2015-03-06T00:00:00.000Z"

	values := map[string]interface{}{
		"observation_id":             4,
		"observation_phenomenontime": phenomenonTime,
		"observation_result":         "!0.5",
		"observation_resulttime":     resultTime,
		"observation_resultquality":  "goed",
		"observation_validtime":      validTime,
		// "observation_parameters": "test",
	}

	// act
	entity, err := observationParamFactory(values)
	entitytype := entity.GetEntityType()
	// todo: how to get the observation??

	// assert
	assert.NotEqual(t, entity, nil)
	// entities..
	assert.Equal(t, err, nil)
	assert.Equal(t, entity.GetID(), 4)
	assert.Equal(t, entitytype, entities.EntityTypeObservation)
	// assert.True(t,*observation.ResultTime == resultTime)
}
