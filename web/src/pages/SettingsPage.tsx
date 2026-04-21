import { useState } from "react";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

export default function SettingsPage() {
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState("");

  return (
    <Card className="max-w-lg space-y-3">
      <h2 className="text-lg font-semibold">账户设置</h2>
      <Input type="password" placeholder="新密码" value={password} onChange={(e) => setPassword(e.target.value)} />
      <Button onClick={async () => { await api.changePassword(password); setMsg("已更新密码"); }}>更新密码</Button>
      {msg && <p className="text-sm text-emerald-600">{msg}</p>}
    </Card>
  );
}
