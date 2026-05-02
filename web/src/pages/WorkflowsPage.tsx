import * as Dialog from "@radix-ui/react-dialog";
import { useEffect, useMemo, useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import ReactFlow, { Background, Controls, Edge, MiniMap, Node, ReactFlowInstance } from "reactflow";
import "reactflow/dist/style.css";
import YAML from "js-yaml";

import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { StatusBadge } from "@/components/ui/status-badge";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/toast";
import type { Workflow, WorkflowRunNode } from "@/types";

type FlowNodeData = {
  name: string;
  config?: Record<string, unknown>;
  label?: string;
};

type WorkflowGraph = {
  nodes: Array<Node<FlowNodeData>>;
  edges: Edge[];
};

type WorkflowSpecNode = {
  id: string;
  name: string;
  action: "apply" | "deploy";
  provider?: string;
  accessId?: string;
  config?: Record<string, unknown>;
};

type WorkflowSpecEdge = {
  source: string;
  target: string;
};

type WorkflowSpec = {
  version: number;
  options?: Record<string, unknown>;
  nodes: WorkflowSpecNode[];
  edges: WorkflowSpecEdge[];
};

type WorkflowTemplateKey = "apply-only" | "apply-aliyun" | "apply-qiniu" | "apply-ssh";

const TEMPLATE_OPTIONS: Array<{ value: WorkflowTemplateKey; label: string }> = [
  { value: "apply-only", label: "仅申请证书" },
  { value: "apply-aliyun", label: "申请 + 部署到 Aliyun CAS" },
  { value: "apply-qiniu", label: "申请 + 部署到 Qiniu" },
  { value: "apply-ssh", label: "申请 + 部署到 SSH 主机" },
];

function makeTemplateSpec(template: WorkflowTemplateKey): WorkflowSpec {
  const base: WorkflowSpec = {
    version: 1,
    options: { failFast: true },
    nodes: [
      {
        id: "apply-1",
        name: "Apply",
        action: "apply",
        provider: "tencentcloud",
        accessId: "",
        config: { domains: "example.com", caProvider: "letsencrypt" },
      },
    ],
    edges: [],
  };

  if (template === "apply-only") return base;

  const provider = template === "apply-aliyun" ? "aliyun-cas" : template === "apply-qiniu" ? "qiniu" : "ssh";
  const deployConfig =
    provider === "ssh"
      ? {
          certPath: "/etc/nginx/ssl/fullchain.pem",
          keyPath: "/etc/nginx/ssl/privkey.pem",
          postCommand: "nginx -s reload",
        }
      : provider === "qiniu"
        ? { certName: "", commonName: "" }
        : { region: "cn-hangzhou" };

  base.nodes.push({ id: "deploy-1", name: "Deploy", action: "deploy", provider, accessId: "", config: deployConfig });
  base.edges.push({ source: "apply-1", target: "deploy-1" });
  return base;
}

function stringify(v: unknown): string {
  if (v === null || v === undefined) return "";
  return String(v).trim();
}

function inferAction(type: string, id: string, name: string, config: Record<string, unknown>): "apply" | "deploy" | "" {
  const t = type.toLowerCase();
  const i = id.toLowerCase();
  const n = name.toLowerCase();
  const provider = stringify(config.provider).toLowerCase();
  if (t === "input" || t === "output" || t === "start" || t === "end") return "";
  if (i === "start" || i === "end" || n === "start" || n === "end") return "";
  if (t.includes("apply") || i.includes("apply") || n.includes("apply")) return "apply";
  if (t.includes("deploy") || i.includes("deploy") || n.includes("deploy")) return "deploy";
  if (provider === "aliyun-cas" || provider === "qiniu" || provider === "ssh") return "deploy";
  if (provider === "aliyun" || provider === "tencentcloud") return "apply";
  return "";
}

function pickWorkflowGraph(workflow: Workflow): unknown {
  const draft = workflow.graphDraft;
  const content = workflow.graphContent;
  const draftNodes = Array.isArray((draft as Record<string, unknown> | undefined)?.nodes) ? ((draft as Record<string, unknown>).nodes as unknown[]) : [];
  const contentNodes = Array.isArray((content as Record<string, unknown> | undefined)?.nodes) ? ((content as Record<string, unknown>).nodes as unknown[]) : [];
  if (draftNodes.length > 0) return draft;
  if (contentNodes.length > 0) return content;
  return draft ?? content ?? {};
}

function graphToSpec(raw: unknown): WorkflowSpec {
  const graph = raw && typeof raw === "object" ? (raw as Record<string, unknown>) : {};
  const rawNodes = Array.isArray(graph.nodes) ? graph.nodes : [];
  const rawEdges = Array.isArray(graph.edges) ? graph.edges : [];
  if (rawNodes.length === 0) return makeTemplateSpec("apply-only");

  const nodes = rawNodes
    .map<WorkflowSpecNode | null>((item, index) => {
      const nodeObj = item && typeof item === "object" ? (item as Record<string, unknown>) : {};
      const id = stringify(nodeObj.id) || `node-${index + 1}`;
      const type = stringify(nodeObj.type);
      const dataObj = nodeObj.data && typeof nodeObj.data === "object" ? (nodeObj.data as Record<string, unknown>) : {};
      const name = stringify(dataObj.name) || id;
      const configRaw = dataObj.config && typeof dataObj.config === "object" ? ({ ...(dataObj.config as Record<string, unknown>) } as Record<string, unknown>) : {};
      const action = inferAction(type, id, name, configRaw);
      if (!action) return null;
      const provider = stringify(configRaw.provider);
      const accessId = stringify(configRaw.accessId || configRaw.accessID || configRaw.access_id || configRaw.providerAccessId);
      delete configRaw.provider;
      delete configRaw.accessId;
      delete configRaw.accessID;
      delete configRaw.access_id;
      delete configRaw.providerAccessId;
      return { id, name, action, provider: provider || undefined, accessId: accessId || undefined, config: configRaw };
    })
    .filter((x): x is WorkflowSpecNode => x !== null);

  const edges = rawEdges
    .map((item) => {
      const edgeObj = item && typeof item === "object" ? (item as Record<string, unknown>) : {};
      const source = stringify(edgeObj.source);
      const target = stringify(edgeObj.target);
      if (!source || !target) return null;
      return { source, target };
    })
    .filter((x): x is WorkflowSpecEdge => Boolean(x));

  return { version: 1, options: { failFast: true }, nodes, edges };
}

function specToGraph(spec: WorkflowSpec): WorkflowGraph {
  const nodes: Array<Node<FlowNodeData>> = spec.nodes.map((node, index) => {
    const cfg: Record<string, unknown> = { ...(node.config || {}), provider: node.provider || "", accessId: node.accessId || "" };
    return {
      id: node.id,
      type: node.action,
      position: { x: 120, y: 80 + index * 140 },
      data: { name: node.name, label: node.name, config: cfg },
      draggable: false,
      selectable: true,
      connectable: false,
    };
  });

  const edges: Edge[] = spec.edges.map((e, index) => ({ id: `e-${e.source}-${e.target}-${index}`, source: e.source, target: e.target, animated: false }));
  return { nodes, edges };
}

function readSpecFromText(text: string): WorkflowSpec {
  const loaded = YAML.load(text);
  const obj = loaded && typeof loaded === "object" ? (loaded as Record<string, unknown>) : {};
  const rawNodes = Array.isArray(obj.nodes) ? obj.nodes : [];
  const rawEdges = Array.isArray(obj.edges) ? obj.edges : [];

  const nodes: WorkflowSpecNode[] = rawNodes.map((item, index) => {
    const node = item && typeof item === "object" ? (item as Record<string, unknown>) : {};
    const actionRaw = stringify(node.action).toLowerCase();
    const action: "apply" | "deploy" = actionRaw === "deploy" ? "deploy" : "apply";
    const config = node.config && typeof node.config === "object" ? (node.config as Record<string, unknown>) : {};
    return {
      id: stringify(node.id) || `node-${index + 1}`,
      name: stringify(node.name) || `Node ${index + 1}`,
      action,
      provider: stringify(node.provider) || undefined,
      accessId: stringify(node.accessId) || undefined,
      config,
    };
  });

  const edges: WorkflowSpecEdge[] = rawEdges
    .map((item) => {
      const edge = item && typeof item === "object" ? (item as Record<string, unknown>) : {};
      const source = stringify(edge.source);
      const target = stringify(edge.target);
      if (!source || !target) return null;
      return { source, target };
    })
    .filter((x): x is WorkflowSpecEdge => Boolean(x));

  return {
    version: Number(obj.version || 1),
    options: obj.options && typeof obj.options === "object" ? (obj.options as Record<string, unknown>) : { failFast: true },
    nodes,
    edges,
  };
}

function dumpSpec(spec: WorkflowSpec): string {
  return YAML.dump(spec, { lineWidth: 120, noRefs: true });
}

function formatDateTime(v?: string) {
  if (!v) return "-";
  const d = new Date(v);
  if (Number.isNaN(d.getTime())) return "-";
  return d.toLocaleString();
}

function nodeStatusColor(status?: string) {
  switch ((status || "").toLowerCase()) {
    case "running":
      return "#ebf5ff";
    case "succeeded":
      return "#e9f9ee";
    case "failed":
      return "#ffecec";
    case "skipped":
      return "#f8f8f8";
    default:
      return "#ffffff";
  }
}

function isActiveRun(status?: string) {
  const s = (status || "").toLowerCase();
  return s === "pending" || s === "processing" || s === "running";
}

export default function WorkflowsPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const { data } = useQuery({ queryKey: ["workflows"], queryFn: api.listWorkflows });

  const [editingID, setEditingID] = useState<string | undefined>(undefined);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [templateKey, setTemplateKey] = useState<WorkflowTemplateKey>("apply-only");
  const [specText, setSpecText] = useState(() => dumpSpec(makeTemplateSpec("apply-only")));
  const [error, setError] = useState("");
  const [runNotice, setRunNotice] = useState<{ type: "success" | "error" | "info"; text: string } | null>(null);
  const [runningWorkflowID, setRunningWorkflowID] = useState<string>("");
  const [viewingRunsWorkflowID, setViewingRunsWorkflowID] = useState<string>("");
  const [selectedRunID, setSelectedRunID] = useState<string>("");
  const [selectedRunNodeID, setSelectedRunNodeID] = useState<string>("");
  const [dialogOpen, setDialogOpen] = useState(false);
  const [flowInstance, setFlowInstance] = useState<ReactFlowInstance | null>(null);
  const seenRunStatusRef = useRef<Record<string, string>>({});

  const parsedSpec = useMemo(() => {
    try {
      return { spec: readSpecFromText(specText), error: "" };
    } catch (e) {
      return { spec: null, error: e instanceof Error ? e.message : "YAML 解析失败" };
    }
  }, [specText]);

  const runsQuery = useQuery({
    queryKey: ["workflow-runs", viewingRunsWorkflowID],
    queryFn: () => api.listWorkflowRuns(viewingRunsWorkflowID),
    enabled: Boolean(viewingRunsWorkflowID),
    refetchInterval: viewingRunsWorkflowID ? 1500 : false,
  });

  const selectedRun = useMemo(() => (runsQuery.data?.items || []).find((x) => x.id === selectedRunID), [runsQuery.data?.items, selectedRunID]);
  const activeRun = isActiveRun(selectedRun?.status);

  const runNodesQuery = useQuery({
    queryKey: ["workflow-run-nodes", viewingRunsWorkflowID, selectedRunID],
    queryFn: () => api.listWorkflowRunNodes(viewingRunsWorkflowID, selectedRunID),
    enabled: Boolean(viewingRunsWorkflowID && selectedRunID),
    refetchInterval: selectedRunID ? (activeRun ? 1200 : 4000) : false,
  });

  const runEventsQuery = useQuery({
    queryKey: ["workflow-run-events", viewingRunsWorkflowID, selectedRunID, selectedRunNodeID],
    queryFn: () => api.listWorkflowRunEvents(viewingRunsWorkflowID, selectedRunID, { nodeId: selectedRunNodeID || undefined, limit: 300 }),
    enabled: Boolean(viewingRunsWorkflowID && selectedRunID),
    refetchInterval: selectedRunID ? (activeRun ? 1200 : 4000) : false,
  });

  useEffect(() => {
    const list = runsQuery.data?.items || [];
    list.forEach((run) => {
      const prev = seenRunStatusRef.current[run.id];
      const now = (run.status || "").toLowerCase();
      if (prev !== now && (now === "succeeded" || now === "failed")) {
        if (now === "succeeded") {
          toast.success(`Run ${run.id.slice(0, 8)} 执行成功`);
        } else {
          toast.error(`Run ${run.id.slice(0, 8)} 执行失败`);
        }
      }
      seenRunStatusRef.current[run.id] = now;
    });
  }, [runsQuery.data?.items, toast]);

  const runNodeByID = useMemo(() => {
    const m = new Map<string, WorkflowRunNode>();
    (runNodesQuery.data?.items || []).forEach((x) => m.set(x.nodeId, x));
    return m;
  }, [runNodesQuery.data?.items]);

  const flowPreview = useMemo(() => {
    const spec = parsedSpec.spec || makeTemplateSpec("apply-only");
    const graph = specToGraph(spec);
    const nodes = graph.nodes.map((n) => {
      const status = runNodeByID.get(n.id)?.status;
      return {
        ...n,
        style: {
          border: n.id === selectedRunNodeID ? "2px solid #171717" : "1px solid #e5e5e5",
          borderRadius: 10,
          padding: 8,
          background: nodeStatusColor(status),
          width: 260,
          boxShadow: "rgba(0,0,0,.04) 0 2px 2px",
        },
        data: { ...n.data, label: `${n.data.name}${status ? ` · ${status}` : ""}` },
      };
    });
    return { nodes, edges: graph.edges };
  }, [parsedSpec.spec, runNodeByID, selectedRunNodeID]);

  useEffect(() => {
    if (!flowInstance) return;
    const timer = window.setTimeout(() => {
      flowInstance.fitView({ padding: 0.15, duration: 250 });
    }, 0);
    return () => window.clearTimeout(timer);
  }, [flowInstance, flowPreview.nodes, flowPreview.edges, dialogOpen]);

  const save = useMutation({
    mutationFn: async () => {
      const spec = readSpecFromText(specText);
      if (spec.nodes.length === 0) throw new Error("至少需要一个节点");
      const graph = specToGraph(spec);
      return api.saveWorkflow({
        id: editingID,
        name,
        description,
        trigger: "manual",
        triggerCron: "",
        enabled: true,
        graphDraft: graph as unknown as Record<string, unknown>,
        graphContent: graph as unknown as Record<string, unknown>,
        hasDraft: true,
        hasContent: true,
      });
    },
    onSuccess: (saved) => {
      setError("");
      setEditingID(saved.id);
      qc.invalidateQueries({ queryKey: ["workflows"] });
      setRunNotice({ type: "success", text: "工作流已保存" });
      toast.success("工作流保存成功");
      setDialogOpen(false);
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : "保存失败";
      setError(msg);
      toast.error(msg);
    },
  });

  const startRun = useMutation({
    mutationFn: async (workflowID: string) => api.startWorkflowRun(workflowID),
    onMutate: (workflowID) => {
      setRunningWorkflowID(workflowID);
      setRunNotice({ type: "info", text: "工作流已提交，正在启动..." });
      toast.info("工作流已触发");
    },
    onSuccess: (res, workflowID) => {
      setRunNotice({ type: "success", text: `已触发运行，Run ID: ${res.runId}` });
      setViewingRunsWorkflowID(workflowID);
      setSelectedRunID(res.runId);
      setSelectedRunNodeID("");
      qc.invalidateQueries({ queryKey: ["workflows"] });
      qc.invalidateQueries({ queryKey: ["workflow-runs", workflowID] });
      qc.invalidateQueries({ queryKey: ["workflow-run-nodes", workflowID, res.runId] });
      qc.invalidateQueries({ queryKey: ["workflow-run-events", workflowID, res.runId] });
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : "触发运行失败";
      setRunNotice({ type: "error", text: msg });
      toast.error(msg);
    },
    onSettled: () => setRunningWorkflowID(""),
  });

  const resetEditor = () => {
    const spec = makeTemplateSpec("apply-only");
    setEditingID(undefined);
    setName("");
    setDescription("");
    setTemplateKey("apply-only");
    setSpecText(dumpSpec(spec));
    setSelectedRunNodeID("");
    setError("");
  };

  const openCreate = () => {
    resetEditor();
    setDialogOpen(true);
  };

  const applyTemplate = () => {
    setSpecText(dumpSpec(makeTemplateSpec(templateKey)));
    setError("");
  };

  const loadWorkflow = (workflow: Workflow) => {
    const spec = graphToSpec(pickWorkflowGraph(workflow));
    setEditingID(workflow.id);
    setName(workflow.name);
    setDescription(workflow.description);
    setSpecText(dumpSpec(spec));
    setError("");
    setDialogOpen(true);
  };

  return (
    <div className="space-y-4">
      <Card>
        <div className="mb-3 flex items-center justify-between">
          <div className="text-sm font-medium">工作流列表</div>
          <Button onClick={openCreate}>新增工作流</Button>
        </div>

        {runNotice ? (
          <div className={`mb-3 rounded-full px-2.5 py-1 text-xs w-fit ${runNotice.type === "error" ? "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]" : runNotice.type === "success" ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]" : "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]"}`}>{runNotice.text}</div>
        ) : null}

        <div className="overflow-x-auto ds-scrollbar">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                <th className="pb-2">名称</th>
                <th className="pb-2">最近运行</th>
                <th className="pb-2">最近时间</th>
                <th className="pb-2 text-right">动作</th>
              </tr>
            </thead>
            <tbody>
              {data?.items.map((w) => (
                <tr key={w.id} className="border-t border-[#f1f1f1]">
                  <td className="py-2">{w.name}</td>
                  <td className="py-2"><StatusBadge status={w.lastRunStatus} /></td>
                  <td className="py-2 text-[#666]">{formatDateTime(w.lastRunTime)}</td>
                  <td className="space-x-2 py-2 text-right">
                    <Button size="sm" variant="outline" onClick={() => loadWorkflow(w)}>编辑</Button>
                    <Button size="sm" disabled={runningWorkflowID === w.id} onClick={() => startRun.mutate(w.id!)}>{runningWorkflowID === w.id ? "触发中..." : "运行"}</Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        const open = viewingRunsWorkflowID !== w.id;
                        setViewingRunsWorkflowID(open ? w.id || "" : "");
                        setSelectedRunID("");
                        setSelectedRunNodeID("");
                      }}
                    >
                      {viewingRunsWorkflowID === w.id ? "收起日志" : "查看日志"}
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={async () => {
                        await api.deleteWorkflow(w.id!);
                        qc.invalidateQueries({ queryKey: ["workflows"] });
                        if (viewingRunsWorkflowID === w.id) {
                          setViewingRunsWorkflowID("");
                          setSelectedRunID("");
                          setSelectedRunNodeID("");
                        }
                        toast.info("工作流已删除");
                      }}
                    >
                      删除
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      {viewingRunsWorkflowID ? (
        <div className="grid grid-cols-1 gap-4 xl:grid-cols-[1fr_1fr]">
          <Card>
            <div className="mb-2 text-sm font-medium">运行记录（最近 30 条）</div>
            <div className="overflow-x-auto ds-scrollbar">
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                    <th className="pb-2">Run ID</th>
                    <th className="pb-2">状态</th>
                    <th className="pb-2">开始</th>
                    <th className="pb-2">结束</th>
                    <th className="pb-2">错误</th>
                  </tr>
                </thead>
                <tbody>
                  {runsQuery.data?.items.length ? (
                    runsQuery.data.items.map((run) => (
                      <tr key={run.id} className={`border-t border-[#f1f1f1] cursor-pointer ${selectedRunID === run.id ? "bg-[#fafafa]" : ""}`} onClick={() => setSelectedRunID(run.id)}>
                        <td className="py-2 font-mono text-xs">{run.id}</td>
                        <td className="py-2"><StatusBadge status={run.status} /></td>
                        <td className="py-2 text-xs text-[#666]">{formatDateTime(run.startedAt)}</td>
                        <td className="py-2 text-xs text-[#666]">{formatDateTime(run.endedAt)}</td>
                        <td className="max-w-[240px] break-all py-2 text-xs text-[var(--ds-danger-fg)]">{run.error || "-"}</td>
                      </tr>
                    ))
                  ) : (
                    <tr><td className="py-3 text-[#808080]" colSpan={5}>{runsQuery.isFetching ? "加载中..." : "暂无运行记录"}</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          </Card>

          <Card>
            <div className="mb-2 flex items-center justify-between">
              <div className="text-sm font-medium">节点详情 {selectedRunID ? `(${selectedRunID.slice(0, 8)})` : ""}</div>
              {selectedRun ? <StatusBadge status={selectedRun.status} /> : null}
            </div>

            <div className="mb-3 max-h-[220px] overflow-auto ds-scrollbar rounded-md border border-[#ebebeb]">
              <table className="w-full text-xs">
                <thead>
                  <tr className="text-left uppercase tracking-wide text-[#808080]">
                    <th className="px-2 py-2">Node</th>
                    <th className="px-2 py-2">状态</th>
                    <th className="px-2 py-2">Provider</th>
                    <th className="px-2 py-2">错误</th>
                  </tr>
                </thead>
                <tbody>
                  {runNodesQuery.data?.items?.length ? (
                    runNodesQuery.data.items.map((n) => (
                      <tr key={n.id} className={`border-t border-[#f1f1f1] cursor-pointer ${selectedRunNodeID === n.nodeId ? "bg-[#fafafa]" : ""}`} onClick={() => setSelectedRunNodeID((prev) => (prev === n.nodeId ? "" : n.nodeId))}>
                        <td className="px-2 py-2 font-mono">{n.nodeId}</td>
                        <td className="px-2 py-2"><StatusBadge status={n.status} /></td>
                        <td className="px-2 py-2">{n.provider || "-"}</td>
                        <td className="px-2 py-2 text-[var(--ds-danger-fg)]">{n.error || "-"}</td>
                      </tr>
                    ))
                  ) : (
                    <tr><td className="px-2 py-3 text-[#808080]" colSpan={4}>{runNodesQuery.isFetching ? "加载中..." : "暂无节点状态"}</td></tr>
                  )}
                </tbody>
              </table>
            </div>

            <div className="text-xs text-[#666]">事件日志 {selectedRunNodeID ? `(Node: ${selectedRunNodeID})` : "(全部节点)"}</div>
            <div className="mt-2 max-h-[260px] overflow-auto ds-scrollbar rounded-md border border-[#ebebeb] p-2">
              {(runEventsQuery.data?.items || []).length ? (
                runEventsQuery.data!.items.map((e) => (
                  <div key={e.id} className="border-b border-[#f1f1f1] py-2 last:border-0">
                    <div className="flex items-center gap-2 text-xs">
                      <span className="font-mono text-[#666]">{formatDateTime(e.createdAt)}</span>
                      <span className="rounded-full bg-[#f5f5f5] px-2 py-0.5">{e.eventType}</span>
                      <span className="font-mono text-[#808080]">{e.nodeId || "-"}</span>
                    </div>
                    <div className="mt-1 text-sm text-[#171717]">{e.message || "-"}</div>
                  </div>
                ))
              ) : (
                <div className="py-4 text-sm text-[#808080]">{runEventsQuery.isFetching ? "加载中..." : "暂无日志"}</div>
              )}
            </div>
          </Card>
        </div>
      ) : null}

      <Dialog.Root open={dialogOpen} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-50 bg-black/35" />
          <Dialog.Content className="fixed left-1/2 top-1/2 z-50 w-[min(1200px,94vw)] max-h-[90vh] -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-xl bg-white p-5 shadow-2xl focus:outline-none">
            <div className="mb-3 flex items-center justify-between">
              <div>
                <Dialog.Title className="text-lg font-semibold tracking-[-0.02em]">{editingID ? "编辑工作流" : "新增工作流"}</Dialog.Title>
                <Dialog.Description className="text-sm text-[#666]">配置优先，流程图用于执行观测。</Dialog.Description>
              </div>
              <StatusBadge status={editingID ? "processing" : "pending"} />
            </div>

            <div className="grid grid-cols-1 gap-4 xl:grid-cols-[1.1fr_0.9fr]">
              <div>
                <div className="mb-3 grid gap-2 md:grid-cols-2">
                  <Input placeholder="工作流名称" value={name} onChange={(e) => setName(e.target.value)} />
                  <Input placeholder="描述" value={description} onChange={(e) => setDescription(e.target.value)} />
                </div>

                <div className="mb-3 flex flex-wrap items-center gap-2">
                  <Button onClick={() => save.mutate()} disabled={save.isPending}>{editingID ? "保存工作流" : "创建工作流"}</Button>
                  <Button variant="outline" onClick={resetEditor}>新建模板</Button>
                  <select
                    className="ds-ring h-9 rounded-md bg-white px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
                    value={templateKey}
                    onChange={(e) => setTemplateKey(e.target.value as WorkflowTemplateKey)}
                  >
                    {TEMPLATE_OPTIONS.map((item) => (
                      <option key={item.value} value={item.value}>{item.label}</option>
                    ))}
                  </select>
                  <Button variant="outline" onClick={applyTemplate}>套用模板</Button>
                </div>

                <Textarea value={specText} onChange={(e) => setSpecText(e.target.value)} className="min-h-[520px] font-mono text-xs" />
                {parsedSpec.error ? <div className="mt-2 rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">YAML 错误: {parsedSpec.error}</div> : null}
                {error ? <div className="mt-2 rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">{error}</div> : null}
              </div>

              <div>
                <div className="mb-2 text-sm font-medium">流程图预览</div>
                <div className="h-[640px] rounded-lg border border-[#ebebeb]">
                  <ReactFlow
                    nodes={flowPreview.nodes}
                    edges={flowPreview.edges}
                    fitView
                    fitViewOptions={{ padding: 0.15, minZoom: 0.2, maxZoom: 1.2 }}
                    onInit={setFlowInstance}
                    onNodeClick={(_, node) => setSelectedRunNodeID(node.id)}
                  >
                    <MiniMap />
                    <Controls />
                    <Background />
                  </ReactFlow>
                </div>
              </div>
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
