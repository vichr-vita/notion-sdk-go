# Notion Go SDK PRD

## Product

Build a Go client library for the Notion API.

Primary goal:

- Let individual Go developers call Notion API endpoints with simple typed methods.
- Use public Notion API reference as implementation source.
- Keep SDK small, readable, and practical.

## Target User

Primary user:

- Go backend developer building personal Notion integrations.

Initial real user:

- Project owner using token auth for personal automation.

## Non-Goals

V0 does not include:

- OAuth flows.
- Admin API.
- Legal holds.
- Exports.
- Managed users.
- Webhook receiver framework.
- View APIs.
- Full perfect typing for every Notion union.
- Code generation.
- Gin or other web frameworks.

## Authentication

Support basic Notion token auth only.

Each request sends:

- `Authorization: Bearer <token>`
- `Notion-Version: 2026-03-11`

API version defaults to `2026-03-11`.

Version can be changed per client:

```go
client := notion.NewClient("secret_xxx", notion.WithVersion("2026-03-11"))
```

## Module

Module path:

```text
github.com/vichr-vita/notion-sdk-go
```

Package:

```go
package notion
```

Import:

```go
import "github.com/vichr-vita/notion-sdk-go"
```

## API Shape

Use resource services on one client.

```go
client := notion.NewClient("secret_xxx")

page, err := client.Pages.Get(ctx, pageID)
blocks, err := client.Blocks.ListChildren(ctx, blockID, nil)
res, err := client.Search.Query(ctx, notion.SearchRequest{Query: "Tasks"})
```

Rules:

- Every method takes `context.Context`.
- Client owns shared config.
- Services share client.
- Client is immutable after creation.
- Client is safe for concurrent use.
- Request structs are caller-owned.
- No package-level request functions.
- No caller-created subclients.

## HTTP Config

Constructor options:

```go
client := notion.NewClient(
    "secret_xxx",
    notion.WithHTTPClient(httpClient),
    notion.WithBaseURL(url),
    notion.WithVersion("2026-03-11"),
)
```

Rules:

- Use stdlib `net/http`.
- Default HTTP client if none given.
- No client-level timeout.
- Caller owns timeout via context or supplied HTTP client.

## V0 Endpoint Scope

Include:

- Pages.
- Page markdown.
- Blocks.
- Data sources.
- Databases.
- Search.
- Users.
- Comments.
- File uploads.
- Custom emojis.

Exclude:

- OAuth.
- Admin endpoints.
- Legal holds.
- Exports.
- Managed users.
- Webhook events.
- Views.

## Data Sources vs Databases

Expose both services:

- `client.DataSources`
- `client.Databases`

Docs should mark `Databases` as legacy or compatibility surface.

Examples should prefer `DataSources`.

No silent aliasing.

## Models

Use typed common fields plus flexible raw payloads.

Rules:

- Strong request structs for core endpoints.
- Strong response shells for common fields.
- Complex Notion unions use discriminator fields plus `json.RawMessage`.
- Major models keep raw JSON when useful.
- ID fields use `string`, not UUID type.
- Optional request fields use `omitempty`.
- Nullable response fields use pointers.
- Struct fields use PascalCase.
- JSON tags use Notion snake_case.

Example:

```go
type Page struct {
    Object         string                     `json:"object"`
    ID             string                     `json:"id"`
    CreatedTime    time.Time                  `json:"created_time"`
    LastEditedTime time.Time                  `json:"last_edited_time"`
    Properties     map[string]json.RawMessage `json:"properties,omitempty"`
    Raw            json.RawMessage            `json:"-"`
}
```

## Request Builders

V0 uses plain structs, not fluent builders.

```go
page, err := client.Pages.Create(ctx, notion.CreatePageRequest{
    Parent: notion.PageParent(parentID),
    Properties: map[string]any{
        "Name": notion.Title("Task"),
    },
})
```

Small helper constructors are allowed for painful common objects:

