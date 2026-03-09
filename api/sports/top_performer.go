package sports

import "github.com/gin-gonic/gin"

func (s *SportsServer) GetTopPerformerHandler(c *gin.Context) {
	sport := c.Param("sport")
	switch sport {
	case "football":
		s.footballServer.GetFootballTopPerformerFunc(c)
	case "cricket":
		s.cricketServer.GetCricketTopPerformerFunc(c)
	default:
		c.JSON(400, gin.H{
			"error": "invalid sport",
		})
	}
}
