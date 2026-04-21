import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { clearToken } from "@/api/client";
import { Button } from "@/components/ui/button";

const items = [
  ["/", "仪表盘"],
  ["/accesses", "授权"],
  ["/workflows", "工作流"],
  ["/certificates", "证书"],
  ["/settings", "设置"],
] as const;

export default function Layout() {
  const navigate = useNavigate();
  return (
    <div className="min-h-screen bg-slate-100">
      <header className="border-b bg-white">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
          <div className="font-semibold">EasySSL</div>
          <div className="flex items-center gap-2">
            {items.map(([to, name]) => (
              <NavLink key={to} to={to} className={({ isActive }) => `rounded px-3 py-1 text-sm ${isActive ? "bg-slate-200" : ""}`}>
                {name}
              </NavLink>
            ))}
          </div>
          <Button variant="outline" onClick={() => { clearToken(); navigate("/login"); }}>退出</Button>
        </div>
      </header>
      <main className="mx-auto max-w-7xl p-4"><Outlet /></main>
    </div>
  );
}
