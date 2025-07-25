package odata

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	godata "github.com/FObersteiner/gosta-odata"
	"github.com/gorilla/mux"
)

// SupportedExpandParameters contains a list of endpoints with their supported expand parameters
var SupportedExpandParameters map[string][]string

// SupportedSelectParameters contains a list of endpoints with their supported select parameters
var SupportedSelectParameters map[string][]string

// QueryOptions extents upon godata.GoDataQuery to implement extra
// odata functions not found in the godata package
type QueryOptions struct {
	godata.GoDataQuery
	Value           *GoDataValueQuery
	Ref             *GoDataRefQuery
	CollectionCount *GoDataCollectionCountQuery
	RawExpand       string
	RawFilter       string
	RawOrderBy      string
}

// ExpandParametersSupported returns if the QueryOptions expand request is supported by the endpoints
func (q *QueryOptions) ExpandParametersSupported(endpoint, expand string) bool {
	return q.checkMap(SupportedExpandParameters, endpoint, expand)
}

// SelectParametersSupported returns if the QueryOptions select request is supported by the endpoints
func (q *QueryOptions) SelectParametersSupported(endpoint, selectString string) bool {
	return q.checkMap(SupportedSelectParameters, endpoint, selectString)
}

func (q *QueryOptions) checkMap(mapToCheck map[string][]string, endpoint, request string) bool {
	// endpoint is not registered
	if _, ok := mapToCheck[endpoint]; !ok {
		return false
	}

	// return true if not found for registered endpoint
	for _, supported := range mapToCheck[endpoint] {
		if strings.EqualFold(supported, request) {
			return true
		}
	}

	return false
}

// GoDataValueQuery true when $value is requested false if not
type GoDataValueQuery bool

// GoDataRefQuery true when $ref is requested false if not
type GoDataRefQuery bool

// GoDataCollectionCountQuery true when addressing the count of a collection
type GoDataCollectionCountQuery bool

// ExpandItemToQueryOptions converts an ExpandItem into QueryOptions
func ExpandItemToQueryOptions(ei *godata.ExpandItem) *QueryOptions {
	qo := QueryOptions{}
	qo.Top = ei.Top
	qo.Filter = ei.Filter
	qo.OrderBy = ei.OrderBy
	qo.Search = ei.Search
	qo.Select = ei.Select
	qo.Skip = ei.Skip
	qo.Expand = ei.Expand

	return &qo
}

// ParseURLQuery parses an incoming url query into QueryOptions
func ParseURLQuery(query url.Values) (*QueryOptions, error) {
	if query == nil || len(query) == 0 {
		return nil, nil
	}

	qo, err := godata.ParseUrlQuery(query)
	if err != nil {
		return nil, err
	}

	result := &QueryOptions{}
	result.GoDataQuery = *qo

	value := query.Get("$value")

	val := GoDataValueQuery(false)
	if value != "" {
		val = GoDataValueQuery(true)
	}

	result.Value = &val

	value = query.Get("$ref")
	ref := GoDataRefQuery(false)

	if value != "" {
		ref = GoDataRefQuery(true)
	}

	result.Ref = &ref

	value = query.Get("$collectioncount")
	cc := GoDataCollectionCountQuery(false)

	if value != "" {
		cc = GoDataCollectionCountQuery(true)
	}

	result.CollectionCount = &cc

	// store raw queries
	result.RawExpand = query.Get("$expand")
	result.RawFilter = query.Get("$filter")
	result.RawOrderBy = query.Get("$orderby")

	return result, err
}

// GetQueryOptions creates QueryOptions based upon the incoming request
// QueryOptions = nil when no query was found, errors != nil if something
// went wrong with parsing the query into QueryOptions and will contain information
// on what went wrong
func GetQueryOptions(r *http.Request, maxEntities int) (*QueryOptions, []error) {
	// If request contains parameters from route wildcard convert it to a select query
	vars := mux.Vars(r)
	value := vars["params"]

	// Encode semicolon to %3B, semicolon is used in $expand with multiple inline queries: $expand=Datastreams/Observations($select=result;$top=2)
	// when r.URL.Query() is called the value will be cut off after semicolon, using r.URL.EscapedPath() won't work	either
	r.URL.RawQuery = strings.ReplaceAll(r.URL.RawQuery, ";", "%3B")

	values := r.URL.Query()

	if len(vars["params"]) > 0 {
		// If $ref found create select query with id
		switch vars["params"] {
		case "$ref":
			value = "id"
			values["$ref"] = []string{"true"}
			values["$select"] = []string{value}
		case "$count":
			values["$collectioncount"] = []string{"true"}
			values["$count"] = []string{"true"}
		default:
			values["$select"] = []string{value}
		}
	}

	if strings.HasSuffix(r.URL.Path, "$value") {
		values["$value"] = []string{"true"}
	}

	if t, ok := r.URL.Query()["$top"]; !ok {
		values["$top"] = []string{strconv.Itoa(maxEntities)}
	} else {
		top, err := strconv.Atoi(t[0])
		if err != nil || top > maxEntities {
			values["$top"] = []string{strconv.Itoa(maxEntities)}
		}
	}

	if _, ok := values["$skip"]; !ok {
		values["$skip"] = []string{"0"}
	}

	qo, e := ParseURLQuery(values)
	if e != nil {
		return nil, []error{e}
	}

	return qo, nil
}
