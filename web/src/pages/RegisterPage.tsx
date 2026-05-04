import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/toast";

export default function RegisterPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const toast = useToast();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password !== confirm) {
      setError(t("register.passwordMismatch"));
      return;
    }
    if (password.length < 8) {
      setError(t("register.passwordMin"));
      return;
    }

    setLoading(true);
    try {
      await api.register(email, password);
      toast.success(t("register.success"));
      navigate("/login");
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-white p-4">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(10,114,239,0.07),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(222,29,141,0.06),transparent_30%)]" />
      <Card className="relative z-10 w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl tracking-[-0.03em]">EasySSL {t("register.title")}</CardTitle>
          <CardDescription>{t("register.description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-3" onSubmit={onSubmit}>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder={t("register.email")}
              required
            />
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder={t("register.password")}
              required
            />
            <Input
              type="password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              placeholder={t("register.confirmPassword")}
              required
            />
            {error && (
              <p className="rounded-md bg-[var(--ds-danger-bg)] px-3 py-2 text-sm text-[var(--ds-danger-fg)]">
                {error}
              </p>
            )}
            <Button className="w-full" type="submit" disabled={loading}>
              {loading ? t("common.loading") : t("register.submit")}
            </Button>
          </form>
          <div className="mt-4 text-center text-sm text-[#666]">
            {t("register.hasAccount")}{" "}
            <Link to="/login" className="text-[#171717] underline underline-offset-2 hover:text-[#555]">
              {t("register.login")}
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
