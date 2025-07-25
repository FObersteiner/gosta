package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointHistoricalLocation(t *testing.T) {
	// arrange
	ep := CreateHistoricalLocationsEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.NotEqual(t, ep, nil)
	assert.Equal(t, ep.GetName(), "yo")
}
