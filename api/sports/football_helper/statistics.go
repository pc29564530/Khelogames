package footballutils

// Safe helper for number conversion
func GetInt32(v interface{}) int32 {
	switch val := v.(type) {
	case nil:
		return 0
	case int:
		return int32(val)
	case int32:
		return val
	case int64:
		return int32(val)
	case float32:
		return int32(val)
	case float64:
		return int32(val)
	default:
		return 0
	}
}

func GetStatisticsUpdateFromIncident(currentStats map[string]interface{}, incidentType string) map[string]interface{} {
	if currentStats == nil {
		currentStats = make(map[string]interface{})
	}

	switch incidentType {

	case "goal":
		currentStats["shots_on_target"] = GetInt32(currentStats["shots_on_target"]) + 1
		currentStats["total_shots"] = GetInt32(currentStats["total_shots"]) + 1

	case "shot_on_target":
		currentStats["shots_on_target"] = GetInt32(currentStats["shots_on_target"]) + 1
		currentStats["total_shots"] = GetInt32(currentStats["total_shots"]) + 1

	case "foul":
		currentStats["fouls"] = GetInt32(currentStats["fouls"]) + 1

	case "corner_kick":
		currentStats["corner_kicks"] = GetInt32(currentStats["corner_kicks"]) + 1

	case "free_kick":
		currentStats["free_kicks"] = GetInt32(currentStats["free_kicks"]) + 1

	case "yellow_card":
		currentStats["yellow_cards"] = GetInt32(currentStats["yellow_cards"]) + 1

	case "red_card":
		currentStats["red_cards"] = GetInt32(currentStats["red_cards"]) + 1

	case "total_shot":
		currentStats["total_shots"] = GetInt32(currentStats["total_shots"]) + 1

	case "penalty":
		currentStats["shots_on_target"] = GetInt32(currentStats["shots_on_target"]) + 1
		currentStats["total_shots"] = GetInt32(currentStats["total_shots"]) + 1

	case "missed_penalty":
		currentStats["goal_keeper_saves"] = GetInt32(currentStats["goal_keeper_saves"]) + 1
		currentStats["shots_on_target"] = GetInt32(currentStats["shots_on_target"]) + 1
		currentStats["total_shots"] = GetInt32(currentStats["total_shots"]) + 1

	case "goal_keeper_saves":
		currentStats["goal_keeper_saves"] = GetInt32(currentStats["goal_keeper_saves"]) + 1
		currentStats["shots_on_target"] = GetInt32(currentStats["shots_on_target"]) + 1
		currentStats["total_shots"] = GetInt32(currentStats["total_shots"]) + 1
	}

	return currentStats
}
