package fpl

// Bootstrap types

type Bootstrap struct {
	Events       []Event       `json:"events"`
	Teams        []Team        `json:"teams"`
	Elements     []Element     `json:"elements"`
	ElementTypes []ElementType `json:"element_types"`
}

type Event struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DeadlineTime string `json:"deadline_time"`
	IsCurrent    bool   `json:"is_current"`
	IsNext       bool   `json:"is_next"`
	IsPrevious   bool   `json:"is_previous"`
	Finished     bool   `json:"finished"`
	DataChecked  bool   `json:"data_checked"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Code      int    `json:"code"`
}

type Element struct {
	ID          int     `json:"id"`
	Code        int     `json:"code"`
	WebName     string  `json:"web_name"`
	FirstName   string  `json:"first_name"`
	SecondName  string  `json:"second_name"`
	Team        int     `json:"team"`
	ElementType int     `json:"element_type"`
	Photo       string  `json:"photo"`
	NowCost     int     `json:"now_cost"`
	EventPoints int     `json:"event_points"`
	TotalPoints int     `json:"total_points"`
	SelectedBy  string  `json:"selected_by_percent"`
	Status      string  `json:"status"`
}

type ElementType struct {
	ID              int    `json:"id"`
	SingularName    string `json:"singular_name"`
	SingularShort   string `json:"singular_name_short"`
	PluralName      string `json:"plural_name"`
}

// Fixture types

type Fixture struct {
	ID                  int            `json:"id"`
	Event               int            `json:"event"`
	TeamH               int            `json:"team_h"`
	TeamA               int            `json:"team_a"`
	TeamHScore          *int           `json:"team_h_score"`
	TeamAScore          *int           `json:"team_a_score"`
	Finished            bool           `json:"finished"`
	FinishedProvisional bool           `json:"finished_provisional"`
	Started             bool           `json:"started"`
	Minutes             int            `json:"minutes"`
	KickoffTime         string         `json:"kickoff_time"`
	Stats               []FixtureStat  `json:"stats"`
}

type FixtureStat struct {
	Identifier string              `json:"identifier"`
	H          []FixtureStatValue  `json:"h"`
	A          []FixtureStatValue  `json:"a"`
}

type FixtureStatValue struct {
	Value   int `json:"value"`
	Element int `json:"element"`
}

// Live types

type LiveResponse struct {
	Elements []LiveElement `json:"elements"`
}

type LiveElement struct {
	ID    int       `json:"id"`
	Stats LiveStats `json:"stats"`
}

type LiveStats struct {
	Minutes      int  `json:"minutes"`
	Starts       int  `json:"starts"`
	GoalsScored  int  `json:"goals_scored"`
	Assists      int  `json:"assists"`
	CleanSheets  int  `json:"clean_sheets"`
	YellowCards  int  `json:"yellow_cards"`
	RedCards     int  `json:"red_cards"`
	Saves        int  `json:"saves"`
	Bonus        int  `json:"bonus"`
	TotalPoints  int  `json:"total_points"`
	InDreamteam  bool `json:"in_dreamteam"`
}

// Entry (team) types

type Entry struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	PlayerFirstName string      `json:"player_first_name"`
	PlayerLastName  string      `json:"player_last_name"`
	CurrentEvent    int         `json:"current_event"`
	Leagues         EntryLeagues `json:"leagues"`
}

type EntryLeagues struct {
	Classic []ClassicLeague `json:"classic"`
}

type ClassicLeague struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	EntryRank  int    `json:"entry_rank"`
	EntryLastRank int `json:"entry_last_rank"`
}

// Picks types

type Picks struct {
	ActiveChip string     `json:"active_chip"`
	Picks      []Pick     `json:"picks"`
}

type Pick struct {
	Element        int  `json:"element"`
	Position       int  `json:"position"`
	Multiplier     int  `json:"multiplier"`
	IsCaptain      bool `json:"is_captain"`
	IsViceCaptain  bool `json:"is_vice_captain"`
}

// League standings types

type LeagueStandings struct {
	League    LeagueInfo         `json:"league"`
	Standings StandingsWrapper   `json:"standings"`
}

type LeagueInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StandingsWrapper struct {
	Results []StandingEntry `json:"results"`
	HasNext bool            `json:"has_next"`
}

type StandingEntry struct {
	Entry      int    `json:"entry"`
	EntryName  string `json:"entry_name"`
	PlayerName string `json:"player_name"`
	Rank       int    `json:"rank"`
	Total      int    `json:"total"`
}
