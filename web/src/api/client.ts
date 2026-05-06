import i18n from "@/i18n";
import type { ApiResp } from "@/types";

const API_BASE = "/api";

type RequestOptions = {
  attachToken?: boolean;
  handleUnauthorizedAsSessionExpired?: boolean;
};

export function getToken() {
  return localStorage.getItem("easyssl_token") || "";
}

export function setToken(token: string) {
  localStorage.setItem("easyssl_token", token);
}

export function clearToken() {
  localStorage.removeItem("easyssl_token");
  localStorage.removeItem("easyssl_role");
}

export function getRole() {
  return localStorage.getItem("easyssl_role") || "";
}

export function setRole(role: string) {
  localStorage.setItem("easyssl_role", role);
}

export async function request<T>(path: string, init?: RequestInit, options?: RequestOptions): Promise<T> {
  const attachToken = options?.attachToken ?? true;
  const handleUnauthorizedAsSessionExpired = options?.handleUnauthorizedAsSessionExpired ?? true;
  const headers = new Headers(init?.headers || {});
  headers.set("Content-Type", "application/json");
  const token = getToken();
  if (attachToken && token) headers.set("Authorization", `Bearer ${token}`);

  const resp = await fetch(`${API_BASE}${path}`, { ...init, headers });
  const data = (await resp.json()) as ApiResp<T>;
  if (data.code !== 0) {
    if (data.code === 401 && handleUnauthorizedAsSessionExpired) {
      clearToken();
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
      throw new Error(i18n.t("common.sessionExpired"));
    }
    throw new Error(data.msg || i18n.t("common.unknownError"));
  }
  return data.data;
}
