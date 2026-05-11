import { useState, useRef, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { api } from "@/api/client";
import { useAppStore } from "@/store/appStore";
import PitchView from "@/components/PitchView";
import PlayerDetailSheet from "@/components/PlayerDetailSheet";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import type { MatchPlayer, MatchResponse } from "@/types/api";

function MatchDetailSkeleton() {
  return (
    <div className="mx-auto max-w-2xl space-y-4 p-4">
      <Skeleton className="h-4 w-24" />
      <Skeleton className="h-20 w-full rounded-xl" />
      <Skeleton className="h-[700px] w-full rounded-xl" />
    </div>
  );
}

// Track previous player points so we can toast when they change.
function usePointsChangeToasts(data: MatchResponse | undefined) {
  const prevPoints = useRef<Record<number, number>>({});

  useEffect(() => {
    if (!data) return;
    const isLive = data.fixture.status === "live";
    if (!isLive) return;

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

export default function MatchDetail() {
  const { id }    = useParams<{ id: string }>();
  const fixtureId = Number(id);

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

  if (matchQuery.isLoading) return <MatchDetailSkeleton />;

  if (matchQuery.isError || !matchQuery.data) {
    return (
      <div className="mx-auto max-w-2xl space-y-4 p-4">
        <Link to="/" className="text-muted-foreground hover:underline text-sm inline-block">← Back</Link>
        <Alert variant="destructive">
          <AlertDescription>Failed to load match. The FPL API may be unavailable — try again shortly.</AlertDescription>
        </Alert>
      </div>
    );
  }

  const { fixture, home, away, totalEntries } = matchQuery.data;

  return (
    <div className="mx-auto max-w-2xl space-y-4 p-4">
      <Link to="/" className="text-muted-foreground hover:underline text-sm inline-block">
        ← Back to matches
      </Link>

      {/* Match header */}
      <div className="bg-card border-border flex items-center justify-between rounded-xl border p-4">
        <div className="flex min-w-0 flex-1 items-center gap-2">
          <img
            src={`https://fantasy.premierleague.com/dist/img/shirts/standard/shirt_${fixture.teamHCode}-110.png`}
            alt={fixture.teamHShort}
            className="h-10 w-10 shrink-0 object-contain"
            onError={(e) => (e.currentTarget.style.display = "none")}
          />
          <span className="truncate font-semibold">{fixture.teamHName}</span>
        </div>

        <div className="shrink-0 px-2 text-center">
          {fixture.scoreH !== null && (
            <span className="text-2xl font-bold tabular-nums">
              {fixture.scoreH} – {fixture.scoreA}
            </span>
          )}
          {fixture.status === "live" && (
            <p className="mt-1 flex items-center justify-center gap-1.5 text-xs font-semibold text-emerald-500">
              <span className="relative flex h-2 w-2">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
                <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
              </span>
              {fixture.minutes}&apos;
            </p>
          )}
          {fixture.status === "finished" && (
            <p className="text-muted-foreground text-xs font-semibold">Full Time</p>
          )}
          {fixture.status === "upcoming" && (
            <p className="text-muted-foreground text-xs">Kick off soon</p>
          )}
        </div>

        <div className="flex min-w-0 flex-1 items-center justify-end gap-2">
          <span className="truncate font-semibold">{fixture.teamAName}</span>
          <img
            src={`https://fantasy.premierleague.com/dist/img/shirts/standard/shirt_${fixture.teamACode}-110.png`}
            alt={fixture.teamAShort}
            className="h-10 w-10 shrink-0 object-contain"
            onError={(e) => (e.currentTarget.style.display = "none")}
          />
        </div>
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
