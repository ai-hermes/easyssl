import { createContext, useContext } from "react";
import { Toaster, toast } from "sonner";

type ToastType = "success" | "error" | "info";

type ToastApi = {
  show: (type: ToastType, text: string) => void;
  success: (text: string) => void;
  error: (text: string) => void;
  info: (text: string) => void;
};

const api: ToastApi = {
  show: (type, text) => {
    if (type === "success") {
      toast.success(text, { className: "!bg-[var(--ds-success-bg)] !text-[var(--ds-success-fg)]" });
      return;
    }
    if (type === "error") {
      toast.error(text, { className: "!bg-[var(--ds-danger-bg)] !text-[var(--ds-danger-fg)]" });
      return;
    }
    toast.message(text, { className: "!bg-[var(--ds-info-bg)] !text-[var(--ds-info-fg)]" });
  },
  success: (text) => toast.success(text, { className: "!bg-[var(--ds-success-bg)] !text-[var(--ds-success-fg)]" }),
  error: (text) => toast.error(text, { className: "!bg-[var(--ds-danger-bg)] !text-[var(--ds-danger-fg)]" }),
  info: (text) => toast.message(text, { className: "!bg-[var(--ds-info-bg)] !text-[var(--ds-info-fg)]" }),
};

const ToastContext = createContext<ToastApi>(api);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  return (
    <ToastContext.Provider value={api}>
      {children}
      <Toaster position="top-right" closeButton duration={2800} />
    </ToastContext.Provider>
  );
}

export function useToast() {
  return useContext(ToastContext);
}
