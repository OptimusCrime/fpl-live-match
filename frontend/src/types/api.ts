export interface League {
  id: number;
  name: string;
  entryRank: number;
}

export interface EntryResponse {
  id: number;
  name: string;
  leagues: League[];
}

export interface Fixture {
  id: number;
  teamH: number;
  teamA: number;
  teamHName: string;
  teamAName: string;
  teamHShort: string;
  teamAShort: string;
  teamHCode: number;
  teamACode: number;
  scoreH: number | null;
  scoreA: number | null;
  status: "finished" | "live" | "upcoming";
  kickoffTime: string;
  minutes: number;
}

export interface FixturesResponse {
  gw: number;
  deadlineTime: string;
  fixtures: Fixture[];
}

export interface MatchPlayer {
  id: number;
  webName: string;
  photoUrl: string;
  position: string;
  livePoints: number;
  minutes: number;
  started: boolean;
  isStarter: boolean;
  subbedOn: boolean;
  subbedOff: boolean;
  ownedCount: number;
  isOwned: boolean;
}

export interface MatchTeam {
  id: number;
  name: string;
  shortName: string;
  code: number;
  players: MatchPlayer[];
}

export interface MatchResponse {
  fixture: Fixture;
  home: MatchTeam;
  away: MatchTeam;
  totalEntries: number;
}

export interface PlayerOwner {
  entryName: string;
  playerName: string;
  rank: number;
}

export interface PlayerOwnersResponse {
  playerId: number;
  owners: PlayerOwner[];
  totalEntries: number;
}
