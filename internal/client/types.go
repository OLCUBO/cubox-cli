package client

import "encoding/json"

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type Folder struct {
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

type Annotation struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Note       string `json:"note"`
	ImageURL   string `json:"image_url"`
	Color      string `json:"color"`
	CardID     string `json:"card_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
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
	Folder       *Folder  `json:"folder,omitempty"`
	URL          string   `json:"url"`
	CreateTime   string   `json:"create_time"`
	UpdateTime   string   `json:"update_time"`
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
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	ArticleTitle string           `json:"article_title"`
	Domain       string           `json:"domain"`
	Read         bool             `json:"read"`
	Starred      bool             `json:"starred"`
	Tags         []string         `json:"tags"`
	Folder       *Folder          `json:"folder,omitempty"`
	URL          string           `json:"url"`
	CreateTime   string           `json:"create_time"`
	UpdateTime   string           `json:"update_time"`
	Content      string           `json:"content"`
	Author       string           `json:"author"`
	Annotations  []Annotation     `json:"annotations"`
	Insight      *Insight         `json:"insight,omitempty"`
}

type CardFilterRequest struct {
	FolderFilters []string `json:"folder_filters,omitempty"`
	TagFilters   []string `json:"tag_filters,omitempty"`
	Starred      *bool    `json:"starred,omitempty"`
	Read         *bool    `json:"read,omitempty"`
	Annotated    *bool    `json:"annotated,omitempty"`
	Archived     *bool    `json:"archived,omitempty"`
	LastCardID   string   `json:"last_card_id,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	Keyword      string   `json:"keyword,omitempty"`
	Page         int      `json:"page,omitempty"`
	StartTime    string   `json:"start_time,omitempty"`
	EndTime      string   `json:"end_time,omitempty"`
}

type SaveCardEntry struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

type SaveCardsRequest struct {
	Cards            []SaveCardEntry `json:"cards"`
	FolderNestedName string          `json:"folder_nested_name,omitempty"`
	TagNestedNames   []string        `json:"tag_nested_names,omitempty"`
}

type CardUpdateRequest struct {
	ID               string   `json:"id"`
	Title            string   `json:"title,omitempty"`
	Description      string   `json:"description,omitempty"`
	Starred          *bool    `json:"starred,omitempty"`
	Read             *bool    `json:"read,omitempty"`
	FolderNestedName *string  `json:"folder_nested_name,omitempty"`
	TagNestedNames   []string `json:"tag_nested_names,omitempty"`
}

type MoveCardsRequest struct {
	FolderID string   `json:"folder_id"`
	CardIDs  []string `json:"card_ids"`
}

type CardAddTagsRequest struct {
	ID               string   `json:"id"`
	AddTagNestedNames []string `json:"add_tag_nested_names,omitempty"`
}

type CardRemoveTagsRequest struct {
	ID                  string   `json:"id"`
	RemoveTagNestedNames []string `json:"remove_tag_nested_names,omitempty"`
}

type RagQueryRequest struct {
	Query string `json:"query"`
}

type TagUpdateRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TagMergeRequest struct {
	SourceTagIDs []string `json:"source_tag_ids"`
	TargetTagID  string   `json:"target_tag_id"`
}

type AnnotationFilterRequest struct {
	Colors           []string `json:"colors,omitempty"`
	LastAnnotationID string   `json:"last_annotation_id,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Keyword          string   `json:"keyword,omitempty"`
	StartTime        string   `json:"start_time,omitempty"`
	EndTime          string   `json:"end_time,omitempty"`
}

