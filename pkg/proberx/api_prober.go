package proberx

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/liushun-ing/integrated_exporter/config"
	"github.com/liushun-ing/integrated_exporter/pkg/constantx"
)

// ProbeApi detect whether an API service exposing monitoring metrics is running properly
// and return its monitoring metrics.
func ProbeApi(as config.ApiService) ([]byte, error) {
	timeout, err := time.ParseDuration(as.Timeout)
	if err != nil {
		log.Printf("Failed to parse timeout duration for probe %s %s: %v", constantx.ApiService, as.Name, err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, as.Address, nil)
	if err != nil {
		log.Printf("Error creating request for probe %s %s: %v", constantx.ApiService, as.Name, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if as.Token != "" {
		req.Header.Set("Authorization", as.Token)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request for probe %s %s: %v", constantx.ApiService, as.Name, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for probe %s %s: %v", constantx.ApiService, as.Name, err)
		return nil, err
	}
	return body, nil
}
