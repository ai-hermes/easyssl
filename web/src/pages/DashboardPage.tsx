import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { api } from "@/api";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/status-badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { formatTime } from "@/lib/time";

function MetricCard({ title, value, hint }: { title: string; value: number; hint?: string }) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wide text-[#808080]">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-3xl font-semibold tracking-[-0.04em] text-[#171717]">{value}</div>
        {hint ? <div className="mt-2 text-xs text-[#666]">{hint}</div> : null}
      </CardContent>
    </Card>
  );
}

export default function DashboardPage() {
  const { t } = useTranslation();
  const { data } = useQuery({ queryKey: ["stats"], queryFn: api.statistics });
  const { data: ws } = useQuery({ queryKey: ["workflow-stats"], queryFn: api.workflowStats });
  const { data: workflows } = useQuery({ queryKey: ["workflows"], queryFn: api.listWorkflows });

  const pendingCount = Array.isArray(ws?.pendingRunIds) ? ws.pendingRunIds.length : 0;
  const processingCount = Array.isArray(ws?.processingRunIds) ? ws.processingRunIds.length : 0;
  const recent = workflows?.items?.slice(0, 6) || [];

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard title={t("dashboard.totalCertificates")} value={data?.certificateTotal ?? 0} />
        <MetricCard title={t("dashboard.expiringSoon")} value={data?.certificateExpiringSoon ?? 0} hint={t("dashboard.expiringSoonHint")} />
        <MetricCard title={t("dashboard.expired")} value={data?.certificateExpired ?? 0} />
        <MetricCard title={t("dashboard.totalWorkflows")} value={data?.workflowTotal ?? 0} hint={t("dashboard.workflowHint", { enabled: data?.workflowEnabled ?? 0, disabled: data?.workflowDisabled ?? 0 })} />
      </div>

      <Card>
        <CardContent className="flex flex-wrap items-center gap-4 p-4">
          <div>
            <div className="text-sm font-medium text-[#171717]">{t("dashboard.schedulerStatus")}</div>
            <div className="text-xs text-[#666]">{t("dashboard.schedulerDescription")}</div>
          </div>
          <div className="ml-auto flex items-center gap-2">
            <Badge variant="outline" className="border-transparent bg-[#f5f5f5] text-[#666]">{t("dashboard.concurrency")}: {ws?.concurrency ?? 0}</Badge>
            <Badge className="border-transparent bg-[var(--ds-warning-bg)] text-[var(--ds-warning-fg)]">{t("dashboard.pending")}: {pendingCount}</Badge>
            <Badge className="border-transparent bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]">{t("dashboard.processing")}: {processingCount}</Badge>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm">{t("dashboard.recentWorkflows")}</CardTitle>
        </CardHeader>
        <CardContent className="pt-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("dashboard.columns.workflow")}</TableHead>
                <TableHead>{t("dashboard.columns.lastStatus")}</TableHead>
                <TableHead>{t("dashboard.columns.lastTime")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {recent.length ? (
                recent.map((w) => (
                  <TableRow key={w.id}>
                    <TableCell>{w.name}</TableCell>
                    <TableCell>
                      <StatusBadge status={w.lastRunStatus} />
                    </TableCell>
                    <TableCell className="text-[#666]">{w.lastRunTime ? formatTime(w.lastRunTime) : "-"}</TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell className="text-[#808080]" colSpan={3}>
                    {t("dashboard.noWorkflowData")}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
