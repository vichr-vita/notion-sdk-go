import { expect, test } from "bun:test";
import { parseEndpoint, parseLinks } from "../src/parser";
import { allows } from "../src/robots";
import { inferSchema } from "../src/schema";

const html = `
<html><body>
<h1>Retrieve a page</h1>
<a href="/reference/post-page">next</a>
<pre>GET /v1/pages/{page_id}</pre>
<p>Gets a Notion page.</p>
<table><tr><th>page_id required</th><td>Identifier for a Notion page</td></tr></table>
<pre>{"object":"page","id":"abc","properties":{}}</pre>
</body></html>`;

test("parse links within reference", () => {
  expect(parseLinks(html, "https://developers.notion.com/reference/intro")).toEqual([
    "https://developers.notion.com/reference/post-page",
  ]);
});

test("parse endpoint method path params examples", () => {
  const endpoint = parseEndpoint(html, "https://developers.notion.com/reference/retrieve-a-page", "now");
  expect(endpoint?.method).toBe("GET");
  expect(endpoint?.path).toBe("/v1/pages/{page_id}");
  expect(endpoint?.parameters[0]?.name).toBe("page_id");
  expect(endpoint?.responseExamples.length).toBe(1);
});

test("robots longest matching rule wins", () => {
  expect(allows("User-agent: *\nDisallow: /reference\nAllow: /reference/intro", "/reference/intro", "x")).toBe(true);
  expect(allows("User-agent: *\nDisallow: /reference", "/reference/foo", "x")).toBe(false);
});

test("infer schema from object", () => {
  expect(inferSchema({ ok: true, tags: ["a"] })).toMatchObject({
    type: "object",
    properties: { ok: { type: "boolean" }, tags: { type: "array" } },
  });
});
