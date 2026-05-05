import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Check, ChevronDown, Clock3, Loader2, Search } from "lucide-react";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/toast";
import { StatusBadge } from "@/components/ui/status-badge";
import { formatTime } from "@/lib/time";
import type { Access, ProviderDefinition, ProviderField } from "@/types";
import { cn } from "@/lib/utils";

const formLabelClassName = "text-sm font-medium leading-5 text-[#171717]";
const formInputClassName =
  "h-9 rounded-md border-0 bg-white px-3 py-1 text-sm text-[#171717] placeholder:text-[#808080] shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px] transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]";
const formTextareaClassName =
  "min-h-[96px] rounded-md border-0 bg-white px-3 py-2 text-sm text-[#171717] placeholder:text-[#808080] shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]";

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
    provider: def?.id ?? "",
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

function ProviderSearchSelect({
  providers,
  value,
  onChange,
  placeholder,
}: {
  providers: ProviderDefinition[];
  value: string;
  onChange: (id: string) => void;
  placeholder: string;
}) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const selected = providers.find((provider) => provider.id === value);

  useEffect(() => {
    if (!open) {
      setQuery("");
      return;
    }
    const frame = window.requestAnimationFrame(() => {
      searchInputRef.current?.focus();
    });
    return () => window.cancelAnimationFrame(frame);
  }, [open]);

  const filteredProviders = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase();
    if (!normalizedQuery) return providers;
    return providers.filter(
      (provider) =>
        provider.label.toLowerCase().includes(normalizedQuery) ||
        provider.id.toLowerCase().includes(normalizedQuery),
    );
  }, [providers, query]);

  return (
    <Popover open={open} onOpenChange={setOpen} modal>
      <PopoverTrigger asChild>
        <button
          type="button"
          className={cn(
            formInputClassName,
            "flex items-center justify-between text-left",
            !selected && "text-[#808080]",
          )}
        >
          <span className="truncate">{selected?.label ?? placeholder}</span>
          <ChevronDown className="ml-2 h-4 w-4 shrink-0 text-[#808080]" />
        </button>
      </PopoverTrigger>
      <PopoverContent
        align="start"
        sideOffset={6}
        className="z-[60] w-[var(--radix-popover-trigger-width)] p-0 shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px,rgba(0,0,0,0.04)_0px_2px_2px,rgba(0,0,0,0.04)_0px_8px_8px_-8px,#fafafa_0px_0px_0px_1px]"
        onOpenAutoFocus={(event) => event.preventDefault()}
      >
        <div className="flex flex-col">
          <div className="flex items-center gap-2 px-3 py-3 shadow-[rgba(0,0,0,0.08)_0px_1px_0px]">
            <Search className="h-4 w-4 shrink-0 text-[#808080]" />
            <input
              ref={searchInputRef}
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              className="w-full bg-transparent text-sm text-[#171717] outline-none placeholder:text-[#808080]"
              placeholder="Search provider..."
            />
          </div>
          <div className="ds-scrollbar max-h-72 overflow-y-auto overscroll-contain p-1">
            {filteredProviders.length === 0 ? (
              <div className="px-3 py-2 text-sm text-[#808080]">No results found.</div>
            ) : (
              filteredProviders.map((provider) => (
                <button
                  key={provider.id}
                  type="button"
                  className={cn(
                    "flex w-full items-center justify-between rounded-md px-3 py-2 text-left text-sm text-[#171717] hover:bg-[#fafafa]",
                    provider.id === value && "bg-[#fafafa] font-medium",
                  )}
                  onClick={() => {
                    onChange(provider.id);
                    setOpen(false);
                  }}
                >
                  <span className="truncate">{provider.label}</span>
                  {provider.id === value ? <Check className="ml-2 h-4 w-4 shrink-0" /> : null}
                </button>
              ))
            )}
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}

