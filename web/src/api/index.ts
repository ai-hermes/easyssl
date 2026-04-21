import { request } from "./client";
import type { Access, Certificate, Workflow } from "@/types";

export const api = {
  login: (email: string, password: string) => request<{ token: string; admin: { id: string; email: string } }>("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) }),
  me: () => request<{ id: string; email: string }>("/auth/me"),
  changePassword: (password: string) => request<{}>("/auth/password", { method: "PUT", body: JSON.stringify({ password }) }),

  listAccesses: () => request<{ items: Access[]; totalItems: number }>("/accesses"),
  saveAccess: (payload: Access) => request<Access>(payload.id ? `/accesses/${payload.id}` : "/accesses", { method: payload.id ? "PUT" : "POST", body: JSON.stringify(payload) }),
  deleteAccess: (id: string) => request<{}>(`/accesses/${id}`, { method: "DELETE" }),

  listWorkflows: () => request<{ items: Workflow[]; totalItems: number }>("/workflows"),
  getWorkflow: (id: string) => request<Workflow>(`/workflows/${id}`),
  saveWorkflow: (payload: Workflow) => request<Workflow>(payload.id ? `/workflows/${payload.id}` : "/workflows", { method: payload.id ? "PUT" : "POST", body: JSON.stringify(payload) }),
  deleteWorkflow: (id: string) => request<{}>(`/workflows/${id}`, { method: "DELETE" }),
  listWorkflowRuns: (id: string) => request<{ items: unknown[]; totalItems: number }>(`/workflows/${id}/runs`),
  startWorkflowRun: (id: string) => request<{ runId: string }>(`/workflows/${id}/runs`, { method: "POST", body: JSON.stringify({ trigger: "manual" }) }),
  cancelWorkflowRun: (id: string, runId: string) => request<{}>(`/workflows/${id}/runs/${runId}/cancel`, { method: "POST" }),
  workflowStats: () => request<{ concurrency: number; pendingRunIds: string[]; processingRunIds: string[] }>("/workflows/stats"),

  listCertificates: () => request<{ items: Certificate[]; totalItems: number }>("/certificates"),
  downloadCertificate: (id: string, format = "PEM") => request<{ fileBytes: string; fileFormat: string }>(`/certificates/${id}/download`, { method: "POST", body: JSON.stringify({ format }) }),
  revokeCertificate: (id: string) => request<{}>(`/certificates/${id}/revoke`, { method: "POST" }),

  statistics: () => request<{ certificateTotal: number; certificateExpiringSoon: number; certificateExpired: number; workflowTotal: number; workflowEnabled: number; workflowDisabled: number }>("/statistics"),
  testNotification: (provider: string, accessId: string) => request<{ sentAt: string }>("/notifications/test", { method: "POST", body: JSON.stringify({ provider, accessId }) }),
};
