import { Badge } from "@/components/ui/badge";
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
  pending: { label: "排队中", cls: "border-transparent bg-secondary text-secondary-foreground" },
  processing: { label: "运行中", cls: "border-transparent bg-blue-100 text-blue-700" },
  running: { label: "运行中", cls: "border-transparent bg-blue-100 text-blue-700" },
  succeeded: { label: "成功", cls: "border-transparent bg-emerald-100 text-emerald-700" },
  failed: { label: "失败", cls: "border-transparent bg-rose-100 text-rose-700" },
  canceled: { label: "已取消", cls: "border-transparent bg-slate-100 text-slate-700" },
  skipped: { label: "已跳过", cls: "border-transparent bg-slate-100 text-slate-700" },
  unknown: { label: "-", cls: "border-transparent bg-slate-100 text-slate-700" },
};

export function StatusBadge({ status, className }: { status?: string; className?: string }) {
  const normalized = normalize(status);
  const meta = statusMap[normalized];
  return <Badge className={cn("px-2.5 py-1 text-xs font-medium", meta.cls, className)}>{meta.label}</Badge>;
}
