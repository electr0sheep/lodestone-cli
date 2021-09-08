package lib

type PvpProfile struct {
	Faction            string
	Rank               string
	Title              string
	TotalXp            string
	Xp                 string
	NextXp             string
	OverallPerformance PvpPerformance
	WeeklyPerformance  PvpPerformance
}
