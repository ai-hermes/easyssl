import * as Dialog from "@radix-ui/react-dialog";
import { useEffect, useMemo, useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import ReactFlow, { Background, Controls, Edge, Handle, MarkerType, Node, NodeProps, NodeTypes, Position, ReactFlowInstance } from "reactflow";
import "reactflow/dist/style.css";
import YAML from "js-yaml";

import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { StatusBadge } from "@/components/ui/status-badge";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/toast";
import type { Workflow, WorkflowRunNode } from "@/types";

type FlowNodeData = {
  name: string;
  action?: "start" | "apply" | "deploy" | "end";
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
  action: "start" | "apply" | "deploy" | "end";
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

function WorkflowPreviewNode({ data }: NodeProps<FlowNodeData>) {
  const action = (data.action || "").toLowerCase();
  return (
    <div className="relative rounded-[inherit] px-2 py-1.5 text-sm text-[#171717]">
      {action !== "start" ? (
        <Handle
          type="target"
          position={Position.Top}
          style={{ left: "50%", top: 0, transform: "translate(-50%, -70%)" }}
        />
      ) : null}
      <div>{data.label || data.name}</div>
      {action !== "end" ? (
        <Handle
          type="source"
          position={Position.Bottom}
          style={{ left: "50%", bottom: 0, transform: "translate(-50%, 70%)" }}
        />
      ) : null}
    </div>
  );
}

const workflowNodeTypes: NodeTypes = {
  start: WorkflowPreviewNode,
  apply: WorkflowPreviewNode,
  deploy: WorkflowPreviewNode,
  end: WorkflowPreviewNode,
};

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
        id: "start",
        name: "Start",
        action: "start",
      },
      {
        id: "apply-1",
        name: "Apply",
        action: "apply",
        provider: "tencentcloud",
        accessId: "",
        config: { domains: "example.com", caProvider: "letsencrypt" },
      },
      {
        id: "end",
        name: "End",
        action: "end",
      },
    ],
    edges: [
      { source: "start", target: "apply-1" },
      { source: "apply-1", target: "end" },
    ],
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

  base.nodes.splice(2, 0, { id: "deploy-1", name: "Deploy", action: "deploy", provider, accessId: "", config: deployConfig });
  base.edges = [
    { source: "start", target: "apply-1" },
    { source: "apply-1", target: "deploy-1" },
    { source: "deploy-1", target: "end" },
  ];
  return base;
}

function stringify(v: unknown): string {
  if (v === null || v === undefined) return "";
  return String(v).trim();
}

