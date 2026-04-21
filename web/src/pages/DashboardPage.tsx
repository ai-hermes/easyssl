import { useQuery } from "@tanstack/react-query";
import { api } from "@/api";
import { Card } from "@/components/ui/card";

export default function DashboardPage() {
  const { data } = useQuery({ queryKey: ["stats"], queryFn: api.statistics });
  const { data: ws } = useQuery({ queryKey: ["workflow-stats"], queryFn: api.workflowStats });

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
      <Card><p className="text-sm text-slate-500">证书总数</p><p className="text-2xl font-semibold">{data?.certificateTotal ?? 0}</p></Card>
      <Card><p className="text-sm text-slate-500">即将过期</p><p className="text-2xl font-semibold">{data?.certificateExpiringSoon ?? 0}</p></Card>
      <Card><p className="text-sm text-slate-500">工作流总数</p><p className="text-2xl font-semibold">{data?.workflowTotal ?? 0}</p></Card>
      <Card className="md:col-span-3"><p className="text-sm text-slate-500">调度并发</p><p className="text-lg">并发: {ws?.concurrency ?? 0} / 等待: {ws?.pendingRunIds.length ?? 0} / 执行中: {ws?.processingRunIds.length ?? 0}</p></Card>
    </div>
  );
}
