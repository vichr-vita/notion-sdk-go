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
