import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/status-badge";
import { useToast } from "@/components/ui/toast";

export default function CertificatesPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const { data } = useQuery({ queryKey: ["certificates"], queryFn: api.listCertificates });

  return (
    <Card>
      <div className="mb-3 text-sm font-medium">证书列表</div>
      <div className="overflow-x-auto ds-scrollbar">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
              <th className="pb-2">SAN</th>
              <th className="pb-2">算法</th>
              <th className="pb-2">到期时间</th>
              <th className="pb-2">状态</th>
              <th className="pb-2 text-right">动作</th>
            </tr>
          </thead>
          <tbody>
            {data?.items.map((c) => (
              <tr key={c.id} className="border-t border-[#f1f1f1]">
                <td className="py-2">{c.subjectAltNames}</td>
                <td>{c.keyAlgorithm || "-"}</td>
                <td>{c.validityNotAfter ? new Date(c.validityNotAfter).toLocaleString() : "-"}</td>
                <td>
                  <StatusBadge status={c.isRevoked ? "failed" : "succeeded"} className="min-w-[56px] justify-center" />
                </td>
                <td className="space-x-2 text-right">
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
                      toast.success("证书文件已开始下载");
                    }}
                  >
                    下载
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={async () => {
                      await api.revokeCertificate(c.id);
                      qc.invalidateQueries({ queryKey: ["certificates"] });
                      toast.info("证书已吊销");
                    }}
                  >
                    吊销
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Card>
  );
}
