import { useState } from "react";
import { useAppStore } from "@/store/appStore";
import { Button } from "@/components/ui/button";

export default function TeamSetup() {
  const [value, setValue] = useState("");
  const setTeamId = useAppStore((s) => s.setTeamId);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = value.trim();
    if (trimmed) setTeamId(trimmed);
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-sm space-y-6 p-6">
        <div className="space-y-2 text-center">
          <h1 className="text-2xl font-bold">FPL Live Match</h1>
          <p className="text-muted-foreground text-sm">
            Enter your FPL team ID to get started. You can find it in the URL
            when viewing your team on the FPL website.
          </p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            placeholder="Team ID (e.g. 1234567)"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            className="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-none"
          />
          <Button type="submit" className="w-full" disabled={!value.trim()}>
            Continue
          </Button>
        </form>
      </div>
    </div>
  );
}