function renderField(
  field: ProviderField,
  value: unknown,
  editing: boolean,
  onChange: (value: unknown) => void,
  t: (key: string) => string,
) {
  const label = `${field.label || field.name}${field.required ? " *" : ""}`;
  const placeholder =
    field.placeholder || (field.secret && editing ? `${label}${t("accesses.leaveBlankHint")}` : label);

  if (field.type === "checkbox") {
    return (
      <div key={field.name} className="flex flex-col gap-2">
        <label className={formLabelClassName} htmlFor={field.name}>
          {label}
        </label>
        <div className="flex h-9 items-center justify-start gap-3">
          <Switch checked={Boolean(value)} onCheckedChange={onChange} id={field.name} />
          <span className="text-sm text-[#666666]">
            {Boolean(value) ? t("common.enabled") : t("common.disabled")}
          </span>
        </div>
      </div>
    );
  }

  if (field.type === "select") {
    return (
      <div key={field.name} className="flex flex-col gap-2">
        <label className={formLabelClassName}>{label}</label>
        <select
          className={cn(formInputClassName, "appearance-none bg-white")}
          value={stringify(value)}
          onChange={(e) => onChange(e.target.value)}
        >
          {(field.options ?? []).map((item) => (
            <option key={item.value} value={item.value}>
              {item.label}
            </option>
          ))}
        </select>
      </div>
    );
  }

  if (field.type === "textarea") {
    return (
      <div key={field.name} className="flex flex-col gap-2 md:col-span-2">
        <label className={formLabelClassName}>{label}</label>
        <Textarea
          className={formTextareaClassName}
          placeholder={placeholder}
          value={stringify(value)}
          onChange={(e) => onChange(e.target.value)}
        />
      </div>
    );
  }

  return (
    <div key={field.name} className="flex flex-col gap-2">
      <label className={formLabelClassName}>{label}</label>
      <Input
        className={formInputClassName}
        type={field.type === "password" ? "password" : field.type === "number" ? "number" : "text"}
        placeholder={placeholder}
        value={stringify(value)}
        onChange={(e) => onChange(field.type === "number" ? e.target.value : e.target.value)}
      />
    </div>
  );
}

