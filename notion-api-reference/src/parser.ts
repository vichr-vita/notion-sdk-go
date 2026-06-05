import * as cheerio from "cheerio";
import type { Confidence, EndpointSpec, HttpMethod, ParameterSpec } from "./types";

const methods = ["GET", "POST", "PATCH", "PUT", "DELETE"] as const;
const pathPattern = "(?:/?v\\d/|/)[A-Za-z0-9_{}./:-]+";
const methodPath = new RegExp(`\\b(GET|POST|PATCH|PUT|DELETE)\\s+(${pathPattern})`, "i");
const pathOnly = new RegExp(`(${pathPattern})`);

export function parseLinks(html: string, baseUrl: string): string[] {
  const $ = cheerio.load(html);
  const urls = new Set<string>();
  $("a[href]").each((_, element) => {
    const href = $(element).attr("href");
    if (!href) return;
    try {
      const url = new URL(href, baseUrl);
      url.hash = "";
      if (url.origin === "https://developers.notion.com" && url.pathname.startsWith("/reference/")) {
        urls.add(url.toString());
      }
    } catch {}
  });
  return [...urls].sort();
}

export function parseTitle(html: string): string | undefined {
  const $ = cheerio.load(html);
  return clean($("h1").first().text()) || clean($("title").text()) || undefined;
}

export function parseEndpoint(html: string, url: string, scrapedAt: string): EndpointSpec | undefined {
  const $ = cheerio.load(html);
  $("script,style,noscript,[data-agent-docs-index]").remove();
  const title = parseTitle(html) ?? slugTitle(url);
  const bodyText = clean($("body").text());
  const endpoint = findMethodPath($, bodyText, url, html);
  if (!endpoint) return undefined;

  const sections = collectSections($);
  const examples = collectCodeBlocks($);
  const jsonExamples = examples.map(parseJsonLoose).filter((value) => value !== undefined);
  const parameters = extractParameters($, endpoint.path, url);
  const warnings: string[] = [];
  if (endpoint.confidence !== "high") warnings.push("method/path inferred from page text; verify before SDK generation");
  if (jsonExamples.length === 0) warnings.push("no JSON examples found");

  return {
    id: `${endpoint.method} ${endpoint.path}`,
    method: endpoint.method,
    path: normalizePath(endpoint.path),
    title,
    summary: firstParagraph($),
    sourceUrl: url,
    scrapedAt,
    confidence: endpoint.confidence,
    parameters,
    requestExamples: jsonExamples.filter((value) => looksLikeRequest(value)),
    responseExamples: jsonExamples.filter((value) => !looksLikeRequest(value)),
    examples,
    sections,
    warnings,
  };
}

function findMethodPath($: cheerio.CheerioAPI, text: string, url: string, html: string): { method: HttpMethod; path: string; confidence: Confidence } | undefined {
  const embedded = html.match(/\\"method\\":\\"(get|post|patch|put|delete)\\"[\s\S]{0,2000}?\\"path\\":\\"([^\\"]+)\\"/i)
    ?? html.match(/"method":"(get|post|patch|put|delete)"[\s\S]{0,2000}?"path":"([^"]+)"/i);
  if (embedded) return { method: embedded[1]!.toUpperCase() as HttpMethod, path: unescapeJson(embedded[2]!), confidence: "high" };

  for (const code of collectCodeBlocks($)) {
    const match = code.match(methodPath);
    if (match) return { method: match[1]!.toUpperCase() as HttpMethod, path: match[2]!, confidence: "high" };
  }
  const textMatch = text.match(methodPath);
  if (textMatch) return { method: textMatch[1]!.toUpperCase() as HttpMethod, path: textMatch[2]!, confidence: "medium" };

  const lowerTitle = parseTitle($.html())?.toLowerCase() ?? "";
  const method = inferMethod(lowerTitle);
  const pathMatch = text.match(pathOnly);
  if (method && pathMatch) return { method, path: pathMatch[1]!, confidence: "low" };

  const slug = new URL(url).pathname.split("/").filter(Boolean).at(-1) ?? "";
  const slugMethod = inferMethod(slug.replace(/-/g, " "));
  if (slugMethod) return { method: slugMethod, path: `/${slug}`, confidence: "low" };
  return undefined;
}

