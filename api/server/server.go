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
	db "khelogames/db/sqlc"
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
	public := router.Group("/auth")
	{
		public.POST("/send_otp", authServer.Otp)
		public.POST("/signup", authServer.CreateSignupFunc)
		public.POST("/users", handlersServer.CreateUserFunc)
		public.POST("/login", authServer.CreateLoginFunc)
		public.DELETE("/removeSession/:username", authServer.DeleteSessionFunc)
		public.POST("/tokens/renew_access", authServer.RenewAccessTokenFunc)
		public.GET("/user/:username", handlersServer.GetUsersFunc)
		public.GET("/getProfile/:owner", handlersServer.GetProfileFunc)
	}
	authRouter := router.Group("/api").Use(authMiddleware(server.tokenMaker))
	{
		authRouter.GET("/ws", messageServer.HandleWebSocket)
		authRouter.GET("/getAllGames", sportsServer.GetGamesFunc)
		authRouter.GET("/getGame/:id", sportsServer.GetGameFunc)
		authRouter.POST("/searchProfile", playersServer.SearchProfileFunc)
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
		authRouter.POST("/createTournament", tournamentServer.AddTournamentFunc)
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
		authRouter.GET("/getPlayerByID", playersServer.GetPlayerFunc)
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

	}
	sportRouter := router.Group("/api/:sport").Use(authMiddleware(server.tokenMaker))
	sportRouter.POST("/createTournamentMatch", tournamentServer.CreateTournamentMatch)
	sportRouter.GET("/getTeamsByGroup", tournamentServer.GetTeamsByGroupFunc)
	//sportRouter.GET("/getTeams/:tournament_id", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getTournamentTeam/:tournament_id", tournamentServer.GetTournamentTeamsFunc)
	sportRouter.GET("/getTournamentsBySport/:game_id", tournamentServer.GetTournamentsBySportFunc)
	sportRouter.GET("/getTournament/:tournament_id", tournamentServer.GetTournamentFunc)

	sportRouter.POST("/addFootballGoalByPlayer", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.GET("/getAllTournamentMatch", tournamentServer.GetTournamentMatch)

	sportRouter.PUT("/updateTournamentStanding", tournamentServer.UpdateTournamentStandingFunc)
	sportRouter.PUT("/updateTournamentDate", tournamentServer.UpdateTournamentDateFunc)

	sportRouter.POST("/createTournamentStanding", tournamentServer.CreateTournamentStandingFunc)
	sportRouter.POST("/createTournamentGroup", tournamentServer.CreateTournamentGroupFunc)
	sportRouter.GET("/getTournamentGroup", tournamentServer.GetTournamentGroupFunc)
	sportRouter.GET("/getTournamentGroups", tournamentServer.GetTournamentGroupsFunc)
	sportRouter.GET("/getTournamentStanding", tournamentServer.GetTournamentStandingFunc)
	sportRouter.POST("/addTournamentTeam", tournamentServer.AddTournamentTeamFunc)
	sportRouter.GET("/getTournamentByLevel", tournamentServer.GetTournamentByLevelFunc)
	sportRouter.PUT("/updateMatchStatus", tournamentServer.UpdateMatchStatusFunc)
	sportRouter.PUT("/updateTournamentStatus", tournamentServer.UpdateTournamentStatusFunc)
	sportRouter.POST("/addGroupTeam", tournamentServer.AddGroupTeamFunc)

	//teams
	sportRouter.POST("/newTeams", teamsServer.AddTeam)
	sportRouter.GET("/getTeam", teamsServer.GetTeamFunc)
	sportRouter.GET("/getTeams", teamsServer.GetTeamsFunc)
	sportRouter.GET("/searchTeams", teamsServer.SearchTeamFunc)
	sportRouter.POST("/addTeamsMemberFunc", teamsServer.AddTeamsMemberFunc)
	sportRouter.GET("/getTeamsMemberFunc", teamsServer.GetTeamsMemberFunc)
	sportRouter.GET("/getTeamsBySport/:game_id", teamsServer.GetTeamsBySportFunc)
	sportRouter.GET("/getMatchByTeamFunc", teamsServer.GetMatchByTeamFunc)
	sportRouter.GET("/getTournamentByTeamFunc", teamsServer.GetTournamentbyTeamFunc)

	//football
	// sportRouter.GET("/getFootballScore", footballServer.GetFootballScore)
	sportRouter.POST("/addFootballIncidents", footballServer.AddFootballIncidents)
	sportRouter.GET("/getFootballIncidents", footballServer.GetFootballIncidents)
	sportRouter.POST("/addFootballIncidentsSubs", footballServer.AddFootballIncidentsSubs)
	sportRouter.PUT("/updateFootballFirstHalfScore", footballServer.UpdateFootballMatchScoreFirstHalfFunc)
	sportRouter.PUT("/updateFootballSecondHalfScore", footballServer.UpdateFootballMatchScoreSecondHalfFunc)
	sportRouter.PUT("/updateFootballMatchScore", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.POST("/addFootballMatchScore", footballServer.AddFootballMatchScoreFunc)

	sportRouter.POST("/addFootballPenalty", footballServer.AddFootballPenaltyFunc)
	sportRouter.POST("/getFootballPenalty", footballServer.GetFootballPenaltyFunc)
	sportRouter.POST("/updateFootballPenaltyScore", footballServer.UpdateFootballPenaltyFunc)

	//football->player
	sportRouter.PUT("/updateCurrentTeamByPlayer", teamsServer.UpdateCurrentTeamByPlayerFunc)
	sportRouter.POST("/addFootballLineUp", footballServer.AddFootballLineUpFunc)
	sportRouter.POST("/addFootballSubstitution", footballServer.AddFootballSubstitionFunc)
	sportRouter.GET("/getFootballLineUp", footballServer.GetFootballLineUpFunc)
	sportRouter.GET("/getFootballSubstitution", footballServer.GetFootballSubstitutionFunc)
	sportRouter.PUT("/updateFootballSubsAndLineUp", footballServer.UpdateFootballSubsAndLineUpFunc)

	sportRouter.POST("/addFootballStatistics", footballServer.AddFootballStatisticsFunc)
	sportRouter.GET("/getFootballStatistics", footballServer.GetFootballStatisticsFunc)
	// sportRouter.PUT("/updateFootballStatistics", footballServer.UpdateFootballStatisticsFunc)

	//cricket
	sportRouter.POST("/addCricketScore", cricketServer.AddCricketScoreFunc)
	sportRouter.POST("/addCricketToss", cricketServer.AddCricketToss)
	sportRouter.GET("/getCricketToss", cricketServer.GetCricketTossFunc)
	sportRouter.PUT("/updateCricketInning", cricketServer.UpdateCricketInningsFunc)
	//cricket->player
	sportRouter.POST("addCricketBatScore", cricketServer.AddCricketBatScoreFunc)
	sportRouter.POST("/addCricketBall", cricketServer.AddCricketBallFunc)
	sportRouter.GET("/getPlayerScoreFunc", cricketServer.GetPlayerScoreFunc)
	sportRouter.GET("/getCricketBowlerFunc", cricketServer.GetCricketBowlerFunc)
	sportRouter.POST("/addCricketWicket", cricketServer.AddCricketWicketFunc)
	sportRouter.PUT("/updateCricketBat", cricketServer.UpdateCricketBatScoreFunc)
	sportRouter.PUT("/updateCricketBall", cricketServer.UpdateCricketBallFunc)
	sportRouter.GET("/getCricketWickets", cricketServer.GetCricketWicketsFunc)

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
