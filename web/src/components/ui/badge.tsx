import { cn } from "@/lib/utils";

export function Badge({ className, ...props }: React.HTMLAttributes<HTMLSpanElement>) {
  return <span className={cn("inline-flex items-center rounded-full bg-[#f5f5f5] px-2.5 py-1 text-xs text-[#4d4d4d]", className)} {...props} />;
}
