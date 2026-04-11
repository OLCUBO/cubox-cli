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

// Card is the shape returned by the card/filter list endpoint.
type Card struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ArticleTitle string   `json:"article_title"`
	Domain       string   `json:"domain"`
	Read         bool     `json:"read"`
	Starred      bool     `json:"starred"`
	Tags         []string `json:"tags"`
	Group        *Group   `json:"group,omitempty"`
	URL          string   `json:"url"`
	CreateTime   string   `json:"create_time"`
	UpdateTime   string   `json:"update_time"`
}

type Highlight struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	ImageURL   string `json:"image_url"`
	CuboxURL   string `json:"cubox_url,omitempty"`
	Note       string `json:"note"`
	Color      string `json:"color"`
	CreateTime string `json:"create_time"`
}

type InsightQA struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type Insight struct {
	ID         string      `json:"id"`
	Summary    string      `json:"summary"`
	QAs        []InsightQA `json:"qas"`
	CreateTime string      `json:"create_time"`
	UpdateTime string      `json:"update_time"`
}

// CardDetail is the full shape returned by the card/detail endpoint.
type CardDetail struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	ArticleTitle string      `json:"article_title"`
	Domain       string      `json:"domain"`
	Read         bool        `json:"read"`
	Starred      bool        `json:"starred"`
	Tags         []string    `json:"tags"`
	Group        *Group      `json:"group,omitempty"`
	URL          string      `json:"url"`
	CreateTime   string      `json:"create_time"`
	UpdateTime   string      `json:"update_time"`
	Content      string      `json:"content"`
	Author       string      `json:"author"`
	Highlights   []Highlight `json:"highlights"`
	Insight      *Insight    `json:"insight,omitempty"`
}

type CardFilterRequest struct {
	GroupFilters []string `json:"group_filters,omitempty"`
	TagFilters   []string `json:"tag_filters,omitempty"`
	Starred      *bool    `json:"starred,omitempty"`
	Read         *bool    `json:"read,omitempty"`
	Annotated    *bool    `json:"annotated,omitempty"`
	LastCardID   string   `json:"last_card_id,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	Keyword      string   `json:"keyword,omitempty"`
	Page         int      `json:"page,omitempty"`
	StartTime    string   `json:"start_time,omitempty"`
	EndTime      string   `json:"end_time,omitempty"`
}

type SaveURLsRequest struct {
	URLs    []string `json:"urls"`
	GroupID string   `json:"group_id,omitempty"`
	TagIDs  []string `json:"tag_ids,omitempty"`
}

type CardUpdateRequest struct {
	ID        string   `json:"id"`
	Starred   *bool    `json:"starred,omitempty"`
	Read      *bool    `json:"read,omitempty"`
	Archive   *bool    `json:"archive,omitempty"`
	GroupID   string   `json:"group_id,omitempty"`
	AddTagIDs []string `json:"add_tag_ids,omitempty"`
}

type MarkFilterRequest struct {
	Colors          []string `json:"colors,omitempty"`
	LastHighlightID string   `json:"last_highlight_id,omitempty"`
	Limit           int      `json:"limit,omitempty"`
	Keyword         string   `json:"keyword,omitempty"`
	StartTime       string   `json:"start_time,omitempty"`
	EndTime         string   `json:"end_time,omitempty"`
}

type Mark struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Note       string `json:"note"`
	ImageURL   string `json:"image_url"`
	Color      string `json:"color"`
	CardID     string `json:"card_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
