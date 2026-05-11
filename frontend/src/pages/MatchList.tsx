import { useSearchParams } from "react-router-dom";
import MatchPanel from "@/components/MatchPanel";

export default function MatchList() {
  const [searchParams, setSearchParams] = useSearchParams();

  const openMatchIds = searchParams.get("matches")?.split(",").map(Number).filter(Boolean) ?? [];

  function closeMatch(id: number) {
    const remaining = openMatchIds.filter((m) => m !== id);
    if (remaining.length === 0) setSearchParams({});
    else setSearchParams({ matches: remaining.join(",") });
  }

  if (openMatchIds.length === 0) {
    return (
      <div className="flex h-[calc(100vh-57px)] items-center justify-center">
        <p className="text-muted-foreground text-sm">Select a match from the header to get started.</p>
      </div>
    );
  }

  return (
    <div className="flex flex-row overflow-x-auto">
      {openMatchIds.map((id) => (
        <div key={id} className="border-border border-r last:border-r-0">
          <MatchPanel fixtureId={id} onClose={() => closeMatch(id)} />
        </div>
      ))}
    </div>
  );
}
