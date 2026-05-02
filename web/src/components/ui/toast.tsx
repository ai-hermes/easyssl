import { createContext, useCallback, useContext, useMemo, useState } from "react";
import { cn } from "@/lib/utils";

type ToastType = "success" | "error" | "info";

type ToastItem = {
  id: string;
  type: ToastType;
  text: string;
};

type ToastApi = {
  show: (type: ToastType, text: string) => void;
  success: (text: string) => void;
  error: (text: string) => void;
  info: (text: string) => void;
};

const ToastContext = createContext<ToastApi | null>(null);

function clsByType(type: ToastType) {
  if (type === "success") return "bg-[var(--ds-success-bg)] text-[var(--ds-success-fg)]";
  if (type === "error") return "bg-[var(--ds-danger-bg)] text-[var(--ds-danger-fg)]";
  return "bg-[var(--ds-info-bg)] text-[var(--ds-info-fg)]";
}

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [items, setItems] = useState<ToastItem[]>([]);

  const show = useCallback((type: ToastType, text: string) => {
    const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
    setItems((prev) => [...prev, { id, type, text }]);
    window.setTimeout(() => {
      setItems((prev) => prev.filter((x) => x.id !== id));
    }, 2800);
  }, []);

  const api = useMemo<ToastApi>(
    () => ({
      show,
      success: (text) => show("success", text),
      error: (text) => show("error", text),
      info: (text) => show("info", text),
    }),
    [show]
  );

  return (
    <ToastContext.Provider value={api}>
      {children}
      <div className="pointer-events-none fixed right-4 top-4 z-[1000] flex max-w-sm flex-col gap-2">
        {items.map((x) => (
          <div key={x.id} className={cn("ds-card rounded-lg px-3 py-2 text-sm", clsByType(x.type))}>
            {x.text}
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within ToastProvider");
  return ctx;
}
