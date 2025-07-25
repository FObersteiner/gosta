package odata

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	godata "github.com/FObersteiner/gosta-odata"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestExpandParametersSupported(t *testing.T) {
	// arrange
	qo := &QueryOptions{}
	SupportedExpandParameters = map[string][]string{"things": {"locations", "datastreams"}}

	// assert
	assert.True(t, qo.ExpandParametersSupported("things", "locations"))
	assert.False(t, qo.ExpandParametersSupported("things", "featuresofinterest"))
	assert.False(t, qo.ExpandParametersSupported("bla", ""))
}

func TestSelectParametersSupported(t *testing.T) {
	// arrange
	qo := &QueryOptions{}
	SupportedSelectParameters = map[string][]string{"things": {"id", "name"}}

	// assert
	assert.True(t, qo.SelectParametersSupported("things", "name"))
	assert.False(t, qo.SelectParametersSupported("things", "nonexistingparam"))
	assert.False(t, qo.SelectParametersSupported("bla", ""))
}

func TestEmptyQuery(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Nil(t, query, "parsing query from localhost/v1.0/things should return nil")
}

func TestWrongQuery(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$count=ok")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Nil(t, query, "parsing query from localhost/v1.0/things?$count=ok should return nil because ok cannot be parsed to bool")
}

func TestParseFilter(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$filter=id%20eq%201&$top=1&$skip=1&$count=true&$expand=Datastreams/Observations,Locations&$orderby=id&$select=name")

	// act
	query, err := ParseURLQuery(uri.Query())

	// assert
	assert.NotNil(t, query, "%v", err)
	assert.NotNil(t, query.Filter)
	assert.NotNil(t, query.Count)
	assert.NotNil(t, query.Expand)
	assert.NotNil(t, query.OrderBy)
	assert.NotNil(t, query.Select)
	assert.NotNil(t, query.Top)
	assert.NotNil(t, query.Skip)
}

func TestParseFilterRef(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$ref=true")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, GoDataRefQuery(true), *query.Ref)
}

func TestParseFilterValue(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$value=true")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, GoDataValueQuery(true), *query.Value)
}

func TestSavingRawQuery(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/things?$filter=id eq 1&$orderby=id desc&$expand=Datastreams/Observations,Locations")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, "id eq 1", query.RawFilter)
	assert.Equal(t, "id desc", query.RawOrderBy)
	assert.Equal(t, "Datastreams/Observations,Locations", query.RawExpand)
}

func TestFilterGeography(t *testing.T) {
	// arrange
	uri, _ := url.Parse("localhost/v1.0/Locations?$filter=geo.intersects(location,geography'LINESTRING(7.5 51.5, 7.5 53.5)')")

	// act
	query, _ := ParseURLQuery(uri.Query())

	// assert
	assert.Equal(t, godata.FilterTokenFunc, query.Filter.Tree.Token.Type)
	assert.Equal(t, godata.FilterTokenLiteral, query.Filter.Tree.Children[0].Token.Type)
	assert.Equal(t, godata.FilterTokenGeography, query.Filter.Tree.Children[1].Token.Type)
}

func TestExpandItemToQueryOptions(t *testing.T) {
	// arrange
	ei := &godata.ExpandItem{Filter: &godata.GoDataFilterQuery{}}

	// act
	qo := ExpandItemToQueryOptions(ei)

	// assert
	assert.Equal(t, ei.Filter, qo.Filter)
}

func TestGetQueryOptions(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things/$value?$top=201", nil)

	// act
	qo, _ := GetQueryOptions(req, 20)

	// assert
	assert.NotEqual(t, qo, nil)
}

func TestReturnNoQueryOptionsOnFailedParse(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things?$count=none", nil)

	// act
	qo, err := GetQueryOptions(req, 20)

	// assert
	assert.Nil(t, qo)
	assert.NotNil(t, err)
}

func TestReadRefFromWildcard(t *testing.T) {
	// arrange
	router := mux.NewRouter()
	router.HandleFunc("/v1.0/Things{id}/{params}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qo, _ := GetQueryOptions(r, 20)
		fmt.Fprintf(w, "%v", *qo.Ref)
	}))

	ts := httptest.NewServer(router)
	defer ts.Close()

	// act
	resp, _ := http.Get(ts.URL + "/v1.0/Things(35)/$ref")

	// assert
	assert.NotEqual(t, resp, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := resp.Body
	result, _ := io.ReadAll(body)
	assert.Equal(t, string(result), "true")
}
