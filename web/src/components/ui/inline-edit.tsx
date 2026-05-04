import * as React from "react";
import { useEffect, useRef, useState } from "react";
import { Pencil, Check, X } from "lucide-react";

import { cn } from "@/lib/utils";
import { Input } from "./input";
import { Button } from "./button";

interface InlineEditProps {
  value: string;
  onSave: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function InlineEdit({ value, onSave, placeholder, className }: InlineEditProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(value);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing) {
      setEditValue(value);
      requestAnimationFrame(() => {
        inputRef.current?.focus();
      });
    }
  }, [isEditing, value]);

  const handleSave = () => {
    onSave(editValue);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditValue(value);
    setIsEditing(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSave();
    }
    if (e.key === "Escape") {
      handleCancel();
    }
  };

  if (isEditing) {
    return (
      <div className={cn("inline-flex items-center gap-2", className)}>
        <Input
          ref={inputRef}
          value={editValue}
          onChange={(e) => setEditValue(e.target.value)}
          onKeyDown={handleKeyDown}
          className="h-8 w-auto min-w-[140px] px-2.5 py-0 text-sm"
        />
        <button
          type="button"
          onClick={handleSave}
          className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-[#171717] text-white transition-colors hover:bg-black"
        >
          <Check className="h-4 w-4" />
        </button>
        <button
          type="button"
          onClick={handleCancel}
          className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-white text-[#171717] transition-colors hover:bg-[#f5f5f5]"
          style={{ boxShadow: "rgba(0,0,0,0.08) 0px 0px 0px 1px" }}
        >
          <X className="h-4 w-4" />
        </button>
      </div>
    );
  }

  return (
    <div
      className={cn(
        "group inline-flex cursor-pointer items-center gap-1.5 rounded-md px-1.5 py-0.5 transition-colors hover:bg-[#f5f5f5]",
        className
      )}
      onClick={() => setIsEditing(true)}
      title="点击编辑"
    >
      <span className={cn("text-sm", !value && "text-[#808080]")}>
        {value || placeholder}
      </span>
      <Pencil className="h-3 w-3 text-[#808080] opacity-0 transition-opacity group-hover:opacity-100" />
    </div>
  );
}
