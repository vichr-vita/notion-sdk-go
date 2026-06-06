# Notion Go SDK Tasks

## Milestone 0: Project Skeleton

- [x] Create `go.mod` with module `github.com/vichr-vita/notion-sdk-go` and `go 1.22`.
- [x] Add root package `notion`.
- [x] Add base files: `client.go`, `request.go`, `errors.go`, `pagination.go`, `models.go`.
- [x] Add service files: `pages.go`, `blocks.go`, `data_sources.go`, `databases.go`, `search.go`, `users.go`, `comments.go`, `file_uploads.go`, `custom_emojis.go`.
- [x] Add test helper for fake `http.RoundTripper`.
- [x] Add `testify` dependency.

## Milestone 1: Client Core

- [x] Implement `NewClient(token string, opts ...Option) *Client`.
- [x] Add immutable shared client config.
- [x] Add services on client: `Pages`, `Blocks`, `DataSources`, `Databases`, `Search`, `Users`, `Comments`, `FileUploads`, `CustomEmojis`.
- [x] Implement `WithHTTPClient(*http.Client)`.
- [x] Implement `WithBaseURL(string)`.
- [x] Implement `WithVersion(string)`.
- [x] Default API version to `2026-03-11`.
- [x] Default HTTP client to `http.DefaultClient`.
- [x] Ensure client safe for concurrent use.
- [x] Test default config.
- [x] Test option overrides.

## Milestone 2: Request Pipeline

- [x] Implement shared request builder.
- [x] Add `Authorization: Bearer <token>` header.
- [x] Add `Notion-Version` header.
- [x] Add JSON request encoding.
- [x] Add JSON response decoding.
- [x] Add empty-body handling.
- [x] Add query parameter encoding for pagination and endpoint filters.
- [x] Add `RequestHook`.
- [x] Add `ResponseHook`.
- [x] Implement `WithRequestHook`.
- [x] Implement `WithResponseHook`.
- [x] Test method, path, query, headers, and body.
- [x] Test hooks fire with expected request/response.

## Milestone 3: Errors

- [x] Create `APIError`.
- [x] Include `StatusCode`.
- [x] Include Notion error `Code`.
- [x] Include Notion error `Message`.
- [x] Include `RequestID`.
- [x] Include response body snippet.
- [x] Parse Notion error JSON when possible.
- [x] Preserve network errors as normal Go errors.
- [x] Preserve JSON decode errors as normal Go errors.
- [x] Test structured Notion error parsing.
- [x] Test malformed error body fallback.
- [x] Test `errors.As` with `*APIError`.

## Milestone 4: Retries

- [x] Add retry config with default max retries = `2`.
- [x] Retry `429`, `500`, `502`, `503`, `504`.
- [x] Honor `Retry-After`.
- [x] Add exponential backoff with jitter.
- [x] Add option to configure retry behavior.
- [x] Ensure non-replayable request bodies handled correctly.
- [x] Test retry count.
- [x] Test no retry for non-retryable status.
- [x] Test retry success after transient failure.
- [x] Test `Retry-After` handling.

## Milestone 5: Pagination

- [x] Create `Pagination` request struct.
- [x] Create generic paginated response shell.
- [x] Support `start_cursor`.
- [x] Support `page_size`.
- [x] Support `has_more`.
- [x] Support `next_cursor`.
- [x] Implement iterator-style helpers for endpoint pages.
- [x] Optionally implement generic `All(ctx, firstPageFn)`.
- [x] Test manual pagination query fields.
- [x] Test iterator stops on `has_more=false`.
- [x] Test iterator returns callback error.

## Milestone 6: Common Models and Helpers

- [x] Define common object fields for pages, blocks, users, comments, data sources, databases, files, custom emojis.
- [x] Use `string` IDs.
- [x] Use pointer fields for nullable responses.
- [x] Use `omitempty` for optional request fields.
- [x] Preserve `json.RawMessage` on major models where useful.
- [x] Add union-friendly discriminator fields.
- [x] Implement `Title(string)`.
- [x] Implement `RichText(string)`.
- [x] Implement `PageParent(string)`.
- [x] Implement `DataSourceParent(string)`.
- [x] Test helper JSON output.
- [x] Test raw JSON preservation.

## Milestone 7: Search Service

- [x] Verify Search endpoint against Notion reference.
- [x] Add `SearchService`.
- [x] Add `Search.Query(ctx, SearchRequest)`.
- [x] Model query, filter, sort, pagination fields.
- [x] Return typed paginated search response with raw objects.
- [x] Test request method, path, body, auth, version.
- [x] Test response decoding.
- [x] Test API error handling through service.

## Milestone 8: Blocks Service

