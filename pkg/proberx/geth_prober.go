package proberx

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/liushun-ing/integrated_exporter/config"
	"github.com/liushun-ing/integrated_exporter/pkg/constantx"
)

// ProbeGeth detect whether a Geth service exposing monitoring metrics is running properly
// and return its monitoring metrics.
func ProbeGeth(gs config.GethService) ([]byte, error) {
	timeout, err := time.ParseDuration(gs.Timeout)
	if err != nil {
		log.Printf("Failed to parse timeout duration for probe %s %s: %v", constantx.GethService, gs.Name, err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, gs.Address, nil)
	if err != nil {
		log.Printf("Error creating request for probe %s %s: %v", constantx.GethService, gs.Name, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if gs.Token != "" {
		req.Header.Set("Authorization", gs.Token)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request for probe %s %s: %v", constantx.GethService, gs.Name, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for probe %s %s: %v", constantx.GethService, gs.Name, err)
		return nil, err
	}
	return body, nil
}
