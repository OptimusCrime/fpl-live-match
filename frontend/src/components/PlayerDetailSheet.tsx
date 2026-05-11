import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/client";
import type { MatchPlayer } from "@/types/api";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

interface Props {
  player: MatchPlayer | null;
  leagueId: number;
  onClose: () => void;
}

export default function PlayerDetailSheet({ player, leagueId, onClose }: Props) {
  const { data, isLoading } = useQuery({
    queryKey: ["playerOwners", player?.id, leagueId],
    queryFn: () => api.getPlayerOwners(player!.id, leagueId),
    enabled: !!player,
  });

  return (
    <Sheet open={!!player} onOpenChange={(open) => { if (!open) onClose(); }}>
      <SheetContent side="bottom" className="max-h-[80vh] overflow-y-auto rounded-t-2xl">
        {player && (
          <>
            <SheetHeader className="flex flex-row items-center gap-4 pb-4">
              <img
                src={player.photoUrl}
                alt={player.webName}
                className="h-16 w-16 rounded-full object-cover object-top border-2 border-border"
                onError={(e) => { e.currentTarget.style.display = "none"; }}
              />
              <div>
                <SheetTitle className="text-left text-lg">{player.webName}</SheetTitle>
                <p className="text-muted-foreground text-sm">
                  {player.livePoints} pts · {player.minutes}&apos; played
                </p>
              </div>
            </SheetHeader>

            {isLoading && (
              <div className="space-y-2 py-4">
                {Array.from({ length: 4 }).map((_, i) => (
                  <div key={i} className="bg-muted h-10 animate-pulse rounded-lg" />
                ))}
              </div>
            )}

            {data && (
              <>
                <p className="text-muted-foreground mb-3 text-sm">
                  Owned by <strong className="text-foreground">{data.owners.length}</strong> of{" "}
                  <strong className="text-foreground">{data.totalEntries}</strong> managers in your league
                </p>

                {data.owners.length === 0 ? (
                  <p className="text-muted-foreground py-4 text-center text-sm">No one in your league owns this player.</p>
                ) : (
                  <div className="divide-border divide-y rounded-lg border">
                    {data.owners.map((owner, i) => (
                      <div key={i} className="flex items-center justify-between px-4 py-3">
                        <div>
                          <p className="text-sm font-medium">{owner.entryName}</p>
                          <p className="text-muted-foreground text-xs">{owner.playerName}</p>
                        </div>
                        <span className="text-muted-foreground text-sm font-semibold">
                          #{owner.rank}
                        </span>
                      </div>
                    ))}
                  </div>
                )}
              </>
            )}
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
