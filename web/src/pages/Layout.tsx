import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { clearToken } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

const items = [
  ["/", "仪表盘"],
  ["/accesses", "授权"],
  ["/workflows", "工作流"],
  ["/certificates", "证书"],
  ["/settings", "设置"],
  ["/docs", "文档"],
] as const;

export default function Layout() {
  const navigate = useNavigate();
  const location = useLocation();
  const title = items.find((x) => x[0] === location.pathname)?.[1] || "EasySSL";
  return (
    <div className="min-h-screen bg-[var(--ds-bg)]">
      <header className="sticky top-0 z-40 bg-white/95 backdrop-blur" style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}>
        <div className="mx-auto flex max-w-[1200px] items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2">
            <div className="text-base font-semibold tracking-[-0.02em] text-[#171717]">EasySSL</div>
            <Badge variant="secondary">Workflow SSL</Badge>
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
      <main className="mx-auto max-w-[1200px] p-4 md:p-6">
        <div className="mb-6 flex items-end justify-between">
          <div>
            <h1 className="text-[28px] font-semibold tracking-[-0.04em] text-[#171717]">{title}</h1>
            <p className="text-sm text-[#666]">清晰配置，实时执行反馈，节点状态可观测。</p>
          </div>
        </div>
        <Outlet />
      </main>
    </div>
  );
}
