package notion

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNullableResponseFieldsDecodeToNilPointers(t *testing.T) {
	var user User
	require.NoError(t, json.Unmarshal([]byte(`{
		"object": "user",
		"id": "user_a",
		"name": null,
		"avatar_url": null
	}`), &user))

	require.Nil(t, user.Name)
	require.Nil(t, user.AvatarURL)

	var page Page
	require.NoError(t, json.Unmarshal([]byte(`{
		"object": "page",
		"id": "page_a",
		"created_time": "2026-03-11T00:00:00Z",
		"last_edited_time": "2026-03-11T00:00:00Z",
		"cover": null,
		"icon": null,
		"parent": {"type": "workspace", "workspace": true},
		"archived": false,
		"in_trash": false,
		"properties": {},
		"url": "https://notion.so/page_a",
		"public_url": null
	}`), &page))

	require.Nil(t, page.Cover)
	require.Nil(t, page.Icon)
	require.Nil(t, page.PublicURL)
}

func TestOptionalRequestFieldsOmitEmptyValues(t *testing.T) {
	body := struct {
		RichText []RichTextObject `json:"rich_text,omitempty"`
	}{
		RichText: []RichTextObject{
			{
				Type: "text",
				Text: &TextContent{Content: "hello"},
			},
		},
	}

	data, err := json.Marshal(body)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"rich_text": [
			{
				"type": "text",
				"text": {
					"content": "hello"
				}
			}
		]
	}`, string(data))
	require.NotContains(t, string(data), `"href":null`)
	require.NotContains(t, string(data), `"link":null`)
}

func TestMajorModelsPreserveRawJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		out  any
	}{
		{
			name: "user",
			data: `{"object":"user","id":"user_a","type":"person","name":"Ada","avatar_url":null,"person":{"email":"ada@example.com"},"extra":true}`,
			out:  &User{},
		},
		{
			name: "page",
			data: `{"object":"page","id":"page_a","created_time":"2026-03-11T00:00:00Z","last_edited_time":"2026-03-11T00:00:00Z","cover":null,"icon":null,"parent":{"type":"workspace","workspace":true},"archived":false,"in_trash":false,"properties":{"Name":{"type":"title"}},"url":"https://notion.so/page_a","public_url":null,"extra":true}`,
			out:  &Page{},
		},
		{
			name: "block",
			data: `{"object":"block","id":"block_a","created_time":"2026-03-11T00:00:00Z","last_edited_time":"2026-03-11T00:00:00Z","has_children":false,"archived":false,"in_trash":false,"type":"paragraph","paragraph":{"rich_text":[]},"extra":true}`,
			out:  &Block{},
		},
		{
			name: "comment",
			data: `{"object":"comment","id":"comment_a","discussion_id":"disc_a","created_time":"2026-03-11T00:00:00Z","last_edited_time":null,"rich_text":[],"extra":true}`,
			out:  &Comment{},
		},
		{
			name: "data source",
			data: `{"object":"data_source","id":"ds_a","created_time":"2026-03-11T00:00:00Z","last_edited_time":"2026-03-11T00:00:00Z","title":[],"description":[],"icon":null,"parent":{"type":"workspace","workspace":true},"archived":false,"in_trash":false,"properties":{},"url":"https://notion.so/ds_a","public_url":null,"extra":true}`,
			out:  &DataSource{},
		},
		{
			name: "database",
			data: `{"object":"database","id":"db_a","created_time":"2026-03-11T00:00:00Z","last_edited_time":"2026-03-11T00:00:00Z","title":[],"description":[],"icon":null,"cover":null,"parent":{"type":"workspace","workspace":true},"archived":false,"in_trash":false,"properties":{},"url":"https://notion.so/db_a","public_url":null,"extra":true}`,
			out:  &Database{},
		},
		{
			name: "custom emoji",
			data: `{"object":"custom_emoji","id":"emoji_a","name":"ship","url":"https://notion.so/emoji.png","extra":true}`,
			out:  &CustomEmoji{},
		},
		{
			name: "rich text",
			data: `{"type":"text","plain_text":"hello","href":null,"text":{"content":"hello"},"extra":true}`,
			out:  &RichTextObject{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(tt.data), tt.out))
			require.JSONEq(t, tt.data, string(rawMessageFromModel(t, tt.out)))
		})
	}
}

func TestBlockPreservesDiscriminatedContent(t *testing.T) {
	var block Block
	require.NoError(t, json.Unmarshal([]byte(`{
		"object": "block",
		"id": "block_a",
		"created_time": "2026-03-11T00:00:00Z",
		"last_edited_time": "2026-03-11T00:00:00Z",
		"has_children": false,
		"archived": false,
		"in_trash": false,
		"type": "paragraph",
		"paragraph": {
			"rich_text": [
				{
					"type": "text",
					"text": {
						"content": "hello"
					}
				}
			]
		}
	}`), &block))

	require.Equal(t, "paragraph", block.Type)
	require.JSONEq(t, `{
		"rich_text": [
			{
				"type": "text",
				"text": {
					"content": "hello"
				}
			}
		]
	}`, string(block.Content))
}

func rawMessageFromModel(t *testing.T, model any) json.RawMessage {
	t.Helper()

	switch v := model.(type) {
	case *User:
		return v.Raw
	case *Page:
		return v.Raw
	case *Block:
		return v.Raw
	case *Comment:
		return v.Raw
	case *DataSource:
		return v.Raw
	case *Database:
		return v.Raw
	case *CustomEmoji:
		return v.Raw
	case *RichTextObject:
		return v.Raw
	default:
		t.Fatalf("missing raw extractor for %T", model)
		return nil
	}
}