- [x] Verify Blocks endpoints against Notion reference.
- [x] Add `Blocks.Get(ctx, blockID)`.
- [x] Add `Blocks.Update(ctx, blockID, request)`.
- [x] Add `Blocks.Delete(ctx, blockID)`.
- [x] Add `Blocks.ListChildren(ctx, blockID, *Pagination)`.
- [x] Add `Blocks.AppendChildren(ctx, blockID, request)`.
- [x] Add `Blocks.ForEachChild(ctx, blockID, *Pagination, func(Block) error)`.
- [x] Test each method path and HTTP verb.
- [x] Test children pagination.
- [x] Test update and append body encoding.

## Milestone 9: Pages Service

- [x] Verify Pages endpoints against Notion reference.
- [x] Add `Pages.Get(ctx, pageID)`.
- [x] Add `Pages.Create(ctx, CreatePageRequest)`.
- [x] Add `Pages.Update(ctx, pageID, UpdatePageRequest)`.
- [x] Add property item retrieval if in V0 reference scope.
- [x] Add page markdown endpoint support.
- [x] Test each method path and HTTP verb.
- [x] Test create parent/property body.
- [x] Test markdown response handling.

## Milestone 10: Data Sources Service

- [x] Verify Data Sources endpoints against Notion reference.
- [x] Add `DataSources.Get(ctx, dataSourceID)`.
- [x] Add `DataSources.Query(ctx, dataSourceID, QueryDataSourceRequest)`.
- [x] Add `DataSources.Create(ctx, CreateDataSourceRequest)`.
- [x] Add `DataSources.Update(ctx, dataSourceID, UpdateDataSourceRequest)`.
- [x] Model filters and sorts flexibly with raw JSON or `any`.
- [x] Test query pagination.
- [x] Test create/update body encoding.
- [x] Test docs/examples prefer data sources.

## Milestone 11: Databases Service

- [x] Verify legacy Databases endpoints against Notion reference.
- [x] Add `Databases.Get(ctx, databaseID)`.
- [x] Add compatibility methods required by current reference.
- [x] Mark service docs as legacy/compatibility surface.
- [x] Avoid silent aliasing to Data Sources.
- [x] Test methods independently from Data Sources.

## Milestone 12: Users Service

- [x] Verify Users endpoints against Notion reference.
- [x] Add `Users.List(ctx, *Pagination)`.
- [x] Add `Users.Get(ctx, userID)`.
- [x] Add `Users.Me(ctx)`.
- [x] Test list pagination query.
- [x] Test get/me paths.

## Milestone 13: Comments Service

- [x] Verify Comments endpoints against Notion reference.
- [x] Add `Comments.List(ctx, request)`.
- [x] Add `Comments.Create(ctx, CreateCommentRequest)`.
- [x] Support page/block discussion targets from reference.
- [x] Test list query/body requirements.
- [x] Test create body encoding.

## Milestone 14: File Uploads Service

- [x] Verify File Upload endpoints against Notion reference.
- [x] Add `FileUploads.List(ctx, *Pagination)`.
- [x] Add `FileUploads.Create(ctx, CreateFileUploadRequest)`.
- [x] Add `FileUploads.Retrieve(ctx, fileUploadID)`.
- [x] Add `FileUploads.Send(ctx, fileUploadID, body)`.
- [x] Add `FileUploads.Complete(ctx, fileUploadID)`.
- [x] Decide minimal upload body abstraction from docs.
- [x] Test all paths and verbs.
- [x] Test send body handling.

## Milestone 15: Custom Emojis Service

- [x] Verify Custom Emojis endpoints against Notion reference.
- [x] Add list/retrieve methods supported by reference.
- [x] Model custom emoji response.
- [x] Test paths, query, response decoding.

## Milestone 16: Examples

- [x] Create `examples/search`.
- [x] Create `examples/page`.
- [x] Create `examples/blocks`.
- [x] Read `NOTION_TOKEN`.
- [ ] Use `client := notion.NewClient(token)`.
- [ ] Prefer `DataSources` over `Databases` where relevant.
- [ ] Ensure examples compile.

## Milestone 17: Documentation

- [ ] Add package overview.
- [ ] Document auth and default version.
- [ ] Document config options.
- [ ] Document pagination patterns.
- [ ] Document `APIError` handling.
- [ ] Document retry defaults and options.
- [ ] Document Data Sources vs Databases legacy split.
- [ ] Document no OAuth/admin/webhook support in V0.

## Milestone 18: Final Verification

- [ ] Run `go test ./...`.
- [ ] Run `go test -race ./...` if test runtime stays reasonable.
- [ ] Run example compile checks.
- [ ] Confirm no live Notion calls in default tests.
- [ ] Confirm no generated SDK code required.
- [ ] Confirm acceptance criteria from PRD pass.
