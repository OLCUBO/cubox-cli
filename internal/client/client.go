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

func (c *Client) ListFolders() ([]Folder, error) {
	data, err := c.get("/c/api/cli/folder/list", nil)
	if err != nil {
		return nil, err
	}
	var folders []Folder
	if err := json.Unmarshal(data, &folders); err != nil {
		return nil, fmt.Errorf("parsing folders: %w", err)
	}
	return folders, nil
}

func (c *Client) ListTags() ([]Tag, error) {
	data, err := c.get("/c/api/cli/tag/list", nil)
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
	data, err := c.post("/c/api/cli/card/filter", req)
	if err != nil {
		return nil, err
	}
	var cards []Card
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("parsing cards: %w", err)
	}
	return cards, nil
}

func (c *Client) GetCardDetail(id string) (*CardDetail, error) {
	data, err := c.get("/c/api/cli/card/detail", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}
	var detail CardDetail
	if err := json.Unmarshal(data, &detail); err != nil {
		return nil, fmt.Errorf("parsing card detail: %w", err)
	}
	return &detail, nil
}

func (c *Client) SaveCards(req *SaveCardsRequest) error {
	_, err := c.post("/c/api/cli/cards/save", req)
	return err
}

func (c *Client) UpdateCard(req *CardUpdateRequest) error {
	_, err := c.post("/c/api/cli/card/update", req)
	return err
}

func (c *Client) AddCardTags(req *CardAddTagsRequest) error {
	_, err := c.post("/c/api/cli/card/add/tags", req)
	return err
}

func (c *Client) RemoveCardTags(req *CardRemoveTagsRequest) error {
	_, err := c.post("/c/api/cli/card/remove/tags", req)
	return err
}

func (c *Client) RagQueryCards(query string) ([]Card, error) {
	data, err := c.post("/c/api/cli/card/rag/query", &RagQueryRequest{Query: query})
	if err != nil {
		return nil, err
	}
	var cards []Card
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("parsing rag results: %w", err)
	}
	return cards, nil
}

func (c *Client) DeleteCards(ids []string) error {
	_, err := c.post("/c/api/cli/cards/delete", ids)
	return err
}

func (c *Client) FilterAnnotations(req *AnnotationFilterRequest) ([]Annotation, error) {
	data, err := c.post("/c/api/cli/annotation/filter", req)
	if err != nil {
		return nil, err
	}
	var annotations []Annotation
	if err := json.Unmarshal(data, &annotations); err != nil {
		return nil, fmt.Errorf("parsing annotations: %w", err)
	}
	return annotations, nil
}
