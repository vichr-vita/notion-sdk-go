# notion-api-reference

Respectful Bun scraper for public Notion API reference.

Outputs:

- `data/raw/*.html` cached pages
- `data/manifest.json` crawl manifest
- `data/intermediate/endpoints.json` normalized endpoint catalog
- `data/agent-spec.json` agent-facing spec
- `data/openapi.json` best-effort OpenAPI 3.1 seed
- `data/report.md` review report

Run:

```bash
bun install
bun run scrape
bun run build:spec
bun test
```

Defaults:

- seed: `https://developers.notion.com/reference/intro`
- one request at a time
- 2–4s jitter between pages
- max 300 pages
- cache unless `--refresh`
- respect robots.txt unless `--ignore-robots`

Flags:

```bash
bun run scrape -- --refresh
bun run scrape -- --delay-min-ms 2000 --delay-max-ms 4000 --max-pages 300
bun run scrape -- --out data --seed https://developers.notion.com/reference/intro
```

Generated OpenAPI = scraped best-effort. Verify low-confidence endpoints and inferred schemas before SDK generation.
