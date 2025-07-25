package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FObersteiner/gosta-server/configuration"
	"github.com/FObersteiner/gosta-server/database/postgis"
	"github.com/FObersteiner/gosta-server/mqtt"
	"github.com/FObersteiner/gosta-server/sensorthings/api"
	"github.com/stretchr/testify/assert"
)

func TestCreateServer(t *testing.T) {
	// arrange
	server := createTestServer(8080, false)

	// act
	server.Stop()

	// assert
	assert.NotNil(t, server)
}

func TestFailRunServerHttp(t *testing.T) {
	// arrange
	server := createTestServer(789456456, false)

	// assert
	assert.Panics(t, func() { server.Start() })
}

func TestFailRunServerHttps(t *testing.T) {
	// arrange
	server := createTestServer(8080, true)

	// assert
	assert.Panics(t, func() { server.Start() })
}

func TestLowerCaseURI(t *testing.T) {
	n := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.Path, "/test")
	})

	ts := httptest.NewServer(LowerCaseURI(n))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/TEST")
	if err == nil && res != nil {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		assert.NotNil(t, b)
	} else {
		t.Fail()
	}
}

func TestPostProcessHandler(t *testing.T) {
	n := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusTeapot)
		rw.Header().Add("Location", "tea location")
		rw.Write([]byte("hello teapot"))
	})

	ts := httptest.NewServer(PostProcessHandler(n, "http://localhost:8080/"))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("X-Forwarded-For", "coffee")

	res, err := client.Do(req)
	if err == nil && res != nil {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		body := string(b)
		assert.NotNil(t, body)
		assert.Equal(t, body, "hello teapot")
		assert.Equal(t, res.StatusCode, http.StatusTeapot)
		assert.Equal(t, res.Header.Get("Location"), "tea location")
	} else {
		t.Fail()
	}
}

func TestIsValidRequestHandler(t *testing.T) {
	n := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// assert.True(t, req.URL.Path == "/test")
	})

	ts := httptest.NewServer(RequestErrorHandler(n))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/sensors?$notexisting eq 'ho'")
	if err == nil && res != nil {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, res.StatusCode, 400)
		assert.NotNil(t, b)
	} else {
		t.Fail()
	}
}

func createTestServer(port int, https bool) Server {
	cfg := configuration.Config{
		Server: configuration.ServerConfig{ExternalURI: "http://localhost:8080/"},
	}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := api.NewAPI(database, cfg, mqttServer)

	return CreateServer("localhost", port, &stAPI, https, "", "")
}
