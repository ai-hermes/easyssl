import { useState } from "react";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/toast";

export default function SettingsPage() {
  const toast = useToast();
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  return (
    <Card className="max-w-lg space-y-3">
      <h2 className="text-lg font-semibold tracking-[-0.02em]">账户设置</h2>
      <p className="text-sm text-[#666]">修改管理员密码后，下次登录即生效。</p>
      <Input type="password" placeholder="新密码" value={password} onChange={(e) => setPassword(e.target.value)} />
      <Button
        onClick={async () => {
          try {
            await api.changePassword(password);
            setMsg({ type: "success", text: "已更新密码" });
            setPassword("");
            toast.success("密码已更新");
          } catch (e) {
            const t = e instanceof Error ? e.message : "更新失败";
            setMsg({ type: "error", text: t });
            toast.error(t);
          }
        }}
      >
        更新密码
      </Button>
      {msg ? (
        <p className={`rounded-md px-3 py-2 text-sm ${msg.type === "success" ? "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]" : "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]"}`}>
          {msg.text}
        </p>
      ) : null}
    </Card>
  );
}
