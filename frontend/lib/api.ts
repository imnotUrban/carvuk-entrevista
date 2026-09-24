import type {
  Boleta,
  BoletaInput,
  BoletaListItem,
  Boilerplate,
  BoilerplateFilters,
  BoilerplateInput,
  PaginatedResponse,
  Producto,
  ProductoFilters,
  ProductoInput,
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

export const productosApi = {
  list: (filters: ProductoFilters = {}) =>
    request<PaginatedResponse<Producto>>(`/productos${buildQuery(filters)}`),
  get: (id: number) => request<Producto>(`/productos/${id}`),
  create: (input: ProductoInput) =>
    request<Producto>("/productos", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  update: (id: number, input: Partial<ProductoInput>) =>
    request<Producto>(`/productos/${id}`, {
      method: "PUT",
      body: JSON.stringify(input),
    }),
  remove: (id: number) => request<void>(`/productos/${id}`, { method: "DELETE" }),
};

export const boletasApi = {
  list: (params: { page?: number; limit?: number; order?: "asc" | "desc" } = {}) =>
    request<PaginatedResponse<BoletaListItem>>(`/boletas${buildQuery(params)}`),
  get: (id: number) => request<Boleta>(`/boletas/${id}`),
  create: (input: BoletaInput) =>
    request<Boleta>("/boletas", {
      method: "POST",
      body: JSON.stringify(input),
    }),
};

export { ApiError };
