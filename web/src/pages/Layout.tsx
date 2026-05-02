import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { clearToken } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

const items = [
  ["/", "仪表盘"],
  ["/accesses", "授权"],
  ["/workflows", "工作流"],
  ["/certificates", "证书"],
  ["/settings", "设置"],
] as const;

export default function Layout() {
  const navigate = useNavigate();
  const location = useLocation();
  const title = items.find((x) => x[0] === location.pathname)?.[1] || "EasySSL";
  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-40 border-b bg-background/90 backdrop-blur">
        <div className="container flex items-center justify-between py-3">
          <div className="flex items-center gap-2">
            <div className="text-base font-semibold tracking-tight">EasySSL</div>
            <Badge variant="secondary">Workflow SSL</Badge>
          </div>
          <div className="flex items-center gap-1">
            {items.map(([to, name]) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  `rounded-md px-3 py-1.5 text-sm transition ${
                    isActive ? "border bg-accent text-accent-foreground font-medium" : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  }`
                }
              >
                {name}
              </NavLink>
            ))}
          </div>
          <Button
            variant="outline"
            onClick={() => {
              clearToken();
              navigate("/login");
            }}
          >
            退出
          </Button>
        </div>
      </header>
      <main className="container py-4 md:py-6">
        <div className="mb-6 flex items-end justify-between">
          <div>
            <h1 className="text-3xl font-semibold tracking-tight">{title}</h1>
            <p className="text-sm text-muted-foreground">清晰配置，实时执行反馈，节点状态可观测。</p>
          </div>
        </div>
        <Separator className="mb-6" />
        <Outlet />
      </main>
    </div>
  );
}
