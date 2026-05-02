import * as React from "react";
import { cn } from "@/lib/utils";

export function Input({ className, ...props }: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={cn(
        "ds-ring flex h-9 w-full rounded-md bg-white px-3 py-1 text-sm text-[#171717] placeholder:text-[#808080] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]",
        className
      )}
      {...props}
    />
  );
}
