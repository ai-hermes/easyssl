import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { ArrowLeft, ChevronDown, ChevronRight, Play, Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Separator } from "@/components/ui/separator";
import { getToken } from "@/api/client";

interface Schema {
  type?: string;
  $ref?: string;
  properties?: Record<string, Schema>;
  items?: Schema;
  description?: string;
  required?: string[];
  additionalProperties?: boolean | Schema;
}

interface Parameter {
  name: string;
  in: string;
  required?: boolean;
  description?: string;
  schema?: Schema;
  type?: string;
}

interface Response {
  description: string;
  schema?: Schema;
}

interface Operation {
  summary: string;
  tags?: string[];
  parameters?: Parameter[];
  responses: Record<string, Response>;
  security?: Array<Record<string, string[]>>;
  consumes?: string[];
  produces?: string[];
}

interface PathItem {
  get?: Operation;
  post?: Operation;
  put?: Operation;
  delete?: Operation;
}

interface SecurityDef {
  type: string;
  in?: string;
  name?: string;
  description?: string;
}

interface OpenAPISpec {
  swagger: string;
  info: { title: string; description?: string; version: string };
  basePath?: string;
  securityDefinitions?: Record<string, SecurityDef>;
  paths: Record<string, PathItem>;
  definitions?: Record<string, Schema>;
}

const methodColors: Record<string, string> = {
  get: "bg-[#e9f9ee] text-[#116329]",
  post: "bg-[#ebf5ff] text-[#0068d6]",
  put: "bg-[#fff4e8] text-[#8f4d00]",
  delete: "bg-[#ffecec] text-[#a01616]",
};

function MethodBadge({ method }: { method: string }) {
  const cls = methodColors[method.toLowerCase()] || "bg-[#f5f5f5] text-[#666]";
  return (
    <span className={`inline-flex min-w-[56px] items-center justify-center rounded px-2 py-0.5 text-xs font-semibold uppercase ${cls}`}>
      {method}
    </span>
  );
}

function resolveRef(ref: string, defs: Record<string, Schema>): Schema | undefined {
  return defs[ref.replace("#/definitions/", "")];
}

function SchemaViewer({ schema, defs, level = 0 }: { schema?: Schema; defs: Record<string, Schema>; level?: number }) {
  const { t } = useTranslation();
  if (!schema) return <span className="text-xs text-[#999]">-</span>;

  let s = schema;
  if (s.$ref) {
    s = resolveRef(s.$ref, defs) || s;
  }

  if (s.type === "object" && s.properties) {
    return (
      <div className={`space-y-1 ${level > 0 ? "ml-3 border-l border-[#eaeaea] pl-2" : ""}`}>
        {Object.entries(s.properties).map(([key, val]) => (
          <div key={key} className="text-xs">
            <div className="flex items-center gap-1.5">
              <code className="font-semibold text-[#171717]">{key}</code>
              <TypeTag schema={val} defs={defs} />
              {s.required?.includes(key) && <Badge variant="destructive" className="h-4 text-[9px] px-1">{t("docs.required")}</Badge>}
              {val.description && <span className="text-[#999]">{val.description}</span>}
            </div>
            {(val.type === "object" || val.type === "array" || val.$ref) && (
              <div className="mt-0.5">
                <SchemaViewer schema={val} defs={defs} level={level + 1} />
              </div>
            )}
          </div>
        ))}
      </div>
    );
  }

  if (s.type === "array" && s.items) {
    return (
      <div className={`${level > 0 ? "ml-3 border-l border-[#eaeaea] pl-2" : ""}`}>
        <div className="text-xs text-[#999] mb-0.5">{t("docs.arrayItems")}</div>
        <SchemaViewer schema={s.items} defs={defs} level={level + 1} />
      </div>
    );
  }

  return <span className="text-xs text-[#999]">{s.type || "object"}</span>;
}

