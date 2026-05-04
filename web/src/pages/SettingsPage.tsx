import { useEffect, useMemo, useState } from "react";
import { Copy } from "lucide-react";
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
    if (!ok) throw new Error("当前浏览器不支持自动复制");
  } finally {
    document.body.removeChild(textarea);
  }
}

export default function SettingsPage() {
  const toast = useToast();
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
      const msg = e instanceof Error ? e.message : "加载 API Key 失败";
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
          <CardTitle>账户设置</CardTitle>
          <CardDescription>修改管理员密码后，下次登录即生效。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <Input type="password" placeholder="新密码" value={password} onChange={(e) => setPassword(e.target.value)} />
          <Button
            onClick={async () => {
              try {
                await api.changePassword(password);
                setPasswordMsg({ type: "success", text: "已更新密码" });
                setPassword("");
                toast.success("密码已更新");
              } catch (e) {
                const t = e instanceof Error ? e.message : "更新失败";
                setPasswordMsg({ type: "error", text: t });
                toast.error(t);
              }
            }}
          >
            更新密码
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
          <CardTitle>API Key</CardTitle>
          <CardDescription>先创建 API Key，再通过 X-API-Key 调用 OpenAPI 申请证书。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-3 md:grid-cols-[1fr_220px_auto]">
            <Input placeholder="Key 名称（例如：ci-prod）" value={apiKeyName} onChange={(e) => setAPIKeyName(e.target.value)} />
            <Input type="datetime-local" value={apiKeyExpiresAt} onChange={(e) => setAPIKeyExpiresAt(e.target.value)} />
            <Button
              disabled={creatingKey}
              onClick={async () => {
                const name = apiKeyName.trim();
                if (!name) {
                  toast.error("请先输入 Key 名称");
                  return;
                }
                setCreatingKey(true);
                try {
                  const expiresAt = apiKeyExpiresAt ? new Date(apiKeyExpiresAt).toISOString() : undefined;
                  const res = await api.createAPIKey({ name, expiresAt });
                  setRevealedToken(res.token);
                  setAPIKeyName("");
                  setAPIKeyExpiresAt("");
                  toast.success("API Key 创建成功（明文仅显示一次）");
                  await loadAPIKeys();
                } catch (e) {
                  const msg = e instanceof Error ? e.message : "创建 API Key 失败";
                  toast.error(msg);
                } finally {
                  setCreatingKey(false);
                }
              }}
            >
              创建 Key
            </Button>
          </div>

          {revealedToken ? (
            <div className="space-y-2 rounded-md bg-[#f8fbff] p-3 shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px]">
              <div className="text-sm font-medium text-[#171717]">新创建的 Key（仅展示一次）</div>
              <div className="flex flex-col gap-2 md:flex-row">
                <Input readOnly value={revealedToken} />
                <Button
                  variant="outline"
                  onClick={async () => {
                    await copyText(revealedToken);
                    toast.success("已复制 API Key");
                  }}
                >
                  复制
                </Button>
              </div>
            </div>
          ) : null}

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>名称</TableHead>
                <TableHead>前缀</TableHead>
                <TableHead>状态</TableHead>
                <TableHead>过期时间</TableHead>
                <TableHead>最后使用</TableHead>
                <TableHead>创建时间</TableHead>
                <TableHead className="text-right">操作</TableHead>
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
                          toast.success("已吊销");
                          await loadAPIKeys();
                        } catch (e) {
                          const msg = e instanceof Error ? e.message : "吊销失败";
                          toast.error(msg);
                        }
                      }}
                    >
                      吊销
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
              {apiKeys.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="py-6 text-center text-sm text-[#666]">
                    {loadingKeys ? "加载中..." : "暂无 API Key"}
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
