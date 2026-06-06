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
//
// Authentication uses an internal integration token. Every request sends
// Authorization: Bearer <token> and a Notion-Version header. The default
// Notion API version is 2026-03-11:
//
//	client := notion.NewClient("secret_xxx")
//
// Override the API version when a workspace or endpoint needs a different
// Notion-Version value:
//
//	client := notion.NewClient("secret_xxx", notion.WithVersion("2026-03-11"))
//
// Configuration uses functional options. WithHTTPClient supplies the
// *http.Client used for all requests; use it for custom transports, proxies, or
// client-wide timeouts. WithBaseURL changes the API root, which is useful for
// tests or compatible gateways. WithRequestHook and WithResponseHook install
// lightweight hooks for request inspection, logging, or test assertions.
//
//	client := notion.NewClient(
//		"secret_xxx",
//		notion.WithHTTPClient(httpClient),
//		notion.WithBaseURL("https://api.notion.com/v1"),
//		notion.WithRequestHook(func(req *http.Request) error {
//			return nil
//		}),
//	)
package notion
