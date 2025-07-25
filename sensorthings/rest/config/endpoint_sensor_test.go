package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointSensor(t *testing.T) {
	// arrange
	ep := CreateSensorsEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.NotEqual(t, ep, nil)
	assert.Equal(t, ep.GetName(), "yo")
}