function inferAction(type: string, id: string, name: string, config: Record<string, unknown>): "start" | "apply" | "deploy" | "end" | "" {
  const t = type.toLowerCase();
  const i = id.toLowerCase();
  const n = name.toLowerCase();
  const provider = stringify(config.provider).toLowerCase();
  if (t === "input" || t === "start" || i === "start" || n === "start") return "start";
  if (t === "output" || t === "end" || i === "end" || n === "end") return "end";
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
      if (action === "start" || action === "end") {
        return { id, name, action };
      }
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
    const cfg: Record<string, unknown> =
      node.action === "start" || node.action === "end"
        ? {}
        : { ...(node.config || {}), provider: node.provider || "", accessId: node.accessId || "" };
    return {
      id: node.id,
      type: node.action,
      position: { x: 120, y: 80 + index * 140 },
      data: { name: node.name, action: node.action, label: node.name, config: cfg },
      draggable: false,
      selectable: true,
      connectable: false,
    };
  });

  const edges: Edge[] = spec.edges.map((e, index) => ({
    id: `e-${e.source}-${e.target}-${index}`,
    source: e.source,
    target: e.target,
    type: "smoothstep",
    animated: false,
    zIndex: 2,
    style: {
      stroke: "#171717",
      strokeWidth: 1.8,
    },
    markerEnd: {
      type: MarkerType.ArrowClosed,
      color: "#171717",
      width: 18,
      height: 18,
    },
  }));
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
    const action: "start" | "apply" | "deploy" | "end" =
      actionRaw === "start"
        ? "start"
        : actionRaw === "deploy"
          ? "deploy"
          : actionRaw === "end"
            ? "end"
            : "apply";
    const config = node.config && typeof node.config === "object" ? (node.config as Record<string, unknown>) : {};
    if (action === "start" || action === "end") {
      return {
        id: stringify(node.id) || `node-${index + 1}`,
        name: stringify(node.name) || `Node ${index + 1}`,
        action,
      };
    }
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
  const [selectedWorkflowID, setSelectedWorkflowID] = useState<string>("");
  const [selectedRunID, setSelectedRunID] = useState<string>("");
  const [selectedRunNodeID, setSelectedRunNodeID] = useState<string>("");
  const [runDrawerOpen, setRunDrawerOpen] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editorFlowInstance, setEditorFlowInstance] = useState<ReactFlowInstance | null>(null);
  const [runFlowInstance, setRunFlowInstance] = useState<ReactFlowInstance | null>(null);
  const seenRunStatusRef = useRef<Record<string, string>>({});

  const parsedSpec = useMemo(() => {
    try {
      return { spec: readSpecFromText(specText), error: "" };
    } catch (e) {
      return { spec: null, error: e instanceof Error ? e.message : "YAML 解析失败" };
    }
  }, [specText]);

  const selectedWorkflow = useMemo(() => (data?.items || []).find((x) => x.id === selectedWorkflowID), [data?.items, selectedWorkflowID]);

  const runsQuery = useQuery({
    queryKey: ["workflow-runs", selectedWorkflowID],
    queryFn: () => api.listWorkflowRuns(selectedWorkflowID),
    enabled: Boolean(selectedWorkflowID),
    refetchInterval: selectedWorkflowID ? 1500 : false,
  });

  const selectedRun = useMemo(() => (runsQuery.data?.items || []).find((x) => x.id === selectedRunID), [runsQuery.data?.items, selectedRunID]);
  const activeRun = isActiveRun(selectedRun?.status);

  const runNodesQuery = useQuery({
    queryKey: ["workflow-run-nodes", selectedWorkflowID, selectedRunID],
    queryFn: () => api.listWorkflowRunNodes(selectedWorkflowID, selectedRunID),
    enabled: Boolean(selectedWorkflowID && selectedRunID),
    refetchInterval: selectedRunID ? (activeRun ? 1200 : 4000) : false,
  });

  const runEventsQuery = useQuery({
    queryKey: ["workflow-run-events", selectedWorkflowID, selectedRunID, selectedRunNodeID],
    queryFn: () => api.listWorkflowRunEvents(selectedWorkflowID, selectedRunID, { nodeId: selectedRunNodeID || undefined, limit: 300 }),
    enabled: Boolean(selectedWorkflowID && selectedRunID && selectedRunNodeID),
    refetchInterval: selectedRunID && selectedRunNodeID ? (activeRun ? 1200 : 4000) : false,
  });

  useEffect(() => {
    const list = runsQuery.data?.items || [];
    const nextSeen: Record<string, string> = {};
    list.forEach((run) => {
      const prev = seenRunStatusRef.current[run.id];
      const now = (run.status || "").toLowerCase();
      if (prev && prev !== now && (now === "succeeded" || now === "failed")) {
        if (now === "succeeded") {
          toast.success(`Run ${run.id.slice(0, 8)} 执行成功`);
        } else {
          toast.error(`Run ${run.id.slice(0, 8)} 执行失败`);
        }
      }
      nextSeen[run.id] = now;
    });
    seenRunStatusRef.current = nextSeen;
  }, [runsQuery.data?.items, toast]);

  useEffect(() => {
    seenRunStatusRef.current = {};
  }, [selectedWorkflowID]);

  useEffect(() => {
    const items = data?.items || [];
    if (items.length > 0 && selectedWorkflowID && !items.some((x) => x.id === selectedWorkflowID)) {
      setSelectedWorkflowID("");
      setSelectedRunID("");
      setSelectedRunNodeID("");
      setRunDrawerOpen(false);
    }
  }, [data?.items, selectedWorkflowID]);

  useEffect(() => {
    const runs = runsQuery.data?.items || [];
    if (selectedRunID && !runs.some((x) => x.id === selectedRunID)) {
      setSelectedRunID("");
      setSelectedRunNodeID("");
      setRunDrawerOpen(false);
    }
  }, [runsQuery.data?.items, selectedRunID]);

  const runNodeByID = useMemo(() => {
    const m = new Map<string, WorkflowRunNode>();
    (runNodesQuery.data?.items || []).forEach((x) => m.set(x.nodeId, x));
    return m;
  }, [runNodesQuery.data?.items]);
  const selectedNodeRun = selectedRunNodeID ? runNodeByID.get(selectedRunNodeID) : undefined;

  const runSpec = useMemo(() => {
    const runGraph = selectedRun?.graph;
    if (runGraph && typeof runGraph === "object") return graphToSpec(runGraph);
    if (selectedWorkflow) return graphToSpec(pickWorkflowGraph(selectedWorkflow));
    return makeTemplateSpec("apply-only");
  }, [selectedRun?.graph, selectedWorkflow]);

  const runFlowPreview = useMemo(() => {
    const graph = specToGraph(runSpec);
    const nodes = graph.nodes.map((n) => {
      const action = (n.data.action || "").toLowerCase();
      const isBoundaryNode = action === "start" || action === "end";
      const status = runNodeByID.get(n.id)?.status || (selectedRunID && !isBoundaryNode ? "unknown" : "");
      return {
        ...n,
        style: {
          border: n.id === selectedRunNodeID ? "2px solid #171717" : "1px solid #e5e5e5",
          borderRadius: 10,
          padding: 8,
          background: nodeStatusColor(status),
          width: 240,
          boxShadow: "rgba(0,0,0,.04) 0 2px 2px",
        },
        data: { ...n.data, label: `${n.data.name}${status ? ` · ${status}` : ""}` },
      };
    });
    return { nodes, edges: graph.edges };
  }, [runSpec, runNodeByID, selectedRunNodeID, selectedRunID]);

  const editorFlowPreview = useMemo(() => {
    const spec = parsedSpec.spec || makeTemplateSpec("apply-only");
    const graph = specToGraph(spec);
    const nodes = graph.nodes.map((n) => {
      return {
        ...n,
        style: {
          border: "1px solid #e5e5e5",
          borderRadius: 10,
          padding: 8,
          background: "#ffffff",
          width: 260,
          boxShadow: "rgba(0,0,0,.04) 0 2px 2px",
        },
        data: { ...n.data, label: n.data.name },
      };
    });
    return { nodes, edges: graph.edges };
  }, [parsedSpec.spec]);

  useEffect(() => {
    if (!editorFlowInstance || !dialogOpen) return;
    const timer = window.setTimeout(() => {
      editorFlowInstance.fitView({ padding: 0.15, duration: 250 });
    }, 0);
    return () => window.clearTimeout(timer);
  }, [editorFlowInstance, editorFlowPreview.nodes, editorFlowPreview.edges, dialogOpen]);

  useEffect(() => {
    if (!runFlowInstance || !runDrawerOpen) return;
    const timer = window.setTimeout(() => {
      runFlowInstance.fitView({ padding: 0.15, duration: 250 });
    }, 0);
    return () => window.clearTimeout(timer);
  }, [runFlowInstance, runFlowPreview.nodes, runFlowPreview.edges, runDrawerOpen]);

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
      setSelectedWorkflowID(workflowID);
      setSelectedRunID("");
      setSelectedRunNodeID("");
      setRunDrawerOpen(false);
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

  const isEditing = Boolean(editingID);
  const selectedWorkflowName = selectedWorkflow?.name || "-";

  return (
    <div className="space-y-4">
      {runNotice ? (
        <div className={`rounded-full px-2.5 py-1 text-xs w-fit ${runNotice.type === "error" ? "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]" : runNotice.type === "success" ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]" : "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]"}`}>{runNotice.text}</div>
      ) : null}

      <div className="flex justify-end">
        <Button onClick={openCreate}>新增工作流</Button>
      </div>

      <Card>
        <CardHeader className="pb-2">
          <div className="text-base font-semibold tracking-[-0.02em]">Workflows</div>
          <div className="text-xs text-[#666]">默认仅展示工作流列表。</div>
        </CardHeader>

        <CardContent className="pt-0">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                <th className="px-2 pb-3 pt-1">名称</th>
                <th className="px-2 pb-3 pt-1">最近运行</th>
                <th className="px-2 pb-3 pt-1">最近时间</th>
                <th className="px-2 pb-3 pt-1 text-right">动作</th>
              </tr>
            </thead>
            <tbody>
              {data?.items.map((w) => {
                const selected = selectedWorkflowID === w.id;
                return (
                  <tr
                    key={w.id}
                    className={`border-t border-[#f1f1f1] cursor-pointer ${selected ? "bg-[#f7f7f5]" : ""}`}
                    onClick={() => {
                      setSelectedWorkflowID(w.id || "");
                      setSelectedRunID("");
                      setSelectedRunNodeID("");
                      setRunDrawerOpen(false);
                    }}
                  >
                    <td className={`px-2 py-3 ${selected ? "text-[#171717] font-medium" : ""}`}>{w.name}</td>
                    <td className="px-2 py-3"><StatusBadge status={w.lastRunStatus} /></td>
                    <td className="px-2 py-3 text-[#666]">{formatDateTime(w.lastRunTime)}</td>
                    <td className="px-2 py-3 space-x-2 text-right">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={(e) => {
                          e.stopPropagation();
                          loadWorkflow(w);
                        }}
                      >
                        编辑
                      </Button>
                      <Button
                        size="sm"
                        disabled={runningWorkflowID === w.id}
                        onClick={(e) => {
                          e.stopPropagation();
                          startRun.mutate(w.id!);
                        }}
                      >
                        {runningWorkflowID === w.id ? "触发中..." : "运行"}
                      </Button>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={async (e) => {
                          e.stopPropagation();
                          await api.deleteWorkflow(w.id!);
                          qc.invalidateQueries({ queryKey: ["workflows"] });
                          if (selectedWorkflowID === w.id) {
                            setSelectedWorkflowID("");
                            setSelectedRunID("");
                            setSelectedRunNodeID("");
                            setRunDrawerOpen(false);
                          }
                          toast.info("工作流已删除");
                        }}
                      >
                        删除
                      </Button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </CardContent>
      </Card>

      {selectedWorkflowID ? (
        <Card className="bg-[#f5f5f7]">
          <div className="mb-2 flex items-center justify-between">
            <div className="text-sm font-medium">执行记录</div>
            <span className="text-xs text-[#666]">{selectedWorkflowName}</span>
          </div>
          <div className="overflow-x-auto ds-scrollbar">
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                  <th className="px-2 pb-3 pt-1">Run ID</th>
                  <th className="px-2 pb-3 pt-1">状态</th>
                  <th className="px-2 pb-3 pt-1">开始</th>
                  <th className="px-2 pb-3 pt-1">结束</th>
                </tr>
              </thead>
              <tbody>
                {(runsQuery.data?.items || []).length ? (
                  runsQuery.data!.items.map((run) => (
                    <tr
                      key={run.id}
                      className="border-t border-[#e9e9ed] cursor-pointer hover:bg-white"
                      onClick={() => {
                        setSelectedRunID(run.id);
                        setSelectedRunNodeID("");
                        setRunDrawerOpen(true);
                      }}
                    >
                      <td className="px-2 py-3 font-mono text-xs">{run.id}</td>
                      <td className="px-2 py-3"><StatusBadge status={run.status} /></td>
                      <td className="px-2 py-3 text-xs text-[#666]">{formatDateTime(run.startedAt)}</td>
                      <td className="px-2 py-3 text-xs text-[#666]">{formatDateTime(run.endedAt)}</td>
                    </tr>
                  ))
                ) : (
                  <tr><td className="px-2 py-3 text-[#808080]" colSpan={4}>{runsQuery.isFetching ? "加载中..." : "暂无运行记录"}</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </Card>
      ) : null}

      <Dialog.Root
        open={runDrawerOpen}
        onOpenChange={(open) => {
          setRunDrawerOpen(open);
          if (!open) setSelectedRunNodeID("");
        }}
      >
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-40 bg-black/30" />
          <Dialog.Content className="fixed right-0 top-0 z-50 h-screen w-[min(960px,96vw)] overflow-y-auto bg-[#f5f5f7] shadow-2xl focus:outline-none">
            <div className="sticky top-0 z-10 border-b border-[#ececec] bg-[#272729] px-5 py-4 text-white">
              <div className="flex items-center justify-between">
                <div>
                  <Dialog.Title className="text-sm font-medium">执行回显 {selectedRunID ? `(${selectedRunID.slice(0, 8)})` : ""}</Dialog.Title>
                  <Dialog.Description className="mt-1 text-xs text-[#c8c8ce]">点击节点可查看该节点执行日志</Dialog.Description>
                </div>
                <div className="flex items-center gap-2">
                  {selectedRun ? <StatusBadge status={selectedRun.status} /> : null}
                  <Button size="sm" variant="outline" onClick={() => setRunDrawerOpen(false)}>关闭</Button>
                </div>
              </div>
            </div>

            <div className="space-y-3 p-5">
              <div className="h-[360px] rounded-lg border border-[#ebebeb] bg-white">
                <ReactFlow
                  nodes={runFlowPreview.nodes}
                  edges={runFlowPreview.edges}
                  nodeTypes={workflowNodeTypes}
                  defaultEdgeOptions={{ markerEnd: { type: MarkerType.ArrowClosed } }}
                  defaultMarkerColor="#171717"
                  fitView
                  fitViewOptions={{ padding: 0.15, minZoom: 0.2, maxZoom: 1.2 }}
                  onInit={setRunFlowInstance}
                  onNodeClick={(_, node) => setSelectedRunNodeID((prev) => (prev === node.id ? "" : node.id))}
                >
                  <Controls />
                  <Background />
                </ReactFlow>
              </div>

              {!selectedRunNodeID ? (
                <div className="rounded-md border border-[#ebebeb] bg-white px-4 py-6 text-sm text-[#666]">
                  请先点击流程图中的节点查看该节点状态与执行日志。
                </div>
              ) : (
                <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
                  <div>
                    <div className="mb-2 text-xs font-medium uppercase tracking-wide text-[#666]">节点状态</div>
                    <div className="rounded-md border border-[#ebebeb] bg-white p-3 text-sm">
                      {selectedNodeRun ? (
                        <div className="space-y-2">
                          <div className="text-xs text-[#808080]">Node ID</div>
                          <div className="font-mono">{selectedNodeRun.nodeId}</div>
                          <div className="text-xs text-[#808080]">状态</div>
                          <div><StatusBadge status={selectedNodeRun.status} /></div>
                          <div className="text-xs text-[#808080]">Provider</div>
                          <div>{selectedNodeRun.provider || "-"}</div>
                          <div className="text-xs text-[#808080]">开始 / 结束</div>
                          <div>{formatDateTime(selectedNodeRun.startedAt)} / {formatDateTime(selectedNodeRun.endedAt)}</div>
                          {selectedNodeRun.error ? (
                            <>
                              <div className="text-xs text-[#808080]">错误</div>
                              <div className="text-[var(--ds-danger-fg)] break-all">{selectedNodeRun.error}</div>
                            </>
                          ) : null}
                        </div>
                      ) : (
                        <div className="text-[#808080]">该节点没有执行记录（可能是 start/end 节点或本次运行未执行到该节点）。</div>
                      )}
                    </div>
                  </div>

                  <div>
                    <div className="mb-2 text-xs font-medium uppercase tracking-wide text-[#666]">事件日志 (Node: {selectedRunNodeID})</div>
                    <div className="max-h-[260px] overflow-auto ds-scrollbar rounded-md border border-[#ebebeb] bg-white p-2">
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
                  </div>
                </div>
              )}
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>

      <Dialog.Root open={dialogOpen} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-50 bg-black/35" />
          <Dialog.Content className="fixed left-1/2 top-1/2 z-50 w-[min(1280px,96vw)] max-h-[92vh] -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-xl bg-white p-0 shadow-2xl focus:outline-none">
            <div className="sticky top-0 z-10 border-b border-[#ececec] bg-white/95 px-5 py-4 backdrop-blur">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <div>
                  <Dialog.Title className="text-lg font-semibold tracking-[-0.02em]">{isEditing ? "编辑工作流" : "新增工作流"}</Dialog.Title>
                  <Dialog.Description className="text-sm text-[#666]">配置优先，流程图用于执行观测。</Dialog.Description>
                </div>
                <div className="flex items-center gap-2">
                  <StatusBadge status={isEditing ? "processing" : "pending"} />
                  <Button onClick={() => save.mutate()} disabled={save.isPending}>{isEditing ? "保存工作流" : "创建工作流"}</Button>
                </div>
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 p-5 xl:grid-cols-[1.15fr_0.85fr]">
              <div className="space-y-4">
                <div className="rounded-lg border border-[#ececec] bg-[#fcfcfc] p-3">
                  <div className="mb-2 text-xs font-medium uppercase tracking-wide text-[#808080]">基本信息</div>
                  <div className="grid gap-2 md:grid-cols-2">
                    <Input placeholder="工作流名称" value={name} onChange={(e) => setName(e.target.value)} />
                    <Input placeholder="描述" value={description} onChange={(e) => setDescription(e.target.value)} />
                  </div>
                  {isEditing ? (
                    <div className="mt-2 text-xs text-[#777]">
                      工作流 ID: <span className="font-mono">{editingID}</span>
                    </div>
                  ) : null}
                </div>

                <div className="rounded-lg border border-[#ececec] p-3">
                  <div className="mb-2 flex items-center justify-between gap-2">
                    <div className="text-xs font-medium uppercase tracking-wide text-[#808080]">YAML 编排</div>
                  </div>
                  {!isEditing ? (
                    <div className="mb-3 flex flex-wrap items-center gap-2">
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
                  ) : null}
                  <Textarea value={specText} onChange={(e) => setSpecText(e.target.value)} className="min-h-[520px] font-mono text-xs" />
                  {parsedSpec.error ? <div className="mt-2 rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">YAML 错误: {parsedSpec.error}</div> : null}
                  {error ? <div className="mt-2 rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">{error}</div> : null}
                </div>
              </div>

              <div className="space-y-2 xl:sticky xl:top-[88px] xl:self-start">
                <div className="flex items-center justify-between">
                  <div className="text-sm font-medium">流程图预览</div>
                  <span className="text-xs text-[#808080]">点击节点可联动运行日志筛选</span>
                </div>
                <div className="h-[680px] rounded-lg border border-[#ebebeb]">
                  <ReactFlow
                    nodes={editorFlowPreview.nodes}
                    edges={editorFlowPreview.edges}
                    nodeTypes={workflowNodeTypes}
                    defaultEdgeOptions={{ markerEnd: { type: MarkerType.ArrowClosed } }}
                    defaultMarkerColor="#171717"
                    fitView
                    fitViewOptions={{ padding: 0.15, minZoom: 0.2, maxZoom: 1.2 }}
                    onInit={setEditorFlowInstance}
                  >
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
