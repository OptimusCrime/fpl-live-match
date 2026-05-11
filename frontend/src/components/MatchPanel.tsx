import { useState, useRef, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { X } from "lucide-react";
import { api } from "@/api/client";
import { useAppStore } from "@/store/appStore";
import PitchView from "@/components/PitchView";
import PlayerDetailSheet from "@/components/PlayerDetailSheet";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import type { MatchPlayer, MatchResponse } from "@/types/api";

function MatchPanelSkeleton({ onClose }: { onClose: () => void }) {
  return (
    <div className="w-auto shrink-0 space-y-4 p-4">
      <div className="flex items-center gap-2">
        <Skeleton className="h-20 flex-1 rounded-xl" />
        <Button variant="ghost" size="icon" onClick={onClose} className="shrink-0 self-start">
          <X className="h-4 w-4" />
        </Button>
      </div>
      <Skeleton className="h-[700px] w-full rounded-xl" />
    </div>
  );
}

function usePointsChangeToasts(data: MatchResponse | undefined) {
  const prevPoints = useRef<Record<number, number>>({});

  useEffect(() => {
    if (!data) return;
    if (data.fixture.status !== "live") return;

    const allPlayers = [...data.home.players, ...data.away.players];
    allPlayers.forEach((p) => {
      const prev = prevPoints.current[p.id];
      if (prev !== undefined && p.livePoints !== prev) {
        const delta = p.livePoints - prev;
        toast(`${p.webName} ${delta > 0 ? `+${delta}` : delta} pts`, {
          description: `Now on ${p.livePoints} pts`,
        });
      }
      prevPoints.current[p.id] = p.livePoints;
    });
  }, [data]);
}

interface Props {
  fixtureId: number;
  onClose: () => void;
}

export default function MatchPanel({ fixtureId, onClose }: Props) {
  const teamId   = useAppStore((s) => s.teamId)!;
  const leagueId = useAppStore((s) => s.leagueId)!;

  const [selectedPlayer, setSelectedPlayer] = useState<MatchPlayer | null>(null);

  const matchQuery = useQuery({
    queryKey: ["leagueMatch", leagueId, fixtureId],
    queryFn:  () => api.getLeagueMatch(leagueId, fixtureId, teamId),
    refetchInterval: (query) => {
      const status = query.state.data?.fixture.status;
      return status === "live" || status === "upcoming" ? 30_000 : false;
    },
  });

  usePointsChangeToasts(matchQuery.data);

  if (matchQuery.isLoading) return <MatchPanelSkeleton onClose={onClose} />;

  if (matchQuery.isError || !matchQuery.data) {
    return (
      <div className="w-auto shrink-0 space-y-4 p-4">
        <div className="flex justify-end">
          <Button variant="ghost" size="icon" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>
        <Alert variant="destructive">
          <AlertDescription>
            Failed to load match. The FPL API may be unavailable — try again shortly.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  const { fixture, home, away, totalEntries } = matchQuery.data;

  return (
    <div className="w-auto shrink-0 space-y-1 p-1">
      {/* Match header */}
      <div className="bg-card border-border relative flex items-center justify-between rounded-xl border p-1">
        <div className="flex min-w-0 flex-1 items-center gap-2">
          <img
            src={`https://fantasy.premierleague.com/dist/img/shirts/standard/shirt_${fixture.teamHCode}-110.png`}
            alt={fixture.teamHShort}
            className="h-5 w-5 shrink-0 object-contain"
            onError={(e) => (e.currentTarget.style.display = "none")}
          />
          <span className="truncate text-sx font-semibold">{fixture.teamHName}</span>
        </div>

        <div className="shrink-0 px-2 flex items-center justify-center gap-1.5">
          {fixture.scoreH !== null && (
            <span className="text-sm font-bold tabular-nums">
              {fixture.scoreH} – {fixture.scoreA}
            </span>
          )}
          {fixture.status === "live" && (
            <span className="flex items-center gap-1 text-xs font-semibold text-emerald-500">
              <span className="relative flex h-2 w-2">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
                <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
              </span>
              {fixture.minutes}&apos;
            </span>
          )}
          {fixture.status === "finished" && (
            <span className="text-muted-foreground text-xs font-semibold">FT</span>
          )}
          {fixture.status === "upcoming" && (
            <span className="text-muted-foreground text-xs">Kick off soon</span>
          )}
        </div>

        <div className="flex min-w-0 flex-1 items-center justify-end gap-2">
          <span className="truncate text-sm font-semibold">{fixture.teamAName}</span>
          <img
            src={`https://fantasy.premierleague.com/dist/img/shirts/standard/shirt_${fixture.teamACode}-110.png`}
            alt={fixture.teamAShort}
            className="h-7 w-7 shrink-0 object-contain"
            onError={(e) => (e.currentTarget.style.display = "none")}
          />
        </div>

        <Button
          variant="ghost"
          size="icon"
          onClick={onClose}
          className="border-border bg-background absolute -right-2 -top-2 h-6 w-6 rounded-full border shadow-sm"
        >
          <X className="h-3 w-3" />
        </Button>
      </div>

      {/* Pitch */}
      <PitchView
        home={home}
        away={away}
        totalEntries={totalEntries}
        onPlayerClick={setSelectedPlayer}
      />

      <PlayerDetailSheet
        player={selectedPlayer}
        leagueId={leagueId}
        onClose={() => setSelectedPlayer(null)}
      />
    </div>
  );
}
