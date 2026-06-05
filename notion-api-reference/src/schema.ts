export function inferSchema(value: unknown): Record<string, unknown> {
  if (value === null) return { type: "null", "x-inferred": true };
  if (Array.isArray(value)) {
    return { type: "array", items: value.length ? mergeSchemas(value.map(inferSchema)) : {}, "x-inferred": true };
  }
  const type = typeof value;
  if (type === "string" || type === "number" || type === "boolean") return { type: type === "number" ? "number" : type, "x-inferred": true };
  if (type === "object") {
    const properties: Record<string, unknown> = {};
    for (const [key, child] of Object.entries(value as Record<string, unknown>)) properties[key] = inferSchema(child);
    return { type: "object", properties, additionalProperties: true, "x-inferred": true };
  }
  return { "x-inferred": true };
}

export function mergeSchemas(schemas: Record<string, unknown>[]): Record<string, unknown> {
  if (schemas.length === 0) return {};
  const types = new Set(schemas.map((schema) => schema.type));
  if (types.size === 1) return schemas[0]!;
  return { oneOf: schemas, "x-inferred": true };
}
