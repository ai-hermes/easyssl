import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Globe, Code } from "lucide-react";
import { getRole } from "@/api/client";
import OpenAPIDocs from "@/components/OpenAPIDocs";

export default function DocsPage() {
  const { t } = useTranslation();
  const isAdmin = getRole() === "admin";
  const [activeTab, setActiveTab] = useState<"openapi" | "api">("openapi");

  const navItems: { key: "openapi" | "api"; label: string; icon: React.ReactNode }[] = [
    { key: "openapi", label: t("layout.nav.openapi"), icon: <Globe size={16} /> },
    ...(isAdmin ? [{ key: "api" as const, label: t("layout.nav.apiDocs"), icon: <Code size={16} /> }] : []),
  ];

  return (
    <div className="grid gap-8 md:grid-cols-[240px_1fr]">
      {/* Left Sidebar */}
      <div
        className="sticky top-[72px] h-fit self-start rounded-lg bg-white p-2"
        style={{ boxShadow: "rgba(0,0,0,0.08) 0px_0px_0px_1px" }}
      >
        <nav className="flex flex-col gap-1">
          {navItems.map((item) => {
            const isActive = activeTab === item.key;
            return (
              <button
                key={item.key}
                onClick={() => setActiveTab(item.key)}
                className={`flex items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition ${
                  isActive
                    ? "bg-white font-medium text-[#171717] shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px]"
                    : "text-[#666] hover:bg-[#fafafa] hover:text-[#171717]"
                }`}
              >
                {item.icon}
                {item.label}
              </button>
            );
          })}
        </nav>
      </div>

      {/* Right Content */}
      <div className="min-w-0 space-y-6">
        {activeTab === "openapi" && (
          <div className="space-y-6">
            <div>
              <h2 className="text-[24px] font-semibold tracking-[-0.04em] text-[#171717]">
                {t("openapiDocs.title")}
              </h2>
            </div>
            <OpenAPIDocs includeTags={["OpenAPI"]} description={t("openapiDocs.description")} />
          </div>
        )}
        {activeTab === "api" && (
          <div className="space-y-6">
            <div>
              <h2 className="text-[24px] font-semibold tracking-[-0.04em] text-[#171717]">
                {t("apiDocs.title")}
              </h2>
            </div>
            <OpenAPIDocs
              includeTags={[
                "Auth", "APIKey", "Access", "Workflow", "Certificate",
                "Providers", "Statistics", "Notification", "User"
              ]}
              description={t("apiDocs.description")}
            />
          </div>
        )}
      </div>
    </div>
  );
}
