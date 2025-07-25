package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointObservedProperties(t *testing.T) {
	// arrange
	ep := CreateObservedPropertiesEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.NotEqual(t, ep, nil)
	assert.Equal(t, ep.GetName(), "yo")
}
