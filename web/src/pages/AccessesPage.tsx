import * as Dialog from "@radix-ui/react-dialog";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/toast";
import type { Access } from "@/types";

const PROVIDERS = [
  { value: "aliyun", label: "Aliyun" },
  { value: "tencentcloud", label: "Tencent Cloud DNSPod" },
  { value: "qiniu", label: "Qiniu" },
  { value: "ssh", label: "SSH" },
];

type AccessFormState = {
  id?: string;
  name: string;
  provider: string;
  accessKeyId: string;
  accessKeySecret: string;
  resourceGroupId: string;
  region: string;
  tencentSecretId: string;
  tencentSecretKey: string;
  tencentSessionToken: string;
  tencentRegion: string;
  qiniuAccessKey: string;
  qiniuSecretKey: string;
  host: string;
  port: string;
  username: string;
  authMethod: "password" | "key";
  password: string;
  key: string;
  keyPassphrase: string;
};

function emptyForm(): AccessFormState {
  return {
    name: "",
    provider: "aliyun",
    accessKeyId: "",
    accessKeySecret: "",
    resourceGroupId: "",
    region: "",
    tencentSecretId: "",
    tencentSecretKey: "",
    tencentSessionToken: "",
    tencentRegion: "ap-guangzhou",
    qiniuAccessKey: "",
    qiniuSecretKey: "",
    host: "",
    port: "22",
    username: "root",
    authMethod: "password",
    password: "",
    key: "",
    keyPassphrase: "",
  };
}

function readAliyunConfig(config: Record<string, unknown> | undefined) {
  const raw = config ?? {};
  return {
    accessKeyId: String(raw.accessKeyId ?? ""),
    resourceGroupId: String(raw.resourceGroupId ?? ""),
    region: String(raw.region ?? ""),
  };
}

function readSSHConfig(config: Record<string, unknown> | undefined) {
  const raw = config ?? {};
  return {
    host: String(raw.host ?? ""),
    port: String(raw.port ?? 22),
    username: String(raw.username ?? "root"),
    authMethod: (String(raw.authMethod ?? "password") === "key" ? "key" : "password") as "password" | "key",
  };
}

function readTencentConfig(config: Record<string, unknown> | undefined) {
  const raw = config ?? {};
  return {
    secretId: String(raw.secretId ?? ""),
    region: String(raw.region ?? "ap-guangzhou"),
  };
}

function readQiniuConfig(config: Record<string, unknown> | undefined) {
  const raw = config ?? {};
  return {
    accessKey: String(raw.accessKey ?? ""),
  };
}

