import type {
  Boilerplate,
  BoilerplateFilters,
  BoilerplateInput,
  PaginatedResponse,
} from "./types";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (!res.ok) {
    let message = res.statusText;
    try {
      const body = await res.json();
      message = body.error ?? message;
    } catch {
      // response had no JSON body
    }
    throw new ApiError(res.status, message);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json() as Promise<T>;
}

function buildQuery(params: object): string {
  const search = new URLSearchParams();
  Object.entries(params as Record<string, unknown>).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== "") {
      search.set(key, String(value));
    }
  });
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export const boilerplateApi = {
  list: (filters: BoilerplateFilters = {}) =>
    request<PaginatedResponse<Boilerplate>>(`/boilerplate${buildQuery(filters)}`),
  get: (id: number) => request<Boilerplate>(`/boilerplate/${id}`),
  create: (input: BoilerplateInput) =>
    request<Boilerplate>("/boilerplate", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  update: (id: number, input: Partial<BoilerplateInput>) =>
    request<Boilerplate>(`/boilerplate/${id}`, {
      method: "PUT",
      body: JSON.stringify(input),
    }),
  remove: (id: number) =>
    request<void>(`/boilerplate/${id}`, { method: "DELETE" }),
};

export { ApiError };
