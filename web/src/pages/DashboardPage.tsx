import { useQuery } from "@tanstack/react-query";
import { api } from "@/api";
import { Card } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/status-badge";

function MetricCard({ title, value, hint }: { title: string; value: number; hint?: string }) {
  return (
    <Card>
      <div className="text-xs uppercase tracking-wide text-[#808080]">{title}</div>
      <div className="mt-2 text-3xl font-semibold tracking-[-0.04em] text-[#171717]">{value}</div>
      {hint ? <div className="mt-2 text-xs text-[#666]">{hint}</div> : null}
    </Card>
  );
}

export default function DashboardPage() {
  const { data } = useQuery({ queryKey: ["stats"], queryFn: api.statistics });
  const { data: ws } = useQuery({ queryKey: ["workflow-stats"], queryFn: api.workflowStats });
  const { data: workflows } = useQuery({ queryKey: ["workflows"], queryFn: api.listWorkflows });

  const pendingCount = Array.isArray(ws?.pendingRunIds) ? ws.pendingRunIds.length : 0;
  const processingCount = Array.isArray(ws?.processingRunIds) ? ws.processingRunIds.length : 0;
  const recent = workflows?.items?.slice(0, 6) || [];

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="证书总数" value={data?.certificateTotal ?? 0} />
        <MetricCard title="即将过期" value={data?.certificateExpiringSoon ?? 0} hint="未来 21 天内" />
        <MetricCard title="已过期" value={data?.certificateExpired ?? 0} />
        <MetricCard title="工作流总数" value={data?.workflowTotal ?? 0} hint={`启用 ${data?.workflowEnabled ?? 0} / 停用 ${data?.workflowDisabled ?? 0}`} />
      </div>

      <Card>
        <div className="flex flex-wrap items-center gap-4">
          <div>
            <div className="text-sm font-medium text-[#171717]">调度器状态</div>
            <div className="text-xs text-[#666]">并发与队列状态实时反映调度压力</div>
          </div>
          <div className="ml-auto flex items-center gap-2">
            <span className="rounded-full bg-[#f5f5f5] px-2.5 py-1 text-xs text-[#666]">并发: {ws?.concurrency ?? 0}</span>
            <span className="rounded-full bg-[var(--ds-warning-bg)] px-2.5 py-1 text-xs text-[var(--ds-warning-fg)]">等待: {pendingCount}</span>
            <span className="rounded-full bg-[var(--ds-info-bg)] px-2.5 py-1 text-xs text-[var(--ds-info-fg)]">执行中: {processingCount}</span>
          </div>
        </div>
      </Card>

      <Card>
        <div className="mb-3 text-sm font-medium text-[#171717]">最近工作流运行</div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                <th className="pb-2">工作流</th>
                <th className="pb-2">最近状态</th>
                <th className="pb-2">最近时间</th>
              </tr>
            </thead>
            <tbody>
              {recent.length ? (
                recent.map((w) => (
                  <tr key={w.id} className="border-t border-[#f1f1f1]">
                    <td className="py-3 text-[#171717]">{w.name}</td>
                    <td className="py-3">
                      <StatusBadge status={w.lastRunStatus} />
                    </td>
                    <td className="py-3 text-[#666]">{w.lastRunTime ? new Date(w.lastRunTime).toLocaleString() : "-"}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td className="py-3 text-[#808080]" colSpan={3}>
                    暂无工作流数据
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
