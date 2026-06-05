import type { AgentSpec, EndpointSpec } from "./types";
import { inferSchema, mergeSchemas } from "./schema";

export function buildOpenApi(agent: AgentSpec): Record<string, unknown> {
  const paths: Record<string, Record<string, unknown>> = {};
  for (const endpoint of agent.endpoints) {
    paths[endpoint.path] ??= {};
    paths[endpoint.path]![endpoint.method.toLowerCase()] = operation(endpoint);
  }

  return {
    openapi: "3.1.0",
    info: {
      title: "Notion API Reference (Scraped Best-Effort)",
      version: "0.0.0",
      description: "Generated from public Notion API reference pages. Verify uncertain inferred schemas before SDK generation.",
    },
    servers: [{ url: "https://api.notion.com" }],
    security: [{ bearerAuth: [] }],
    paths,
    components: {
      securitySchemes: {
        bearerAuth: { type: "http", scheme: "bearer" },
      },
      parameters: {
        NotionVersion: {
          name: "Notion-Version",
          in: "header",
          required: true,
          schema: { type: "string" },
          description: "Notion API version header. Value should come from Notion docs/versioning policy.",
        },
      },
    },
    "x-generation-notes": agent.notes,
  };
}

function operation(endpoint: EndpointSpec): Record<string, unknown> {
  const requestSchemas = endpoint.requestExamples.map(inferSchema);
  const responseSchemas = endpoint.responseExamples.map(inferSchema);
  return {
    summary: endpoint.title,
    description: endpoint.summary,
    parameters: [
      { $ref: "#/components/parameters/NotionVersion" },
      ...endpoint.parameters.map((parameter) => ({
        name: parameter.name,
        in: parameter.in,
        required: parameter.in === "path" ? true : parameter.required ?? false,
        description: parameter.description,
        schema: parameter.schema ?? { type: "string" },
        "x-confidence": parameter.confidence,
        "x-source-url": parameter.sourceUrl,
      })),
    ],
    requestBody: requestSchemas.length ? {
      required: true,
      content: { "application/json": { schema: mergeSchemas(requestSchemas), examples: examples(endpoint.requestExamples) } },
    } : undefined,
    responses: {
      "200": {
        description: "Successful response. Schema inferred from docs examples where available.",
        content: responseSchemas.length ? { "application/json": { schema: mergeSchemas(responseSchemas), examples: examples(endpoint.responseExamples) } } : undefined,
      },
      default: { description: "Error response. See Notion API error docs." },
    },
    "x-source-url": endpoint.sourceUrl,
    "x-scraped-at": endpoint.scrapedAt,
    "x-confidence": endpoint.confidence,
    "x-warnings": endpoint.warnings,
  };
}

function examples(values: unknown[]): Record<string, unknown> {
  return Object.fromEntries(values.slice(0, 5).map((value, index) => [`example${index + 1}`, { value }]));
}
