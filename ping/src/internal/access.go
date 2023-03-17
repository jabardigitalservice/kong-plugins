package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Kong/go-pdk"
)

func New() interface{} {
	return &Config{}
}

const (
	STATUS_OK     = "OK"
	STATUS_FAILED = "FAILED"
)

type (
	Response struct {
		Services []map[string]interface{} `json:"services"`
		Hostname string                   `json:"hostname"`
	}
)

func (conf *Config) Access(kong *pdk.PDK) {
	var (
		status     = "OK"
		statusCode = http.StatusOK
		ctx        = context.Background()
		response   Response
	)

	for _, service := range conf.Services {
		hostname, err := os.Hostname()
		if err != nil {
			kong.Response.Exit(500, err.Error(), nil)
		}

		response.Hostname = hostname

		resp, err := pingService(ctx, service["url"], service["method"], 1000)
		if err != nil {
			if resp != nil {
				statusCode = resp.StatusCode
			} else {
				statusCode = http.StatusFailedDependency
			}

			status = STATUS_FAILED
		}
		response.Services = append(response.Services, map[string]interface{}{
			"name":   service["name"],
			"code":   statusCode,
			"status": status,
			"url":    service["url"],
		})
	}

	jsonResponseString, err := json.Marshal(response)
	if err != nil {
		kong.Response.Exit(500, err.Error(), nil)
	}

	kong.Response.Exit(200, string(jsonResponseString), nil)
}

func pingService(ctx context.Context, url string, method string, timeout int64) (*http.Response, error) {
	inquiryTimeout := time.Duration(timeout) * time.Second

	ctx, cancel := context.WithTimeout(ctx, inquiryTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("pingService: failed to http.NewRequestWithContext %w", err)
	}

	var client = &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pingService: failed to http.DefaultClient.Do %w", err)
	}
	defer res.Body.Close()

	return res, nil
}
