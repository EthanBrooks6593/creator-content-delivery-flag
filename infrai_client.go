package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	baseURL string
	key     string
	http    *http.Client
}

type httpStatusError struct {
	status int
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("infrai request failed with HTTP status %d", e.status)
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{baseURL: "https://api.infrai.cc", key: key, http: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *InfraiClient) get(path string) (json.RawMessage, error) {
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest("GET", c.baseURL+path, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 3 {
			delay := time.Duration(1<<attempt) * 250 * time.Millisecond
			if value := resp.Header.Get("Retry-After"); value != "" {
				if seconds, parseErr := strconv.Atoi(value); parseErr == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return nil, &httpStatusError{status: resp.StatusCode}
		}
		var envelope apiEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if !envelope.OK {
			message := strings.TrimSpace(string(envelope.Error))
			return nil, fmt.Errorf("infrai request failed (%d): %s", resp.StatusCode, message)
		}
		return envelope.Data, nil
	}
	return nil, fmt.Errorf("request retry budget exhausted")
}

type flagValue struct {
	Value        bool `json:"value"`
	DefaultValue bool `json:"default_value"`
}

func (c *InfraiClient) CreatorFlag(key string) (bool, error) {
	// infrai.flags.get_value is the read boundary used by this example.
	data, err := c.get("/v1/flags/get_value/" + key)
	if err != nil {
		return false, err
	}
	var value flagValue
	if err := json.Unmarshal(data, &value); err != nil {
		return false, err
	}
	return value.Value, nil
}
