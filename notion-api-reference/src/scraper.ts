import { join, relative } from "node:path";
import { readFile } from "node:fs/promises";
import type { AgentSpec, CrawlConfig, EndpointSpec, Manifest, PageRecord } from "./types";
import { cachePath, ensureDir, readTextIfExists, writeJson, writeText } from "./fs";
import { parseEndpoint, parseLinks, parseTitle } from "./parser";
import { robotsAllows } from "./robots";
import { buildOpenApi } from "./openapi";

export async function scrape(config: CrawlConfig): Promise<void> {
  await ensureDir(config.outDir);
  if (!(await robotsAllows(config))) {
    await writeJson(join(config.outDir, "manifest.json"), emptyManifest(config, "robots disallow seed path"));
    throw new Error("robots disallow seed path; pass --ignore-robots only if user explicitly approves");
  }

  const scrapedAt = new Date().toISOString();
  const queue = [config.seed];
  const seen = new Set<string>();
  const pages: PageRecord[] = [];
  let stoppedReason: string | undefined;

  while (queue.length && pages.length < config.maxPages) {
    const url = queue.shift()!;
    if (seen.has(url)) continue;
    seen.add(url);
    const page = await fetchPage(config, url);
    pages.push(page);
    for (const link of page.links) if (!seen.has(link) && !queue.includes(link)) queue.push(link);
    if (page.status === 429 || page.status === 403) {
      stoppedReason = `stopped on HTTP ${page.status}`;
      break;
    }
    if (queue.length) await sleep(jitter(config.delayMinMs, config.delayMaxMs));
  }
  if (queue.length && pages.length >= config.maxPages) stoppedReason = `max pages reached: ${config.maxPages}`;

  const manifest: Manifest = { seed: config.seed, scrapedAt, userAgent: config.userAgent, pages, stoppedReason };
  await writeJson(join(config.outDir, "manifest.json"), manifest);
  await build(config.outDir);
}

export async function build(outDir: string): Promise<void> {
  const manifest = JSON.parse(await readFile(join(outDir, "manifest.json"), "utf8")) as Manifest;
  const endpoints: EndpointSpec[] = [];
  for (const page of manifest.pages) {
    if (page.status < 200 || page.status >= 300) continue;
    const html = await readFile(page.cachePath, "utf8");
    const endpoint = parseEndpoint(html, page.url, manifest.scrapedAt);
    if (endpoint) endpoints.push(endpoint);
  }
  endpoints.sort((a, b) => a.path.localeCompare(b.path) || a.method.localeCompare(b.method));

  const agent: AgentSpec = {
    title: "Notion API Reference Agent Spec (Scraped Best-Effort)",
    generatedAt: new Date().toISOString(),
    source: manifest.seed,
    auth: { type: "bearer", requiredHeaders: ["Authorization", "Notion-Version"] },
    endpoints,
    notes: [
      "Static HTML scrape only. No browser rendering used.",
      "Schemas inferred from examples/tables where possible. Verify before SDK generation.",
      "Unknown object shapes allow additionalProperties.",
    ],
  };

  await writeJson(join(outDir, "intermediate", "endpoints.json"), endpoints);
  await writeJson(join(outDir, "agent-spec.json"), agent);
  await writeJson(join(outDir, "openapi.json"), buildOpenApi(agent));
  await writeReport(outDir, manifest, endpoints);
}

async function fetchPage(config: CrawlConfig, url: string): Promise<PageRecord> {
  const path = cachePath(config.outDir, url);
  const cached = config.refresh ? undefined : await readTextIfExists(path);
  if (cached) {
    return { url, cachePath: path, status: 200, fetchedAt: "cached", title: parseTitle(cached), links: parseLinks(cached, url) };
  }
  const response = await fetch(url, { headers: { "user-agent": config.userAgent, accept: "text/html,application/xhtml+xml" } });
  const text = await response.text();
  await writeText(path, text);
  return { url, cachePath: path, status: response.status, fetchedAt: new Date().toISOString(), title: parseTitle(text), links: response.ok ? parseLinks(text, url) : [] };
}

async function writeReport(outDir: string, manifest: Manifest, endpoints: EndpointSpec[]): Promise<void> {
  const lines = [
    "# Scrape Report",
    "",
    `Seed: ${manifest.seed}`,
    `Pages: ${manifest.pages.length}`,
    `Endpoints extracted: ${endpoints.length}`,
    manifest.stoppedReason ? `Stopped: ${manifest.stoppedReason}` : "Stopped: complete/max queue empty",
    "",
    "## Low confidence endpoints",
    ...endpoints.filter((endpoint) => endpoint.confidence !== "high").map((endpoint) => `- ${endpoint.method} ${endpoint.path} (${endpoint.confidence}) ← ${endpoint.sourceUrl}`),
    "",
    "## Pages without endpoints",
    ...manifest.pages.filter((page) => !endpoints.some((endpoint) => endpoint.sourceUrl === page.url)).map((page) => `- ${page.title ?? "untitled"} ← ${page.url}`),
  ];
  await writeText(join(outDir, "report.md"), `${lines.join("\n")}\n`);
}

function emptyManifest(config: CrawlConfig, stoppedReason: string): Manifest {
  return { seed: config.seed, scrapedAt: new Date().toISOString(), userAgent: config.userAgent, pages: [], stoppedReason };
}

function jitter(min: number, max: number): number {
  return min + Math.floor(Math.random() * Math.max(1, max - min));
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
