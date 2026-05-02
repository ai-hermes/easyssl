import { cn } from "@/lib/utils";

type StatusType = "pending" | "processing" | "running" | "succeeded" | "failed" | "canceled" | "skipped" | "unknown";

function normalize(status?: string): StatusType {
  const s = (status || "").toLowerCase();
  if (s === "pending") return "pending";
  if (s === "processing" || s === "running") return "running";
  if (s === "succeeded" || s === "success") return "succeeded";
  if (s === "failed" || s === "error") return "failed";
  if (s === "canceled") return "canceled";
  if (s === "skipped") return "skipped";
  return "unknown";
}

const statusMap: Record<StatusType, { label: string; cls: string }> = {
  pending: { label: "排队中", cls: "bg-[#f5f5f5] text-[#666]" },
  processing: { label: "运行中", cls: "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]" },
  running: { label: "运行中", cls: "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]" },
  succeeded: { label: "成功", cls: "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]" },
  failed: { label: "失败", cls: "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]" },
  canceled: { label: "已取消", cls: "bg-[#f5f5f5] text-[#666]" },
  skipped: { label: "已跳过", cls: "bg-[#f8f8f8] text-[#666]" },
  unknown: { label: "-", cls: "bg-[#f5f5f5] text-[#666]" },
};

export function StatusBadge({ status, className }: { status?: string; className?: string }) {
  const normalized = normalize(status);
  const meta = statusMap[normalized];
  return <span className={cn("inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium", meta.cls, className)}>{meta.label}</span>;
}