export default function AccessesPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const { data } = useQuery({ queryKey: ["accesses"], queryFn: api.listAccesses });
  const [form, setForm] = useState<AccessFormState>(emptyForm());
  const [notice, setNotice] = useState<{ type: "error" | "success" | "info"; text: string } | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const editing = Boolean(form.id);
  const isAliyun = form.provider === "aliyun";
  const isTencent = form.provider === "tencentcloud";
  const isQiniu = form.provider === "qiniu";
  const isSSH = form.provider === "ssh";

  const save = useMutation({
    mutationFn: async () => {
      setNotice(null);
      const payload: Access = {
        id: form.id,
        name: form.name.trim(),
        provider: form.provider,
        config: isAliyun
          ? {
              accessKeyId: form.accessKeyId.trim(),
              accessKeySecret: form.accessKeySecret.trim(),
              resourceGroupId: form.resourceGroupId.trim(),
              region: form.region.trim(),
            }
          : isTencent
            ? {
                secretId: form.tencentSecretId.trim(),
                secretKey: form.tencentSecretKey.trim(),
                sessionToken: form.tencentSessionToken.trim(),
                region: form.tencentRegion.trim(),
              }
            : isQiniu
              ? {
                  accessKey: form.qiniuAccessKey.trim(),
                  secretKey: form.qiniuSecretKey.trim(),
                }
              : {
                  host: form.host.trim(),
                  port: Number(form.port || 22),
                  username: form.username.trim(),
                  authMethod: form.authMethod,
                  password: form.password,
                  key: form.key,
                  keyPassphrase: form.keyPassphrase,
                },
      };
      return api.saveAccess(payload);
    },
    onSuccess: () => {
      setForm(emptyForm());
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
    setForm(emptyForm());
    setDialogOpen(true);
    setNotice(null);
  };

  const startEdit = (access: Access) => {
    if (access.provider === "aliyun") {
      const aliyun = readAliyunConfig(access.config);
      setForm({
        id: access.id,
        name: access.name,
        provider: access.provider,
        accessKeyId: aliyun.accessKeyId,
        accessKeySecret: "",
        resourceGroupId: aliyun.resourceGroupId,
        region: aliyun.region,
        tencentSecretId: "",
        tencentSecretKey: "",
        tencentSessionToken: "",
        tencentRegion: "ap-guangzhou",
        qiniuAccessKey: "",
        qiniuSecretKey: "",
        host: "",
        port: "22",
        username: "root",
        authMethod: "password",
        password: "",
        key: "",
        keyPassphrase: "",
      });
    } else if (access.provider === "tencentcloud") {
      const tencent = readTencentConfig(access.config);
      setForm({
        id: access.id,
        name: access.name,
        provider: access.provider,
        accessKeyId: "",
        accessKeySecret: "",
        resourceGroupId: "",
        region: "",
        tencentSecretId: tencent.secretId,
        tencentSecretKey: "",
        tencentSessionToken: "",
        tencentRegion: tencent.region,
        qiniuAccessKey: "",
        qiniuSecretKey: "",
        host: "",
        port: "22",
        username: "root",
        authMethod: "password",
        password: "",
        key: "",
        keyPassphrase: "",
      });
    } else if (access.provider === "qiniu") {
      const qiniu = readQiniuConfig(access.config);
      setForm({
        id: access.id,
        name: access.name,
        provider: access.provider,
        accessKeyId: "",
        accessKeySecret: "",
        resourceGroupId: "",
        region: "",
        tencentSecretId: "",
        tencentSecretKey: "",
        tencentSessionToken: "",
        tencentRegion: "ap-guangzhou",
        qiniuAccessKey: qiniu.accessKey,
        qiniuSecretKey: "",
        host: "",
        port: "22",
        username: "root",
        authMethod: "password",
        password: "",
        key: "",
        keyPassphrase: "",
      });
    } else {
      const ssh = readSSHConfig(access.config);
      setForm({
        id: access.id,
        name: access.name,
        provider: access.provider,
        accessKeyId: "",
        accessKeySecret: "",
        resourceGroupId: "",
        region: "",
        tencentSecretId: "",
        tencentSecretKey: "",
        tencentSessionToken: "",
        tencentRegion: "ap-guangzhou",
        qiniuAccessKey: "",
        qiniuSecretKey: "",
        host: ssh.host,
        port: ssh.port,
        username: ssh.username,
        authMethod: ssh.authMethod,
        password: "",
        key: "",
        keyPassphrase: "",
      });
    }
    setDialogOpen(true);
    setNotice(null);
  };

  return (
    <div className="space-y-4">
      <Card>
        <div className="mb-3 flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold tracking-[-0.02em]">授权列表</h2>
            <p className="text-sm text-[#666]">先新增授权，再在工作流节点中引用对应 accessId。</p>
          </div>
          <Button onClick={openCreate}>新增授权</Button>
        </div>

        {notice ? (
          <div
            className={`mb-3 rounded-md px-3 py-2 text-sm ${
              notice.type === "error"
                ? "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]"
                : notice.type === "success"
                  ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]"
                  : "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]"
            }`}
          >
            {notice.text}
          </div>
        ) : null}

        <div className="overflow-x-auto ds-scrollbar">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-xs uppercase tracking-wide text-[#808080]">
                <th className="pb-2">名称</th>
                <th className="pb-2">ID</th>
                <th className="pb-2">Provider</th>
                <th className="pb-2">配置</th>
                <th className="pb-2 text-right">动作</th>
              </tr>
            </thead>
            <tbody>
              {data?.items.map((x) => (
                <tr key={x.id} className="border-t border-[#f1f1f1]">
                  <td className="py-2">{x.name}</td>
                  <td className="font-mono text-xs text-[#666]">{x.id}</td>
                  <td>{x.provider}</td>
                  <td className="text-xs text-[#666]">
                    {x.provider === "aliyun"
                      ? (() => {
                          const cfg = readAliyunConfig(x.config);
                          const rg = cfg.resourceGroupId ? `, rg=${cfg.resourceGroupId}` : "";
                          return `ak=${cfg.accessKeyId}${rg}`;
                        })()
                      : x.provider === "tencentcloud"
                        ? (() => {
                            const cfg = readTencentConfig(x.config);
                            return `sid=${cfg.secretId}, region=${cfg.region}`;
                          })()
                        : x.provider === "qiniu"
                          ? (() => {
                              const cfg = readQiniuConfig(x.config);
                              return `ak=${cfg.accessKey}`;
                            })()
                          : (() => {
                              const cfg = readSSHConfig(x.config);
                              return `${cfg.username}@${cfg.host}:${cfg.port} (${cfg.authMethod})`;
                            })()}
                  </td>
                  <td className="space-x-2 text-right">
                    <Button size="sm" variant="outline" onClick={() => startEdit(x)}>
                      编辑
                    </Button>
                    <Button size="sm" variant="outline" disabled={testAccess.isPending} onClick={() => testAccess.mutate(x.id!)}>
                      测试
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
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      <Dialog.Root open={dialogOpen} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-50 bg-black/35" />
          <Dialog.Content className="fixed left-1/2 top-1/2 z-50 w-[min(860px,92vw)] max-h-[86vh] -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-xl bg-white p-5 shadow-2xl focus:outline-none">
            <div className="mb-4">
              <Dialog.Title className="text-lg font-semibold tracking-[-0.02em]">{editing ? "编辑授权" : "新增授权"}</Dialog.Title>
              <Dialog.Description className="text-sm text-[#666]">配置完成后可用于证书申请或部署节点。</Dialog.Description>
            </div>

            <div className="grid gap-2 md:grid-cols-2">
              <Input placeholder="授权名称" value={form.name} onChange={(e) => setForm((v) => ({ ...v, name: e.target.value }))} />
              <select
                className="ds-ring h-9 rounded-md bg-white px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
                value={form.provider}
                onChange={(e) => setForm((v) => ({ ...v, provider: e.target.value }))}
              >
                {PROVIDERS.map((item) => (
                  <option key={item.value} value={item.value}>
                    {item.label}
                  </option>
                ))}
              </select>

              {isAliyun ? (
                <>
                  <Input placeholder="AccessKey ID" value={form.accessKeyId} onChange={(e) => setForm((v) => ({ ...v, accessKeyId: e.target.value }))} />
                  <Input placeholder={editing ? "AccessKey Secret（留空表示不修改）" : "AccessKey Secret"} type="password" value={form.accessKeySecret} onChange={(e) => setForm((v) => ({ ...v, accessKeySecret: e.target.value }))} />
                  <Input placeholder="Resource Group ID（可选）" value={form.resourceGroupId} onChange={(e) => setForm((v) => ({ ...v, resourceGroupId: e.target.value }))} />
                  <Input placeholder="Region（可选）" value={form.region} onChange={(e) => setForm((v) => ({ ...v, region: e.target.value }))} />
                </>
              ) : null}

              {isTencent ? (
                <>
                  <Input placeholder="SecretId" value={form.tencentSecretId} onChange={(e) => setForm((v) => ({ ...v, tencentSecretId: e.target.value }))} />
                  <Input placeholder={editing ? "SecretKey（留空表示不修改）" : "SecretKey"} type="password" value={form.tencentSecretKey} onChange={(e) => setForm((v) => ({ ...v, tencentSecretKey: e.target.value }))} />
                  <Input placeholder={editing ? "SessionToken（留空表示不修改）" : "SessionToken（可选）"} type="password" value={form.tencentSessionToken} onChange={(e) => setForm((v) => ({ ...v, tencentSessionToken: e.target.value }))} />
                  <Input placeholder="Region（默认 ap-guangzhou）" value={form.tencentRegion} onChange={(e) => setForm((v) => ({ ...v, tencentRegion: e.target.value }))} />
                </>
              ) : null}

              {isQiniu ? (
                <>
                  <Input placeholder="Qiniu AccessKey" value={form.qiniuAccessKey} onChange={(e) => setForm((v) => ({ ...v, qiniuAccessKey: e.target.value }))} />
                  <Input type="password" placeholder={editing ? "Qiniu SecretKey（留空表示不修改）" : "Qiniu SecretKey"} value={form.qiniuSecretKey} onChange={(e) => setForm((v) => ({ ...v, qiniuSecretKey: e.target.value }))} />
                </>
              ) : null}

              {isSSH ? (
                <>
                  <Input placeholder="Host" value={form.host} onChange={(e) => setForm((v) => ({ ...v, host: e.target.value }))} />
                  <Input placeholder="Port" value={form.port} onChange={(e) => setForm((v) => ({ ...v, port: e.target.value }))} />
                  <Input placeholder="Username" value={form.username} onChange={(e) => setForm((v) => ({ ...v, username: e.target.value }))} />
                  <select
                    className="ds-ring h-9 rounded-md bg-white px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]"
                    value={form.authMethod}
                    onChange={(e) => setForm((v) => ({ ...v, authMethod: (e.target.value === "key" ? "key" : "password") as "password" | "key" }))}
                  >
                    <option value="password">Password</option>
                    <option value="key">Private Key</option>
                  </select>
                  {form.authMethod === "password" ? (
                    <Input type="password" placeholder={editing ? "Password（留空表示不修改）" : "Password"} value={form.password} onChange={(e) => setForm((v) => ({ ...v, password: e.target.value }))} />
                  ) : (
                    <>
                      <Input type="password" placeholder={editing ? "Private Key（留空表示不修改）" : "Private Key"} value={form.key} onChange={(e) => setForm((v) => ({ ...v, key: e.target.value }))} />
                      <Input type="password" placeholder={editing ? "Key Passphrase（留空表示不修改）" : "Key Passphrase（可选）"} value={form.keyPassphrase} onChange={(e) => setForm((v) => ({ ...v, keyPassphrase: e.target.value }))} />
                    </>
                  )}
                </>
              ) : null}
            </div>

            <div className="mt-5 flex justify-end gap-2">
              <Button variant="outline" onClick={() => setDialogOpen(false)}>
                取消
              </Button>
              <Button disabled={save.isPending} onClick={() => save.mutate()}>
                {editing ? "保存修改" : "新增授权"}
              </Button>
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
