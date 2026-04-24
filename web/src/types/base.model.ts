export interface ApiError {
  code: number;
  message: string;
}

export interface BaseResponse<T> {
  data: T | null;
  error: ApiError | null;
  meta?: Record<string, unknown>;
}

export interface CursorResponse<T> {
  data: T[];
  cursor?: string;
}