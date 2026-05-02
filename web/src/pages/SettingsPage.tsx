import { useState } from "react";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/toast";

export default function SettingsPage() {
  const toast = useToast();
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  return (
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
      </CardContent>
    </Card>
  );
}
