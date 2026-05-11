import type { MatchPlayer } from "@/types/api";

function ownershipColor(pct: number): string {
  if (pct >= 30) return "bg-red-500 text-white";
  if (pct >= 15) return "bg-orange-400 text-white";
  if (pct >= 5)  return "bg-amber-400 text-gray-900";
  return "bg-white/70 text-gray-600";
}

interface Props {
  player: MatchPlayer;
  totalEntries?: number;
  onClick?: (player: MatchPlayer) => void;
}

export default function PlayerCard({ player, totalEntries, onClick }: Props) {
  const pct = totalEntries && totalEntries > 0
    ? (player.ownedCount / totalEntries) * 100
    : null;

  return (
    <button
      onClick={() => onClick?.(player)}
      className="flex flex-col items-center gap-0.5"
    >
      {/* Photo with ownership badge overlay */}
      <div className="relative">
        <div
          className={`h-10 w-10 overflow-hidden rounded-full bg-white/10 shadow ${
            player.isOwned
              ? "ring-2 ring-yellow-400 ring-offset-1 ring-offset-transparent"
              : "border-2 border-white/30"
          }`}
        >
          <img
            src={player.photoUrl}
            alt={player.webName}
            className="h-full w-full object-cover object-top"
            onError={(e) => { e.currentTarget.style.display = "none"; }}
          />
        </div>

        {pct !== null && (
          <span className={`absolute -top-1 -right-3 rounded px-1 py-0.5 text-[9px] font-bold shadow-sm ${ownershipColor(pct)}`}>
            {Math.round(pct)}%
          </span>
        )}

        {player.subbedOn && (
          <span className="absolute -bottom-1 -left-1 flex h-4 w-4 items-center justify-center rounded-full bg-green-500 shadow-sm">
            <svg viewBox="0 0 10 10" className="h-2.5 w-2.5 fill-white">
              <path d="M5 1 L8.5 6 H1.5 Z" />
            </svg>
          </span>
        )}

        {player.subbedOff && (
          <span className="absolute -bottom-1 -left-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 shadow-sm">
            <svg viewBox="0 0 10 10" className="h-2.5 w-2.5 fill-white">
              <path d="M5 9 L8.5 4 H1.5 Z" />
            </svg>
          </span>
        )}
      </div>

      {/* Name */}
      <span className="max-w-[60px] truncate rounded bg-black/50 px-1 py-0.5 text-[10px] font-medium text-white">
        {player.webName}
      </span>

      {/* Points and minutes */}
      <span className="rounded bg-white/90 px-1 py-0.5 text-[10px] font-bold text-gray-900 shadow-sm whitespace-nowrap">
        {player.livePoints} pts · {player.minutes}&apos;
      </span>
    </button>
  );
}
