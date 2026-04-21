import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Table } from "@/components/ui/table";

const starterGraph = {
  nodes: [
    { id: "start", type: "start", data: { name: "Start" } },
    { id: "apply", type: "bizApply", data: { name: "Apply", config: { domains: "example.com" } } },
    { id: "deploy", type: "bizDeploy", data: { name: "Deploy", config: { provider: "ssh" } } },
    { id: "notify", type: "bizNotify", data: { name: "Notify", config: { provider: "webhook" } } },
    { id: "end", type: "end", data: { name: "End" } }
  ]
};

export default function WorkflowsPage() {
  const qc = useQueryClient();
  const { data } = useQuery({ queryKey: ["workflows"], queryFn: api.listWorkflows });
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [graphText, setGraphText] = useState(JSON.stringify(starterGraph, null, 2));

  const save = useMutation({
    mutationFn: async () => api.saveWorkflow({
      name,
      description,
      trigger: "manual",
      triggerCron: "",
      enabled: true,
      graphDraft: JSON.parse(graphText),
      graphContent: JSON.parse(graphText),
      hasDraft: true,
      hasContent: true,
    }),
    onSuccess: () => { setName(""); setDescription(""); qc.invalidateQueries({ queryKey: ["workflows"] }); }
  });

  return (
    <Card>
      <div className="mb-4 space-y-2">
        <Input placeholder="工作流名称" value={name} onChange={(e) => setName(e.target.value)} />
        <Input placeholder="描述" value={description} onChange={(e) => setDescription(e.target.value)} />
        <Textarea value={graphText} onChange={(e) => setGraphText(e.target.value)} />
        <Button onClick={() => save.mutate()}>创建工作流</Button>
      </div>
      <Table>
        <thead><tr><th className="text-left">名称</th><th className="text-left">状态</th><th className="text-left">动作</th></tr></thead>
        <tbody>
          {data?.items.map((w) => (
            <tr key={w.id} className="border-t">
              <td className="py-2">{w.name}</td>
              <td>{w.enabled ? "启用" : "停用"}</td>
              <td className="space-x-2">
                <Button size="sm" onClick={async () => { await api.startWorkflowRun(w.id!); }}>运行</Button>
                <Button size="sm" variant="outline" onClick={async () => { await api.deleteWorkflow(w.id!); qc.invalidateQueries({ queryKey: ["workflows"] }); }}>删除</Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Card>
  );
}