function TypeTag({ schema, defs }: { schema?: Schema; defs: Record<string, Schema> }) {
  if (!schema) return null;
  if (schema.$ref) {
    const name = schema.$ref.replace("#/definitions/", "");
    return <Badge variant="outline" className="h-4 text-[9px] px-1">{name}</Badge>;
  }
  if (schema.type === "array" && schema.items?.$ref) {
    return <Badge variant="outline" className="h-4 text-[9px] px-1">{schema.items.$ref.replace("#/definitions/", "")}[]</Badge>;
  }
  return <Badge variant="outline" className="h-4 text-[9px] px-1">{schema.type || "object"}{schema.type === "array" ? "[]" : ""}</Badge>;
}

function smartString(key: string): string {
  const k = key.toLowerCase();
  if (k.includes("email")) return "admin@example.com";
  if (k.includes("password")) return "1234567890";
  if (k.includes("provider")) return "tencentcloud";
  if (k.includes("domain")) return "example.com";
  if (k.includes("url") || k.includes("endpoint")) return "https://example.com";
  if (k.includes("format")) return "PEM";
  if (k.includes("trigger") && k.includes("cron")) return "0 0 * * *";
  if (k.includes("trigger")) return "manual";
  if (k.includes("caprovider") || k.includes("ca_provider")) return "letsencrypt";
  if (k.includes("keyalgorithm") || k.includes("key_algorithm")) return "RSA2048";
  if (k.includes("contactemail") || k.includes("contact_email")) return "admin@example.com";
  if (k.includes("accessid") || k.includes("access_id")) return "access-uuid-1234";
  if (k.includes("workflowid") || k.includes("workflow_id")) return "workflow-uuid-5678";
  if (k.includes("runid") || k.includes("run_id")) return "run-uuid-abcd";
  if (k.includes("certificateid") || k.includes("certificate_id")) return "cert-uuid-efgh";
  if (k.includes("id")) return "550e8400-e29b-41d4-a716-446655440000";
  if (k.includes("name") && k.includes("access")) return "my-dns-access";
  if (k.includes("name") && k.includes("workflow")) return "ssl-renewal";
  if (k.includes("name") && k.includes("key")) return "prod-api-key";
  if (k.includes("name")) return "example-name";
  if (k.includes("prefix")) return "ek_";
  if (k.includes("token")) return "eyJhbGciOiJIUzI1NiIs...";
  if (k.includes("status")) return "active";
  if (k.includes("role")) return "admin";
  if (k.includes("description")) return "This is an example description.";
  if (k.includes("message")) return "Operation completed successfully.";
  if (k.includes("eventtype") || k.includes("event_type")) return "started";
  if (k.includes("error")) return "";
  if (k.includes("nodeid") || k.includes("node_id")) return "node-1";
  if (k.includes("action")) return "apply";
  if (k.includes("graph")) return "{}";
  if (k.includes("output")) return "{}";
  if (k.includes("payload")) return "{}";
  if (k.includes("config")) return "{}";
  if (k.includes("reserve")) return "";
  if (k.includes("subject")) return "*.example.com";
  if (k.includes("serial")) return "00:01:02:03:04:05:06:07";
  if (k.includes("issuer")) return "Let's Encrypt";
  if (k.includes("source")) return "acme";
  return "string";
}

function smartNumber(key: string): number {
  const k = key.toLowerCase();
  if (k.includes("timeout")) return 60;
  if (k.includes("ttl")) return 300;
  if (k.includes("port")) return 443;
  if (k.includes("total")) return 10;
  if (k.includes("propagation")) return 120;
  return 1;
}

function buildExample(schema: Schema, defs: Record<string, Schema>, key = ""): unknown {
  let s: Schema | undefined = schema;
  if (s?.$ref) s = resolveRef(s.$ref, defs);
  if (!s) return null;

  if (s.type === "object") {
    if (s.properties) {
      const obj: Record<string, unknown> = {};
      Object.entries(s.properties).forEach(([k, v]) => {
        obj[k] = buildExample(v, defs, k);
      });
      return obj;
    }
    // additionalProperties (e.g. map[string]interface{})
    if (s.additionalProperties) {
      return { key: "value" };
    }
    return {};
  }
  if (s.type === "array") {
    const item = s.items ? buildExample(s.items, defs, key) : null;
    if (key.toLowerCase().includes("domain")) return ["ssl1.example.com", "*.example.com"];
    return item !== null ? [item] : [];
  }
  if (s.type === "string") return smartString(key);
  if (s.type === "integer" || s.type === "number") return smartNumber(key);
  if (s.type === "boolean") return false;
  return null;
}

