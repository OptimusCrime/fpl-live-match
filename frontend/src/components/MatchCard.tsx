import { useSearchParams } from "react-router-dom";
import { format, parseISO } from "date-fns";
import type { Fixture } from "@/types/api";

function kitUrl(teamCode: number) {
  return `https://fantasy.premierleague.com/dist/img/shirts/standard/shirt_${teamCode}-110.png`;
}

function StatusBadge({ status, minutes }: { status: Fixture["status"]; minutes: number }) {
  if (status === "live") {
    return (
      <span className="flex items-center gap-1.5 rounded-full bg-emerald-500/15 px-2.5 py-0.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400">
        <span className="relative flex h-2 w-2">
          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
          <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
        </span>
        {minutes}&apos;
      </span>
    );
  }
  if (status === "finished") {
    return (
      <span className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold text-muted-foreground">
        FT
      </span>
    );
  }
  return null;
}

function KickoffTime({ kickoffTime }: { kickoffTime: string }) {
  const date = parseISO(kickoffTime);
  return (
    <div className="text-center">
      <p className="text-sm font-semibold tabular-nums">{format(date, "HH:mm")}</p>
      <p className="text-muted-foreground text-xs">{format(date, "EEE d MMM")}</p>
    </div>
  );
}

interface Props {
  fixture: Fixture;
}

export default function MatchCard({ fixture }: Props) {
  const [searchParams, setSearchParams] = useSearchParams();

  const openIds = searchParams.get("matches")?.split(",").map(Number).filter(Boolean) ?? [];
  const isOpen = openIds.includes(fixture.id);

  function handleClick() {
    if (isOpen) return;
    setSearchParams({ matches: [...openIds, fixture.id].join(",") });
  }

  return (
    <button
      onClick={handleClick}
      className={`border-border bg-card hover:bg-accent/50 flex w-full items-center gap-4 rounded-xl border p-4 text-left transition-colors ${isOpen ? "ring-primary ring-2" : ""}`}
    >
      {/* Home team */}
      <div className="flex flex-1 items-center justify-end gap-3">
        <span className="text-sm font-semibold">{fixture.teamHName}</span>
        <img
          src={kitUrl(fixture.teamHCode)}
          alt={fixture.teamHShort}
          className="h-9 w-9 object-contain"
          onError={(e) => (e.currentTarget.style.display = "none")}
        />
      </div>

      {/* Score / time */}
      <div className="flex w-28 shrink-0 flex-row items-center justify-center gap-2">
        {fixture.status === "upcoming" ? (
          <KickoffTime kickoffTime={fixture.kickoffTime} />
        ) : (
          <p className="text-sm font-bold tabular-nums">
            {fixture.scoreH ?? "–"} – {fixture.scoreA ?? "–"}
          </p>
        )}
        <StatusBadge status={fixture.status} minutes={fixture.minutes} />
      </div>

      {/* Away team */}
      <div className="flex flex-1 items-center gap-3">
        <img
          src={kitUrl(fixture.teamACode)}
          alt={fixture.teamAShort}
          className="h-9 w-9 object-contain"
          onError={(e) => (e.currentTarget.style.display = "none")}
        />
        <span className="text-sm font-semibold">{fixture.teamAName}</span>
      </div>
    </button>
  );
}
