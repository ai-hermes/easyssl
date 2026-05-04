import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { api } from "@/api";
import { setToken, setRole } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

export default function LoginPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const resp = await api.login(email, password);
      setToken(resp.token);
      setRole(resp.user.role);
      navigate("/");
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-white p-4">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(10,114,239,0.07),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(222,29,141,0.06),transparent_30%)]" />
      <Card className="relative z-10 w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl tracking-[-0.03em]">EasySSL {t("login.title")}</CardTitle>
          <CardDescription>{t("login.description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-3" onSubmit={onSubmit}>
            <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder={t("login.email")} />
            <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder={t("login.password")} />
            {error && <p className="rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">{error}</p>}
            <Button className="w-full" type="submit">{t("login.submit")}</Button>
          </form>
          <div className="mt-4 text-center text-sm text-[#666]">
            {t("login.noAccount")}{" "}
            <Link to="/register" className="text-[#171717] underline underline-offset-2 hover:text-[#555]">
              {t("login.register")}
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
