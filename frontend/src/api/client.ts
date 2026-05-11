import type {
  EntryResponse,
  FixturesResponse,
  MatchResponse,
  PlayerOwnersResponse,
} from "@/types/api";

const BASE_URL = import.meta.env.PROD ? "https://fpl.optimuscrime.net" : "";

async function get<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, init);
  if (!res.ok) {
    throw new Error(`${res.status} ${res.statusText} — ${path}`);
  }
  return res.json() as Promise<T>;
}

export const api = {
  getEntry: (teamId: string) =>
    get<EntryResponse>(`/api/team/${teamId}/leagues`),

  getFixtures: () =>
    get<FixturesResponse>("/api/fixtures"),

  getLeagueMatch: (leagueId: number, matchId: number, profileId: string) =>
    get<MatchResponse>(`/api/league/${leagueId}/match/${matchId}`, {
      headers: { "X-CURRENT-PROFILE-ID": profileId },
    }),

  getPlayerOwners: (playerId: number, leagueId: number) =>
    get<PlayerOwnersResponse>(`/api/player/${playerId}/owners?league=${leagueId}`),
};
