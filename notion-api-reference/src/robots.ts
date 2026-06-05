import type { CrawlConfig } from "./types";

export async function robotsAllows(config: CrawlConfig): Promise<boolean> {
  if (config.ignoreRobots) return true;
  const seed = new URL(config.seed);
  const robotsUrl = `${seed.origin}/robots.txt`;
  const response = await fetch(robotsUrl, { headers: { "user-agent": config.userAgent } });
  if (!response.ok) return true;
  const text = await response.text();
  return allows(text, seed.pathname, config.userAgent);
}

export function allows(robotsText: string, path: string, userAgent: string): boolean {
  const groups: { agents: string[]; rules: { allow: boolean; path: string }[] }[] = [];
  let current: { agents: string[]; rules: { allow: boolean; path: string }[] } | undefined;

  for (const rawLine of robotsText.split(/\r?\n/)) {
    const line = rawLine.replace(/#.*/, "").trim();
    if (!line.includes(":")) continue;
    const [fieldRaw, ...rest] = line.split(":");
    const field = fieldRaw!.trim().toLowerCase();
    const value = rest.join(":").trim();
    if (field === "user-agent") {
      current = { agents: [value.toLowerCase()], rules: [] };
      groups.push(current);
    } else if ((field === "allow" || field === "disallow") && current) {
      current.rules.push({ allow: field === "allow", path: value });
    }
  }

  const ua = userAgent.toLowerCase();
  const relevant = groups.filter((group) => group.agents.some((agent) => agent === "*" || ua.includes(agent)));
  const matchingRules = relevant.flatMap((group) => group.rules).filter((rule) => rule.path && path.startsWith(rule.path));
  if (matchingRules.length === 0) return true;
  matchingRules.sort((a, b) => b.path.length - a.path.length);
  return matchingRules[0]!.allow;
}
