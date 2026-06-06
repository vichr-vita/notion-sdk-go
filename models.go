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

// TitleProperty is a page property value containing title rich text.
type TitleProperty struct {
	Title []RichTextObject `json:"title"`
}

// RichTextProperty is a page property value containing rich text.
type RichTextProperty struct {
	RichText []RichTextObject `json:"rich_text"`
}

// Title creates a title property value with one plain text object.
func Title(content string) TitleProperty {
	return TitleProperty{Title: []RichTextObject{plainRichText(content)}}
}

// RichText creates a rich text property value with one plain text object.
func RichText(content string) RichTextProperty {
	return RichTextProperty{RichText: []RichTextObject{plainRichText(content)}}
}

// PageParent creates a parent reference to a page.
func PageParent(pageID string) Parent {
	return Parent{Type: "page_id", PageID: pageID}
}

// BlockParent creates a parent reference to a block.
func BlockParent(blockID string) Parent {
	return Parent{Type: "block_id", BlockID: blockID}
}

// DataSourceParent creates a parent reference to a data source.
func DataSourceParent(dataSourceID string) Parent {
	return Parent{Type: "data_source_id", DataSourceID: dataSourceID}
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
	Content        json.RawMessage `json:"-"`
	Raw            json.RawMessage `json:"-"`
}

// Comment is a Notion comment response object.
type Comment struct {
	Object         string              `json:"object"`
	ID             string              `json:"id"`
	Parent         *Parent             `json:"parent,omitempty"`
	DiscussionID   string              `json:"discussion_id,omitempty"`
	CreatedTime    time.Time           `json:"created_time"`
	LastEditedTime *time.Time          `json:"last_edited_time"`
	CreatedBy      *User               `json:"created_by,omitempty"`
	RichText       []RichTextObject    `json:"rich_text,omitempty"`
	DisplayName    *CommentDisplayName `json:"display_name,omitempty"`
	Attachments    []CommentAttachment `json:"attachments,omitempty"`
	Raw            json.RawMessage     `json:"-"`
}

// CommentDisplayName is a resolved comment author display name.
type CommentDisplayName struct {
	Type         string          `json:"type,omitempty"`
	ResolvedName string          `json:"resolved_name,omitempty"`
	Raw          json.RawMessage `json:"-"`
}

// CommentAttachment is a file attachment on a comment.
type CommentAttachment struct {
	File File            `json:"file"`
	Raw  json.RawMessage `json:"-"`
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
	Href        *string         `json:"href,omitempty"`
	Annotations *Annotations    `json:"annotations,omitempty"`
	Text        *TextContent    `json:"text,omitempty"`
	Mention     json.RawMessage `json:"mention,omitempty"`
	Equation    json.RawMessage `json:"equation,omitempty"`
	Raw         json.RawMessage `json:"-"`
}

// TextContent contains plain rich text content.
type TextContent struct {
	Content string `json:"content,omitempty"`
	Link    *Link  `json:"link,omitempty"`
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

func plainRichText(content string) RichTextObject {
	return RichTextObject{
		Type: "text",
		Text: &TextContent{Content: content},
	}
}

func (u *User) UnmarshalJSON(data []byte) error {
	type userAlias User
	var v userAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*u = User(v)
	return nil
}

func (p *Page) UnmarshalJSON(data []byte) error {
	type pageAlias Page
	var v pageAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*p = Page(v)
	return nil
}

func (b *Block) UnmarshalJSON(data []byte) error {
	type blockAlias Block
	var v blockAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	if v.Type != "" {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			return err
		}
		v.Content = append(v.Content[:0], fields[v.Type]...)
	}
	*b = Block(v)
	return nil
}

func (c *Comment) UnmarshalJSON(data []byte) error {
	type commentAlias Comment
	var v commentAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*c = Comment(v)
	return nil
}

func (d *CommentDisplayName) UnmarshalJSON(data []byte) error {
	type commentDisplayNameAlias CommentDisplayName
	var v commentDisplayNameAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*d = CommentDisplayName(v)
	return nil
}

func (a *CommentAttachment) UnmarshalJSON(data []byte) error {
	type commentAttachmentAlias CommentAttachment
	var v commentAttachmentAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*a = CommentAttachment(v)
	return nil
}

func (d *DataSource) UnmarshalJSON(data []byte) error {
	type dataSourceAlias DataSource
	var v dataSourceAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*d = DataSource(v)
	return nil
}

func (d *Database) UnmarshalJSON(data []byte) error {
	type databaseAlias Database
	var v databaseAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*d = Database(v)
	return nil
}

func (c *CustomEmoji) UnmarshalJSON(data []byte) error {
	type customEmojiAlias CustomEmoji
	var v customEmojiAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*c = CustomEmoji(v)
	return nil
}

func (r *RichTextObject) UnmarshalJSON(data []byte) error {
	type richTextObjectAlias RichTextObject
	var v richTextObjectAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*r = RichTextObject(v)
	return nil
}

func unmarshalWithRaw(data []byte, out any, raw *json.RawMessage) error {
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}
	*raw = append((*raw)[:0], data...)
	return nil
}
