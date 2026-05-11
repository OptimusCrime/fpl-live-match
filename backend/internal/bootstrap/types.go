package bootstrap

type Team struct {
	ID        int
	Name      string
	ShortName string
	Code      int
}

type Player struct {
	ID          int
	Code        int
	WebName     string
	Team        int
	ElementType int
	Photo       string
	EventPoints int
}

type GameweekInfo struct {
	ID           int    `json:"id"`
	DeadlineTime string `json:"deadline_time"`
}