function TryItOut({
  method,
  path,
  op,
  basePath,
  defs,
}: {
  method: string;
  path: string;
  op: Operation;
  basePath: string;
  defs: Record<string, Schema>;
}) {
  const { t } = useTranslation();
  const [params, setParams] = useState<Record<string, string>>({});
  const [body, setBody] = useState("{}");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ status: number; data: unknown } | null>(null);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const fullPath = basePath + path;
  const pathParams = op.parameters?.filter((p) => p.in === "path") || [];
  const queryParams = op.parameters?.filter((p) => p.in === "query") || [];
  const bodyParam = op.parameters?.find((p) => p.in === "body");

  useEffect(() => {
    if (bodyParam?.schema) {
      try {
        const ex = buildExample(bodyParam.schema, defs);
        setBody(JSON.stringify(ex, null, 2));
      } catch {
        setBody("{}");
      }
    }
  }, [bodyParam, defs]);

  async function send() {
    setLoading(true);
    setError("");
    setResult(null);
    try {
      let url = fullPath;
      pathParams.forEach((p) => {
        url = url.replace(`{${p.name}}`, encodeURIComponent(params[p.name] || ""));
      });
      const q = new URLSearchParams();
      queryParams.forEach((p) => {
        if (params[p.name]) q.set(p.name, params[p.name]);
      });
      if (q.toString()) url += `?${q.toString()}`;

      const headers: Record<string, string> = {};
      if (op.consumes?.includes("application/json")) {
        headers["Content-Type"] = "application/json";
      }
      const token = getToken();
      if (token) headers["Authorization"] = `Bearer ${token}`;

      const init: RequestInit = {
        method: method.toUpperCase(),
        headers,
      };
      if (["post", "put"].includes(method.toLowerCase()) && bodyParam) {
        init.body = body;
      }

      const resp = await fetch(url, init);
      let data: unknown;
      try {
        data = await resp.json();
      } catch {
        data = await resp.text();
      }
      setResult({ status: resp.status, data });
    } catch (e) {
      setError(e instanceof Error ? e.message : t("docs.requestFailed"));
    } finally {
      setLoading(false);
    }
  }

  const allParams = [...pathParams, ...queryParams];

  return (
    <div className="space-y-3 pt-2">
      <Separator />
      <div className="flex items-center gap-2">
        <Play className="h-3.5 w-3.5 text-[#0068d6]" />
        <span className="text-xs font-medium text-[#171717]">{t("docs.tryItOut")}</span>
      </div>

      {allParams.length > 0 && (
        <div className="space-y-2">
          <div className="text-xs font-medium text-[#999]">{t("docs.parameters")}</div>
          <div className="grid gap-2">
            {allParams.map((p) => (
              <div key={p.name} className="grid grid-cols-[120px_1fr] items-center gap-2">
                <label className="text-xs text-[#666]">
                  {p.name} <span className="text-[#999]">({p.in})</span>
                  {p.required && <span className="text-[#a01616]">*</span>}
                </label>
                <Input
                  size-condensed
                  placeholder={p.description || p.name}
                  value={params[p.name] || ""}
                  onChange={(e) => setParams((prev) => ({ ...prev, [p.name]: e.target.value }))}
                  className="h-7 text-xs"
                />
              </div>
            ))}
          </div>
        </div>
      )}

      {bodyParam && (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-[#999]">{t("docs.requestBody")}</span>
            <Button
              variant="ghost"
              size="sm"
              className="h-6 text-[10px]"
              onClick={() => {
                if (bodyParam.schema) {
                  try {
                    setBody(JSON.stringify(buildExample(bodyParam.schema, defs), null, 2));
                  } catch { /* noop */ }
                }
              }}
            >
              {t("docs.resetExample")}
            </Button>
          </div>
          <Textarea
            value={body}
            onChange={(e) => setBody(e.target.value)}
            className="min-h-[120px] font-mono text-xs"
          />
        </div>
      )}

      <Button size="sm" onClick={send} disabled={loading} className="text-xs">
        {loading ? t("docs.sending") : t("docs.sendRequest")}
      </Button>

      {error && (
        <div className="rounded-md bg-[#ffecec] px-3 py-2 text-xs text-[#a01616]">{error}</div>
      )}

      {result && (
        <div className="space-y-1.5">
          <div className="flex items-center gap-2">
            <Badge variant={result.status >= 200 && result.status < 300 ? "secondary" : "destructive"} className="text-[10px]">
              {result.status}
            </Badge>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={async () => {
                await navigator.clipboard.writeText(JSON.stringify(result.data, null, 2));
                setCopied(true);
                setTimeout(() => setCopied(false), 1500);
              }}
            >
              {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
            </Button>
          </div>
          <pre className="max-h-[300px] overflow-auto rounded-md bg-[#fafafa] p-3 text-xs font-mono text-[#171717] shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px]">
            {JSON.stringify(result.data, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}

export default function DocsPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [spec, setSpec] = useState<OpenAPISpec | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});

  useEffect(() => {
    fetch("/api/openapi.json")
      .then((r) => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`);
        return r.json();
      })
      .then((d) => setSpec(d))
      .catch((e) => setError(e instanceof Error ? e.message : t("common.unknownError")))
      .finally(() => setLoading(false));
  }, [t]);

  const grouped = useMemo(() => {
    if (!spec) return [] as Array<{ tag: string; items: Array<{ id: string; method: string; path: string; op: Operation }> }>;
    const map = new Map<string, Array<{ id: string; method: string; path: string; op: Operation }>>();
    Object.entries(spec.paths).forEach(([path, item]) => {
      (Object.keys(item) as Array<keyof PathItem>).forEach((method) => {
        const op = item[method];
        if (!op) return;
        const tag = op.tags?.[0] || "Default";
        if (!map.has(tag)) map.set(tag, []);
        map.get(tag)!.push({ id: `${method}-${path}`, method, path, op });
      });
    });
    return Array.from(map.entries()).map(([tag, items]) => ({ tag, items }));
  }, [spec]);

  const toggle = (id: string) => setExpanded((p) => ({ ...p, [id]: !p[id] }));

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center text-sm text-[#666]">{t("docs.loading")}</div>
    );
  }

  if (error || !spec) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4 text-sm text-[#a01616]">
        <p>{t("docs.loadFailed", { error: error || t("common.unknownError") })}</p>
        <Button variant="outline" onClick={() => navigate(-1)}>{t("common.back")}</Button>
      </div>
    );
  }

  const defs = spec.definitions || {};

  return (
    <div className="min-h-screen bg-[var(--ds-bg)]">
      <div className="sticky top-0 z-50 border-b bg-white/95 backdrop-blur" style={{ boxShadow: "rgba(0,0,0,0.08)_0px_0px_0px_1px" }}>
        <div className="mx-auto flex max-w-[1200px] items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="sm" onClick={() => navigate("/settings")}>
              <ArrowLeft className="mr-1 h-4 w-4" />
              {t("common.back")}
            </Button>
            <span className="text-base font-semibold text-[#171717]">{t("docs.pageTitle", { title: spec.info.title })}</span>
            <Badge variant="secondary">v{spec.info.version}</Badge>
          </div>
        </div>
      </div>

      <main className="mx-auto max-w-[1200px] space-y-6 p-4 md:p-6">
        {spec.info.description && (
          <Card>
            <CardContent className="py-4 text-sm leading-relaxed text-[#666] whitespace-pre-line">
              {spec.info.description}
            </CardContent>
          </Card>
        )}

        {spec.securityDefinitions && (
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{t("docs.authMethods")}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {Object.entries(spec.securityDefinitions).map(([key, def]) => (
                <div key={key} className="flex items-center gap-2 text-sm">
                  <Badge variant="outline">{key}</Badge>
                  <span className="text-[#666]">
                    {def.description}（{def.type}，Header: <code className="rounded bg-[#f5f5f5] px-1 py-0.5 text-xs">{def.name}</code>）
                  </span>
                </div>
              ))}
            </CardContent>
          </Card>
        )}

        {grouped.map(({ tag, items }) => (
          <Card key={tag}>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">{tag}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-1">
              {items.map(({ id, method, path, op }) => {
                const isOpen = expanded[id];
                const fullPath = (spec.basePath || "") + path;
                return (
                  <div key={id} className="rounded-md border border-[#f1f1f1]">
                    <button
                      onClick={() => toggle(id)}
                      className="flex w-full items-center gap-3 px-3 py-2.5 text-left transition hover:bg-[#fafafa]"
                    >
                      {isOpen ? <ChevronDown className="h-4 w-4 text-[#999]" /> : <ChevronRight className="h-4 w-4 text-[#999]" />}
                      <MethodBadge method={method} />
                      <code className="text-sm text-[#171717]">{fullPath}</code>
                      <span className="ml-auto text-xs text-[#666]">{op.summary}</span>
                    </button>
                    {isOpen && (
                      <div className="border-t border-[#f1f1f1] bg-[#fafafa] px-3 py-3 text-sm space-y-4">
                        {op.consumes && (
                          <div className="flex items-center gap-2">
                            <span className="text-xs font-medium text-[#999]">{t("docs.consumes")}</span>
                            <span className="text-xs text-[#666]">{op.consumes.join(", ")}</span>
                          </div>
                        )}
                        {op.produces && (
                          <div className="flex items-center gap-2">
                            <span className="text-xs font-medium text-[#999]">{t("docs.produces")}</span>
                            <span className="text-xs text-[#666]">{op.produces.join(", ")}</span>
                          </div>
                        )}

                        {op.parameters && op.parameters.length > 0 && (
                          <div>
                            <div className="mb-1.5 text-xs font-medium text-[#999]">{t("docs.parameters")}</div>
                            <div className="rounded-md border border-[#eaeaea] bg-white divide-y divide-[#f1f1f1]">
                              {op.parameters.map((p, i) => (
                                <div key={i} className="px-3 py-2">
                                  <div className="flex items-center gap-2 text-xs">
                                    <code className="min-w-[100px] font-semibold text-[#171717]">{p.name}</code>
                                    <Badge variant="outline" className="text-[10px]">{p.in}</Badge>
                                    {p.required && <Badge variant="destructive" className="text-[10px]">{t("docs.required")}</Badge>}
                                    <TypeTag schema={p.schema} defs={defs} />
                                  </div>
                                  {p.description && <div className="mt-0.5 text-[11px] text-[#999]">{p.description}</div>}
                                  {(p.schema?.type === "object" || p.schema?.$ref) && p.in === "body" && (
                                    <div className="mt-1.5">
                                      <SchemaViewer schema={p.schema} defs={defs} />
                                    </div>
                                  )}
                                </div>
                              ))}
                            </div>
                          </div>
                        )}

                        <div>
                          <div className="mb-1.5 text-xs font-medium text-[#999]">{t("docs.responses")}</div>
                          <div className="rounded-md border border-[#eaeaea] bg-white divide-y divide-[#f1f1f1]">
                            {Object.entries(op.responses).map(([code, res]) => (
                              <div key={code} className="px-3 py-2">
                                <div className="flex items-center gap-2 text-xs">
                                  <Badge variant="secondary" className="text-[10px]">{code}</Badge>
                                  <span className="text-[#666]">{res.description}</span>
                                  {res.schema && <TypeTag schema={res.schema} defs={defs} />}
                                </div>
                                {res.schema && (
                                  <div className="mt-1.5">
                                    <SchemaViewer schema={res.schema} defs={defs} />
                                  </div>
                                )}
                              </div>
                            ))}
                          </div>
                        </div>

                        <TryItOut method={method} path={path} op={op} basePath={spec.basePath || ""} defs={defs} />
                      </div>
                    )}
                  </div>
                );
              })}
            </CardContent>
          </Card>
        ))}
      </main>
    </div>
  );
}
