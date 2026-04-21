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
    <div className="flex min-h-screen items-center justify-center bg-slate-100 p-4">
      <Card className="w-full max-w-md">
        <h1 className="mb-4 text-xl font-semibold">EasySSL 登录</h1>
        <form className="space-y-3" onSubmit={onSubmit}>
          <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="邮箱" />
          <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="密码" />
          {error && <p className="text-sm text-red-500">{error}</p>}
          <Button className="w-full" type="submit">登录</Button>
        </form>
      </Card>
    </div>
  );
}
