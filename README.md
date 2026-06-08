# notion-sdk-go

Go client for the Notion API.

This package is intentionally small: one client, typed request and response models,
context-aware methods, pagination helpers, retry handling, and access to raw JSON
for flexible Notion objects.

## Install

```bash
go get github.com/vichr-vita/notion-sdk-go
```

## Quick start

Create an internal Notion integration, give it access to the pages or data
sources you want to use, and pass the integration token to the client.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	notion "github.com/vichr-vita/notion-sdk-go"
)

func main() {
	client := notion.NewClient(os.Getenv("NOTION_TOKEN"))

	resp, err := client.Search.Query(context.Background(), notion.SearchRequest{
		Query:    "Tasks",
		PageSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("%s %s\n", result.Object, result.ID)
	}
}
```

## Supported APIs

- Pages
- Blocks
- Data sources
- Legacy databases
- Search
- Users
- Comments
- File uploads
- Custom emojis

The default `Notion-Version` header is `2026-03-11`.

Override it when needed:

```go
client := notion.NewClient(
	"secret_xxx",
	notion.WithVersion("2026-03-11"),
)
```

## Data sources

Use `DataSources` for current Notion database-style content.

```go
rows, err := client.DataSources.Query(ctx, dataSourceID, notion.QueryDataSourceRequest{
	PageSize: 100,
})
```

`Databases` remains available for legacy database container endpoints. It does
not proxy or alias data source operations.

## Pagination

Paginated endpoints accept `*notion.Pagination` for manual cursor control.

```go
resp, err := client.Blocks.ListChildren(ctx, blockID, &notion.Pagination{
	PageSize: 100,
})
```

Helpers are available for cursor iteration:

```go
err := client.Blocks.ForEachChild(ctx, blockID, &notion.Pagination{
	PageSize: 100,
}, func(block notion.Block) error {
	return nil
})
```

Generic helpers work with endpoints that return `PaginatedResponse[T]`:

```go
blocks, err := notion.All(ctx, func(ctx context.Context, p *notion.Pagination) (*notion.PaginatedResponse[notion.Block], error) {
	return client.Blocks.ListChildren(ctx, blockID, p)
})
```

## Retries

The client retries `429`, `500`, `502`, `503`, and `504` responses up to two
times by default. It honors `Retry-After` when Notion sends it, then uses
exponential backoff with jitter.

```go
client := notion.NewClient("secret_xxx", notion.WithRetryConfig(notion.RetryConfig{
	MaxRetries: 4,
	Delay:      250 * time.Millisecond,
	MaxDelay:   2 * time.Second,
	Jitter:     0.1,
}))
```

## Errors

Non-2xx API responses return `*notion.APIError`.

```go
page, err := client.Pages.Get(ctx, pageID)
if err != nil {
	var apiErr *notion.APIError
	if errors.As(err, &apiErr) {
		log.Printf("notion status=%d code=%s request_id=%s: %s",
			apiErr.StatusCode,
			apiErr.Code,
			apiErr.RequestID,
			apiErr.Message,
		)
	}
	return err
}
```

## Examples

Runnable examples live in `examples/`:

- `examples/search`
- `examples/page`
- `examples/blocks`
- `examples/data_source`

Most examples require:

```bash
export NOTION_TOKEN=secret_xxx
```

Some examples also require IDs such as `NOTION_PAGE_ID` or
`NOTION_DATA_SOURCE_ID`.

## Scope

This is a v0 SDK. It intentionally does not include OAuth flows, admin
endpoints, managed users, webhook receiver helpers, webhook event models, legal
holds, exports, or view APIs.

Use an internal integration token with workspace-granted access.

## API reference generator

`notion-api-reference/` contains the scraper and generated reference artifacts
used while building the SDK. It stays in the public repository so the source
mapping is auditable. Local raw HTML cache and dependency folders are ignored.

## Development

```bash
go test ./...
go vet ./...
```

## Release

This module follows Go module versioning. For the first public release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## License

MIT
