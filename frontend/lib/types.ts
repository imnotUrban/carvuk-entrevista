export type BoilerplateStatus = "draft" | "active" | "archived";

export interface Boilerplate {
  id: number;
  name: string;
  code: string;
  status: BoilerplateStatus;
  quantity: number;
  amount: number;
  created_at: string;
  updated_at: string;
}

export interface Producto {
  id: number;
  nombre: string;
  precio: number;
  stock: number;
  created_at: string;
  updated_at: string;
}

export interface ProductoInput {
  nombre: string;
  precio: number;
  stock: number;
}

export interface ProductoFilters {
  nombre?: string;
  page?: number;
  limit?: number;
  sort_by?: "id" | "nombre" | "precio" | "stock" | "created_at";
  order?: "asc" | "desc";
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface BoilerplateInput {
  name: string;
  code: string;
  status: BoilerplateStatus;
  quantity: number;
  amount: number;
}

export interface BoilerplateFilters {
  name?: string;
  code?: string;
  status?: BoilerplateStatus;
  min_amount?: number;
  max_amount?: number;
  page?: number;
  limit?: number;
}

export interface BoletaItem {
  producto_id: number;
  nombre: string;
  precio_unitario: number;
  cantidad: number;
  subtotal: number;
}

export interface Boleta {
  id: number;
  compra_id: number;
  valor_bruto: number;
  impuesto: number;
  valor_neto: number;
  porcentaje_impuesto: number;
  usuario: { id: number; nombre: string; correo: string };
  items: BoletaItem[];
  created_at: string;
}

export interface BoletaInput {
  items: { producto_id: number; cantidad: number }[];
}

export interface BoletaListItem {
  id: number;
  created_at: string;
  valor_bruto: number;
  cantidad_items: number;
  usuario: { id: number; nombre: string; correo: string };
}
