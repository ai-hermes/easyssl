import * as Dialog from "@radix-ui/react-dialog";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/toast";
import type { Access, ProviderDefinition, ProviderField } from "@/types";

type AccessFormState = {
  id?: string;
  name: string;
  provider: string;
  config: Record<string, unknown>;
};

function defaultValue(field: ProviderField): unknown {
  if (field.default !== undefined) return field.default;
  if (field.type === "checkbox") return false;
  if (field.type === "number") return "";
  return "";
}

function buildDefaultConfig(def?: ProviderDefinition) {
  const config: Record<string, unknown> = {};
  for (const field of def?.accessFields ?? []) {
    const value = defaultValue(field);
    if (value !== "" && value !== undefined) config[field.name] = value;
  }
  return config;
}

function emptyForm(def?: ProviderDefinition): AccessFormState {
  return {
    name: "",
    provider: def?.id ?? "aliyun",
    config: buildDefaultConfig(def),
  };
}

function stringify(v: unknown) {
  if (v === null || v === undefined) return "";
  return String(v);
}

function readConfigValue(config: Record<string, unknown>, field: ProviderField) {
  const value = config[field.name];
  if (value === undefined || value === null) return defaultValue(field);
  return value;
}

function updateConfig(config: Record<string, unknown>, field: ProviderField, value: unknown) {
  const next = { ...config };
  next[field.name] = value;
  return next;
}

function renderField(field: ProviderField, value: unknown, editing: boolean, onChange: (value: unknown) => void) {
  const label = `${field.label || field.name}${field.required ? " *" : ""}`;
  const placeholder = field.placeholder || (field.secret && editing ? `${label}（留空表示不修改）` : label);
  if (field.type === "checkbox") {
    return (
      <label key={field.name} className="ds-ring flex h-9 items-center gap-2 rounded-md bg-white px-3 text-sm text-[#171717]">
        <input type="checkbox" checked={Boolean(value)} onChange={(e) => onChange(e.target.checked)} />
        {label}
      </label>
    );
  }
  if (field.type === "select") {
    return (
      <select
        key={field.name}
        className="ds-ring h-9 rounded-md bg-white px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
        value={stringify(value)}
        onChange={(e) => onChange(e.target.value)}
      >
        {(field.options ?? []).map((item) => (
          <option key={item.value} value={item.value}>
            {item.label}
          </option>
        ))}
      </select>
    );
  }
  if (field.type === "textarea") {
    return (
      <textarea
        key={field.name}
        className="ds-ring min-h-24 rounded-md bg-white px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
        placeholder={placeholder}
        value={stringify(value)}
        onChange={(e) => onChange(e.target.value)}
      />
    );
  }
  return (
    <Input
      key={field.name}
      type={field.type === "password" ? "password" : field.type === "number" ? "number" : "text"}
      placeholder={placeholder}
      value={stringify(value)}
      onChange={(e) => onChange(field.type === "number" ? e.target.value : e.target.value)}
    />
  );
}

