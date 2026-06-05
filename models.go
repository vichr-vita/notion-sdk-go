package notion

import (
	"encoding/json"
	"time"
)

// Parent identifies the Notion object that owns another object.
type Parent struct {
	Type         string `json:"type,omitempty"`
	PageID       string `json:"page_id,omitempty"`
	BlockID      string `json:"block_id,omitempty"`
	DatabaseID   string `json:"database_id,omitempty"`
	DataSourceID string `json:"data_source_id,omitempty"`
	Workspace    *bool  `json:"workspace,omitempty"`
}

// User is a Notion user response object.
type User struct {
	Object    string          `json:"object"`
	ID        string          `json:"id"`
	Type      string          `json:"type,omitempty"`
	Name      *string         `json:"name"`
	AvatarURL *string         `json:"avatar_url"`
	Person    *PersonUser     `json:"person,omitempty"`
	Bot       *BotUser        `json:"bot,omitempty"`
	Raw       json.RawMessage `json:"-"`
}

// PersonUser contains person-specific user fields.
type PersonUser struct {
	Email string `json:"email,omitempty"`
}

// BotUser contains bot-specific user fields.
type BotUser struct {
	Owner         json.RawMessage `json:"owner,omitempty"`
	WorkspaceName *string         `json:"workspace_name,omitempty"`
}

// Page is a Notion page response object.
type Page struct {
	Object         string                     `json:"object"`
	ID             string                     `json:"id"`
	CreatedTime    time.Time                  `json:"created_time"`
	LastEditedTime time.Time                  `json:"last_edited_time"`
	CreatedBy      *User                      `json:"created_by,omitempty"`
	LastEditedBy   *User                      `json:"last_edited_by,omitempty"`
	Cover          *File                      `json:"cover"`
	Icon           *Icon                      `json:"icon"`
	Parent         Parent                     `json:"parent"`
	Archived       bool                       `json:"archived"`
	InTrash        bool                       `json:"in_trash"`
	Properties     map[string]json.RawMessage `json:"properties,omitempty"`
	URL            *string                    `json:"url,omitempty"`
	PublicURL      *string                    `json:"public_url"`
	Raw            json.RawMessage            `json:"-"`
}

// Block is a Notion block response object.
type Block struct {
	Object         string          `json:"object"`
	ID             string          `json:"id"`
	Parent         *Parent         `json:"parent,omitempty"`
	CreatedTime    time.Time       `json:"created_time"`
	LastEditedTime time.Time       `json:"last_edited_time"`
	CreatedBy      *User           `json:"created_by,omitempty"`
	LastEditedBy   *User           `json:"last_edited_by,omitempty"`
	HasChildren    bool            `json:"has_children"`
	Archived       bool            `json:"archived"`
	InTrash        bool            `json:"in_trash"`
	Type           string          `json:"type,omitempty"`
	Raw            json.RawMessage `json:"-"`
}

// Comment is a Notion comment response object.
type Comment struct {
	Object         string           `json:"object"`
	ID             string           `json:"id"`
	Parent         *Parent          `json:"parent,omitempty"`
	DiscussionID   string           `json:"discussion_id,omitempty"`
	CreatedTime    time.Time        `json:"created_time"`
	LastEditedTime *time.Time       `json:"last_edited_time"`
	CreatedBy      *User            `json:"created_by,omitempty"`
	RichText       []RichTextObject `json:"rich_text,omitempty"`
	Raw            json.RawMessage  `json:"-"`
}

// DataSource is a Notion data source response object.
type DataSource struct {
	Object         string                     `json:"object"`
	ID             string                     `json:"id"`
	CreatedTime    time.Time                  `json:"created_time"`
	LastEditedTime time.Time                  `json:"last_edited_time"`
	CreatedBy      *User                      `json:"created_by,omitempty"`
	LastEditedBy   *User                      `json:"last_edited_by,omitempty"`
	Title          []RichTextObject           `json:"title,omitempty"`
	Description    []RichTextObject           `json:"description,omitempty"`
	Icon           *Icon                      `json:"icon"`
	Parent         Parent                     `json:"parent"`
	Archived       bool                       `json:"archived"`
	InTrash        bool                       `json:"in_trash"`
	Properties     map[string]json.RawMessage `json:"properties,omitempty"`
	URL            *string                    `json:"url,omitempty"`
	PublicURL      *string                    `json:"public_url"`
	Raw            json.RawMessage            `json:"-"`
}

// Database is a legacy Notion database response object.
type Database struct {
	Object         string                     `json:"object"`
	ID             string                     `json:"id"`
	CreatedTime    time.Time                  `json:"created_time"`
	LastEditedTime time.Time                  `json:"last_edited_time"`
	CreatedBy      *User                      `json:"created_by,omitempty"`
	LastEditedBy   *User                      `json:"last_edited_by,omitempty"`
	Title          []RichTextObject           `json:"title,omitempty"`
	Description    []RichTextObject           `json:"description,omitempty"`
	Icon           *Icon                      `json:"icon"`
	Cover          *File                      `json:"cover"`
	Parent         Parent                     `json:"parent"`
	Archived       bool                       `json:"archived"`
	InTrash        bool                       `json:"in_trash"`
	Properties     map[string]json.RawMessage `json:"properties,omitempty"`
	URL            *string                    `json:"url,omitempty"`
	PublicURL      *string                    `json:"public_url"`
	Raw            json.RawMessage            `json:"-"`
}

// File is a Notion file-like object.
type File struct {
	Type     string        `json:"type,omitempty"`
	File     *HostedFile   `json:"file,omitempty"`
	External *ExternalFile `json:"external,omitempty"`
}

// HostedFile contains a Notion-hosted file URL and optional expiry.
type HostedFile struct {
	URL        string     `json:"url,omitempty"`
	ExpiryTime *time.Time `json:"expiry_time,omitempty"`
}

// ExternalFile contains an externally hosted file URL.
type ExternalFile struct {
	URL string `json:"url,omitempty"`
}

// Icon is a Notion icon object.
type Icon struct {
	Type     string        `json:"type,omitempty"`
	Emoji    string        `json:"emoji,omitempty"`
	File     *HostedFile   `json:"file,omitempty"`
	External *ExternalFile `json:"external,omitempty"`
	Custom   *CustomEmoji  `json:"custom_emoji,omitempty"`
}

// CustomEmoji is a Notion custom emoji response object.
type CustomEmoji struct {
	Object string          `json:"object"`
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	URL    string          `json:"url"`
	Raw    json.RawMessage `json:"-"`
}

// RichTextObject is a Notion rich text object.
type RichTextObject struct {
	Type        string          `json:"type,omitempty"`
	PlainText   string          `json:"plain_text,omitempty"`
	Href        *string         `json:"href"`
	Annotations *Annotations    `json:"annotations,omitempty"`
	Text        *TextContent    `json:"text,omitempty"`
	Mention     json.RawMessage `json:"mention,omitempty"`
	Equation    json.RawMessage `json:"equation,omitempty"`
	Raw         json.RawMessage `json:"-"`
}

// TextContent contains plain rich text content.
type TextContent struct {
	Content string `json:"content,omitempty"`
	Link    *Link  `json:"link"`
}

// Link contains a rich text link URL.
type Link struct {
	URL string `json:"url,omitempty"`
}

// Annotations contains rich text formatting flags.
type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color,omitempty"`
}
