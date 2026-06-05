import { parseArgs } from "./config";
import { build, scrape } from "./scraper";

try {
  const { command, config } = parseArgs(Bun.argv.slice(2));
  if (command === "scrape") await scrape(config);
  else await build(config.outDir);
  console.log(`${command} done → ${config.outDir}`);
} catch (error) {
  console.error(error instanceof Error ? error.message : error);
  process.exit(1);
}
