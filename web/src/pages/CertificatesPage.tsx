import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Table } from "@/components/ui/table";

export default function CertificatesPage() {
  const qc = useQueryClient();
  const { data } = useQuery({ queryKey: ["certificates"], queryFn: api.listCertificates });

  return (
    <Card>
      <Table>
        <thead><tr><th className="text-left">SAN</th><th>算法</th><th>状态</th><th></th></tr></thead>
        <tbody>
          {data?.items.map((c) => (
            <tr key={c.id} className="border-t">
              <td className="py-2">{c.subjectAltNames}</td>
              <td>{c.keyAlgorithm}</td>
              <td>{c.isRevoked ? "已吊销" : "有效"}</td>
              <td className="space-x-2 text-right">
                <Button size="sm" onClick={async () => {
                  const r = await api.downloadCertificate(c.id, "PEM");
                  const blob = new Blob([r.fileBytes], { type: "text/plain" });
                  const a = document.createElement("a");
                  a.href = URL.createObjectURL(blob);
                  a.download = `${c.id}.pem`;
                  a.click();
                }}>下载</Button>
                <Button size="sm" variant="outline" onClick={async () => { await api.revokeCertificate(c.id); qc.invalidateQueries({ queryKey: ["certificates"] }); }}>吊销</Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Card>
  );
}
