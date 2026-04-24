export interface ApiError {
  code: number;
  message: string;
}

export interface BaseResponse<T> {
  data: T | null;
  error: ApiError | null;
  meta?: {
    limit: number,
    lastKey: string,
    hasMore: boolean
  };
}