export default function AccessesPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const { data } = useQuery({ queryKey: ["accesses"], queryFn: api.listAccesses });
  const { data: providerData } = useQuery({ queryKey: ["providers", "access"], queryFn: () => api.listProviders("access") });
  const providers = providerData?.items ?? [];
  const firstProvider = providers[0];
  const providerMap = useMemo(() => new Map(providers.map((item) => [item.id, item])), [providers]);
  const [form, setForm] = useState<AccessFormState>(() => emptyForm());
  const [notice, setNotice] = useState<{ type: "error" | "success" | "info"; text: string } | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const editing = Boolean(form.id);
  const selectedProvider = providerMap.get(form.provider) ?? firstProvider;

  const save = useMutation({
    mutationFn: async () => {
      setNotice(null);
      const payload: Access = {
        id: form.id,
        name: form.name.trim(),
        provider: form.provider,
        config: form.config,
      };
      return api.saveAccess(payload);
    },
    onSuccess: () => {
      setForm(emptyForm(firstProvider));
      setDialogOpen(false);
      qc.invalidateQueries({ queryKey: ["accesses"] });
      setNotice({ type: "success", text: "授权保存成功" });
      toast.success("授权保存成功");
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : "保存失败";
      setNotice({ type: "error", text: msg });
      toast.error(msg);
    },
  });

  const testAccess = useMutation({
    mutationFn: (id: string) => api.testAccess(id),
    onSuccess: () => {
      setNotice({ type: "success", text: "授权测试成功" });
      toast.success("授权测试成功");
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : "测试失败";
      setNotice({ type: "error", text: msg });
      toast.error(msg);
    },
  });

  const openCreate = () => {
    setForm(emptyForm(firstProvider));
    setDialogOpen(true);
    setNotice(null);
  };

  const startEdit = (access: Access) => {
    const def = providerMap.get(access.provider);
    setForm({
      id: access.id,
      name: access.name,
      provider: access.provider,
      config: { ...buildDefaultConfig(def), ...(access.config ?? {}) },
    });
    setDialogOpen(true);
    setNotice(null);
  };

  const onProviderChange = (provider: string) => {
    const def = providerMap.get(provider);
    setForm((value) => ({ ...value, provider, config: buildDefaultConfig(def) }));
  };

  return (
    <div className="space-y-5">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-[-0.04em] text-[#171717]">授权管理</h1>
          <p className="mt-1 text-sm text-[#666]">管理 DNS、云厂商、主机与面板授权，可用于证书申请或部署。</p>
        </div>
        <Button onClick={openCreate}>新增授权</Button>
      </div>

      {notice ? <div className={`rounded-md px-3 py-2 text-sm ${notice.type === "error" ? "bg-red-50 text-red-700" : notice.type === "success" ? "bg-green-50 text-green-700" : "bg-blue-50 text-blue-700"}`}>{notice.text}</div> : null}

      <Card>
        <CardHeader>
          <CardTitle>Accesses</CardTitle>
          <CardDescription>Provider credentials are stored server-side. Secret fields are masked after save.</CardDescription>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="text-xs uppercase tracking-wide text-[#777]">
              <tr>
                <th className="px-2 pb-3 pt-1">Name</th>
                <th className="px-2 pb-3 pt-1">Provider</th>
                <th className="px-2 pb-3 pt-1">Config</th>
                <th className="px-2 pb-3 pt-1 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#eee]">
              {(data?.items ?? []).map((x) => {
                const def = providerMap.get(x.provider);
                return (
                  <tr key={x.id}>
                    <td className="px-2 py-3 font-medium text-[#171717]">{x.name}</td>
                    <td className="px-2 py-3 font-mono text-xs text-[#444]">{def?.label ?? x.provider}</td>
                    <td className="px-2 py-3 text-xs text-[#777]">{Object.keys(x.config ?? {}).join(", ") || "-"}</td>
                    <td className="px-2 py-3">
                      <div className="flex justify-end gap-2">
                        <Button size="sm" variant="outline" disabled={testAccess.isPending} onClick={() => testAccess.mutate(x.id!)}>
                          测试
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => startEdit(x)}>
                          编辑
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={async () => {
                            await api.deleteAccess(x.id!);
                            qc.invalidateQueries({ queryKey: ["accesses"] });
                            toast.info("授权已删除");
                          }}
                        >
                          删除
                        </Button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </CardContent>
      </Card>

      <Dialog.Root open={dialogOpen} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-50 bg-black/35" />
          <Dialog.Content className="fixed left-1/2 top-1/2 z-50 max-h-[86vh] w-[min(900px,92vw)] -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-xl bg-white p-5 shadow-2xl focus:outline-none">
            <div className="mb-4">
              <Dialog.Title className="text-lg font-semibold tracking-[-0.02em]">{editing ? "编辑授权" : "新增授权"}</Dialog.Title>
              <Dialog.Description className="text-sm text-[#666]">配置完成后可用于证书申请或部署节点。新增 provider 表单由后端元数据自动生成。</Dialog.Description>
            </div>

            <div className="grid gap-2 md:grid-cols-2">
              <Input placeholder="授权名称" value={form.name} onChange={(e) => setForm((v) => ({ ...v, name: e.target.value }))} />
              <select
                className="ds-ring h-9 rounded-md bg-white px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
                value={form.provider}
                onChange={(e) => onProviderChange(e.target.value)}
              >
                {providers.map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.label} ({item.id})
                  </option>
                ))}
              </select>

              {(selectedProvider?.accessFields ?? []).map((field) =>
                renderField(field, readConfigValue(form.config, field), editing, (value) => setForm((current) => ({ ...current, config: updateConfig(current.config, field, value) }))),
              )}
            </div>

            <div className="mt-3 text-xs text-[#777]">
              Provider: <span className="font-mono">{selectedProvider?.id ?? form.provider}</span>
              {selectedProvider?.accessProviderId && selectedProvider.accessProviderId !== selectedProvider.id ? <span> · Access: <span className="font-mono">{selectedProvider.accessProviderId}</span></span> : null}
            </div>

            <div className="mt-5 flex justify-end gap-2">
              <Button variant="outline" onClick={() => setDialogOpen(false)}>
                取消
              </Button>
              <Button disabled={save.isPending || providers.length === 0} onClick={() => save.mutate()}>
                {editing ? "保存修改" : "新增授权"}
              </Button>
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
