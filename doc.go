// Package notion provides a small Go client for the Notion API.
//
// Create one client with NewClient, then call endpoint groups through its
// services:
//
//	client := notion.NewClient("secret_xxx")
//
//	page, err := client.Pages.Get(ctx, pageID)
//	blocks, err := client.Blocks.ListChildren(ctx, blockID, nil)
//	results, err := client.Search.Query(ctx, notion.SearchRequest{Query: "Tasks"})
//
// The client includes services for pages, blocks, data sources, databases,
// search, users, comments, file uploads, and custom emojis. Request methods
// accept context.Context and return typed response structs with raw JSON fields
// where Notion objects are intentionally flexible.
package notion
