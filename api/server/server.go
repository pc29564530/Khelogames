package server

import (
	"khelogames/api/auth"
	"khelogames/api/handlers"
	"khelogames/api/messenger"
	"khelogames/api/players"
	"khelogames/api/sports"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/api/teams"
	"khelogames/api/tournaments"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/token"
	util "khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config        util.Config
	store         *db.Store
	tokenMaker    token.Maker
	logger        *logger.Logger
	router        *gin.Engine
	messageServer *messenger.MessageServer
}

func NewServer(config util.Config,
	store *db.Store,
	tokenMaker token.Maker,
	logger *logger.Logger,
	authServer *auth.AuthServer,
	handlersServer *handlers.HandlersServer,
	tournamentServer *tournaments.TournamentServer,

	footballServer *football.FootballServer,
	cricketServer *cricket.CricketServer,
	teamsServer *teams.TeamsServer,
	messageServer *messenger.MessageServer,
	playersServer *players.PlayerServer,
	sportsServer *sports.SportsServer,
	router *gin.Engine,
) (*Server, error) {

	server := &Server{
		config:        config,
		store:         store,
		tokenMaker:    tokenMaker,
		logger:        logger,
		router:        router,
		messageServer: messageServer,
	}

	router.Use(corsHandle())
	router.StaticFS("/api/images", http.Dir("/Users/pawan/database/Khelogames/images"))
	router.StaticFS("/api/videos", http.Dir("/Users/pawan/database/Khelogames/videos"))
	router.StaticFS("/media", http.Dir("/tmp/khelogames_media_uploads"))

	public := router.Group("/auth")
	{
		public.POST("/send_otp", authServer.Otp)
		public.POST("/users", handlersServer.CreateUserFunc)
		public.DELETE("/removeSession/:username", authServer.DeleteSessionFunc)
		public.POST("/tokens/renew_access", authServer.RenewAccessTokenFunc)
		public.GET("/user/:username", handlersServer.GetUsersFunc)
		public.GET("/getProfile/:owner", handlersServer.GetProfileFunc)
		public.POST("/google/createGoogleSignUp", authServer.CreateGoogleSignUp)
		public.POST("/google/createGoogleSignIn", authServer.CreateGoogleSignIn)
		public.GET("/google/handleGoogleRedirect", authServer.HandleGoogleRedirect)
		public.POST("/createMobileSignup", authServer.CreateMobileSignUp)
		public.POST("/createMobileSignin", authServer.CreateMobileSignIn)
		public.GET("/getUserByMobileNumber", handlersServer.GetUserByMobileNumber)
		public.GET("/getUserByGmail", handlersServer.GetUserByGmail)
	}

	authRouter := router.Group("/api").Use(authMiddleware(server.tokenMaker))
	{
		// added the funcitonality for the matches by player
		authRouter.GET("/getGroups", tournamentServer.GetGroupsFunc)
		authRouter.GET("/isFollowing", handlersServer.IsFollowingFunc)
		authRouter.GET("/checkConnection", handlersServer.CheckConnectionFunc)
		authRouter.PUT("/updateProfile", handlersServer.UpdateProfileFunc)
		authRouter.GET("/ws", messageServer.HandleWebSocket)
		authRouter.GET("/getAllGames", sportsServer.GetGamesFunc)
		authRouter.GET("/getGame/:id", sportsServer.GetGameFunc)
		// authRouter.POST("/searchProfile", playersServer.SearchProfileFunc)
		authRouter.POST("/addJoinCommunity", handlersServer.AddJoinCommunityFunc)
		authRouter.GET("/getUserByCommunity/:community_name", handlersServer.GetUserByCommunityFunc)
		authRouter.GET("/getCommunityByUser", handlersServer.GetCommunityByUserFunc)
		authRouter.GET("/user_list", handlersServer.ListUsersFunc)
		authRouter.POST("/communities", handlersServer.CreateCommunitesFunc)
		//authRouter.GET("/communities/:id", server.GetCommunitiesFunc)
		authRouter.GET("/community/:id", handlersServer.GetCommunityFunc)
		authRouter.GET("/get_all_communities", handlersServer.GetAllCommunitiesFunc)
		authRouter.GET("/getCommunityByCommunityName/:communities_name", handlersServer.GetCommunityByCommunityNameFunc)
		authRouter.POST("/create_thread", handlersServer.CreateThreadFunc)
		authRouter.GET("/getThread/:id", handlersServer.GetThreadFunc)
		authRouter.PUT("/update_like", handlersServer.UpdateThreadLikeFunc)
		authRouter.GET("/all_threads", handlersServer.GetAllThreadsFunc)
		authRouter.GET("/getAllThreadByCommunity/:communities_name", handlersServer.GetAllThreadsByCommunitiesFunc)
		authRouter.GET("/get_communities_member/:communities_name", handlersServer.GetCommunitiesMemberFunc)
		authRouter.POST("/create_follow/:following_owner", handlersServer.CreateFollowingFunc)
		authRouter.GET("/getFollower", handlersServer.GetAllFollowerFunc)
		authRouter.GET("/getFollowing", handlersServer.GetAllFollowingFunc)
		authRouter.POST("/createComment/:threadId", handlersServer.CreateCommentFunc)
		authRouter.GET("/getComment/:thread_id", handlersServer.GetAllCommentFunc)
		authRouter.GET("/getCommentByUser/:username", handlersServer.GetCommentByUserFunc)
		authRouter.DELETE("/unFollow/:following_owner", handlersServer.DeleteFollowingFunc)
		authRouter.POST("/createLikeThread/:thread_id", handlersServer.CreateLikeFunc)
		authRouter.GET("/countLike/:thread_id", handlersServer.CountLikeFunc)
		authRouter.GET("/checkLikeByUser/:thread_id", handlersServer.CheckLikeByUserFunc)
		authRouter.POST("/createProfile", handlersServer.CreateProfileFunc)
		authRouter.PUT("/editProfile", handlersServer.UpdateProfileFunc)
		authRouter.PUT("/updateAvatar", handlersServer.UpdateAvatarUrlFunc)
		authRouter.PUT("/updateFullName", handlersServer.UpdateFullNameFunc)
		authRouter.PUT("/updateBio", handlersServer.UpdateBioFunc)
		authRouter.GET("getThreadByUser/:username", handlersServer.GetThreadByUserFunc)
		authRouter.GET("/getMessage/:receiver_username", messageServer.GetMessageByReceiverFunc)
		authRouter.PUT("/updateAvatarUrl", handlersServer.UpdateAvatarUrlFunc)
		authRouter.GET("/getMessagedUser", messageServer.GetUserByMessageSendFunc)
		authRouter.POST("/createUploadMedia", messageServer.CreateUploadMediaFunc)
		authRouter.POST("/createMessageMedia", messageServer.CreateMessageMediaFunc)
		authRouter.POST("/createCommunityMessage", messageServer.CreateCommunityMessageFunc)
		authRouter.GET("/getCommunityMessage", messageServer.GetCommunityByMessageFunc)
		authRouter.GET("/getCommunityByMessage", messageServer.GetCommunityByMessageFunc)
		authRouter.GET("/GetAllThreadDetailFunc", handlersServer.GetAllThreadDetailFunc)
		authRouter.GET("/GetAllThreadsByCommunityDetailsFunc/:communities_name", handlersServer.GetAllThreadsByCommunityDetailsFunc)
		//player
		authRouter.POST("/newPlayer", playersServer.NewPlayerFunc)
		authRouter.GET("/getPlayerByCountry", playersServer.GetPlayerByCountry)
		authRouter.GET("/getPlayersBySport", playersServer.GetPlayersBySportFunc)
		authRouter.GET("/getPlayerByID", playersServer.GetPlayerFunc)
		authRouter.GET("/getPlayerByProfileID", playersServer.GetPlayerByProfileIDFunc)
		authRouter.GET("/getAllPlayers", playersServer.GetAllPlayerFunc)
		authRouter.GET("/getPlayerSearch", playersServer.GetPlayerSearchFunc)
		authRouter.GET("/updatePlayerMedia", playersServer.UpdatePlayerMediaFunc)
		authRouter.GET("/updatePlayerPosition", playersServer.UpdatePlayerPositionFunc)

		authRouter.PUT("/inActiveUserFromCommunity", handlersServer.InActiveUserFromCommunityFunc)
		authRouter.PUT("/updateDeleteMessage", messageServer.UpdateDeleteMessageFunc)
		authRouter.DELETE("/deleteScheduleMessage", messageServer.DeleteScheduleMessageFunc)
		authRouter.DELETE("/deleteCommentByUser", handlersServer.DeleteCommentByUserFunc)
		authRouter.DELETE("/deleteAdmin", handlersServer.DeleteAdminFunc)
		authRouter.PUT("/updateCommunityByDescription", handlersServer.UpdateCommunityByDescriptionFunc)
		authRouter.PUT("/updateCommunityByCommunityName", handlersServer.UpdateCommunityByCommunityNameFunc)
		authRouter.GET("/getRoles", handlersServer.GetRolesFunc)
		authRouter.POST("/addUserRole", handlersServer.AddUserRoleFunc)
		authRouter.POST("/applyForVerification", handlersServer.AddUserVerificationFunc)
		authRouter.GET("/getPlayerCricketStats", playersServer.GetPlayerCricketStatsByMatchTypeFunc)
		authRouter.GET("/getFootballPlayerStats/:playerID", playersServer.GetFootballPlayerStatsFunc)
		authRouter.POST("/createUploadChunks", handlersServer.CreateUploadMediaFunc)
		authRouter.POST("/completedChunkUpload", handlersServer.CompletedChunkUploadFunc)
	}
	sportRouter := router.Group("/api/:sport").Use(authMiddleware(server.tokenMaker))
	sportRouter.POST("/createTournamentMatch", tournamentServer.CreateTournamentMatch)
	sportRouter.POST("/createTournament", tournamentServer.AddTournamentFunc)
	//sportRouter.GET("/getTeamsByGroup", tournamentServer.GetTeamsByGroupFunc)
	//sportRouter.GET("/getTeams/:tournament_id", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getTournamentTeam/:tournament_id", tournamentServer.GetTournamentTeamsFunc)
	sportRouter.GET("/getTournamentsBySport/:game_id", tournamentServer.GetTournamentsBySportFunc)
	sportRouter.GET("/getTournament/:tournament_id", tournamentServer.GetTournamentFunc)

	sportRouter.POST("/addFootballGoalByPlayer", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.GET("/getAllTournamentMatch", tournamentServer.GetTournamentMatch)
	sportRouter.GET("/getFootballStanding", tournamentServer.GetFootballStandingFunc)
	sportRouter.GET("/getCricketStanding", tournamentServer.GetCricketStandingFunc)
	sportRouter.PUT("/updateFootballStanding", tournamentServer.UpdateFootballStandingFunc)
	sportRouter.PUT("/updateCricketStanding", tournamentServer.UpdateCricketStandingFunc)
	sportRouter.PUT("/updateTournamentDate", tournamentServer.UpdateTournamentDateFunc)

	sportRouter.POST("/createTournamentStanding", tournamentServer.CreateTournamentStandingFunc)
	sportRouter.POST("/addTournamentTeam", tournamentServer.AddTournamentTeamFunc)
	sportRouter.GET("/getTournamentByLevel", tournamentServer.GetTournamentByLevelFunc)
	sportRouter.PUT("/updateMatchStatus", tournamentServer.UpdateMatchStatusFunc)
	sportRouter.PUT("/updateMatchResult", tournamentServer.UpdateMatchResultFunc)
	sportRouter.PUT("/updateTournamentStatus", tournamentServer.UpdateTournamentStatusFunc)
	sportRouter.GET("/getMatchByMatchID", handlersServer.GetMatchByMatchIDFunc)

	//teams
	sportRouter.POST("/newTeams", teamsServer.AddTeam)
	sportRouter.GET("/getTeam", teamsServer.GetTeamFunc)
	sportRouter.GET("/getTeams", teamsServer.GetTeamsFunc)
	sportRouter.GET("/searchTeams", teamsServer.SearchTeamFunc)
	sportRouter.POST("/addTeamsMemberFunc", teamsServer.AddTeamsMemberFunc)
	sportRouter.PUT("/removePlayerFromTeam", teamsServer.RemovePlayerFromTeamFunc)
	sportRouter.GET("/getTeamsMemberFunc", teamsServer.GetTeamsMemberFunc)
	sportRouter.GET("/getTeamsBySport/:game_id", teamsServer.GetTeamsBySportFunc)
	sportRouter.GET("/getMatchByTeamFunc", teamsServer.GetMatchByTeamFunc)
	sportRouter.GET("/getMatchesByTeam", teamsServer.GetMatchesByTeamFunc)
	sportRouter.GET("/getTournamentByTeamFunc", teamsServer.GetTournamentbyTeamFunc)

	sportRouter.GET("/getAllMatches", handlersServer.GetAllMatchesFunc)

	//football
	// sportRouter.GET("/getFootballScore", footballServer.GetFootballScore)
	sportRouter.POST("/addFootballIncidents", footballServer.AddFootballIncidents)
	sportRouter.GET("/getFootballIncidents", footballServer.GetFootballIncidents)
	sportRouter.POST("/addFootballIncidentsSubs", footballServer.AddFootballIncidentsSubs)
	sportRouter.PUT("/updateFootballFirstHalfScore", footballServer.UpdateFootballMatchScoreFirstHalfFunc)
	sportRouter.PUT("/updateFootballSecondHalfScore", footballServer.UpdateFootballMatchScoreSecondHalfFunc)
	sportRouter.PUT("/updateFootballMatchScore", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.POST("/addFootballMatchScore", footballServer.AddFootballMatchScoreFunc)

	//football->player
	sportRouter.POST("/addFootballLineUp", footballServer.AddFootballLineUpFunc)
	sportRouter.POST("/addFootballSubstitution", footballServer.AddFootballSubstitionFunc)
	sportRouter.GET("/getFootballLineUp", footballServer.GetFootballLineUpFunc)
	sportRouter.GET("/getFootballMatchSquad", footballServer.GetFootballMatchSquadFunc)
	sportRouter.GET("/getFootballSubstitution", footballServer.GetFootballSubstitutionFunc)
	sportRouter.PUT("/updateFootballSubsAndLineUp", footballServer.UpdateFootballSubsAndLineUpFunc)

	sportRouter.POST("/addFootballStatistics", footballServer.AddFootballStatisticsFunc)
	sportRouter.GET("/getFootballStatistics", footballServer.GetFootballStatisticsFunc)
	sportRouter.POST("/addFootballMatchSquad", footballServer.AddFootballSquadFunc)
	sportRouter.GET("/getFootballTournamentPlayerGoal/:id", tournamentServer.GetFootballTournamentPlayersGoalsFunc)
	sportRouter.GET("/getFootballTournamentPlayerYellowCard/:id", tournamentServer.GetFootballTournamentPlayersYellowCardFunc)
	sportRouter.GET("/getFootballTournamentPlayerRedCard/:id", tournamentServer.GetFootballTournamentPlayersRedCardFunc)

	// sportRouter.PUT("/updateFootballStatistics", footballServer.UpdateFootballStatisticsFunc)

	//cricket
	sportRouter.POST("/addCricketScore", cricketServer.AddCricketScoreFunc)
	sportRouter.POST("/addCricketToss", cricketServer.AddCricketToss)
	sportRouter.GET("/getCricketToss", cricketServer.GetCricketTossFunc)
	sportRouter.PUT("/updateCricketInning", cricketServer.UpdateCricketInningsFunc)
	sportRouter.PUT("/updateCricketEndInning", cricketServer.UpdateCricketEndInningsFunc)
	sportRouter.PUT("/updateCricketNoBall", cricketServer.UpdateNoBallsRunsFunc)
	sportRouter.PUT("/updateCricketWide", cricketServer.UpdateWideBallFunc)
	sportRouter.PUT("/updateCricketRegularScore", cricketServer.UpdateInningScoreFunc)
	sportRouter.GET("/getCurrentBatsman", cricketServer.GetCurrentBatsmanFunc)
	sportRouter.GET("/getCurrentBowler", cricketServer.GetCurrentBowlerFunc)
	//squad
	sportRouter.POST("/addCricketSquad", cricketServer.AddCricketSquadFunc)
	sportRouter.GET("/getCricketMatchSquad", cricketServer.GetCricketMatchSquadFunc)
	//tournament data
	sportRouter.GET("/getCricketTournamentMostRuns/:id", tournamentServer.GetCricketTournamentMostRunsFunc)
	sportRouter.GET("/getCricketTournamentHighestRuns/:id", tournamentServer.GetCricketTournamentHighestRunsFunc)
	sportRouter.GET("getCricketTournamentMostSixes/:id", tournamentServer.GetCricketTournamentMostSixesFunc)
	sportRouter.GET("/getCricketTournamentMostFours/:id", tournamentServer.GetCricketTournamentMostFoursFunc)
	sportRouter.GET("/getCricketTournamentMostFifties/:id", tournamentServer.GetCricketTournamentMostFiftiesFunc)
	sportRouter.GET("/getCricketTournamentMostHundreds/:id", tournamentServer.GetCricketTournamentMostHundredsFunc)
	sportRouter.GET("/getCricketTournamentBowlingStrike/:id", tournamentServer.GetCricketTournamentBowlingStrikeRateFunc)
	sportRouter.GET("/getCricketTournamentBowlingEconomy/:id", tournamentServer.GetCricketTournamentBowlingEconomyRateFunc)
	sportRouter.GET("/getCricketTournamentBowlingAverage/:id", tournamentServer.GetCricketTournamentBowlingAverageFunc)
	sportRouter.GET("/getCricketTournamentMostWickets/:id", tournamentServer.GetCricketTournamentMostWicketsFunc)
	sportRouter.GET("/getCricketTournamentFiveWicketsHaul/:id", tournamentServer.GetCricketTournamentBowlingFiveWicketHaulFunc)
	sportRouter.GET("/getCricketTournamentBattingAverage/:id", tournamentServer.GetCricketTournamentBattingAverageFunc)
	sportRouter.GET("/getCricketTournamentBattingStrike/:id", tournamentServer.GetCricketTournamentBattingStrikeFunc)
	//cricket->player
	sportRouter.POST("/addCricketBatScore", cricketServer.AddCricketBatScoreFunc)
	sportRouter.POST("/addCricketBall", cricketServer.AddCricketBallFunc)
	sportRouter.GET("/getPlayerScoreFunc", cricketServer.GetPlayerScoreFunc)
	sportRouter.GET("/getCricketBowlerFunc", cricketServer.GetCricketBowlerFunc)
	sportRouter.PUT("/updateCricketBat", cricketServer.UpdateCricketBatScoreFunc)
	sportRouter.PUT("/updateCricketBall", cricketServer.UpdateCricketBallFunc)
	sportRouter.GET("/getCricketWickets", cricketServer.GetCricketWicketsFunc)
	sportRouter.POST("/wickets", cricketServer.AddCricketWicketsFunc)
	sportRouter.PUT("/updateBowlingBowlerStatus", cricketServer.UpdateBowlingBowlerFunc)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	go server.messageServer.StartWebSocketHub()
	return server.router.Run(address)
}

func corsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
