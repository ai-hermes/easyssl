import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "@/api";
import { setToken } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

export default function LoginPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("admin@easyssl.local");
  const [password, setPassword] = useState("1234567890");
  const [error, setError] = useState("");

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const resp = await api.login(email, password);
      setToken(resp.token);
      navigate("/");
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-white p-4">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(10,114,239,0.07),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(222,29,141,0.06),transparent_30%)]" />
      <Card className="relative z-10 w-full max-w-md">
        <h1 className="mb-1 text-2xl font-semibold tracking-[-0.03em]">EasySSL 登录</h1>
        <p className="mb-4 text-sm text-[#666]">使用管理员账号进入证书自动化控制台。</p>
        <form className="space-y-3" onSubmit={onSubmit}>
          <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="邮箱" />
          <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="密码" />
          {error && <p className="rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">{error}</p>}
          <Button className="w-full" type="submit">登录</Button>
        </form>
      </Card>
    </div>
  );
}
