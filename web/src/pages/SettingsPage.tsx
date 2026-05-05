import { useEffect, useState } from "react";
import { CalendarIcon, Copy, KeyRound, User, GitBranch } from "lucide-react";
import { useTranslation } from "react-i18next";
import { format } from "date-fns";
import i18n from "@/i18n";
import { api } from "@/api";
import { getRole } from "@/api/client";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useToast } from "@/components/ui/toast";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import type { APIKeyItem } from "@/types";
import { cn } from "@/lib/utils";
import { formatTime } from "@/lib/time";

async function copyText(text: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text);
    return;
  }
  const textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "");
  textarea.style.position = "fixed";
  textarea.style.left = "-9999px";
  textarea.style.top = "0";
  document.body.appendChild(textarea);
  textarea.focus();
  textarea.select();
  try {
    const ok = document.execCommand("copy");
    if (!ok) throw new Error(i18n.t("common.copyFailed"));
  } finally {
    document.body.removeChild(textarea);
  }
}

export default function SettingsPage() {
  const toast = useToast();
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<"profile" | "openapi">("profile");

  const [password, setPassword] = useState("");
  const [passwordMsg, setPasswordMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const [apiKeyName, setAPIKeyName] = useState("");
  const [apiKeyExpiresAt, setAPIKeyExpiresAt] = useState("");
  const [apiKeys, setAPIKeys] = useState<APIKeyItem[]>([]);
  const [loadingKeys, setLoadingKeys] = useState(false);
  const [creatingKey, setCreatingKey] = useState(false);
  const [revealedToken, setRevealedToken] = useState("");
  const [version, setVersion] = useState("");
  const [versionDetail, setVersionDetail] = useState("");
  const [versionUrl, setVersionUrl] = useState("");
  const isAdmin = getRole() === "admin";

  async function loadAPIKeys() {
    setLoadingKeys(true);
    try {
      const res = await api.listAPIKeys();
      setAPIKeys(res.items || []);
    } catch (e) {
      const msg = e instanceof Error ? e.message : t("settings.apiKey.loadFailed");
      toast.error(msg);
    } finally {
      setLoadingKeys(false);
    }
  }

  useEffect(() => {
    void loadAPIKeys();
    if (isAdmin) {
      api.version().then((res) => {
        setVersion(res.version);
        setVersionDetail(res.detail);
        const url = res.tagUrl || res.branchUrl || res.commitUrl;
        setVersionUrl(url);
      }).catch(() => {});
    }
  }, []);

  const navItems: { key: "profile" | "openapi"; label: string; icon: React.ReactNode }[] = [
    { key: "profile", label: t("settings.nav.profile"), icon: <User size={16} /> },
    { key: "openapi", label: t("settings.nav.openapi"), icon: <KeyRound size={16} /> },
  ];

  return (
    <div className="grid gap-8 md:grid-cols-[240px_1fr]">
      {/* Left Sidebar */}
      <div
        className="h-fit rounded-lg bg-white p-2"
        style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}
      >
        <nav className="flex flex-col gap-1">
          {navItems.map((item) => {
            const isActive = activeTab === item.key;
            return (
              <button
                key={item.key}
                onClick={() => setActiveTab(item.key)}
                className={`flex items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition ${
                  isActive
                    ? "bg-white font-medium text-[#171717] shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px]"
                    : "text-[#666] hover:bg-[#fafafa] hover:text-[#171717]"
                }`}
              >
                {item.icon}
                {item.label}
              </button>
            );
          })}
        </nav>
        {isAdmin && version ? (
          <div className="mt-auto border-t border-[#ebebeb] pt-3">
            <div className="flex items-center gap-1.5 text-xs text-[#666]">
              <GitBranch size={12} />
              {versionUrl ? (
                <a
                  href={versionUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  title={versionDetail}
                  className="font-mono hover:text-[#171717] hover:underline"
                >
                  {version}
                </a>
              ) : (
                <span className="font-mono" title={versionDetail}>{version}</span>
              )}
            </div>
          </div>
        ) : null}
      </div>

      {/* Right Content */}
      <div className="min-w-0 space-y-6">
        {activeTab === "profile" && (
          <div className="space-y-6">
            <div>
              <h2 className="text-[24px] font-semibold tracking-[-0.04em] text-[#171717]">
                {t("settings.account.title")}
              </h2>
              <p className="mt-1 text-sm text-[#666]">{t("settings.account.desc")}</p>
            </div>

            <div
              className="rounded-lg bg-white p-6"
              style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}
            >
              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-[#171717]">{t("settings.account.newPassword")}</label>
                  <Input
                    type="password"
                    placeholder={t("settings.account.newPassword")}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                </div>
                <Button
                  onClick={async () => {
                    try {
                      await api.changePassword(password);
                      setPasswordMsg({ type: "success", text: t("settings.account.success") });
                      setPassword("");
                      toast.success(t("settings.account.success"));
                    } catch (e) {
                      const errMsg = e instanceof Error ? e.message : t("settings.account.error");
                      setPasswordMsg({ type: "error", text: errMsg });
                      toast.error(errMsg);
                    }
                  }}
                >
                  {t("settings.account.update")}
                </Button>
                {passwordMsg ? (
                  <p
                    className={`rounded-md px-3 py-2 text-sm ${
                      passwordMsg.type === "success"
                        ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]"
                        : "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]"
                    }`}
                  >
                    {passwordMsg.text}
                  </p>
                ) : null}
              </div>
            </div>
          </div>
        )}

        {activeTab === "openapi" && (
          <div className="space-y-6">
            <div>
              <h2 className="text-[24px] font-semibold tracking-[-0.04em] text-[#171717]">
                {t("settings.apiKey.title")}
              </h2>
              <p className="mt-1 text-sm text-[#666]">{t("settings.apiKey.desc")}</p>
            </div>

            <div
              className="rounded-lg bg-white p-6"
              style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}
            >
              <div className="space-y-6">
                <div className="grid gap-3 md:grid-cols-[1fr_auto_auto]">
                  <Input
                    placeholder={t("settings.apiKey.namePlaceholder")}
                    value={apiKeyName}
                    onChange={(e) => setAPIKeyName(e.target.value)}
                  />
                  <Popover>
                    <PopoverTrigger asChild>
                      <Button
                        variant={"outline"}
                        className={cn(
                          "w-[220px] justify-start text-left font-normal",
                          !apiKeyExpiresAt && "text-muted-foreground"
                        )}
                      >
                        <CalendarIcon className="mr-2 h-4 w-4" />
                        {apiKeyExpiresAt ? format(new Date(apiKeyExpiresAt), "yyyy-MM-dd") : t("settings.apiKey.expiresAt")}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={apiKeyExpiresAt ? new Date(apiKeyExpiresAt) : undefined}
                        onSelect={(date) => {
                          if (date) {
                            // Set time to end of day (23:59:59)
                            date.setHours(23, 59, 59, 0);
                            setAPIKeyExpiresAt(date.toISOString().slice(0, 16));
                          } else {
                            setAPIKeyExpiresAt("");
                          }
                        }}
                        initialFocus
                      />
                    </PopoverContent>
                  </Popover>
                  <Button
                    disabled={creatingKey}
                    onClick={async () => {
                      const name = apiKeyName.trim();
                      if (!name) {
                        toast.error(t("settings.apiKey.enterName"));
                        return;
                      }
                      setCreatingKey(true);
                      try {
                        const expiresAt = apiKeyExpiresAt ? new Date(apiKeyExpiresAt).toISOString() : undefined;
                        const res = await api.createAPIKey({ name, expiresAt });
                        setRevealedToken(res.token);
                        setAPIKeyName("");
                        setAPIKeyExpiresAt("");
                        toast.success(t("settings.apiKey.createSuccess"));
                        await loadAPIKeys();
                      } catch (e) {
                        const msg = e instanceof Error ? e.message : t("settings.apiKey.createFailed");
                        toast.error(msg);
                      } finally {
                        setCreatingKey(false);
                      }
                    }}
                  >
                    {t("settings.apiKey.create")}
                  </Button>
                </div>

                {revealedToken ? (
                  <div
                    className="space-y-2 rounded-md bg-[#f8fbff] p-3"
                    style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}
                  >
                    <div className="text-sm font-medium text-[#171717]">{t("settings.apiKey.newToken")}</div>
                    <div className="flex flex-col gap-2 md:flex-row">
                      <Input readOnly value={revealedToken} />
                      <Button
                        variant="outline"
                        onClick={async () => {
                          await copyText(revealedToken);
                          toast.success(t("settings.apiKey.copySuccess"));
                        }}
                      >
                        <Copy size={14} className="mr-1" />
                        {t("common.copy")}
                      </Button>
                    </div>
                  </div>
                ) : null}

                <Separator style={{ backgroundColor: "#ebebeb" }} />

                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t("settings.apiKey.columns.name")}</TableHead>
                      <TableHead>{t("settings.apiKey.columns.prefix")}</TableHead>
                      <TableHead>{t("settings.apiKey.columns.status")}</TableHead>
                      <TableHead>{t("settings.apiKey.columns.expiresAt")}</TableHead>
                      <TableHead>{t("settings.apiKey.columns.lastUsed")}</TableHead>
                      <TableHead>{t("settings.apiKey.columns.createdAt")}</TableHead>
                      <TableHead className="text-right">{t("settings.apiKey.columns.actions")}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {apiKeys.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell>{item.name}</TableCell>
                        <TableCell className="font-mono text-xs">{item.prefix}</TableCell>
                        <TableCell>
                          <Badge variant={item.status === "active" ? "secondary" : "destructive"}>{item.status}</Badge>
                        </TableCell>
                        <TableCell>{formatTime(item.expiresAt)}</TableCell>
                        <TableCell>{formatTime(item.lastUsedAt)}</TableCell>
                        <TableCell>{formatTime(item.createdAt)}</TableCell>
                        <TableCell className="text-right">
                          <Button
                            size="sm"
                            variant="outline"
                            disabled={item.status !== "active"}
                            onClick={async () => {
                              try {
                                await api.revokeAPIKey(item.id);
                                toast.success(t("settings.apiKey.revoked"));
                                await loadAPIKeys();
                              } catch (e) {
                                const msg = e instanceof Error ? e.message : t("settings.apiKey.revokeFailed");
                                toast.error(msg);
                              }
                            }}
                          >
                            {t("settings.apiKey.revoke")}
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                    {apiKeys.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="py-6 text-center text-sm text-[#666]">
                          {loadingKeys ? t("common.loading") : t("settings.apiKey.noKeys")}
                        </TableCell>
                      </TableRow>
                    ) : null}
                  </TableBody>
                </Table>
              </div>
            </div>
          </div>
        )}

      </div>
    </div>
  );
}
