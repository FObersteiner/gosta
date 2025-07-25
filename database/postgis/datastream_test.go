package postgis

import (
	entities "github.com/gost/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatastreamParamFactory(t *testing.T) {
	// arrange
	values := map[string]interface{}{
		"datastream_id":          4,
		"datastream_name":        "name",
		"datastream_description": "desc",
	}

	// act
	entity, err := datastreamParamFactory(values)
	entitytype := entity.GetEntityType()

	// assert
	assert.NotEqual(t, entity, nil)
	// entities..
	assert.Equal(t, err, nil)
	assert.Equal(t, entity.GetID(), 4)
	assert.Equal(t, entitytype, entities.EntityTypeDatastream)
}
