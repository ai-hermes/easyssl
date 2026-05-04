import { useTranslation } from "react-i18next";
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

const statusCls: Record<StatusType, string> = {
  pending: "border-transparent bg-secondary text-secondary-foreground",
  processing: "border-transparent bg-blue-100 text-blue-700",
  running: "border-transparent bg-blue-100 text-blue-700",
  succeeded: "border-transparent bg-emerald-100 text-emerald-700",
  failed: "border-transparent bg-rose-100 text-rose-700",
  canceled: "border-transparent bg-slate-100 text-slate-700",
  skipped: "border-transparent bg-slate-100 text-slate-700",
  unknown: "border-transparent bg-slate-100 text-slate-700",
};

export function StatusBadge({ status, className }: { status?: string; className?: string }) {
  const { t } = useTranslation();
  const normalized = normalize(status);
  const label = t(`status.${normalized}`);
  return <Badge className={cn("px-2.5 py-1 text-xs font-medium", statusCls[normalized], className)}>{label}</Badge>;
}
