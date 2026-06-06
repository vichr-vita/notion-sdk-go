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
//
// Paginated endpoints accept a *Pagination value for manual cursor control:
//
//	resp, err := client.Blocks.ListChildren(ctx, blockID, &notion.Pagination{
//		PageSize: 100,
//	})
//	if err != nil {
//		return err
//	}
//	if resp.HasMore && resp.NextCursor != nil {
//		resp, err = client.Blocks.ListChildren(ctx, blockID, &notion.Pagination{
//			StartCursor: *resp.NextCursor,
//			PageSize:    100,
//		})
//	}
//
// Convenience helpers iterate through cursor-based responses. Service-specific
// helpers, such as Blocks.ForEachChild, preserve endpoint-specific arguments:
//
//	err := client.Blocks.ForEachChild(ctx, blockID, &notion.Pagination{
//		PageSize: 100,
//	}, func(block notion.Block) error {
//		return nil
//	})
//
// Generic helpers are available for endpoints that return PaginatedResponse:
//
//	blocks, err := notion.All(ctx, func(ctx context.Context, p *notion.Pagination) (*notion.PaginatedResponse[notion.Block], error) {
//		return client.Blocks.ListChildren(ctx, blockID, p)
//	})
//
// Non-2xx API responses return *APIError. Use errors.As to inspect Notion's
// structured error code, message, request ID, HTTP status, and body snippet.
// Network and JSON decode failures remain normal Go errors.
//
//	page, err := client.Pages.Get(ctx, pageID)
//	if err != nil {
//		var apiErr *notion.APIError
//		if errors.As(err, &apiErr) {
//			log.Printf("notion status=%d code=%s request_id=%s: %s",
//				apiErr.StatusCode,
//				apiErr.Code,
//				apiErr.RequestID,
//				apiErr.Message,
//			)
//		}
//		return err
//	}
//	_ = page
package notion
