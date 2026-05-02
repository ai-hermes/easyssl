import type { ApiResp } from "@/types";

const API_BASE = "/api";

export function getToken() {
  return localStorage.getItem("easyssl_token") || "";
}

export function setToken(token: string) {
  localStorage.setItem("easyssl_token", token);
}

export function clearToken() {
  localStorage.removeItem("easyssl_token");
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers || {});
  headers.set("Content-Type", "application/json");
  const token = getToken();
  if (token) headers.set("Authorization", `Bearer ${token}`);

  const resp = await fetch(`${API_BASE}${path}`, { ...init, headers });
  const data = (await resp.json()) as ApiResp<T>;
  if (data.code !== 0) {
    if (data.code === 401) {
      clearToken();
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
      throw new Error("登录状态已过期，请重新登录");
    }
    throw new Error(data.msg);
  }
  return data.data;
}
