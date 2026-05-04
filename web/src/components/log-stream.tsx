import * as React from "react";
import { cn } from "@/lib/utils";

type LogLevel = "success" | "error" | "warning" | "info" | "neutral";

function inferLevel(eventType: string): LogLevel {
  const t = eventType.toLowerCase();
  if (t === "started" || t === "succeeded" || t === "completed" || t === "success") return "success";
  if (t === "failed" || t === "error" || t === "failure") return "error";
  if (t === "warn" || t === "warning") return "warning";
  if (t === "log" || t === "info" || t === "debug") return "info";
  return "neutral";
}

const levelStyles: Record<LogLevel, { bg: string; text: string; dot: string }> = {
  success: { bg: "bg-[#e9f9ee]", text: "text-[#116329]", dot: "bg-[#116329]" },
  error: { bg: "bg-[#ffecec]", text: "text-[#a01616]", dot: "bg-[#a01616]" },
  warning: { bg: "bg-[#fff4e8]", text: "text-[#8f4d00]", dot: "bg-[#8f4d00]" },
  info: { bg: "bg-[#ebf5ff]", text: "text-[#0068d6]", dot: "bg-[#0068d6]" },
  neutral: { bg: "bg-[#f5f5f5]", text: "text-[#666]", dot: "bg-[#999]" },
};

function formatLogTime(v?: string) {
  if (!v) return "--:--:--";
  const d = new Date(v);
  if (Number.isNaN(d.getTime())) return "--:--:--";
  const hh = String(d.getHours()).padStart(2, "0");
  const mm = String(d.getMinutes()).padStart(2, "0");
  const ss = String(d.getSeconds()).padStart(2, "0");
  const ms = String(d.getMilliseconds()).padStart(3, "0");
  return `${hh}:${mm}:${ss}.${ms}`;
}

export interface LogStreamItem {
  id: string;
  createdAt: string;
  eventType: string;
  message: string;
  nodeId?: string;
}

export interface LogStreamProps {
  items: LogStreamItem[];
  className?: string;
  emptyText?: string;
  loading?: boolean;
}

export const LogStream = React.forwardRef<HTMLDivElement, LogStreamProps>(
  ({ items, className, emptyText = "暂无日志", loading = false }, ref) => {
    const scrollRef = React.useRef<HTMLDivElement>(null);
    const bottomRef = React.useRef<HTMLDivElement>(null);

    React.useImperativeHandle(ref, () => scrollRef.current as HTMLDivElement);

    React.useEffect(() => {
      if (bottomRef.current) {
        bottomRef.current.scrollIntoView({ behavior: "smooth", block: "end" });
      }
    }, [items.length]);

    if (loading && items.length === 0) {
      return (
        <div className={cn("flex items-center justify-center py-10 text-sm text-[#808080]", className)}>
          加载中...
        </div>
      );
    }

    if (items.length === 0) {
      return (
        <div className={cn("flex items-center justify-center py-10 text-sm text-[#808080]", className)}>
          {emptyText}
        </div>
      );
    }

    return (
      <div
        ref={scrollRef}
        className={cn("max-h-[420px] overflow-auto ds-scrollbar font-mono text-[13px] leading-[1.6]", className)}
      >
        {items.map((item) => {
          const level = inferLevel(item.eventType);
          const style = levelStyles[level];
          return (
            <div
              key={item.id}
              className="group flex gap-3 px-4 py-2.5 hover:bg-[#fafafa]"
            >
              <div className="shrink-0 pt-[3px]">
                <span className="inline-block text-[11px] text-[#999] tabular-nums">
                  {formatLogTime(item.createdAt)}
                </span>
              </div>

              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span
                    className={cn(
                      "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide",
                      style.bg,
                      style.text
                    )}
                  >
                    <span className={cn("h-1 w-1 rounded-full", style.dot)} />
                    {item.eventType}
                  </span>
                  {item.nodeId ? (
                    <span className="text-[11px] text-[#b3b3b3]">{item.nodeId}</span>
                  ) : null}
                </div>
                <div className="mt-0.5 whitespace-pre-wrap text-[13px] text-[#171717]">
                  {item.message || "-"}
                </div>
              </div>
            </div>
          );
        })}
        <div ref={bottomRef} />
      </div>
    );
  }
);
LogStream.displayName = "LogStream";
