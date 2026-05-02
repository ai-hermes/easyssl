import * as React from "react";
import { cn } from "@/lib/utils";

export function Textarea({ className, ...props }: React.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      className={cn(
        "ds-ring ds-scrollbar min-h-24 w-full rounded-md bg-white px-3 py-2 text-sm text-[#171717] placeholder:text-[#808080] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)]",
        className
      )}
      {...props}
    />
  );
}
