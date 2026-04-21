import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Table } from "@/components/ui/table";

export default function AccessesPage() {
  const qc = useQueryClient();
  const { data } = useQuery({ queryKey: ["accesses"], queryFn: api.listAccesses });
  const [name, setName] = useState("");
  const [provider, setProvider] = useState("cloudflare");

  const save = useMutation({
    mutationFn: () => api.saveAccess({ name, provider, config: {} }),
    onSuccess: () => { setName(""); qc.invalidateQueries({ queryKey: ["accesses"] }); },
  });

  return (
    <Card>
      <div className="mb-4 flex gap-2">
        <Input placeholder="授权名称" value={name} onChange={(e) => setName(e.target.value)} />
        <Input placeholder="provider" value={provider} onChange={(e) => setProvider(e.target.value)} />
        <Button onClick={() => save.mutate()}>新增</Button>
      </div>
      <Table>
        <thead><tr><th className="text-left">名称</th><th className="text-left">Provider</th><th></th></tr></thead>
        <tbody>
          {data?.items.map((x) => (
            <tr key={x.id} className="border-t"><td className="py-2">{x.name}</td><td>{x.provider}</td><td className="text-right"><Button variant="ghost" onClick={async () => { await api.deleteAccess(x.id!); qc.invalidateQueries({ queryKey: ["accesses"] }); }}>删除</Button></td></tr>
          ))}
        </tbody>
      </Table>
    </Card>
  );
}
