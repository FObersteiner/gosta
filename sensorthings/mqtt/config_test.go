package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTopics(t *testing.T) {
	// arrange
	// act
	topics := CreateTopics("GOST")
	// assert
	assert.Positive(t, len(topics), "Must have more than zero topics")
}
