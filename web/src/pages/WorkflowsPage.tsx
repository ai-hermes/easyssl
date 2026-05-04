import * as Dialog from "@radix-ui/react-dialog";
import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
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
import { Switch } from "@/components/ui/switch";
import { InlineEdit } from "@/components/ui/inline-edit";
import { useToast } from "@/components/ui/toast";
import { LogStream } from "@/components/log-stream";
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
    <div className="relative flex h-full w-full items-center justify-center rounded-[inherit] px-3 py-2 text-sm text-[#171717]">
      {action !== "start" ? (
        <Handle
          type="target"
          position={Position.Top}
          style={{ left: "50%", top: -5, transform: "translateX(-50%)" }}
        />
      ) : null}
      <div className="truncate">{data.label || data.name}</div>
      {action !== "end" ? (
        <Handle
          type="source"
          position={Position.Bottom}
          style={{ left: "50%", bottom: -5, transform: "translateX(-50%)" }}
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

function getTemplateOptions(t: (key: string) => string): Array<{ value: WorkflowTemplateKey; label: string }> {
  return [
    { value: "apply-only", label: t("workflows.templates.applyOnly") },
    { value: "apply-aliyun", label: t("workflows.templates.applyAliyun") },
    { value: "apply-qiniu", label: t("workflows.templates.applyQiniu") },
    { value: "apply-ssh", label: t("workflows.templates.applySSH") },
  ];
}

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
  const { t } = useTranslation();
  const { data } = useQuery({ queryKey: ["workflows"], queryFn: api.listWorkflows });

  const [editingID, setEditingID] = useState<string | undefined>(undefined);
  const [editingEnabled, setEditingEnabled] = useState(true);
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
      return { spec: null, error: e instanceof Error ? e.message : t("workflows.yamlParseError") };
    }
  }, [specText, t]);

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
          toast.success(t("workflows.runSuccess", { runId: run.id.slice(0, 8) }));
        } else {
          toast.error(t("workflows.runFailed", { runId: run.id.slice(0, 8) }));
        }
      }
      nextSeen[run.id] = now;
    });
    seenRunStatusRef.current = nextSeen;
  }, [runsQuery.data?.items, toast, t]);

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
          borderRadius: 10,
          background: nodeStatusColor(status),
          width: 240,
          boxShadow:
            n.id === selectedRunNodeID
              ? "rgb(23, 23, 23) 0px 0px 0px 2px, rgba(0,0,0,.04) 0px 2px 2px"
              : "rgb(235, 235, 235) 0px 0px 0px 1px, rgba(0,0,0,.04) 0px 2px 2px",
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
          borderRadius: 10,
          background: "#ffffff",
          width: 260,
          boxShadow: "rgb(235, 235, 235) 0px 0px 0px 1px, rgba(0,0,0,.04) 0px 2px 2px",
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
      if (spec.nodes.length === 0) throw new Error(t("workflows.minOneNode"));
      const graph = specToGraph(spec);
      return api.saveWorkflow({
        id: editingID,
        name,
        description,
        trigger: "manual",
        triggerCron: "",
        enabled: editingEnabled,
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
      setRunNotice({ type: "success", text: t("workflows.saveSuccess") });
      toast.success(t("workflows.saveSuccess"));
      setDialogOpen(false);
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : t("common.saveFailed");
      setError(msg);
      toast.error(msg);
    },
  });

  const startRun = useMutation({
    mutationFn: async (workflowID: string) => api.startWorkflowRun(workflowID),
    onMutate: (workflowID) => {
      setRunningWorkflowID(workflowID);
      setRunNotice({ type: "info", text: t("workflows.submitting") });
      toast.info(t("workflows.triggered"));
    },
    onSuccess: (res, workflowID) => {
      setRunNotice({ type: "success", text: t("workflows.runTriggered", { runId: res.runId }) });
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
      const msg = e instanceof Error ? e.message : t("workflows.triggerFailed");
      setRunNotice({ type: "error", text: msg });
      toast.error(msg);
    },
    onSettled: () => setRunningWorkflowID(""),
  });

  const resetEditor = () => {
    const spec = makeTemplateSpec("apply-only");
    setEditingID(undefined);
    setEditingEnabled(true);
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
    setEditingEnabled(workflow.enabled);
    setName(workflow.name);
    setDescription(workflow.description);
    setSpecText(dumpSpec(spec));
    setError("");
    setDialogOpen(true);
  };

  const isEditing = Boolean(editingID);
  const selectedWorkflowName = selectedWorkflow?.name || "-";
  const templateOptions = useMemo(() => getTemplateOptions(t), [t]);

  return (
    <div className="space-y-4">
      {runNotice ? (
        <div className={`rounded-full px-2.5 py-1 text-xs w-fit ${runNotice.type === "error" ? "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]" : runNotice.type === "success" ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]" : "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]"}`}>{runNotice.text}</div>
      ) : null}

      <div className="flex justify-end">
        <Button onClick={openCreate}>{t("workflows.addWorkflow")}</Button>
      </div>

      <Card>
        <CardHeader className="pb-2">
          <div className="text-base font-semibold tracking-[-0.02em]">Workflows</div>
          <div className="text-xs text-[#666]">{t("workflows.listDescription")}</div>
        </CardHeader>

        <CardContent className="pt-0">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                <th className="px-2 pb-3 pt-1">{t("workflows.columns.name")}</th>
                <th className="px-2 pb-3 pt-1">{t("workflows.columns.status")}</th>
                <th className="px-2 pb-3 pt-1">{t("workflows.columns.lastRun")}</th>
                <th className="px-2 pb-3 pt-1">{t("workflows.columns.lastTime")}</th>
                <th className="px-2 pb-3 pt-1 text-right">{t("workflows.columns.actions")}</th>
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
                    <td className="px-2 py-3">
                      <span className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${w.enabled ? "border-transparent bg-emerald-100 text-emerald-700" : "border-transparent bg-slate-100 text-slate-700"}`}>
                        {w.enabled ? t("common.enabled") : t("common.disabled")}
                      </span>
                    </td>
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
                        {t("common.edit")}
                      </Button>
                      <Button
                        size="sm"
                        disabled={runningWorkflowID === w.id}
                        onClick={(e) => {
                          e.stopPropagation();
                          startRun.mutate(w.id!);
                        }}
                      >
                        {runningWorkflowID === w.id ? t("common.triggering") : t("common.run")}
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
                          toast.info(t("workflows.deleteSuccess"));
                        }}
                      >
                        {t("common.delete")}
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
            <div className="text-sm font-medium">{t("workflows.executionHistory")}</div>
            <span className="text-xs text-[#666]">{selectedWorkflowName}</span>
          </div>
          <div className="overflow-x-auto ds-scrollbar">
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                  <th className="px-2 pb-3 pt-1">{t("workflows.columns.runId")}</th>
                  <th className="px-2 pb-3 pt-1">{t("workflows.columns.status")}</th>
                  <th className="px-2 pb-3 pt-1">{t("workflows.columns.startedAt")}</th>
                  <th className="px-2 pb-3 pt-1">{t("workflows.columns.endedAt")}</th>
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
                  <tr><td className="px-2 py-3 text-[#808080]" colSpan={4}>{runsQuery.isFetching ? t("common.loading") : t("workflows.noRuns")}</td></tr>
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
          <Dialog.Content className="fixed right-0 top-0 z-50 h-screen w-[min(960px,96vw)] overflow-y-auto bg-white shadow-2xl focus:outline-none">
            <div className="sticky top-0 z-10 border-b border-[#f1f1f1] bg-white/95 px-6 py-4 backdrop-blur">
              <div className="flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <Dialog.Title className="text-sm font-semibold text-[#171717]">{t("workflows.executionEcho")} {selectedRunID ? `(${selectedRunID.slice(0, 8)})` : ""}</Dialog.Title>
                    {selectedRun ? <StatusBadge status={selectedRun.status} /> : null}
                  </div>
                  <Dialog.Description className="mt-1 text-xs text-[#666]">{t("workflows.clickNodeForLogs")}</Dialog.Description>
                </div>
                <div className="flex items-center">
                  <Button size="sm" variant="outline" onClick={() => setRunDrawerOpen(false)}>{t("common.close")}</Button>
                </div>
              </div>
            </div>

            <div className="space-y-6 p-6">
              <div className="h-[360px] rounded-xl bg-white ds-card">
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
                <div className="rounded-xl bg-white ds-card px-6 py-8 text-center text-sm text-[#666]">
                  {t("workflows.selectNodePrompt")}
                </div>
              ) : (
                <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
                  <div>
                    <div className="mb-2 text-[11px] font-semibold uppercase tracking-wider text-[#808080]">{t("workflows.nodeStatus")}</div>
                    <div className="rounded-xl bg-white ds-card p-5 text-sm">
                      {selectedNodeRun ? (
                        <div className="space-y-4">
                          <div className="grid grid-cols-3 gap-4">
                            <div>
                              <div className="text-[11px] font-medium uppercase tracking-wider text-[#808080]">Node ID</div>
                              <div className="mt-1 font-mono text-sm text-[#171717]">{selectedNodeRun.nodeId}</div>
                            </div>
                            <div>
                              <div className="text-[11px] font-medium uppercase tracking-wider text-[#808080]">{t("workflows.status")}</div>
                              <div className="mt-1"><StatusBadge status={selectedNodeRun.status} /></div>
                            </div>
                            <div>
                              <div className="text-[11px] font-medium uppercase tracking-wider text-[#808080]">Provider</div>
                              <div className="mt-1 font-mono text-sm text-[#171717]">{selectedNodeRun.provider || "-"}</div>
                            </div>
                          </div>
                          <div className="grid grid-cols-2 gap-4">
                            <div>
                              <div className="text-[11px] font-medium uppercase tracking-wider text-[#808080]">{t("workflows.startTime")}</div>
                              <div className="mt-1 text-sm text-[#171717]">{formatDateTime(selectedNodeRun.startedAt)}</div>
                            </div>
                            <div>
                              <div className="text-[11px] font-medium uppercase tracking-wider text-[#808080]">{t("workflows.endTime")}</div>
                              <div className="mt-1 text-sm text-[#171717]">{formatDateTime(selectedNodeRun.endedAt)}</div>
                            </div>
                          </div>
                          {selectedNodeRun.error ? (
                            <div>
                              <div className="text-[11px] font-medium uppercase tracking-wider text-[#808080]">{t("workflows.error")}</div>
                              <div className="mt-1 text-sm text-[var(--ds-danger-fg)] break-all">{selectedNodeRun.error}</div>
                            </div>
                          ) : null}
                        </div>
                      ) : (
                        <div className="text-sm text-[#808080]">{t("workflows.noNodeRecords")}</div>
                      )}
                    </div>
                  </div>

                  <div>
                    <div className="mb-2 text-[11px] font-semibold uppercase tracking-wider text-[#808080]">{t("workflows.eventLogs")} <span className="font-normal normal-case tracking-normal text-[#999]">(Node: {selectedRunNodeID})</span></div>
                    <div className="rounded-xl bg-white ds-card p-0 overflow-hidden">
                      <LogStream
                        items={runEventsQuery.data?.items || []}
                        loading={runEventsQuery.isFetching}
                        emptyText={runEventsQuery.isFetching ? t("common.loading") : t("workflows.noLogs")}
                      />
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
                  <Dialog.Title className="text-lg font-semibold tracking-[-0.02em]">{isEditing ? (name || t("workflows.editTitle")) : t("workflows.createTitle")}</Dialog.Title>
                  <Dialog.Description className="text-sm text-[#666]">{isEditing ? (description || t("workflows.dialogDescription")) : t("workflows.dialogDescription")}</Dialog.Description>
                </div>
                <div className="flex items-center gap-2">
                  <Button onClick={() => save.mutate()} disabled={save.isPending}>{t("common.save")}</Button>
                  <Dialog.Close asChild>
                    <Button variant="outline" onClick={() => setDialogOpen(false)}>{t("common.cancel")}</Button>
                  </Dialog.Close>
                </div>
              </div>
            </div>

            <div className="space-y-4 p-5">
              {/* Section 1: 基本信息 */}
              <div className="rounded-lg p-4" style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 2px, #fafafa 0px 0px 0px 1px" }}>
                <div className="mb-3 text-xs font-medium uppercase tracking-wide text-[#808080]">{t("workflows.basicInfo")}</div>
                {isEditing ? (
                  <div className="mb-4 flex items-center gap-3">
                    <span className="w-20 shrink-0 whitespace-nowrap text-sm font-medium text-[#666]">{t("workflows.workflowId")}:</span>
                    <span className="font-mono text-sm text-[#4d4d4d]">{editingID}</span>
                  </div>
                ) : null}
                <div className="space-y-4">
                  <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                    <div className="flex items-center gap-3">
                      <label className="w-20 shrink-0 whitespace-nowrap text-sm font-medium text-[#666]">{t("workflows.nameLabel")}:</label>
                      <InlineEdit value={name} onSave={setName} placeholder={t("workflows.namePlaceholder")} />
                    </div>
                    <div className="flex items-center gap-3">
                      <label className="w-20 shrink-0 whitespace-nowrap text-sm font-medium text-[#666]">{t("workflows.descriptionPlaceholder")}:</label>
                      <InlineEdit value={description} onSave={setDescription} placeholder={t("workflows.descriptionPlaceholder")} />
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <label className="w-20 shrink-0 whitespace-nowrap text-sm font-medium text-[#666]">{t("workflows.statusLabel")}:</label>
                    <div className="flex h-8 items-center gap-2">
                      <Switch checked={editingEnabled} onCheckedChange={setEditingEnabled} />
                      <span className="text-sm text-[#4d4d4d]">{editingEnabled ? t("common.enabled") : t("common.disabled")}</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Section 2: 流程编排 */}
              <div className="rounded-lg p-4" style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 2px, #fafafa 0px 0px 0px 1px" }}>
                <div className="mb-3 flex items-center justify-between gap-2">
                  <div className="text-xs font-medium uppercase tracking-wide text-[#808080]">{t("workflows.yamlOrchestration")}</div>
                  {!isEditing ? (
                    <div className="flex flex-wrap items-center gap-2">
                      <Button variant="outline" onClick={resetEditor}>{t("workflows.newTemplate")}</Button>
                      <select
                        className="ds-ring h-9 rounded-md bg-white px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
                        value={templateKey}
                        onChange={(e) => setTemplateKey(e.target.value as WorkflowTemplateKey)}
                      >
                        {templateOptions.map((item) => (
                          <option key={item.value} value={item.value}>{item.label}</option>
                        ))}
                      </select>
                      <Button variant="outline" onClick={applyTemplate}>{t("workflows.applyTemplate")}</Button>
                    </div>
                  ) : null}
                </div>
                <div className="grid grid-cols-1 gap-4 xl:grid-cols-[1fr_1fr]">
                  <div className="space-y-2">
                    <Textarea value={specText} onChange={(e) => setSpecText(e.target.value)} className="h-[520px] min-h-0 font-mono text-xs" />
                    {parsedSpec.error ? <div className="rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">{t("workflows.yamlError")}: {parsedSpec.error}</div> : null}
                    {error ? <div className="rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">{error}</div> : null}
                  </div>
                  <div className="h-[520px] rounded-lg" style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}>
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
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
