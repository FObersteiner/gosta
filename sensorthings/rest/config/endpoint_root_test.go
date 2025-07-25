package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointRoot(t *testing.T) {
	// arrange
	ep := CreateRootEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.NotEqual(t, ep, nil)
	assert.Equal(t, ep.GetName(), "yo")
}
