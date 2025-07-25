package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointVersion(t *testing.T) {
	// arrange
	ep := CreateVersionEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.Equal(t, ep.GetName(), "yo")
}
