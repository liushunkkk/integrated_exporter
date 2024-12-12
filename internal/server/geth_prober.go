package server

import (
	"integrated-exporter/config"
	"io"
	"log"
	"net/http"
	"time"
)

func probeGeth(gs config.GethService) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, gs.Address, nil)
	if err != nil {
		log.Printf("Error creating request for probe %s %s: %v", GethService, gs.Name, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if gs.Token != "" {
		req.Header.Set("Authorization", gs.Token)
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request for probe %s %s: %v", GethService, gs.Name, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for probe %s %s: %v", GethService, gs.Name, err)
		return nil, err
	}
	return body, nil
}
