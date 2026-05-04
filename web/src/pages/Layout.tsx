import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { clearToken, getRole } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export default function Layout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t, i18n } = useTranslation();

  const isAdmin = getRole() === "admin";
  const items = [
    ["/", t("layout.nav.dashboard")],
    ["/accesses", t("layout.nav.accesses")],
    ["/workflows", t("layout.nav.workflows")],
    ["/certificates", t("layout.nav.certificates")],
    ...(isAdmin ? [["/users", t("layout.nav.users")] as [string, string]] : []),
    ["/settings", t("layout.nav.settings")],
    ["/docs", t("layout.nav.docs")],
  ] as [string, string][];

  const title = items.find((x) => x[0] === location.pathname)?.[1] || "EasySSL";

  return (
    <div className="min-h-screen bg-[var(--ds-bg)]">
      <header className="sticky top-0 z-40 bg-white/95 backdrop-blur" style={{ boxShadow: "rgba(0,0,0,0.08)_0px_0px_0px_1px" }}>
        <div className="mx-auto flex max-w-[1200px] items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2">
            <img src="/logo-v1.png" alt="EasySSL" className="h-6 w-6 rounded" />
            <div className="text-base font-semibold tracking-[-0.02em] text-[#171717]">{t("layout.brand")}</div>
            <Badge variant="secondary">{t("layout.tagline")}</Badge>
          </div>
          <div className="flex items-center gap-1">
            {items.map(([to, name]) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  `rounded-md px-3 py-1.5 text-sm transition ${
                    isActive ? "bg-white text-[#171717] font-medium shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px]" : "text-[#666] hover:bg-[#fafafa] hover:text-[#171717]"
                  }`
                }
              >
                {name}
              </NavLink>
            ))}
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => i18n.changeLanguage(i18n.language === "zh" ? "en" : "zh")}
            >
              {i18n.language === "zh" ? "EN" : "中"}
            </Button>
            <Button
              variant="outline"
              onClick={() => {
                clearToken();
                navigate("/login");
              }}
            >
              {t("layout.logout")}
            </Button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-[1200px] p-4 md:p-6">
        <div className="mb-6 flex items-end justify-between">
          <div>
            <h1 className="text-[28px] font-semibold tracking-[-0.04em] text-[#171717]">{title}</h1>
          </div>
        </div>
        <Outlet />
      </main>
    </div>
  );
}
