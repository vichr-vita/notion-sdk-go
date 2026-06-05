import type { CrawlConfig } from "./types";

const defaults: CrawlConfig = {
  seed: "https://developers.notion.com/reference/intro",
  outDir: "data",
  refresh: false,
  delayMinMs: 2000,
  delayMaxMs: 4000,
  maxPages: 300,
  ignoreRobots: false,
  userAgent: "notion-api-reference-scraper/0.1 (educational; respectful cache)",
};

export function parseArgs(argv: string[]): { command: "scrape" | "build"; config: CrawlConfig } {
  const [commandRaw = "scrape", ...args] = argv;
  if (commandRaw !== "scrape" && commandRaw !== "build") {
    throw new Error(`unknown command: ${commandRaw}`);
  }

  const config = { ...defaults };
  for (let index = 0; index < args.length; index += 1) {
    const arg = args[index];
    if (arg === "--refresh") config.refresh = true;
    else if (arg === "--ignore-robots") config.ignoreRobots = true;
    else if (arg === "--seed") config.seed = requireValue(args, ++index, arg);
    else if (arg === "--out") config.outDir = requireValue(args, ++index, arg);
    else if (arg === "--delay-min-ms") config.delayMinMs = toInt(requireValue(args, ++index, arg), arg);
    else if (arg === "--delay-max-ms") config.delayMaxMs = toInt(requireValue(args, ++index, arg), arg);
    else if (arg === "--max-pages") config.maxPages = toInt(requireValue(args, ++index, arg), arg);
    else throw new Error(`unknown flag: ${arg}`);
  }
  return { command: commandRaw, config };
}

function requireValue(args: string[], index: number, flag: string): string {
  const value = args[index];
  if (!value) throw new Error(`${flag} needs value`);
  return value;
}

function toInt(value: string, flag: string): number {
  const parsed = Number.parseInt(value, 10);
  if (!Number.isFinite(parsed)) throw new Error(`${flag} needs integer`);
  return parsed;
}
