package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	configStr = `{
      "echo_config": {
        "endpoint": "/echo"
      }
	}`
	configStrMinimal = `{
      "echo_config": {}
	}`
	configStrEmpty   = `{}`
	configStrApiEcho = `{
      "echo_config": {
        "endpoint": "/api/echo"
      }
	}`
	configStrFaulty = `{
      "echo_config": {
        "endpoint": 123
      }
	}`
)

type EmptyHandler struct {
}

func (h EmptyHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(404)
}

func TestSampleRequest(t *testing.T) {
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStr), &extraConfig)

	handlers, _ := HandlerRegisterer.registerHandlers(context.Background(), extraConfig, EmptyHandler{})
	ts := httptest.NewUnstartedServer(handlers)
	ts.Start()
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/echo", bytes.NewBuffer([]byte("Test Body")))
	req.Header.Set("Key", "Value")
	client := &http.Client{}

	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bodyBytes)

	assert.Contains(t, body, "Key\":[\"Value\"]")
	assert.Contains(t, body, "\"body\":\"Test Body\"")
}

func TestDifferentEndpoint(t *testing.T) {
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStrApiEcho), &extraConfig)

	handlers, _ := HandlerRegisterer.registerHandlers(context.Background(), extraConfig, EmptyHandler{})
	ts := httptest.NewUnstartedServer(handlers)
	ts.Start()
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/echo", bytes.NewBuffer([]byte("Test Body")))
	req.Header.Set("Key", "Value")
	client := &http.Client{}

	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bodyBytes)

	assert.Contains(t, body, "Key\":[\"Value\"]")
	assert.Contains(t, body, "\"body\":\"Test Body\"")
}

func TestWhitelistingIsWorkingCorrectly(t *testing.T) {
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStr), &extraConfig)

	handlers, _ := HandlerRegisterer.registerHandlers(context.Background(), extraConfig, EmptyHandler{})
	ts := httptest.NewUnstartedServer(handlers)
	ts.Start()
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/otherUrl", nil)
	client := &http.Client{}

	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 404)
}

func Test_parseConfig_returns_config_when_valid(t *testing.T) {
	// given
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStrApiEcho), &extraConfig)
	var cfg config

	// when
	err := parseConfig(extraConfig, &cfg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "/api/echo", cfg.Endpoint)
}

func Test_parseConfig_minimal_returns_config_when_valid(t *testing.T) {
	// given
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStrMinimal), &extraConfig)
	var cfg config

	// when
	err := parseConfig(extraConfig, &cfg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "/echo", cfg.Endpoint)
}

func Test_parseConfig_minimal_returns_config_when_empty(t *testing.T) {
	// given
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStrEmpty), &extraConfig)
	var cfg config

	// when
	err := parseConfig(extraConfig, &cfg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "/echo", cfg.Endpoint)
}

func Test_parseConfig_returns_error_when_invalid(t *testing.T) {
	// given
	var extraConfig map[string]interface{}
	json.Unmarshal([]byte(configStrFaulty), &extraConfig)
	var cfg config

	// when
	err := parseConfig(extraConfig, &cfg)

	// then
	assert.Error(t, err)
}
