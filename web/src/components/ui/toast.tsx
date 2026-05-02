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
      toast.success(text);
      return;
    }
    if (type === "error") {
      toast.error(text);
      return;
    }
    toast.message(text);
  },
  success: (text) => toast.success(text),
  error: (text) => toast.error(text),
  info: (text) => toast.message(text),
};

const ToastContext = createContext<ToastApi>(api);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  return (
    <ToastContext.Provider value={api}>
      {children}
      <Toaster position="top-right" richColors closeButton duration={3000} />
    </ToastContext.Provider>
  );
}

export function useToast() {
  return useContext(ToastContext);
}
