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
] as const;

export default function Layout() {
  const navigate = useNavigate();
  const location = useLocation();
  const title = items.find((x) => x[0] === location.pathname)?.[1] || "EasySSL";
  return (
    <div className="min-h-screen bg-white">
      <header className="sticky top-0 z-40 border-b border-[#ebebeb] bg-white/90 backdrop-blur">
        <div className="mx-auto flex max-w-[1200px] items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2">
            <div className="text-base font-semibold tracking-[-0.02em]">EasySSL</div>
            <Badge className="bg-[#ebf5ff] text-[#0068d6]">Workflow SSL</Badge>
          </div>
          <div className="flex items-center gap-1">
            {items.map(([to, name]) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  `rounded-md px-3 py-1.5 text-sm transition ${
                    isActive ? "ds-ring bg-white text-[#171717] font-medium" : "text-[#666] hover:bg-[#fafafa] hover:text-[#171717]"
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