function inferMethod(text: string): HttpMethod | undefined {
  if (/^(retrieve|get|query|list|search)\b/.test(text)) return "GET";
  if (/^(create|append|duplicate|send)\b/.test(text)) return "POST";
  if (/^(update|patch)\b/.test(text)) return "PATCH";
  if (/^(delete|archive)\b/.test(text)) return "DELETE";
  return undefined;
}

function collectSections($: cheerio.CheerioAPI): Record<string, string> {
  const sections: Record<string, string> = {};
  $("h2,h3").each((_, heading) => {
    const name = clean($(heading).text());
    if (!name) return;
    const chunks: string[] = [];
    let node = $(heading).next();
    while (node.length && !/^h[23]$/i.test(node.prop("tagName") ?? "")) {
      const text = clean(node.text());
      if (text) chunks.push(text);
      node = node.next();
    }
    if (chunks.length) sections[name] = chunks.join("\n").slice(0, 8000);
  });
  return sections;
}

export function collectCodeBlocks($: cheerio.CheerioAPI): string[] {
  const blocks: string[] = [];
  $("pre, code").each((_, element) => {
    const text = clean($(element).text());
    if (text && !blocks.includes(text)) blocks.push(text);
  });
  return blocks;
}

function extractParameters($: cheerio.CheerioAPI, path: string, sourceUrl: string): ParameterSpec[] {
  const params = new Map<string, ParameterSpec>();
  for (const match of path.matchAll(/\{([^}]+)\}/g)) {
    params.set(match[1]!, { name: match[1]!, in: "path", required: true, schema: { type: "string" }, confidence: "high", sourceUrl });
  }
  $("table tr").each((_, row) => {
    const cells = $(row).find("th,td").map((__, cell) => clean($(cell).text())).get();
    if (cells.length < 2) return;
    const name = cells[0]!.replace(/required|optional/ig, "").trim();
    if (!/^[A-Za-z_][\w.-]*$/.test(name)) return;
    const required = cells.join(" ").toLowerCase().includes("required");
    if (!params.has(name)) {
      params.set(name, { name, in: "query", required, description: cells.slice(1).join(" "), schema: { type: "string" }, confidence: "low", sourceUrl });
    }
  });
  return [...params.values()];
}

function firstParagraph($: cheerio.CheerioAPI): string | undefined {
  const text = clean($("p").first().text());
  return text || undefined;
}

function parseJsonLoose(text: string): unknown | undefined {
  const trimmed = text.trim();
  if (!trimmed.startsWith("{") && !trimmed.startsWith("[")) return undefined;
  try { return JSON.parse(trimmed); } catch { return undefined; }
}

function looksLikeRequest(value: unknown): boolean {
  if (!value || typeof value !== "object" || Array.isArray(value)) return false;
  const object = value as Record<string, unknown>;
  if (typeof object.object === "string" || typeof object.id === "string") return false;
  const keys = Object.keys(object);
  return keys.some((key) => ["parent", "properties", "children", "filter", "sorts"].includes(key));
}

function normalizePath(path: string): string {
  const cleaned = path.replace(/^v\d\//, "/v1/").replace(/\?.*$/, "");
  return cleaned.startsWith("/") ? cleaned : `/${cleaned}`;
}

function unescapeJson(value: string): string {
  return value.replace(/\\\\\//g, "/").replace(/\\\\"/g, '"');
}

function slugTitle(url: string): string {
  return (new URL(url).pathname.split("/").at(-1) ?? url).replace(/-/g, " ");
}

function clean(value: string): string {
  return value.replace(/\s+/g, " ").trim();
}
