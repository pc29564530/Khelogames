package footballutils

type StatisticsUpdate struct {
	Penalty         int32
	ShotsOnTarget   int32
	TotalShots      int32
	CornerKicks     int32
	Fouls           int32
	GoalkeeperSaves int32
	FreeKicks       int32
	YellowCards     int32
	RedCards        int32
}

func GetStatisticsUpdateFromIncident(incidentType string) StatisticsUpdate {
	switch incidentType {
	case "goal":
		return StatisticsUpdate{
			ShotsOnTarget: 1,
			TotalShots:    1,
		}
	case "fouls":
		return StatisticsUpdate{
			Fouls:     1,
			FreeKicks: 1,
		}
	case "yellow_cards":
		return StatisticsUpdate{
			YellowCards: 1,
		}
	case "red_cards":
		return StatisticsUpdate{
			RedCards: 1,
		}
	case "goalkeeper_saves":
		return StatisticsUpdate{
			GoalkeeperSaves: 1,
			ShotsOnTarget:   1,
			TotalShots:      1,
		}
	case "corner_kicks":
		return StatisticsUpdate{
			CornerKicks: 1,
		}
	case "total_shots":
		return StatisticsUpdate{
			TotalShots: 1,
		}
	case "shots_on_target":
		return StatisticsUpdate{
			ShotsOnTarget: 1,
		}
	case "penalty":
		return StatisticsUpdate{
			ShotsOnTarget: 1,
			TotalShots:    1,
		}
	case "missed_penalty":
		return StatisticsUpdate{
			ShotsOnTarget:   1,
			TotalShots:      1,
			GoalkeeperSaves: 1,
		}
	default:
		return StatisticsUpdate{}
	}
}