export default function AccessesPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const { t } = useTranslation();
  const { data } = useQuery({ queryKey: ["accesses"], queryFn: api.listAccesses });
  const { data: providerData } = useQuery({
    queryKey: ["providers", "access"],
    queryFn: () => api.listProviders("access"),
  });
  const providers = providerData?.items ?? [];
  const firstProvider = providers[0];
  const providerMap = useMemo(() => new Map(providers.map((item) => [item.id, item])), [providers]);
  const [form, setForm] = useState<AccessFormState>(() => emptyForm());
  const [notice, setNotice] = useState<{ type: "error" | "success" | "info"; text: string } | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const editing = Boolean(form.id);
  const selectedProvider = providerMap.get(form.provider);

  useEffect(() => {
    if (!dialogOpen || editing || form.provider || !firstProvider) return;
    setForm((current) => ({ ...current, provider: firstProvider.id, config: buildDefaultConfig(firstProvider) }));
  }, [dialogOpen, editing, firstProvider, form.provider]);

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
      setNotice({ type: "success", text: t("accesses.saveSuccess") });
      toast.success(t("accesses.saveSuccess"));
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : t("common.saveFailed");
      setNotice({ type: "error", text: msg });
      toast.error(msg);
    },
  });

  const testAccess = useMutation({
    mutationFn: (id: string) => api.testAccess(id),
    onMutate: () => {
      setNotice(null);
    },
    onSuccess: () => {
      toast.success(t("accesses.testSuccess"));
      qc.invalidateQueries({ queryKey: ["accesses"] });
    },
    onError: (e) => {
      const msg = e instanceof Error ? e.message : t("accesses.testFailed");
      toast.error(msg);
      qc.invalidateQueries({ queryKey: ["accesses"] });
    },
  });
  const testingAccessId = testAccess.isPending ? testAccess.variables : undefined;

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
          <h1 className="text-2xl font-semibold tracking-[-0.04em] text-[#171717]">{t("accesses.title")}</h1>
          <p className="mt-1 text-sm text-[#666]">{t("accesses.description")}</p>
        </div>
        <Button onClick={openCreate}>{t("accesses.addAccess")}</Button>
      </div>

      {notice ? (
        <div
          className={`rounded-md px-3 py-2 text-sm ${
            notice.type === "error"
              ? "bg-red-50 text-red-700"
              : notice.type === "success"
                ? "bg-green-50 text-green-700"
                : "bg-blue-50 text-blue-700"
          }`}
        >
          {notice.text}
        </div>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle>{t("accesses.title")}</CardTitle>
          <CardDescription>{t("accesses.listDescription")}</CardDescription>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="text-xs uppercase tracking-wide text-[#777]">
              <tr>
                <th className="px-2 pb-3 pt-1">{t("accesses.columns.name")}</th>
                <th className="px-2 pb-3 pt-1">{t("accesses.columns.provider")}</th>
                <th className="px-2 pb-3 pt-1">{t("accesses.columns.config")}</th>
                <th className="px-2 pb-3 pt-1">{t("accesses.columns.testStatus")}</th>
                <th className="px-2 pb-3 pt-1 text-right">{t("accesses.columns.actions")}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#eee]">
              {(data?.items ?? []).map((x) => {
                const def = providerMap.get(x.provider);
                return (
                  <tr key={x.id}>
                    <td className="px-2 py-3 font-medium text-[#171717]">{x.name}</td>
                    <td className="px-2 py-3 font-mono text-xs text-[#444]">{def?.label ?? x.provider}</td>
                    <td className="px-2 py-3 text-xs text-[#777]">
                      {Object.keys(x.config ?? {}).join(", ") || "-"}
                    </td>
                    <td className="px-2 py-3 align-top">
                      <div className="min-w-[176px]">
                        {x.lastTestResult ? (
                          <div className="inline-flex min-w-[176px] items-center justify-between gap-3 py-1">
                            <StatusBadge status={x.lastTestResult} className="w-fit shrink-0 px-2 py-0 text-[11px] leading-5" />
                            {x.lastTestedAt ? (
                              <span className="flex min-w-0 items-center justify-end gap-1 text-[11px] text-[#666666] [font-variant-numeric:tabular-nums]">
                                <Clock3 className="h-3 w-3 shrink-0 text-[#808080]" />
                                <span className="truncate">{formatTime(x.lastTestedAt)}</span>
                              </span>
                            ) : null}
                          </div>
                        ) : (
                          <div className="inline-flex min-w-[176px] items-center py-1">
                            <span className="font-mono text-[11px] uppercase tracking-[0.08em] text-[#808080]">
                              {t("accesses.testStatus.notTested")}
                            </span>
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="px-2 py-3">
                      <div className="flex justify-end gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          disabled={testAccess.isPending}
                          onClick={() => testAccess.mutate(x.id!)}
                        >
                          {testingAccessId === x.id ? (
                            <>
                              <Loader2 className="h-3.5 w-3.5 animate-spin" />
                              {t("common.loading")}
                            </>
                          ) : (
                            t("common.test")
                          )}
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => startEdit(x)}>
                          {t("common.edit")}
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={async () => {
                            await api.deleteAccess(x.id!);
                            qc.invalidateQueries({ queryKey: ["accesses"] });
                            toast.info(t("accesses.deleteSuccess"));
                          }}
                        >
                          {t("common.delete")}
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

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-h-[88vh] max-w-2xl overflow-hidden border-0 bg-white p-0 shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px,rgba(0,0,0,0.04)_0px_2px_2px,rgba(0,0,0,0.04)_0px_8px_8px_-8px,#fafafa_0px_0px_0px_1px] sm:rounded-lg">
          <DialogHeader className="px-6 pb-4 pt-6 text-left shadow-[rgba(0,0,0,0.08)_0px_1px_0px]">
            <DialogTitle className="text-[24px] font-semibold leading-8 tracking-[-0.96px] text-[#171717]">
              {editing ? t("accesses.editTitle") : t("accesses.createTitle")}
            </DialogTitle>
            <DialogDescription className="max-w-2xl text-sm leading-6 text-[#666666]">
              {t("accesses.dialogDescription")}
            </DialogDescription>
          </DialogHeader>

          <div className="ds-scrollbar max-h-[calc(88vh-152px)] overflow-y-auto px-6 py-6">
            <div className="space-y-6">
              <section className="space-y-4">
                <div className="flex flex-col gap-2">
                  <label className={formLabelClassName}>Provider</label>
                  <ProviderSearchSelect
                    providers={providers}
                    value={form.provider}
                    onChange={onProviderChange}
                    placeholder={t("accesses.selectProvider")}
                  />
                </div>

              </section>

              {selectedProvider ? (
                <>
                  <section className="space-y-4">
                    <div className="flex flex-col gap-2">
                      <label className={formLabelClassName}>Alias</label>
                      <Input
                        className={formInputClassName}
                        placeholder={t("accesses.namePlaceholder")}
                        value={form.name}
                        onChange={(e) => setForm((v) => ({ ...v, name: e.target.value }))}
                      />
                    </div>
                  </section>

                  <section className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2">
                      {(selectedProvider.accessFields ?? []).map((field) =>
                        renderField(
                          field,
                          readConfigValue(form.config, field),
                          editing,
                          (value) =>
                            setForm((current) => ({
                              ...current,
                              config: updateConfig(current.config, field, value),
                            })),
                          t,
                        ),
                      )}
                    </div>
                  </section>
                </>
              ) : null}
            </div>
          </div>

          <DialogFooter className="bg-white px-6 py-4 shadow-[rgba(0,0,0,0.08)_0px_-1px_0px] sm:justify-end sm:space-x-2">
            <Button
              variant="outline"
              className="h-9 rounded-md px-3 text-sm font-medium shadow-[rgb(235,235,235)_0px_0px_0px_1px]"
              onClick={() => setDialogOpen(false)}
            >
              {t("common.cancel")}
            </Button>
            <Button
              className="h-9 rounded-md bg-[#171717] px-4 text-sm font-medium text-white hover:bg-black"
              disabled={save.isPending || providers.length === 0 || !form.provider || !form.name.trim()}
              onClick={() => save.mutate()}
            >
              {editing ? t("common.saveChanges") : t("accesses.create")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
