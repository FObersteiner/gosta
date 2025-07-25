package endpoint

import (
	"sort"
	"testing"

	entities "github.com/FObersteiner/gosta-core"
	"github.com/FObersteiner/gosta-server/sensorthings/models"
	"github.com/stretchr/testify/assert"
)

func TestEndPointGetNameShouldReturnCorrectName(t *testing.T) {
	// arrange
	endpoint := Endpoint{}
	endpoint.Name = "test"
	endpoint.URL = "http://www.nu.nl"

	// act
	name := endpoint.GetName()
	output := endpoint.ShowOutputInfo()
	url := endpoint.GetURL()
	ops := endpoint.GetOperations()
	expand := endpoint.GetSupportedExpandParams()
	sel := endpoint.GetSupportedSelectParams()
	// point.AreQueryOptionsSupported()

	// assert
	assert.Equal(t, "test", name, "name should be correct")
	assert.False(t, output)
	assert.Equal(t, "http://www.nu.nl", url)
	assert.Equal(t, len(ops), 0)
	assert.Equal(t, len(expand), 0)
	assert.Equal(t, len(sel), 0)
}

func TestIsDynamic(t *testing.T) {
	// arrange
	urlDynamic := "http://www.{}.nl"
	urlNotDynamic := "http://www.nu.nl"

	// act
	resultNotDynamic := isDynamic(urlNotDynamic)
	resultDynamic := isDynamic(urlDynamic)

	// assert
	assert.False(t, resultNotDynamic)
	assert.True(t, resultDynamic)
}

func TestEndPointSortDuplicate(t *testing.T) {
	// arrange
	ep1 := &EndpointWrapper{Operation: models.EndpointOperation{Path: "ep1", OperationType: models.HTTPOperationGet}}
	ep2 := &EndpointWrapper{Operation: models.EndpointOperation{Path: "ep1", OperationType: models.HTTPOperationGet}}
	eps := SortedEndpoints{ep1, ep2}

	// assert
	assert.Panics(t, func() { sort.Sort(eps) })
}

func TestEndPointSort(t *testing.T) {
	// arrange
	endpoints := map[entities.EntityType]models.Endpoint{}
	endpoints[entities.EntityTypeDatastream] = &Endpoint{
		Operations: []models.EndpointOperation{
			{OperationType: models.HTTPOperationGet, Path: "ep1"},
			{OperationType: models.HTTPOperationPost, Path: "ep2"},
			{OperationType: models.HTTPOperationGet, Path: "{c:.*}ep3"},
			{OperationType: models.HTTPOperationGet, Path: "ep4{c:.*}"},
			{OperationType: models.HTTPOperationGet, Path: "{c:.*}ep5{test}"},
			{OperationType: models.HTTPOperationGet, Path: "ep6{test}"},
			{OperationType: models.HTTPOperationGet, Path: "ep7"},
		},
	}

	// act
	eps := EndpointsToSortedList(&endpoints)

	// assert
	assert.Equal(t, len(eps), 7, "Number of Endpoints should be 7")
	// post becomes first after sorting
	assert.Equal(t, eps[0].Operation.Path, "ep2")
	assert.Equal(t, eps[1].Operation.Path, "ep1")
	assert.Equal(t, eps[2].Operation.Path, "ep7")
	assert.Equal(t, eps[3].Operation.Path, "ep6{test}")
	assert.Equal(t, eps[4].Operation.Path, "{c:.*}ep3")
	assert.Equal(t, eps[5].Operation.Path, "ep4{c:.*}")
	assert.Equal(t, eps[6].Operation.Path, "{c:.*}ep5{test}")
}

func TestEndPointSortDynamic(t *testing.T) {
	// arrange
	httpep1 := &EndpointWrapper{}
	httpep1.Operation.Path = "ep1{}"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &EndpointWrapper{}
	httpep2.Operation.Path = "ep2{}longer"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := SortedEndpoints{httpep1, httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.Equal(t, len(eps), 2, "Number of Endpoints should be 2")
	// when both urls are dynamic, the longer path comes first
	assert.Equal(t, eps[0].Operation.Path, "ep2{}longer")
	assert.Equal(t, eps[1].Operation.Path, "ep1{}")
}

func TestEndPointSortlength(t *testing.T) {
	// arrange
	httpep1 := &EndpointWrapper{}
	httpep1.Operation.Path = "ep1"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &EndpointWrapper{}
	httpep2.Operation.Path = "ep2longer"
	httpep2.Operation.OperationType = models.HTTPOperationPost

	eps := SortedEndpoints{httpep1, httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.Equal(t, len(eps), 2, "Number of Endpoints should be 2")
	// when both urls are dynamic, the longer path comes first
	assert.Equal(t, eps[0].Operation.Path, "ep2longer")
	assert.Equal(t, eps[1].Operation.Path, "ep1")
}

func TestEndPointNotDynamic(t *testing.T) {
	// arrange
	httpep1 := &EndpointWrapper{}
	httpep1.Operation.Path = "ep1 {c:.*}"
	httpep1.Operation.OperationType = models.HTTPOperationGet
	httpep2 := &EndpointWrapper{}
	httpep2.Operation.Path = "ep2longer"
	httpep2.Operation.OperationType = models.HTTPOperationGet
	eps := SortedEndpoints{httpep1, httpep2}

	// act
	sort.Sort(eps)

	// assert
	assert.Equal(t, eps[0].Operation.Path, "ep2longer")
	assert.Equal(t, eps[1].Operation.Path, "ep1 {c:.*}")
}
