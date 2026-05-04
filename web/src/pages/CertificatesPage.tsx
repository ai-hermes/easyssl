import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/status-badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useToast } from "@/components/ui/toast";

export default function CertificatesPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const { t } = useTranslation();
  const { data } = useQuery({ queryKey: ["certificates"], queryFn: api.listCertificates });

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm">{t("certificates.title")}</CardTitle>
      </CardHeader>
      <CardContent className="pt-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("certificates.columns.san")}</TableHead>
              <TableHead>{t("certificates.columns.algorithm")}</TableHead>
              <TableHead>{t("certificates.columns.expiresAt")}</TableHead>
              <TableHead>{t("certificates.columns.status")}</TableHead>
              <TableHead className="text-right">{t("certificates.columns.actions")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data?.items.map((c) => (
              <TableRow key={c.id}>
                <TableCell>{c.subjectAltNames}</TableCell>
                <TableCell>{c.keyAlgorithm || "-"}</TableCell>
                <TableCell>{c.validityNotAfter ? new Date(c.validityNotAfter).toLocaleString() : "-"}</TableCell>
                <TableCell>
                  <StatusBadge status={c.isRevoked ? "failed" : "succeeded"} className="min-w-[56px] justify-center" />
                </TableCell>
                <TableCell className="space-x-2 text-right">
                  <Button
                    size="sm"
                    onClick={async () => {
                      const r = await api.downloadCertificate(c.id, "PEM");
                      const bin = atob(r.fileBytesBase64);
                      const bytes = new Uint8Array(bin.length);
                      for (let i = 0; i < bin.length; i += 1) bytes[i] = bin.charCodeAt(i);
                      const blob = new Blob([bytes], { type: r.mimeType || "application/octet-stream" });
                      const a = document.createElement("a");
                      const url = URL.createObjectURL(blob);
                      a.href = url;
                      a.download = r.fileName || `${c.id}.zip`;
                      a.click();
                      URL.revokeObjectURL(url);
                      toast.success(t("certificates.downloadStarted"));
                    }}
                  >
                    {t("common.download")}
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={async () => {
                      await api.revokeCertificate(c.id);
                      qc.invalidateQueries({ queryKey: ["certificates"] });
                      toast.info(t("certificates.revokeSuccess"));
                    }}
                  >
                    {t("common.revoke")}
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
