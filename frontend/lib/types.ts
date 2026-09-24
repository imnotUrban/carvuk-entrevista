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
