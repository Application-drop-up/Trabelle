import { z, ZodError } from "zod";

const apiErrorSchema = z.object({ message: z.string() });

export const voidSchema = z.void();

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

function parseOrThrow<T extends z.ZodTypeAny>(schema: T, path: string, data: unknown): z.infer<T> {
  try {
    return schema.parse(data);
  } catch (err) {
    if (err instanceof ZodError) {
      throw new Error(`Unexpected response from server (${path})`, { cause: err });
    }
    throw err;
  }
}

async function request<T extends z.ZodTypeAny>(
  schema: T,
  path: string,
  options?: RequestInit,
): Promise<z.infer<T>> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (res.status === 204) return parseOrThrow(schema, path, undefined);

  if (!res.ok) {
    const json = await res.json().catch(() => null);
    const parsedError = apiErrorSchema.safeParse(json);
    throw new Error(parsedError.success ? parsedError.data.message : res.statusText);
  }

  const json = await res.json();
  return parseOrThrow(schema, path, json);
}

export const apiClient = {
  get: <T extends z.ZodTypeAny>(schema: T, path: string) => request(schema, path),
  post: <T extends z.ZodTypeAny>(schema: T, path: string, body: unknown) =>
    request(schema, path, { method: "POST", body: JSON.stringify(body) }),
  patch: <T extends z.ZodTypeAny>(schema: T, path: string, body: unknown) =>
    request(schema, path, { method: "PATCH", body: JSON.stringify(body) }),
  delete: (path: string) => request(voidSchema, path, { method: "DELETE" }),
};
