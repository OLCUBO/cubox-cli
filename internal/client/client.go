package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) get(path string, params map[string]string) (json.RawMessage, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) post(path string, body interface{}) (json.RawMessage, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling request body: %w", err)
	}
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) (json.RawMessage, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if apiResp.Code != 200 {
		return nil, fmt.Errorf("API error %d: %s", apiResp.Code, apiResp.Message)
	}
	return apiResp.Data, nil
}

func (c *Client) ListGroups() ([]Group, error) {
	data, err := c.get("/c/api/third-party/group/list", nil)
	if err != nil {
		return nil, err
	}
	var groups []Group
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("parsing groups: %w", err)
	}
	return groups, nil
}

func (c *Client) ListTags() ([]Tag, error) {
	data, err := c.get("/c/api/third-party/tag/list", nil)
	if err != nil {
		return nil, err
	}
	var tags []Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("parsing tags: %w", err)
	}
	return tags, nil
}

func (c *Client) FilterCards(req *CardFilterRequest) ([]Card, error) {
	data, err := c.post("/c/api/third-party/card/filter", req)
	if err != nil {
		return nil, err
	}
	var cards []Card
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("parsing cards: %w", err)
	}
	return cards, nil
}

func (c *Client) GetCardContent(id string) (string, error) {
	data, err := c.get("/c/api/third-party/card/content", map[string]string{"id": id})
	if err != nil {
		return "", err
	}
	var content string
	if err := json.Unmarshal(data, &content); err != nil {
		return "", fmt.Errorf("parsing content: %w", err)
	}
	return content, nil
}
