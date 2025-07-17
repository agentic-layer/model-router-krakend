package main

import (
	"github.com/agentic-layer/model-router-krakend/lib/header"
	"github.com/agentic-layer/model-router-krakend/lib/logging"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const (
	pluginName      = "echo"
	configKey       = "echo_config"
	defaultEndpoint = "/echo"
)

type config struct {
	Endpoint string `json:"endpoint"`
}

type registerer string

// HandlerRegisterer is the name of the symbol krakend looks up to try and register plugins
var HandlerRegisterer = registerer(pluginName)
var logger = logging.New(pluginName)

func main() {}

func init() {
	logger.Info("loaded")
}

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
	logger.Info("registered")
}

/*
 *	Most parts of this function are copied from https://github.com/traefik/whoami/ under the Apache License Version 2.0
 */
func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	var cfg config
	err := parseConfig(extra, &cfg)
	if err != nil {
		return nil, err
	}
	logger.Info("configuration loaded successfully")

	return http.HandlerFunc(r.handleRequest(cfg, handler)), nil
}

func (r registerer) handleRequest(cfg config, handler http.Handler) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == cfg.Endpoint {
			hostname, _ := os.Hostname()

			data := struct {
				Hostname string      `json:"hostname,omitempty"`
				IP       []string    `json:"ip,omitempty"`
				Headers  http.Header `json:"headers,omitempty"`
				URL      string      `json:"url,omitempty"`
				Host     string      `json:"host,omitempty"`
				Method   string      `json:"method,omitempty"`
				Name     string      `json:"name,omitempty"`
				Body     string      `json:"body,omitempty"`
			}{
				Hostname: hostname,
				IP:       []string{},
				Headers:  req.Header,
				URL:      req.URL.RequestURI(),
				Host:     req.Host,
				Method:   req.Method,
				Body:     getRequestBody(req),
			}

			ifaces, _ := net.Interfaces()
			for _, i := range ifaces {
				addrs, _ := i.Addrs()
				// handle err
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip != nil {
						data.IP = append(data.IP, ip.String())
					}
				}
			}

			w.Header().Set(header.ContentType, "application/json")
			if err := json.NewEncoder(w).Encode(data); err != nil {
				logger.Error("unable to parse response: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		handler.ServeHTTP(w, req)
	}
}

func getRequestBody(req *http.Request) string {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Warn("request  body could not be parsed: %s", err)
		return ""
	}
	return string(bodyBytes)
}

func parseConfig(extra map[string]interface{}, config *config) error {
	if extra[configKey] == nil {
		config.Endpoint = defaultEndpoint
		logger.Info("using default %s.endpoint %v", configKey, defaultEndpoint)
		return nil
	}

	authConfig, ok := extra[configKey].(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot read extra_config.%s", configKey)
	}

	raw, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("cannot marshall extra config back to JSON: %s", err.Error())
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return fmt.Errorf("cannot parse extra config: %s", err.Error())
	}

	if config.Endpoint == "" {
		config.Endpoint = defaultEndpoint
		logger.Info("using default %s.endpoint %v", configKey, defaultEndpoint)
	}

	return nil
}
