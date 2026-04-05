import type { BaseResponse } from "./base.model";
import { snakeToCamel } from "./case";
import { getCookie } from "./cookie";

type ApiFetchOptions = {
  baseUrl?: string;
  withCredentials?: boolean;
  query?: Record<string, any>;
} & RequestInit;

export async function apiFetch<T = any>(
  url: string,
  options?: ApiFetchOptions
): Promise<BaseResponse<T>> {

  try {
    const {
      withCredentials = false,
      baseUrl = import.meta.env.VITE_API_URL,
      query,
      ...fetchOptions
    } = options || {};

    if (!baseUrl) {
      throw new Error(
        "Server VITE_API_URL is not configured. Please set VITE_API_URL environment variable."
      );
    }

    const headers: Record<string, any> = {
      ...fetchOptions?.headers,
      apikey: import.meta.env.VITE_API_KEY || "",
    };

    if (withCredentials) {
      const accessToken = getCookie("access_token");
      if (accessToken) {
        headers["Authorization"] = `Bearer ${accessToken}`;
      }
      // If no token, continue without auth (some endpoints may not require it)
    }

    let queryString = "";
    if (query && Object.keys(query).length > 0) {
      const searchParams = new URLSearchParams();
      Object.entries(query).forEach(([key, value]) => {
        if (
          value !== undefined &&
          value !== null &&
          value !== "" &&
          (Array.isArray(value) ? value.length > 0 : true)
        ) {
          if (Array.isArray(value)) {
            value.forEach((v) => searchParams.append(key, String(v)));
          } else {
            searchParams.append(key, String(value));
          }
        }
      });

      queryString = `?${searchParams.toString()}`;
    }

    const fullUrl = `${baseUrl}${url}${queryString}`;

    const response = await fetch(fullUrl, {
      method: fetchOptions.method || "GET",
      ...fetchOptions,
      headers
    });

    if (!response.ok) {
      let message = "Unknown error";
      try {
        const errorData = await response.json();
        message = errorData.message || message;
      } catch (_) { }
      return {
        data: null,
        error: {
          code: response.status,
          message,
        },
        meta: undefined,
      } as BaseResponse<T>;
    }

    let rawData: any;
    const contentType = response.headers.get("content-type");

    if (contentType && contentType.includes("application/json")) {
      const text = await response.text();
      rawData = text ? JSON.parse(text) : {};
    } else {
      rawData = {};
    }

    const data = snakeToCamel<any>(rawData);

    if (data.pagination || data.meta) {
      return {
        data: data.data ? (data.data as T) : ([] as T),
        error: null,
        meta: {
          pagination: data.pagination,
          ...data.meta,
        },
      } as BaseResponse<T>;
    }
    return {
      data: data.data !== undefined ? (data.data as T) : (data as T),
      error: null,
    } as BaseResponse<T>;
  } catch (error: any) {
    console.error(error);
    return {
      data: null,
      error: {
        code: error.code || 500,
        message: error.message || "Unknown error",
      },
      meta: undefined,
    } as BaseResponse<T>;
  }
}