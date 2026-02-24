package server

import (
	"context"
	"database/sql"
	"khelogames/api/auth"
	"khelogames/api/handlers"
	"khelogames/api/messenger"
	"khelogames/api/players"
	"khelogames/api/sports"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/api/teams"
	apiToken "khelogames/api/token"
	"khelogames/api/tournaments"
	"khelogames/core/token"
	coreToken "khelogames/core/token"
	db "khelogames/database"
	"khelogames/hub"
	"khelogames/logger"
	util "khelogames/util"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config         util.Config
	store          *db.Store
	tokenMaker     token.Maker
	logger         *logger.Logger
	router         *gin.Engine
	messageServer  *messenger.MessageServer
	dbConn         *sql.DB
	rabbitChan     *amqp091.Channel
	isShuttingDown *atomic.Bool
	httpServer     *http.Server
}

const (
	PermAdmin                 = "ADMIN"
	PermUpdateMatch           = "UPDATE_MATCH"
	PermUpdateTournament      = "UPDATE_TOURNAMENT"
	PermUpdateTournamentAdmin = "UPDATE_TOURNAMENT_ADMIN"
	PermUpdateTeam            = "UPDATE_TEAM"
	PermUpdateCommunity       = "UPDATE_COMMUNITY"
)

func (server *Server) setupStaticFiles() {
	mediaBase := server.config.MediaBasePath
	if mediaBase == "" {
		mediaBase = "./media"
	}

	server.router.Static("/media/images", mediaBase+"/images")
	server.router.Static("/media/videos", mediaBase+"/videos")
	server.router.Static("/media/uploads", mediaBase+"/uploads")
}

