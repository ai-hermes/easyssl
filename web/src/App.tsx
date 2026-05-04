import { Navigate, Route, Routes } from "react-router-dom";
import { getToken } from "@/api/client";
import Layout from "@/pages/Layout";
import LoginPage from "@/pages/LoginPage";
import RegisterPage from "@/pages/RegisterPage";
import DashboardPage from "@/pages/DashboardPage";
import AccessesPage from "@/pages/AccessesPage";
import WorkflowsPage from "@/pages/WorkflowsPage";
import CertificatesPage from "@/pages/CertificatesPage";
import SettingsPage from "@/pages/SettingsPage";
import UsersPage from "@/pages/UsersPage";
import DocsPage from "@/pages/DocsPage";

function Guard({ children }: { children: JSX.Element }) {
  return getToken() ? children : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/" element={<Guard><Layout /></Guard>}>
        <Route index element={<DashboardPage />} />
        <Route path="accesses" element={<AccessesPage />} />
        <Route path="workflows" element={<WorkflowsPage />} />
        <Route path="certificates" element={<CertificatesPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="users" element={<UsersPage />} />
        <Route path="docs" element={<DocsPage />} />
      </Route>
    </Routes>
  );
}
