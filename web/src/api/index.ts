import { request } from "./client";
import type { APIKeyItem, Access, Certificate, ProviderDefinition, User, Workflow, WorkflowRun, WorkflowRunEvent, WorkflowRunNode } from "@/types";

export const api = {
  login: (email: string, password: string) => request<{ token: string; user: { id: string; email: string; role: string; status: string } }>("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) }),
  register: (email: string, password: string) => request<{ id: string; email: string; role: string; status: string }>("/auth/register", { method: "POST", body: JSON.stringify({ email, password }) }),
  me: () => request<{ id: string; email: string; role: string; status: string }>("/auth/me"),
  changePassword: (password: string) => request<{}>("/auth/password", { method: "PUT", body: JSON.stringify({ password }) }),
  createAPIKey: (payload: { name: string; expiresAt?: string }) =>
    request<{ id: string; name: string; prefix: string; status: string; expiresAt?: string; createdAt: string; token: string }>("/auth/api-keys", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
  listAPIKeys: () => request<{ items: APIKeyItem[]; totalItems: number }>("/auth/api-keys"),
  revokeAPIKey: (id: string) => request<{}>(`/auth/api-keys/${id}`, { method: "DELETE" }),

  listProviders: (kind?: "access" | "dns" | "deploy") => {
    const suffix = kind ? `?kind=${encodeURIComponent(kind)}` : "";
    return request<{ items: ProviderDefinition[]; totalItems: number }>(`/providers${suffix}`);
  },

  listAccesses: () => request<{ items: Access[]; totalItems: number }>("/accesses"),
  saveAccess: (payload: Access) => request<Access>(payload.id ? `/accesses/${payload.id}` : "/accesses", { method: payload.id ? "PUT" : "POST", body: JSON.stringify(payload) }),
  deleteAccess: (id: string) => request<{}>(`/accesses/${id}`, { method: "DELETE" }),
  testAccess: (id: string) => request<{ testedAt: string }>(`/accesses/${id}/test`, { method: "POST" }),

  listWorkflows: () => request<{ items: Workflow[]; totalItems: number }>("/workflows"),
  getWorkflow: (id: string) => request<Workflow>(`/workflows/${id}`),
  saveWorkflow: (payload: Workflow) => request<Workflow>(payload.id ? `/workflows/${payload.id}` : "/workflows", { method: payload.id ? "PUT" : "POST", body: JSON.stringify(payload) }),
  deleteWorkflow: (id: string) => request<{}>(`/workflows/${id}`, { method: "DELETE" }),
  listWorkflowRuns: (id: string) => request<{ items: WorkflowRun[]; totalItems: number }>(`/workflows/${id}/runs`),
  listWorkflowRunNodes: (id: string, runId: string) => request<{ items: WorkflowRunNode[]; totalItems: number }>(`/workflows/${id}/runs/${runId}/nodes`),
  listWorkflowRunEvents: (id: string, runId: string, params?: { nodeId?: string; since?: string; limit?: number }) => {
    const query = new URLSearchParams();
    if (params?.nodeId) query.set("nodeId", params.nodeId);
    if (params?.since) query.set("since", params.since);
    if (params?.limit) query.set("limit", String(params.limit));
    const suffix = query.toString() ? `?${query.toString()}` : "";
    return request<{ items: WorkflowRunEvent[]; totalItems: number }>(`/workflows/${id}/runs/${runId}/events${suffix}`);
  },
  startWorkflowRun: (id: string) => request<{ runId: string }>(`/workflows/${id}/runs`, { method: "POST", body: JSON.stringify({ trigger: "manual" }) }),
  cancelWorkflowRun: (id: string, runId: string) => request<{}>(`/workflows/${id}/runs/${runId}/cancel`, { method: "POST" }),
  workflowStats: () => request<{ concurrency: number; pendingRunIds: string[]; processingRunIds: string[] }>("/workflows/stats"),

  listCertificates: () => request<{ items: Certificate[]; totalItems: number }>("/certificates"),
  downloadCertificate: (id: string, format = "PEM") =>
    request<{ fileName: string; fileFormat: string; mimeType: string; fileBytesBase64: string }>(`/certificates/${id}/download`, { method: "POST", body: JSON.stringify({ format }) }),
  revokeCertificate: (id: string) => request<{}>(`/certificates/${id}/revoke`, { method: "POST" }),

  statistics: () => request<{ certificateTotal: number; certificateExpiringSoon: number; certificateExpired: number; workflowTotal: number; workflowEnabled: number; workflowDisabled: number }>("/statistics"),
  testNotification: (provider: string, accessId: string) => request<{ sentAt: string }>("/notifications/test", { method: "POST", body: JSON.stringify({ provider, accessId }) }),

  listUsers: () => request<{ items: User[]; totalItems: number }>("/admin/users"),
  updateUserStatus: (id: string, status: string) => request<{}>(`/admin/users/${id}/status`, { method: "PUT", body: JSON.stringify({ status }) }),
  version: () => request<{ version: string }>("/admin/version"),
};
