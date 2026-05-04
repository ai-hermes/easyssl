import { lazy, Suspense } from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import { getToken } from "@/api/client";
import Layout from "@/pages/Layout";
import LoginPage from "@/pages/LoginPage";
import DashboardPage from "@/pages/DashboardPage";
import AccessesPage from "@/pages/AccessesPage";
import WorkflowsPage from "@/pages/WorkflowsPage";
import CertificatesPage from "@/pages/CertificatesPage";
import SettingsPage from "@/pages/SettingsPage";
const DocsPage = lazy(() => import("@/pages/DocsPage"));

function Guard({ children }: { children: JSX.Element }) {
  return getToken() ? children : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/docs" element={<Suspense fallback={<div className="flex h-screen items-center justify-center text-sm text-[#666]">加载文档…</div>}><DocsPage /></Suspense>} />
      <Route path="/" element={<Guard><Layout /></Guard>}>
        <Route index element={<DashboardPage />} />
        <Route path="accesses" element={<AccessesPage />} />
        <Route path="workflows" element={<WorkflowsPage />} />
        <Route path="certificates" element={<CertificatesPage />} />
        <Route path="settings" element={<SettingsPage />} />
      </Route>
    </Routes>
  );
}
