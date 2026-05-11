import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AppState {
  teamId: string | null;
  leagueId: number | null;
  leagueName: string | null;
  setTeamId: (id: string) => void;
  setLeague: (id: number, name: string) => void;
  clearTeam: () => void;
  clearLeague: () => void;
}

export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      teamId: null,
      leagueId: null,
      leagueName: null,
      setTeamId: (id) => set({ teamId: id }),
      setLeague: (id, name) => set({ leagueId: id, leagueName: name }),
      clearTeam: () => set({ teamId: null, leagueId: null, leagueName: null }),
      clearLeague: () => set({ leagueId: null, leagueName: null }),
    }),
    { name: "fpl-live-match" }
  )
);
