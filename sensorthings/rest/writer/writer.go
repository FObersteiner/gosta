package writer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	gostErrors "github.com/FObersteiner/gosta-server/errors"
	"github.com/FObersteiner/gosta-server/sensorthings/models"
	"github.com/FObersteiner/gosta-server/sensorthings/odata"
)

// SendJSONResponse sends the desired message to the user
// the message will be marshalled into JSON
func SendJSONResponse(w http.ResponseWriter, status int, data interface{}, qo *odata.QueryOptions, indentJSON bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if data != nil {
		b, err := JSONMarshal(data, true, indentJSON)
		if err != nil {
			panic(err)
		}

		if qo != nil {
			// $count for collection is requested url/$count and not the query ?$count=true, ToDo: move to API code?e
			if qo.CollectionCount != nil && bool(*qo.CollectionCount) {
				b, err = convertForCountResponse(b, qo)
			} else if qo.Value != nil && bool(*qo.Value) {
				b, err = convertForValueResponse(b, qo)
			}

			if err != nil {
				SendError(w, []error{err}, indentJSON)
			}
		}

		w.WriteHeader(status)
		w.Write(b)
	}
}

// JSONMarshal converts the data and converts special characters such as &
func JSONMarshal(data interface{}, safeEncoding, indentJSON bool) ([]byte, error) {
	var b []byte

	var err error
	if indentJSON {
		b, err = json.MarshalIndent(data, "", "   ")
	} else {
		b, err = json.Marshal(data)
	}

	// This code is needed if the response contains special characters like &, <, >,
	// and those characters must not be converted to safe encoding.
	if safeEncoding {
		b = bytes.ReplaceAll(b, []byte("\\u003c"), []byte("<"))
		b = bytes.ReplaceAll(b, []byte("\\u003e"), []byte(">"))
		b = bytes.ReplaceAll(b, []byte("\\u0026"), []byte("&"))
	}

	return b, err
}

// SendError creates an ErrorResponse message and sets it to the user
// using SendJSONResponse
func SendError(w http.ResponseWriter, error []error, indentJSON bool) {
	// errs cannot be marshalled, create strings
	errs := make([]string, len(error))
	for idx, value := range error {
		errs[idx] = value.Error()
	}

	// Set the status code, default   for error, check if there is an ApiError an get
	statusCode := http.StatusInternalServerError

	if len(error) > 0 {
		// if there is Encoding type error, sends bad request (400 range)
		if strings.Contains(errs[0], "Encoding not supported") ||
			strings.Contains(errs[0], "No matching token") ||
			strings.Contains(errs[0], "invalid input syntax") ||
			strings.Contains(errs[0], "Error executing query") {
			statusCode = http.StatusBadRequest
		}

		{
			var e gostErrors.APIError
			switch {
			case errors.As(error[0], &e):
				statusCode = e.GetHTTPErrorStatusCode()
			}
		}
	}

	statusText := http.StatusText(statusCode)
	errorResponse := models.ErrorResponse{
		Error: models.ErrorContent{
			StatusText: statusText,
			StatusCode: statusCode,
			Messages:   errs,
		},
	}

	SendJSONResponse(w, statusCode, errorResponse, nil, indentJSON)
}

func convertForValueResponse(b []byte, qo *odata.QueryOptions) ([]byte, error) {
	// $value is requested only send back the value
	errMessage := fmt.Errorf("Unable to retrieve $value for %v", qo.Select.SelectItems)

	var m map[string]json.RawMessage

	err := json.Unmarshal(b, &m)
	if err != nil || qo.Select == nil || qo.Select.SelectItems == nil || len(qo.Select.SelectItems) == 0 {
		return nil, gostErrors.NewRequestInternalServerError(errMessage)
	}

	// if selected equals the key in json add to mVal
	mVal := []byte{}

	for k, v := range m {
		if strings.ToLower(k) == qo.Select.SelectItems[0].Segments[0].Value {
			mVal = v
		}
	}

	if len(mVal) == 0 {
		return nil, gostErrors.NewBadRequestError(errMessage)
	}

	value := string(mVal)
	value = strings.TrimPrefix(value, "\"")
	value = strings.TrimSuffix(value, "\"")

	b = []byte(value)

	return b, nil
}

func convertForCountResponse(b []byte, qo *odata.QueryOptions) ([]byte, error) {
	var m map[string]json.RawMessage

	json.Unmarshal(b, &m)

	if count, ok := m["@iot.count"]; ok {
		b = []byte(string(count))
	} else {
		return nil, gostErrors.NewBadRequestError(errors.New("/$count not available for endpoint"))
	}

	return b, nil
}
