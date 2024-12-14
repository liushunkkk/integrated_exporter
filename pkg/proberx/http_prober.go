package proberx

import (
	"bytes"
	"errors"
	"fmt"
	"integrated-exporter/config"
	"integrated-exporter/pkg/constantx"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func ProbeHttp(hs config.HttpService) error {
	timeout, err := time.ParseDuration(hs.Timeout)
	if err != nil {
		log.Printf("Failed to parse timeout duration for probe %s %s: %v", constantx.HttpService, hs.Name, err)
		return err
	}

	url := hs.Address
	reqBody := []byte(hs.Body)

	req, err := http.NewRequest(hs.Method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating request for probe %s %s: %v", constantx.HttpService, hs.Name, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if hs.Token != "" {
		req.Header.Set("Authorization", hs.Token)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request for probe %s %s: %v", constantx.HttpService, hs.Name, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for probe %s %s: %v", constantx.HttpService, hs.Name, err)
		return err
	}
	if hs.Response != "" {
		if !strings.Contains(string(body), hs.Response) {
			return errors.New(fmt.Sprintf("%s %s probe response does not contain %s", constantx.HttpService, hs.Name, hs.Response))
		}
	}
	return nil
}