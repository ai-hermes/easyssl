import { useEffect, useMemo, useState } from "react";
import { Copy } from "lucide-react";
import { useTranslation } from "react-i18next";
import i18n from "@/i18n";
import { api } from "@/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useToast } from "@/components/ui/toast";
import type { APIKeyItem } from "@/types";

function fmtTime(raw?: string) {
  if (!raw) return "-";
  const t = new Date(raw);
  if (Number.isNaN(t.getTime())) return raw;
  return t.toLocaleString();
}

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
  const [password, setPassword] = useState("");
  const [passwordMsg, setPasswordMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const [apiKeyName, setAPIKeyName] = useState("");
  const [apiKeyExpiresAt, setAPIKeyExpiresAt] = useState("");
  const [apiKeys, setAPIKeys] = useState<APIKeyItem[]>([]);
  const [loadingKeys, setLoadingKeys] = useState(false);
  const [creatingKey, setCreatingKey] = useState(false);
  const [revealedToken, setRevealedToken] = useState("");


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
  }, []);

  return (
    <div className="space-y-6">
      <Card className="max-w-lg">
        <CardHeader>
          <CardTitle>{t("settings.account.title")}</CardTitle>
          <CardDescription>{t("settings.account.desc")}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <Input type="password" placeholder={t("settings.account.newPassword")} value={password} onChange={(e) => setPassword(e.target.value)} />
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
                passwordMsg.type === "success" ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]" : "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]"
              }`}
            >
              {passwordMsg.text}
            </p>
          ) : null}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t("settings.apiKey.title")}</CardTitle>
          <CardDescription>{t("settings.apiKey.desc")}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-3 md:grid-cols-[1fr_220px_auto]">
            <Input placeholder={t("settings.apiKey.namePlaceholder")} value={apiKeyName} onChange={(e) => setAPIKeyName(e.target.value)} />
            <Input type="datetime-local" value={apiKeyExpiresAt} onChange={(e) => setAPIKeyExpiresAt(e.target.value)} />
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
            <div className="space-y-2 rounded-md bg-[#f8fbff] p-3 shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px]">
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
                  {t("common.copy")}
                </Button>
              </div>
            </div>
          ) : null}

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
                  <TableCell>{fmtTime(item.expiresAt)}</TableCell>
                  <TableCell>{fmtTime(item.lastUsedAt)}</TableCell>
                  <TableCell>{fmtTime(item.createdAt)}</TableCell>
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


        </CardContent>
      </Card>
    </div>
  );
}
