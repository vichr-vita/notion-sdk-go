export type HttpMethod = "GET" | "POST" | "PATCH" | "PUT" | "DELETE";
export type Confidence = "high" | "medium" | "low";

export interface CrawlConfig {
  seed: string;
  outDir: string;
  refresh: boolean;
  delayMinMs: number;
  delayMaxMs: number;
  maxPages: number;
  ignoreRobots: boolean;
  userAgent: string;
}

export interface PageRecord {
  url: string;
  cachePath: string;
  status: number;
  fetchedAt: string;
  title?: string;
  links: string[];
}

export interface Manifest {
  seed: string;
  scrapedAt: string;
  userAgent: string;
  pages: PageRecord[];
  stoppedReason?: string;
}

export interface ParameterSpec {
  name: string;
  in: "path" | "query" | "header";
  required?: boolean;
  description?: string;
  schema?: Record<string, unknown>;
  confidence?: Confidence;
  sourceUrl?: string;
}

export interface EndpointSpec {
  id: string;
  method: HttpMethod;
  path: string;
  title: string;
  summary?: string;
  sourceUrl: string;
  scrapedAt: string;
  confidence: Confidence;
  parameters: ParameterSpec[];
  requestExamples: unknown[];
  responseExamples: unknown[];
  examples: string[];
  sections: Record<string, string>;
  warnings: string[];
}

export interface AgentSpec {
  title: string;
  generatedAt: string;
  source: string;
  auth: {
    type: "bearer";
    requiredHeaders: string[];
  };
  endpoints: EndpointSpec[];
  notes: string[];
}
