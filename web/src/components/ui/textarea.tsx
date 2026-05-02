import * as React from "react";

import { cn } from "@/lib/utils";

const Textarea = React.forwardRef<HTMLTextAreaElement, React.ComponentProps<"textarea">>(({ className, ...props }, ref) => {
  return (
    <textarea
      className={cn(
        "ds-scrollbar flex min-h-[96px] w-full rounded-md bg-white px-3 py-2 text-sm text-[#171717] shadow-[rgba(0,0,0,0.08)_0px_0px_0px_1px] placeholder:text-[#808080] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ds-focus)] disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
      ref={ref}
      {...props}
    />
  );
});
Textarea.displayName = "Textarea";

export { Textarea };
