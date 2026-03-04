package threadsync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	DefaultBaseURL = "https://api.threadsync.io/v1"
	Version        = "0.1.0"
)

type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
	Connections *ConnectionsService
	Sync        *SyncService
}

type Connection struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type SyncConfig struct {
	Source      Endpoint `json:"source"`
	Destination Endpoint `json:"destination"`
	Schedule    string   `json:"schedule"`
}

type Endpoint struct {
	Connection string `json:"connection"`
	Object     string `json:"object,omitempty"`
	Table      string `json:"table,omitempty"`
}

type SyncResult struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	RecordsSynced int    `json:"records_synced,omitempty"`
}

func New(token string) *Client {
	if token == "" {
		token = os.Getenv("THREADSYNC_API_TOKEN")
	}
	c := &Client{
		token:      token,
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{},
	}
	c.Connections = &ConnectionsService{client: c}
	c.Sync = &SyncService{client: c}
	return c
}

func (c *Client) request(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("threadsync: marshal error: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("threadsync: request error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "threadsync-go/"+Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("threadsync: request failed: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("threadsync: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("threadsync: API error %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// ConnectionsService handles connection-related API calls.
type ConnectionsService struct{ client *Client }

func (s *ConnectionsService) Create(provider string, options map[string]interface{}) (*Connection, error) {
	body := map[string]interface{}{"provider": provider}
	for k, v := range options {
		body[k] = v
	}
	data, err := s.client.request("POST", "/connections", body)
	if err != nil {
		return nil, err
	}
	var conn Connection
	return &conn, json.Unmarshal(data, &conn)
}

func (s *ConnectionsService) Get(id string) (*Connection, error) {
	data, err := s.client.request("GET", "/connections/"+id, nil)
	if err != nil {
		return nil, err
	}
	var conn Connection
	return &conn, json.Unmarshal(data, &conn)
}

// SyncService handles sync-related API calls.
type SyncService struct{ client *Client }

func (s *SyncService) Create(config *SyncConfig) (*SyncResult, error) {
	data, err := s.client.request("POST", "/syncs", config)
	if err != nil {
		return nil, err
	}
	var sync SyncResult
	return &sync, json.Unmarshal(data, &sync)
}

func (s *SyncService) Get(id string) (*SyncResult, error) {
	data, err := s.client.request("GET", "/syncs/"+id, nil)
	if err != nil {
		return nil, err
	}
	var sync SyncResult
	return &sync, json.Unmarshal(data, &sync)
}
