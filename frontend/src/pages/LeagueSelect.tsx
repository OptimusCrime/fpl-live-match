import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/client";
import { useAppStore } from "@/store/appStore";
import { Button } from "@/components/ui/button";

export default function LeagueSelect() {
  const teamId = useAppStore((s) => s.teamId)!;
  const setLeague = useAppStore((s) => s.setLeague);
  const clearTeam = useAppStore((s) => s.clearTeam);

  const { data, isLoading, isError } = useQuery({
    queryKey: ["entry", teamId],
    queryFn: () => api.getEntry(teamId),
  });

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-muted-foreground text-sm">Loading leagues…</p>
      </div>
    );
  }

  if (isError || !data) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="space-y-4 text-center">
          <p className="text-destructive text-sm">
            Could not load team. Check your team ID.
          </p>
          <Button variant="outline" onClick={clearTeam}>
            Try again
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-sm space-y-6 p-6">
        <div className="space-y-1 text-center">
          <h1 className="text-2xl font-bold">Select a League</h1>
          <p className="text-muted-foreground text-sm">{data.name}</p>
        </div>
        <div className="space-y-2">
          {data.leagues.map((league) => (
            <button
              key={league.id}
              onClick={() => setLeague(league.id, league.name)}
              className="border-border bg-card hover:bg-accent w-full rounded-lg border p-4 text-left transition-colors"
            >
              <p className="font-medium">{league.name}</p>
              <p className="text-muted-foreground text-xs">
                Rank #{league.entryRank.toLocaleString()}
              </p>
            </button>
          ))}
        </div>
        <Button variant="ghost" size="sm" className="w-full" onClick={clearTeam}>
          Change team ID
        </Button>
      </div>
    </div>
  );
}