func (server *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (server *Server) readinessCheck(ctx *gin.Context) {

	if server.isShuttingDown != nil && server.isShuttingDown.Load() {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "shutting_down",
		})
		return
	}

	if server.dbConn != nil {
		ctxPing, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := server.dbConn.PingContext(ctxPing); err != nil {
			server.logger.Error("DB health check failed", err)
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"reason": "database unreachable",
			})
			return
		}
	}

	if server.rabbitChan != nil && server.rabbitChan.IsClosed() {
		server.logger.Error("RabbitMQ health check failed: channel closed")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"reason": "rabbitmq channel close",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func NewServer(config util.Config,
	store *db.Store,
	tokenMaker coreToken.Maker,
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
	tokenServer *apiToken.TokenServer,
	hub *hub.Hub,
	conn *sql.DB,
	rabbitChan *amqp091.Channel,
	httpServer *http.Server,
) (*Server, error) {
	isShuttingDown := &atomic.Bool{}
	isShuttingDown.Store(false)

	server := &Server{
		config:         config,
		store:          store,
		tokenMaker:     tokenMaker,
		logger:         logger,
		router:         router,
		messageServer:  messageServer,
		dbConn:         conn,
		rabbitChan:     rabbitChan,
		isShuttingDown: isShuttingDown,
		httpServer:     httpServer,
	}

	server.setupStaticFiles()

	router.GET("/health", server.healthCheck)
	router.GET("/ready", server.readinessCheck)
	router.Use(server.corsHandle())

	public := router.Group("/auth")
	{
		public.DELETE("/removeSession/:public_id", authServer.DeleteSessionFunc)
		public.POST("/tokens/renew_access", tokenServer.RenewAccessTokenFunc)
		public.GET("/getProfile/:public_id", handlersServer.GetProfileFunc)
		public.GET("/getProfileByPublicID/:profile_public_id", handlersServer.GetProfileByPublicIDFunc)
		public.POST("/google", authServer.CreateGoogleSignUpFunc)
		// if server.config.AppEnv == "development" {
		public.POST("/google/createEmailSignUp", authServer.CreateEmailSignUpFunc)
		public.POST("/google/createEmailSignIn", authServer.CreateEmailSignInFunc)
		// }
	}

	authRouter := router.Group("/api").Use(authMiddleware(server.tokenMaker, server.logger))
	{
		// added the funcitonality for the matches by player
		authRouter.GET("/getPlayerWithProfile/:public_id", handlersServer.GetPlayerWithProfileFunc)
		authRouter.GET("/getGroups", tournamentServer.GetGroupsFunc)
		authRouter.GET("/isFollowing/:target_public_id", handlersServer.IsFollowingFunc)
		// authRouter.GET("/checkConnection", handlersServer.CheckConnectionFunc)
		authRouter.PUT("/updateProfile", handlersServer.UpdateProfileFunc)
		authRouter.GET("/ws", hub.HandleWebSocket)
		authRouter.GET("/getAllGames", sportsServer.GetGamesFunc)
		authRouter.GET("/getGame/:id", sportsServer.GetGameFunc)
		authRouter.POST("/search-player", playersServer.SearchPlayerFunc)
		authRouter.POST("/search-user", handlersServer.SearchUserFunc)
		authRouter.POST("/addJoinCommunity/:community_public_id", handlersServer.AddJoinCommunityFunc)
		authRouter.GET("/getCommunityByUser", handlersServer.GetCommunityByUserFunc)
		// authRouter.GET("/user_list", handlersServer.ListUsersFunc)
		authRouter.POST("/createCommunity", handlersServer.CreateCommunitesFunc)
		//authRouter.GET("/communities/:id", server.GetCommunitiesFunc)
		authRouter.GET("/community/:public_id", handlersServer.GetCommunityFunc)
		authRouter.GET("/getAllCommunities", handlersServer.GetAllCommunitiesFunc)
		authRouter.GET("/getCommunityByCommunityName/:communities_name", handlersServer.GetCommunityByCommunityNameFunc)
		authRouter.POST("/create_thread", handlersServer.CreateThreadFunc)
		authRouter.GET("/getThread/:public_id", handlersServer.GetThreadFunc)
		authRouter.GET("/all_threads", handlersServer.GetAllThreadsFunc)
		authRouter.GET("/getAllThreadByCommunity/:communities_name", handlersServer.GetAllThreadsByCommunitiesFunc)
		authRouter.GET("/getCommunityMember/:community_public_id", handlersServer.GetCommunitiesMemberFunc)
		authRouter.POST("/create_follow/:target_public_id", handlersServer.CreateUserConnectionsFunc)
		authRouter.GET("/getFollower", handlersServer.GetAllFollowerFunc)
		authRouter.GET("/getFollowing", handlersServer.GetAllFollowingFunc)
		authRouter.GET("/getFollowerCount/:public_id", handlersServer.GetFollowerCountFunc)
		authRouter.GET("/getFollowingCount/:public_id", handlersServer.GetFollowingCountFunc)
		authRouter.POST("/createComment/:thread_public_id", handlersServer.CreateCommentFunc)
		authRouter.GET("/getComments/:public_id", handlersServer.GetAllCommentFunc)
		authRouter.DELETE("/unFollow/:target_public_id", handlersServer.DeleteFollowingFunc)
		authRouter.POST("/likeThread/:thread_public_id", handlersServer.LikeThreadFunc)
		// authRouter.POST("/createProfile", handlersServer.CreateProfileFunc)
		authRouter.PUT("/editProfile", handlersServer.UpdateProfileFunc)
		// authRouter.PUT("/updateAvatar", handlersServer.UpdateAvatarUrlFunc)
		// authRouter.PUT("/updateFullName", handlersServer.UpdateFullNameFunc)
		// authRouter.PUT("/updateBio", handlersServer.UpdateBioFunc)
		authRouter.GET("/getThreadByUser/:public_id", handlersServer.GetThreadByUserFunc)
		authRouter.GET("/getMessage/:receiver_public_id", messageServer.GetMessageByReceiverFunc)
		// authRouter.PUT("/updateAvatarUrl", handlersServer.UpdateAvatarUrlFunc)
		authRouter.GET("/getMessagedUser", messageServer.GetUserByMessageSendFunc)
		// authRouter.POST("/createUploadMedia", messageServer.CreateUploadMediaFunc)
		// authRouter.POST("/createMessageMedia", messageServer.CreateMessageMediaFunc)
		authRouter.POST("/createCommunityMessage", messageServer.CreateCommunityMessageFunc)
		authRouter.GET("/getCommunityMessage", messageServer.GetCommunityByMessageFunc)
		authRouter.GET("/getCommunityByMessage", messageServer.GetCommunityByMessageFunc)
		authRouter.GET("/GetAllThreadDetailFunc", handlersServer.GetAllThreadDetailFunc)
		// authRouter.GET("/GetAllThreadsByCommunityDetailsFunc/:communities_name", handlersServer.GetAllThreadsByCommunityDetailsFunc)
		//player
		authRouter.POST("/newPlayer", playersServer.NewPlayerFunc)
		authRouter.GET("/getPlayerByCountry", playersServer.GetPlayerByCountry)
		authRouter.GET("/getPlayersBySport/:game_id", playersServer.GetPlayersBySportFunc)
		authRouter.GET("/getPlayer/:public_id", playersServer.GetPlayerFunc)
		authRouter.GET("/getPlayerByProfile/:profile_public_id", playersServer.GetPlayerByProfilePublicIDFunc)
		// authRouter.GET("/getPlayerByProfileID", playersServer.GetPlayerByProfileIDFunc)
		authRouter.GET("/getAllPlayers", playersServer.GetAllPlayerFunc)
		authRouter.GET("/getPlayerSearch", playersServer.GetPlayerSearchFunc)
		// authRouter.GET("/updatePlayerMedia", playersServer.UpdatePlayerMediaFunc)
		// authRouter.GET("/updatePlayerPosition", playersServer.UpdatePlayerPositionFunc)

		// authRouter.PUT("/updateDeleteMessage", messageServer.UpdateDeleteMessageFunc)
		authRouter.DELETE("/deleteScheduleMessage", messageServer.DeleteScheduleMessageFunc)
		authRouter.DELETE("/deleteCommentByUser", handlersServer.DeleteCommentByUserFunc)
		// authRouter.DELETE("/deleteAdmin", handlersServer.DeleteAdminFunc)
		authRouter.PUT("/updateCommunityByDescription/:community_public_id", handlersServer.UpdateCommunityByDescriptionFunc)
		authRouter.PUT("/updateCommunityByCommunityName", handlersServer.UpdateCommunityByCommunityNameFunc)
		authRouter.GET("/getRoles", handlersServer.GetRolesFunc)
		authRouter.POST("/addUserRole", handlersServer.AddUserRoleFunc)
		authRouter.POST("/applyForVerification", handlersServer.AddUserVerificationFunc)
		authRouter.GET("/getPlayerCricketStats", playersServer.GetPlayerCricketStatsByMatchTypeFunc)
		authRouter.GET("/getFootballPlayerStats/:player_public_id", playersServer.GetFootballPlayerStatsFunc)
		authRouter.POST("/createUploadChunks", handlersServer.CreateUploadMediaFunc)
		authRouter.POST("/completedChunkUpload", handlersServer.CompletedChunkUploadFunc)
		//authRouter.PUT("/updateThreadCommentCount/:public_id", handlersServer.UpdateThreadCommentCountFunc)
		authRouter.GET("/getPlayerByTeam/:team_public_id", teamsServer.GetPlayersByTeamFunc)
		authRouter.GET("/getTeamByPlayer/:player_public_id", teamsServer.GetTeamsByPlayerFunc)
		authRouter.POST("/uploadMatchMedia/:match_public_id", server.RequiredPermission(PermUpdateTournament), handlersServer.CreateMatchMediaFunc)
		authRouter.GET("/getMatchMedia/:match_public_id", handlersServer.GetMatchMediaFunc)
		authRouter.PUT("/update-user-location", handlersServer.UpdateUserLocationFunc)
		authRouter.POST("/add-location", handlersServer.AddLocationFunc)
		authRouter.POST("/assign-role", handlersServer.AssignRoleFunc)
		authRouter.DELETE("/revoke-role", handlersServer.RevokeRoleFunc)
		authRouter.GET("/get-my-roles", handlersServer.GetMyRolesFunc)
		authRouter.GET("/get-user-roles/:profile_public_id", handlersServer.GetUserRolesFunc)
		authRouter.GET("/get-resource-roles", handlersServer.GetResourceRolesFunc)
		authRouter.GET("/get-all-roles", handlersServer.GetAllRolesFunc)
	}
	sportRouter := router.Group("/api/:sport").Use(authMiddleware(server.tokenMaker, server.logger))
	//tournament
	sportRouter.GET("/get-tournament-by-location", tournamentServer.GetTournamentByLocationFunc)
	sportRouter.POST("/createTournamentUserRole/:tournament_public_id", tournamentServer.AddTournamentUserRolesFunc)
	sportRouter.POST("/createTournamentMatch/:tournament_public_id", server.RequiredPermission(PermUpdateTournament), tournamentServer.CreateTournamentMatch)
	sportRouter.POST("/createTournament", tournamentServer.AddTournamentFunc)
	//sportRouter.GET("/getTeamsByGroup", tournamentServer.GetTeamxsByGroupFunc)
	//sportRouter.GET("/getTeams/:tournament_id", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getTournamentTeam/:tournament_public_id", tournamentServer.GetTournamentTeamsFunc)
	sportRouter.GET("/getTournamentsBySport/:game_id", tournamentServer.GetTournamentsBySportFunc)
	sportRouter.GET("/getTournament/:tournament_public_id", tournamentServer.GetTournamentFunc)
	sportRouter.GET("/getAllTournamentMatch/:tournament_public_id", tournamentServer.GetTournamentMatch)
	sportRouter.PUT("/updateMatchSubStatus/:match_public_id", server.RequiredPermission(PermUpdateMatch), tournamentServer.UpdateMatchSubStatusFunc)
	sportRouter.GET("/get-matches-by-location", handlersServer.GetMatchesByLocationFunc)
	///
	// sportRouter.POST("/addFootballGoalByPlayer", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.GET("/getFootballStanding/:tournament_public_id", tournamentServer.GetFootballStandingFunc)
	sportRouter.GET("/getCricketStanding/:tournament_public_id", tournamentServer.GetCricketStandingFunc)
	// sportRouter.PUT("/updateFootballStanding", tournamentServer.UpdateFootballStandingFunc)
	// sportRouter.PUT("/updateCricketStanding", tournamentServer.UpdateCricketStandingFunc)
	//sportRouter.PUT("/updateTournamentDate/:tournament_public_id", tournamentServer.UpdateTournamentDateFunc)

	sportRouter.POST("/createTournamentStanding", server.RequiredPermission(PermUpdateTournament), tournamentServer.CreateTournamentStandingFunc)
	// sportRouter.POST("/addTournamentTeam", tournamentServer.AddTournamentTeamFunc)
	sportRouter.GET("/getTournamentByLevel", tournamentServer.GetTournamentByLevelFunc)
	sportRouter.PUT("/updateMatchStatus/:match_public_id", server.RequiredPermission(PermUpdateMatch), tournamentServer.UpdateMatchStatusFunc)
	sportRouter.GET("/getCricketCurrentInning/:match_public_id", cricketServer.GetCricketCurrentInningFunc)
	sportRouter.PUT("/updateMatchResult", tournamentServer.UpdateMatchResultFunc)
	sportRouter.PUT("/updateTournamentStatus/:tournament_public_id", server.RequiredPermission(PermUpdateTournament), tournamentServer.UpdateTournamentStatusFunc)
	sportRouter.GET("/getMatchByMatchID/:match_public_id", handlersServer.GetMatchByMatchIDFunc)
	sportRouter.GET("getTournamentParticipants/:tournament_public_id", tournamentServer.GetTournamentParticipantsFunc)
	sportRouter.POST("/addTournamentParticipants", server.RequiredPermission(PermUpdateTournament), tournamentServer.AddTournamentParticipantsFunc)

	//teams //teams database update completed
	sportRouter.PUT("/update-team-location/:team_public_id", server.RequiredPermission(PermUpdateTeam), teamsServer.UpdateTeamLocationFunc)
	sportRouter.POST("/newTeams", teamsServer.AddTeam)
	//sportRouter.GET("/getTeam/:public_id", teamsServer.GetTeamFunc)
	sportRouter.GET("/getTeams", teamsServer.GetTeamsFunc)
	sportRouter.GET("/searchTeams", teamsServer.SearchTeamFunc)
	sportRouter.POST("/addTeamsMemberFunc", server.RequiredPermission(PermUpdateTeam), teamsServer.AddTeamsMemberFunc)
	sportRouter.PUT("/removePlayerFromTeam", server.RequiredPermission(PermUpdateTeam), teamsServer.RemovePlayerFromTeamFunc)
	sportRouter.GET("/getTeamsMemberFunc/:team_public_id", teamsServer.GetTeamsMemberFunc)
	sportRouter.GET("/getTeamsBySport/:game_id", teamsServer.GetTeamsBySportFunc)
	// sportRouter.GET("/getMatchByTeamFunc", teamsServer.GetMatchByTeamFunc)
	sportRouter.GET("/getMatchesByTeam/:team_public_id", teamsServer.GetMatchesByTeamFunc)
	//sportRouter.GET("/getTournamentByTeamFunc/:team_public_id", teamsServer.GetTournamentbyTeamFunc)

	//matches
	sportRouter.GET("/getAllMatches", handlersServer.GetAllMatchesFunc)

	//football
	// sportRouter.GET("/getFootballScore", footballServer.GetFootballScore)
	sportRouter.POST("/addFootballIncidents", server.RequiredPermission(PermUpdateMatch), footballServer.AddFootballIncidentsFunc)
	sportRouter.GET("/getFootballIncidents/:match_public_id", footballServer.GetFootballIncidentsFunc)
	sportRouter.POST("/addFootballIncidentsSubs", server.RequiredPermission(PermUpdateMatch), footballServer.AddFootballIncidentsSubs)
	// sportRouter.PUT("/updateFootballFirstHalfScore", footballServer.UpdateFootballMatchScoreFirstHalfFunc)
	// sportRouter.PUT("/updateFootballSecondHalfScore", footballServer.UpdateFootballMatchScoreSecondHalfFunc)
	// sportRouter.PUT("/updateFootballMatchScore", footballServer.UpdateFootballMatchScoreFunc)
	//sportRouter.POST("/addFootballMatchScore", footballServer.AddFootballMatchScoreFunc)

	//football->player
	// sportRouter.POST("/addFootballLineUp", footballServer.AddFootballLineUpFunc)
	// sportRouter.POST("/addFootballSubstitution", footballServer.AddFootballSubstitionFunc)
	// sportRouter.GET("/getFootballLineUp", footballServer.GetFootballLineUpFunc)
	sportRouter.GET("/getFootballMatchSquad", footballServer.GetFootballMatchSquadFunc)
	// sportRouter.GET("/getFootballSubstitution", footballServer.GetFootballSubstitutionFunc)
	// sportRouter.PUT("/updateFootballSubsAndLineUp", footballServer.UpdateFootballSubsAndLineUpFunc)

	// sportRouter.POST("/addFootballStatistics", footballServer.AddFootballStatisticsFunc)
	//sportRouter.GET("/getFootballStatistics", footballServer.GetFootballStatisticsFunc)
	sportRouter.POST("/addFootballMatchSquad", server.RequiredPermission(PermUpdateMatch), footballServer.AddFootballSquadFunc)
	sportRouter.GET("/getFootballTournamentPlayerGoal/:tournament_public_id", tournamentServer.GetFootballTournamentPlayersGoalsFunc)
	sportRouter.GET("/getFootballTournamentPlayerYellowCard/:tournament_public_id", tournamentServer.GetFootballTournamentPlayersYellowCardFunc)
	sportRouter.GET("/getFootballTournamentPlayerRedCard/:tournament_public_id", tournamentServer.GetFootballTournamentPlayersRedCardFunc)

	// sportRouter.PUT("/updateFootballStatistics", footballServer.UpdateFootballStatisticsFunc)

	//cricket
	sportRouter.POST("/addCricketScore", server.RequiredPermission(PermUpdateMatch), cricketServer.AddCricketScoreFunc)
	sportRouter.POST("/addCricketToss", server.RequiredPermission(PermUpdateMatch), cricketServer.AddCricketTossFunc)
	sportRouter.GET("/getCricketToss/:match_public_id", cricketServer.GetCricketTossFunc)
	// sportRouter.PUT("/updateCricketInning", cricketServer.UpdateCricketInningsFunc)
	sportRouter.PUT("/updateCricketEndInning", server.RequiredPermission(PermUpdateMatch), cricketServer.UpdateCricketEndInningsFunc)
	sportRouter.PUT("/updateCricketNoBall", server.RequiredPermission(PermUpdateMatch), cricketServer.UpdateNoBallsRunsFunc)
	sportRouter.PUT("/updateCricketWide", server.RequiredPermission(PermUpdateMatch), cricketServer.UpdateWideBallFunc)
	sportRouter.PUT("/updateCricketRegularScore", server.RequiredPermission(PermUpdateMatch), cricketServer.UpdateInningScoreFunc)
	sportRouter.GET("/getCurrentBatsman", cricketServer.GetCurrentBatsmanFunc)
	sportRouter.GET("/getCurrentBowler", cricketServer.GetCurrentBowlerFunc)
	//squad
	sportRouter.POST("/addCricketSquad", server.RequiredPermission(PermUpdateMatch), cricketServer.AddCricketSquadFunc)
	sportRouter.GET("/getCricketMatchSquad", cricketServer.GetCricketMatchSquadFunc)
	//tournament data
	sportRouter.GET("/getCricketTournamentMostRuns/:tournament_public_id", tournamentServer.GetCricketTournamentMostRunsFunc)
	sportRouter.GET("/getCricketTournamentHighestRuns/:tournament_public_id", tournamentServer.GetCricketTournamentHighestRunsFunc)
	sportRouter.GET("getCricketTournamentMostSixes/:tournament_public_id", tournamentServer.GetCricketTournamentMostSixesFunc)
	sportRouter.GET("/getCricketTournamentMostFours/:tournament_public_id", tournamentServer.GetCricketTournamentMostFoursFunc)
	sportRouter.GET("/getCricketTournamentMostFifties/:tournament_public_id", tournamentServer.GetCricketTournamentMostFiftiesFunc)
	sportRouter.GET("/getCricketTournamentMostHundreds/:tournament_public_id", tournamentServer.GetCricketTournamentMostHundredsFunc)
	sportRouter.GET("/getCricketTournamentBowlingStrike/:tournament_public_id", tournamentServer.GetCricketTournamentBowlingStrikeRateFunc)
	sportRouter.GET("/getCricketTournamentBowlingEconomy/:tournament_public_id", tournamentServer.GetCricketTournamentBowlingEconomyRateFunc)
	sportRouter.GET("/getCricketTournamentBowlingAverage/:tournament_public_id", tournamentServer.GetCricketTournamentBowlingAverageFunc)
	sportRouter.GET("/getCricketTournamentMostWickets/:tournament_public_id", tournamentServer.GetCricketTournamentMostWicketsFunc)
	sportRouter.GET("/getCricketTournamentFiveWicketsHaul/:tournament_public_id", tournamentServer.GetCricketTournamentBowlingFiveWicketHaulFunc)
	sportRouter.GET("/getCricketTournamentBattingAverage/:tournament_public_id", tournamentServer.GetCricketTournamentBattingAverageFunc)
	sportRouter.GET("/getCricketTournamentBattingStrike/:tournament_public_id", tournamentServer.GetCricketTournamentBattingStrikeFunc)
	//cricket->player
	sportRouter.POST("/addCricketBatScore", server.RequiredPermission(PermUpdateMatch), cricketServer.AddCricketBatScoreFunc)
	sportRouter.POST("/addCricketBall", server.RequiredPermission(PermUpdateMatch), cricketServer.AddCricketBallFunc)
	sportRouter.GET("/getPlayerScoreFunc", cricketServer.GetPlayerScoreFunc)
	sportRouter.GET("/getCricketBowlerFunc", cricketServer.GetCricketBowlerFunc)
	// sportRouter.PUT("/updateCricketBat", cricketServer.UpdateCricketBatScoreFunc)
	// sportRouter.PUT("/updateCricketBall", cricketServer.UpdateCricketBallFunc)
	sportRouter.GET("/getCricketWickets", cricketServer.GetCricketWicketsFunc)
	sportRouter.POST("/wickets", server.RequiredPermission(PermUpdateMatch), cricketServer.AddCricketWicketsFunc)
	sportRouter.PUT("/updateBowlingBowlerStatus", server.RequiredPermission(PermUpdateMatch), cricketServer.UpdateBowlingBowlerFunc)

	sportRouter.GET("/getLiveMatches", handlersServer.GetLiveMatchesFunc)

	server.router = router
	return server, nil
}

// Middleware for large file upload
func TimeoutMiddleware(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (server *Server) Start(address string) error {
	server.httpServer = &http.Server{
		Addr:    address,
		Handler: server.router,

		ReadTimeout:       0 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      0 * time.Second,
		IdleTimeout:       0 * time.Second,
	}
	return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logger.Info("Shutting down HTTP server...")

	// Set shutting down flag first
	server.SetShuttingDown(true)

	// Gracefully shutdown the HTTP server
	if server.httpServer != nil {
		return server.httpServer.Shutdown(ctx)
	}

	return nil
}

func (server *Server) SetShuttingDown(isShuttingDown bool) {
	if server.isShuttingDown != nil {
		server.isShuttingDown.Store(isShuttingDown)
	}
}

func (server *Server) corsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := server.config.AllowedOrigins
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
