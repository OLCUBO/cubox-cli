package client

import "encoding/json"

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type Group struct {
	ID            string  `json:"id"`
	NestedName    string  `json:"nested_name"`
	Name          string  `json:"name"`
	ParentID      *string `json:"parent_id"`
	Uncategorized bool    `json:"uncategorized,omitempty"`
}

type Tag struct {
	ID         string  `json:"id"`
	NestedName string  `json:"nested_name"`
	Name       string  `json:"name"`
	ParentID   *string `json:"parent_id"`
}

type Highlight struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	ImageURL   string `json:"image_url"`
	CuboxURL   string `json:"cubox_url"`
	Note       string `json:"note"`
	Color      string `json:"color"`
	CreateTime string `json:"create_time"`
}

type Card struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	ArticleTitle string      `json:"article_title"`
	Domain       string      `json:"domain"`
	Type         string      `json:"type"`
	Tags         []string    `json:"tags"`
	URL          string      `json:"url"`
	CuboxURL     string      `json:"cubox_url"`
	WordsCount   int         `json:"words_count"`
	CreateTime   string      `json:"create_time"`
	UpdateTime   string      `json:"update_time"`
	Highlights   []Highlight `json:"highlights"`
}

type CardFilterRequest struct {
	GroupFilters   []string `json:"group_filters,omitempty"`
	TypeFilters    []string `json:"type_filters,omitempty"`
	TagFilters     []string `json:"tag_filters,omitempty"`
	Starred        *bool    `json:"starred,omitempty"`
	Read           *bool    `json:"read,omitempty"`
	Annotated      *bool    `json:"annotated,omitempty"`
	LastCardID     string   `json:"last_card_id,omitempty"`
	LastCardUpdate string   `json:"last_card_update_time,omitempty"`
	Limit          int      `json:"limit,omitempty"`
}
