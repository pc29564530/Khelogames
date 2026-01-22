package errorhandler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validationFieldAliases = map[string]string{
	"TournamentPublicID": "tournament_public_id",
	"AwayTeamPublicID":   "away_team_public_id",
	"HomeTeamPublicID":   "home_team_public_id",
	"StartTimestamp":     "start_timestamp",
	"EndTimestamp":       "end_timestamp",
	"MatchFormat":        "match_format",
	"KnockoutLevelID":    "knockout_level_id",
	"CommunityName":      "name",
	"CommunityPublicID":  "community_public_id",
	"Title":              "title",
	"MediaURL":           "media_url",
	"HasKnockout":        "has_knockout",
	"Level":              "level",
	"GameID":             "game_id",
	"GroupCount":         "group_count",
	"MaxGroupTeams":      "max_group_teams",
	"City":               "city",
	"State":              "State",
	"Country":            "country",
	"Stage":              "stage",
	"Name":               "name",
	"Gender":             "gender",
	"PlayerCount":        "player_count",
	"Type":               "type",
	"CommentText":        "comment_text",
	"FullName":           "full_name",
}

func ExtractValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	fmt.Println("Line no 12 error; ", err)
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		fmt.Println("Line no 26: ", errors)
		return errors
	}

	for _, fe := range ve {
		field := fe.Field()

		if alias, ok := validationFieldAliases[field]; ok {
			fmt.Println("Alias: ", alias)
			field = alias
		} else {
			field = strings.ToLower(field)
		}
		fmt.Println("Fields : ", fe.Tag())
		switch fe.Tag() {

		case "required":
			errors[field] = "This field is required"

		case "email", "email_regex":
			errors[field] = "Enter a valid email address"

		case "min":
			errors[field] = fmt.Sprintf(
				"Must be at least %s characters",
				fe.Param(),
			)

		case "max":
			errors[field] = fmt.Sprintf(
				"Must be at most %s characters",
				fe.Param(),
			)

		case "strong_password":
			errors[field] = "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, one number, and one special character"

		case "oneof":
			errors[field] = "Invalid selection"

		case "url":
			errors[field] = "Enter a valid URL"

		default:
			errors[field] = "Invalid value"
		}
	}
	fmt.Println("Errors: ", errors)
	return errors
}
