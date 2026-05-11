import type { MatchPlayer, MatchTeam } from "@/types/api";
import PlayerCard from "./PlayerCard";

function median(values: number[]): number {
  if (values.length === 0) return 0;
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  return sorted.length % 2 !== 0
    ? sorted[mid]
    : Math.round((sorted[mid - 1] + sorted[mid]) / 2);
}

// FPL position order from goalkeeper to forward
const POSITION_ORDER = ["GKP", "DEF", "MID", "FWD"];

interface PositionRow {
  position: string;
  starters: MatchPlayer[];  // began the match in the XI
  subs: MatchPlayer[];      // came on as substitutes
  reserves: MatchPlayer[];  // did not play (league-owned, shown on side)
}

function groupByPosition(players: MatchPlayer[]): PositionRow[] {
  const rows: PositionRow[] = POSITION_ORDER.map((pos) => ({
    position: pos,
    starters:  players.filter((p) => p.position === pos && p.started &&  p.isStarter),
    subs:      players.filter((p) => p.position === pos && p.started && !p.isStarter),
    reserves:  players.filter((p) => p.position === pos && !p.started),
  })).filter((r) => r.starters.length > 0 || r.subs.length > 0 || r.reserves.length > 0);

  // Pre-match: nobody has started yet — treat everyone as a starter
  if (rows.every((r) => r.starters.length === 0 && r.subs.length === 0)) {
    return POSITION_ORDER.map((pos) => ({
      position: pos,
      starters: players.filter((p) => p.position === pos),
      subs: [],
      reserves: [],
    })).filter((r) => r.starters.length > 0);
  }

  return rows;
}

interface RowProps {
  row: PositionRow;
  totalEntries: number;
  avgPoints: number;
  onPlayerClick: (player: MatchPlayer) => void;
}

function PositionRowView({ row, totalEntries, avgPoints, onPlayerClick }: RowProps) {
  const card = (player: MatchPlayer) => (
    <PlayerCard
      key={player.id}
      player={player}
      totalEntries={totalEntries}
      avgPoints={avgPoints}
      onClick={onPlayerClick}
    />
  );

  return (
    <div className="flex items-center justify-center gap-2 px-2">
      {row.starters.map((p) => card(p))}
      {row.subs.length > 0 && (
        <div className="flex gap-2 pl-4 opacity-75">
          {row.subs.map((p) => card(p))}
        </div>
      )}
      {row.reserves.length > 0 && (
        <div className="flex gap-2 pl-4 opacity-55">
          {row.reserves.map((p) => card(p))}
        </div>
      )}
    </div>
  );
}

interface HalfProps {
  team: MatchTeam;
  reversed?: boolean;
  totalEntries: number;
  avgPoints: number;
  onPlayerClick: (player: MatchPlayer) => void;
}

function TeamHalf({ team, reversed = false, totalEntries, avgPoints, onPlayerClick }: HalfProps) {
  const rows = groupByPosition(team.players);
  const orderedRows = reversed ? [...rows].reverse() : rows;

  return (
    <div className="flex flex-1 flex-col justify-evenly py-1 px-4 gap-3">
      {orderedRows.map((row) => (
        <PositionRowView
          key={row.position}
          row={row}
          totalEntries={totalEntries}
          avgPoints={avgPoints}
          onPlayerClick={onPlayerClick}
        />
      ))}
    </div>
  );
}

interface Props {
  home: MatchTeam;
  away: MatchTeam;
  totalEntries: number;
  onPlayerClick: (player: MatchPlayer) => void;
}

export default function PitchView({ home, away, totalEntries, onPlayerClick }: Props) {
  const allPoints = [...home.players, ...away.players].map((p) => p.livePoints);
  const avgPoints = median(allPoints);

  return (
    <div className="w-full overflow-x-auto rounded-xl">
      <div
        className="relative flex flex-col"
        style={{
          background: "linear-gradient(180deg, #2d7a3a 0%, #3a9e4a 25%, #2d7a3a 50%, #3a9e4a 75%, #2d7a3a 100%)",
          height: 831,
          minWidth: 320,
        }}
      >
        <PitchMarkings />

        {/* Home team — top half: GKP at top, FWD toward centre */}
        <div className="relative z-10 flex flex-1 flex-col border-b border-white/20">
          <TeamHalf
            team={home}
            totalEntries={totalEntries}
            avgPoints={avgPoints}
            onPlayerClick={onPlayerClick}
          />
        </div>

        {/* Away team — bottom half: FWD toward centre, GKP at bottom */}
        <div className="relative z-10 flex flex-1 flex-col">
          <TeamHalf
            team={away}
            reversed
            totalEntries={totalEntries}
            avgPoints={avgPoints}
            onPlayerClick={onPlayerClick}
          />
        </div>
      </div>
    </div>
  );
}

function PitchMarkings() {
  // All vertical values are in px, calibrated to the fixed 831px field height (scale = 831/960).
  // Horizontal values use percentages so the markings stretch with any field width.
  // The center circle uses fixed px for both axes to stay circular at any width.
  return (
    <div className="pointer-events-none absolute inset-0" style={{ opacity: 0.22 }}>
      {/* Outer border */}
      <div className="absolute border-2 border-white" style={{ inset: 10 }} />

      {/* Halfway line */}
      <div className="absolute bg-white" style={{ top: '50%', left: 10, right: 10, height: 2, marginTop: -1 }} />

      {/* Center circle — fixed px so it stays circular */}
      <div
        className="absolute rounded-full border-2 border-white"
        style={{ top: '50%', left: '50%', width: 110, height: 110, marginTop: -55, marginLeft: -55 }}
      />

      {/* Center spot */}
      <div className="absolute rounded-full bg-white" style={{ top: '50%', left: '50%', width: 6, height: 6, marginTop: -3, marginLeft: -3 }} />

      {/* Top penalty area */}
      <div
        className="absolute border-white"
        style={{ top: 10, left: '25%', right: '25%', height: 78, borderWidth: '0 2px 2px 2px', borderStyle: 'solid' }}
      />

      {/* Top goal area */}
      <div
        className="absolute border-white"
        style={{ top: 10, left: '37.5%', right: '37.5%', height: 35, borderWidth: '0 2px 2px 2px', borderStyle: 'solid' }}
      />

      {/* Top penalty spot */}
      <div className="absolute rounded-full bg-white" style={{ top: 120, left: '50%', width: 6, height: 6, marginLeft: -3 }} />

      {/* Bottom penalty area */}
      <div
        className="absolute border-white"
        style={{ bottom: 10, left: '25%', right: '25%', height: 78, borderWidth: '2px 2px 0 2px', borderStyle: 'solid' }}
      />

      {/* Bottom goal area */}
      <div
        className="absolute border-white"
        style={{ bottom: 10, left: '37.5%', right: '37.5%', height: 35, borderWidth: '2px 2px 0 2px', borderStyle: 'solid' }}
      />

      {/* Bottom penalty spot */}
      <div className="absolute rounded-full bg-white" style={{ bottom: 120, left: '50%', width: 6, height: 6, marginLeft: -3 }} />
    </div>
  );
}
