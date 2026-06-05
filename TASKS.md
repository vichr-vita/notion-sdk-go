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
- [ ] Support `page_size`.
- [ ] Support `has_more`.
- [ ] Support `next_cursor`.
- [ ] Implement iterator-style helpers for endpoint pages.
- [ ] Optionally implement generic `All(ctx, firstPageFn)`.
- [ ] Test manual pagination query fields.
- [ ] Test iterator stops on `has_more=false`.
- [ ] Test iterator returns callback error.

## Milestone 6: Common Models and Helpers

- [ ] Define common object fields for pages, blocks, users, comments, data sources, databases, files, custom emojis.
- [ ] Use `string` IDs.
- [ ] Use pointer fields for nullable responses.
- [ ] Use `omitempty` for optional request fields.
- [ ] Preserve `json.RawMessage` on major models where useful.
- [ ] Add union-friendly discriminator fields.
- [ ] Implement `Title(string)`.
- [ ] Implement `RichText(string)`.
- [ ] Implement `PageParent(string)`.
- [ ] Implement `DataSourceParent(string)`.
- [ ] Test helper JSON output.
- [ ] Test raw JSON preservation.

## Milestone 7: Search Service

- [ ] Verify Search endpoint against Notion reference.
- [ ] Add `SearchService`.
- [ ] Add `Search.Query(ctx, SearchRequest)`.
- [ ] Model query, filter, sort, pagination fields.
- [ ] Return typed paginated search response with raw objects.
- [ ] Test request method, path, body, auth, version.
- [ ] Test response decoding.
- [ ] Test API error handling through service.

## Milestone 8: Blocks Service

- [ ] Verify Blocks endpoints against Notion reference.
- [ ] Add `Blocks.Get(ctx, blockID)`.
- [ ] Add `Blocks.Update(ctx, blockID, request)`.
- [ ] Add `Blocks.Delete(ctx, blockID)`.
- [ ] Add `Blocks.ListChildren(ctx, blockID, *Pagination)`.
- [ ] Add `Blocks.AppendChildren(ctx, blockID, request)`.
- [ ] Add `Blocks.ForEachChild(ctx, blockID, *Pagination, func(Block) error)`.
- [ ] Test each method path and HTTP verb.
- [ ] Test children pagination.
- [ ] Test update and append body encoding.

## Milestone 9: Pages Service

- [ ] Verify Pages endpoints against Notion reference.
- [ ] Add `Pages.Get(ctx, pageID)`.
- [ ] Add `Pages.Create(ctx, CreatePageRequest)`.
- [ ] Add `Pages.Update(ctx, pageID, UpdatePageRequest)`.
- [ ] Add property item retrieval if in V0 reference scope.
- [ ] Add page markdown endpoint support.
- [ ] Test each method path and HTTP verb.
- [ ] Test create parent/property body.
- [ ] Test markdown response handling.

## Milestone 10: Data Sources Service

- [ ] Verify Data Sources endpoints against Notion reference.
- [ ] Add `DataSources.Get(ctx, dataSourceID)`.
- [ ] Add `DataSources.Query(ctx, dataSourceID, QueryDataSourceRequest)`.
- [ ] Add `DataSources.Create(ctx, CreateDataSourceRequest)`.
- [ ] Add `DataSources.Update(ctx, dataSourceID, UpdateDataSourceRequest)`.
- [ ] Model filters and sorts flexibly with raw JSON or `any`.
- [ ] Test query pagination.
- [ ] Test create/update body encoding.
- [ ] Test docs/examples prefer data sources.

## Milestone 11: Databases Service

- [ ] Verify legacy Databases endpoints against Notion reference.
- [ ] Add `Databases.Get(ctx, databaseID)`.
- [ ] Add compatibility methods required by current reference.
- [ ] Mark service docs as legacy/compatibility surface.
- [ ] Avoid silent aliasing to Data Sources.
- [ ] Test methods independently from Data Sources.

## Milestone 12: Users Service

- [ ] Verify Users endpoints against Notion reference.
- [ ] Add `Users.List(ctx, *Pagination)`.
- [ ] Add `Users.Get(ctx, userID)`.
- [ ] Add `Users.Me(ctx)`.
- [ ] Test list pagination query.
- [ ] Test get/me paths.

## Milestone 13: Comments Service

- [ ] Verify Comments endpoints against Notion reference.
- [ ] Add `Comments.List(ctx, request)`.
- [ ] Add `Comments.Create(ctx, CreateCommentRequest)`.
- [ ] Support page/block discussion targets from reference.
- [ ] Test list query/body requirements.
- [ ] Test create body encoding.

## Milestone 14: File Uploads Service

- [ ] Verify File Upload endpoints against Notion reference.
- [ ] Add `FileUploads.List(ctx, *Pagination)`.
- [ ] Add `FileUploads.Create(ctx, CreateFileUploadRequest)`.
- [ ] Add `FileUploads.Retrieve(ctx, fileUploadID)`.
- [ ] Add `FileUploads.Send(ctx, fileUploadID, body)`.
- [ ] Add `FileUploads.Complete(ctx, fileUploadID)`.
- [ ] Decide minimal upload body abstraction from docs.
- [ ] Test all paths and verbs.
- [ ] Test send body handling.

## Milestone 15: Custom Emojis Service

- [ ] Verify Custom Emojis endpoints against Notion reference.
- [ ] Add list/retrieve methods supported by reference.
- [ ] Model custom emoji response.
- [ ] Test paths, query, response decoding.

## Milestone 16: Examples

- [ ] Create `examples/search`.
- [ ] Create `examples/page`.
- [ ] Create `examples/blocks`.
- [ ] Read `NOTION_TOKEN`.
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
