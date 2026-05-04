import { useTranslation } from "react-i18next";
import OpenAPIDocs from "@/components/OpenAPIDocs";

export default function DocsPage() {
  const { t } = useTranslation();
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-[24px] font-semibold tracking-[-0.04em] text-[#171717]">
          {t("docs.title")}
        </h2>
      </div>
      <OpenAPIDocs />
    </div>
  );
}