- `notion.Title("Name")`
- `notion.RichText("text")`
- `notion.PageParent(id)`
- `notion.DataSourceParent(id)`

## Pagination

Expose manual and convenience pagination.

Manual:

```go
resp, err := client.Blocks.ListChildren(ctx, blockID, &notion.Pagination{
    PageSize: 100,
})
```

Convenience:

```go
err := client.Blocks.ForEachChild(ctx, blockID, nil, func(block notion.Block) error {
    return nil
})
```

Generic helper may exist:

```go
notion.All(ctx, firstPageFn)
```

## Errors

Non-2xx responses return `*notion.APIError`.

```go
var apiErr *notion.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.StatusCode)
    fmt.Println(apiErr.Code)
    fmt.Println(apiErr.Message)
    fmt.Println(apiErr.RequestID)
}
```

Rules:

- Parse Notion error JSON when possible.
- Keep status code.
- Keep Notion error code.
- Keep Notion message.
- Keep request ID when response header provides it.
- Keep response body snippet for debugging.
- Network and JSON errors stay normal Go errors.

## Retries

Retry by default:

- `429`
- `500`
- `502`
- `503`
- `504`

Rules:

- Honor `Retry-After`.
- Use exponential backoff with jitter.
- Default max retries = `2`.
- Retry behavior is configurable.
- No global rate limiter in v0.

## File Uploads

V0 supports file upload endpoint methods:

- `List`
- `Create`
- `Retrieve`
- `Send`
- `Complete`

No multipart abstraction in v0 unless needed during implementation.

## Hooks

No logger dependency.

Support lightweight hooks:

```go
type RequestHook func(*http.Request)
type ResponseHook func(*http.Response)
```

Options:

```go
notion.WithRequestHook(hook)
notion.WithResponseHook(hook)
```

## Go Version

Minimum Go version:

```text
go 1.22
```

Rules:

- Avoid unnecessary dependencies.
- Use generics only where they improve public API clarity.

## Repo Layout

Root package files:

```text
go.mod
client.go
request.go
errors.go
pagination.go

pages.go
blocks.go
databases.go
data_sources.go
search.go
users.go
comments.go
file_uploads.go
custom_emojis.go

models.go
*_test.go
```

Reference and examples:

```text
notion-api-reference/
spec/
examples/search/
examples/page/
examples/blocks/
```

## Spec Workflow

Use scraped Notion reference as lookup source only.

Do not generate SDK code from OpenAPI.

Rules:

- Keep scraper in `notion-api-reference/`.
- Create curated spec files only if useful.
- Read API reference before implementing each endpoint.
- Ignore low-confidence scraped endpoints unless manually verified.
- Fix behavior from docs, not from inferred scraped paths.

## Testing

Use TDD:

1. Write failing test.
2. Implement behavior.
3. Make test pass.

Test stack:

- `github.com/stretchr/testify`
- procedural tests.
- use `require` mostly.
- no table-driven tests.

Test with fake `http.RoundTripper`.

Cover:

- method.
- path.
- query.
- body.
- auth header.
- version header.
- error parsing.
- retry behavior.
- pagination loop.

No live Notion tests by default.

Optional later:

- integration tests behind `integration` build tag.
- require `NOTION_TOKEN`.

## Examples

Create minimal runnable examples:

- `examples/search`
- `examples/page`
- `examples/blocks`

Examples use:

```text
NOTION_TOKEN
```

No examples for OAuth, admin, or webhooks in v0.

## Acceptance Criteria

V0 is acceptable when:

- Go module imports from repo root.
- Client can call included endpoint groups.
- Basic token auth works.
- Notion version header defaults correctly.
- User can override version, base URL, and HTTP client.
- Non-2xx responses return `APIError`.
- Pagination supports manual and iterator-style flows.
- Retry behavior works for selected status codes.
- Tests cover core request behavior.
- Examples compile.
- No generator needed to build or maintain SDK.
