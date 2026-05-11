import { BrowserRouter, Routes, Route, Link, useSearchParams } from "react-router-dom";
import { QueryClient, QueryClientProvider, useQuery } from "@tanstack/react-query";
import { format, parseISO, isPast } from "date-fns";
import { ChevronDown, Check } from "lucide-react";
import { useAppStore } from "@/store/appStore";
import { api } from "@/api/client";
import TeamSetup from "@/pages/TeamSetup";
import LeagueSelect from "@/pages/LeagueSelect";
import MatchList from "@/pages/MatchList";
import { Toaster } from "@/components/ui/sonner";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import type { Fixture } from "@/types/api";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: 1, staleTime: 30_000 },
  },
});

function FixtureRow({ fixture, isOpen, onToggle }: { fixture: Fixture; isOpen: boolean; onToggle: () => void }) {
  return (
    <button
      onClick={onToggle}
      className="hover:bg-accent flex w-full items-center gap-3 rounded-md px-2 py-1.5 text-sm transition-colors"
    >
      <span className="flex w-4 shrink-0 items-center justify-center">
        {isOpen && <Check className="h-3.5 w-3.5 text-primary" />}
      </span>
      <div className="flex flex-1 items-center justify-between gap-2 min-w-0">
        <span className="truncate font-medium">{fixture.teamHShort}</span>
        <span className="text-muted-foreground shrink-0 tabular-nums text-xs">
          {fixture.status === "upcoming"
            ? format(parseISO(fixture.kickoffTime), "HH:mm")
            : `${fixture.scoreH ?? "–"} – ${fixture.scoreA ?? "–"}`}
        </span>
        <span className="truncate text-right font-medium">{fixture.teamAShort}</span>
      </div>
      {fixture.status === "live" && (
        <span className="shrink-0 text-xs font-semibold text-emerald-500">{fixture.minutes}&apos;</span>
      )}
      {fixture.status === "finished" && (
        <span className="text-muted-foreground shrink-0 text-xs">FT</span>
      )}
    </button>
  );
}

function MatchSelectorPopover() {
  const [searchParams, setSearchParams] = useSearchParams();
  const openIds = searchParams.get("matches")?.split(",").map(Number).filter(Boolean) ?? [];

  const { data } = useQuery({
    queryKey: ["fixtures"],
    queryFn: () => api.getFixtures(),
    refetchInterval: (query) =>
      query.state.data?.fixtures.some((f) => f.status === "live") ? 30_000 : 60_000,
  });

  function toggleMatch(id: number) {
    const next = openIds.includes(id)
      ? openIds.filter((m) => m !== id)
      : [...openIds, id];
    if (next.length === 0) setSearchParams({});
    else setSearchParams({ matches: next.join(",") });
  }

  const live     = data?.fixtures.filter((f) => f.status === "live")     ?? [];
  const upcoming = data?.fixtures.filter((f) => f.status === "upcoming") ?? [];
  const finished = data?.fixtures.filter((f) => f.status === "finished") ?? [];

  const label = data ? `GW ${data.gw}` : "Matches";

  return (
    <Popover>
      <PopoverTrigger render={<Button variant="ghost" size="sm" className="gap-1" />}>
        {label}
        <ChevronDown className="h-3.5 w-3.5 opacity-60" />
      </PopoverTrigger>
      <PopoverContent align="start" className="w-64 p-2">
        {data && (
          <p className="text-muted-foreground mb-2 px-2 text-xs">
            {(() => {
              const date = parseISO(data.deadlineTime);
              const passed = isPast(date);
              return `${passed ? "Deadline passed" : "Deadline"} ${format(date, "EEE d MMM, HH:mm")}`;
            })()}
          </p>
        )}

        {live.length > 0 && (
          <div className="mb-1">
            <p className="text-muted-foreground mb-0.5 px-2 text-[10px] font-semibold uppercase tracking-wider">Live</p>
            {live.map((f) => (
              <FixtureRow key={f.id} fixture={f} isOpen={openIds.includes(f.id)} onToggle={() => toggleMatch(f.id)} />
            ))}
          </div>
        )}

        {upcoming.length > 0 && (
          <div className="mb-1">
            <p className="text-muted-foreground mb-0.5 px-2 text-[10px] font-semibold uppercase tracking-wider">Upcoming</p>
            {upcoming.map((f) => (
              <FixtureRow key={f.id} fixture={f} isOpen={openIds.includes(f.id)} onToggle={() => toggleMatch(f.id)} />
            ))}
          </div>
        )}

        {finished.length > 0 && (
          <div>
            <p className="text-muted-foreground mb-0.5 px-2 text-[10px] font-semibold uppercase tracking-wider">Finished</p>
            {finished.map((f) => (
              <FixtureRow key={f.id} fixture={f} isOpen={openIds.includes(f.id)} onToggle={() => toggleMatch(f.id)} />
            ))}
          </div>
        )}

        {!data && (
          <p className="text-muted-foreground px-2 py-4 text-center text-sm">Loading fixtures…</p>
        )}
      </PopoverContent>
    </Popover>
  );
}

function NavBar() {
  const leagueName = useAppStore((s) => s.leagueName);
  const clearLeague = useAppStore((s) => s.clearLeague);
  const clearTeam   = useAppStore((s) => s.clearTeam);

  return (
    <header className="border-border bg-background flex items-center justify-between border-b px-4 py-2">
      <div className="flex items-center gap-2">
        <Link to="/" className="font-semibold tracking-tight hover:opacity-80 transition-opacity">
          FPL Live
        </Link>
        <MatchSelectorPopover />
      </div>
      <div className="flex items-center gap-1">
        {leagueName && (
          <span className="text-muted-foreground mr-1 hidden max-w-[160px] truncate text-xs sm:block">
            {leagueName}
          </span>
        )}
        <Button variant="ghost" size="sm" onClick={clearLeague}>
          Switch league
        </Button>
        <Button variant="ghost" size="sm" onClick={clearTeam}>
          Switch team
        </Button>
      </div>
    </header>
  );
}

function AppRoutes() {
  const teamId   = useAppStore((s) => s.teamId);
  const leagueId = useAppStore((s) => s.leagueId);

  if (!teamId)   return <TeamSetup />;
  if (!leagueId) return <LeagueSelect />;

  return (
    <div className="flex min-h-screen flex-col">
      <NavBar />
      <main className="flex-1">
        <Routes>
          <Route path="/" element={<MatchList />} />
        </Routes>
      </main>
      <footer className="border-border text-muted-foreground border-t px-4 py-3 text-center text-xs">
        <span>&copy; {new Date().getFullYear()} </span>
        <a
          href="https://github.com/OptimusCrime/fpl-live-match"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-foreground underline underline-offset-2 transition-colors"
        >
          OptimusCrime/fpl-live-match
        </a>
        <span className="mx-2 opacity-40">&middot;</span>
        <span>Not affiliated with Fantasy Premier League</span>
      </footer>
    </div>
  );
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AppRoutes />
        <Toaster position="bottom-center" richColors />
      </BrowserRouter>
    </QueryClientProvider>
  );
}
