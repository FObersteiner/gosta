package odata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryValidator(t *testing.T) {
	// arrange
	assert.True(t, IsValidOdataQuery("$filter=name eq 'ho'"))
	assert.True(t, IsValidOdataQuery(fmt.Sprintf("%sfilter=name eq 'ho'", "%24")))
	assert.False(t, IsValidOdataQuery("$notexisting=name eq 'ho'"))
}
