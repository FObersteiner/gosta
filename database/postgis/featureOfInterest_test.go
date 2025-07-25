package postgis

import (
	entities "github.com/gost/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFoiParamFactory(t *testing.T) {
	// arrange
	values := map[string]interface{}{
		"featureofinterest_id":          4,
		"featureofinterest_name":        "name",
		"featureofinterest_description": "desc",
	}
	// todo: encodingtype + feature

	// act
	entity, err := featureOfInterestParamFactory(values)
	entitytype := entity.GetEntityType()

	// assert
	assert.NotEqual(t, entity, nil)
	// entities..
	assert.Equal(t, err, nil)
	assert.Equal(t, entity.GetID(), 4)
	assert.Equal(t, entitytype, entities.EntityTypeFeatureOfInterest)
}
