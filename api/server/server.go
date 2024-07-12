package server

import (
	"khelogames/api/auth"
	"khelogames/api/clubs"

	"khelogames/api/handlers"
	"khelogames/api/messenger"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
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
	clubServer *clubs.ClubServer,
	messageServer *messenger.MessageServer,
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
		authRouter.GET("/ws", messageServer.)
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
		authRouter.PUT("/updateCover", handlersServer.UpdateCoverUrlFunc)
		authRouter.PUT("/updateFullName", handlersServer.UpdateFullNameFunc)
		authRouter.PUT("/updateBio", handlersServer.UpdateBioFunc)
		authRouter.GET("getThreadByUser/:username", handlersServer.GetThreadByUserFunc)
		authRouter.GET("/getMessage/:receiver_username", messageServer.GetMessageByReceiverFunc)
		authRouter.PUT("/updateAvatarUrl", handlersServer.UpdateAvatarUrlFunc)
		authRouter.PUT("/updateClubSport", clubServer.UpdateClubSportFunc)
		authRouter.POST("/addClubMember", clubServer.AddClubMemberFunc)
		authRouter.POST("/createTournament", tournamentServer.CreateTournamentFunc)
		authRouter.GET("/getPlayerProfile", handlersServer.GetPlayerProfileFunc)
		authRouter.GET("/getAllPlayerProfile", handlersServer.GetAllPlayerProfileFunc)
		authRouter.POST("/addPlayerProfile", handlersServer.AddPlayerProfileFunc)
		authRouter.PUT("/updatePlayerProfileAvatarUrl", handlersServer.UpdatePlayerProfileAvatarFunc)
		authRouter.POST("/addGroupTeam", tournamentServer.AddGroupTeamFunc)
		//authRouter.POST("/createTournamentOrganization", tournamentOrganizerServer.CreateTournamentOrganizationFunc)
		authRouter.GET("/getMessagedUser", messageServer.GetUserByMessageSendFunc)
		authRouter.POST("/createUploadMedia", messageServer.CreateUploadMediaFunc)
		authRouter.POST("/createMessageMedia", messageServer.CreateMessageMediaFunc)
		authRouter.POST("/createCommunityMessage", messageServer.CreateCommunityMessageFunc)
		authRouter.GET("/getCommunityMessage", messageServer.GetCommunityByMessageFunc)
		authRouter.GET("/getCommunityByMessage", messageServer.GetCommunityByMessageFunc)
		authRouter.POST("/createOrganizer", tournamentServer.CreateOrganizerFunc)
		authRouter.GET("/getOrganizer", tournamentServer.GetOrganizerFunc)
		authRouter.POST("/createClub", clubServer.CreateClubFunc)
		authRouter.GET("/GetAllThreadDetailFunc", handlersServer.GetAllThreadDetailFunc)
		authRouter.GET("/GetAllThreadsByCommunityDetailsFunc/:communities_name", handlersServer.GetAllThreadsByCommunityDetailsFunc)
	}
	sportRouter := router.Group("/api/:sport").Use(authMiddleware(server.tokenMaker))
	sportRouter.POST("/createTournamentMatch", tournamentServer.CreateTournamentMatch)
	sportRouter.GET("/getTeamsByGroup", groupTeamServer.GetTeamsByGroupFunc)
	sportRouter.GET("/getTeams/:tournament_id", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getTeam/:team_id", tournamentServer.GetTeamFunc)
	sportRouter.GET("/getTournamentsBySport", tournamentServer.GetTournamentsBySportFunc)
	sportRouter.GET("/getTournament/:tournament_id", tournamentServer.GetTournamentFunc)
	sportRouter.POST("/addFootballMatchScore", footballServer.AddFootballMatchScoreFunc)
	//sportRouter.GET("/getFootballMatchScore", FootballServer.GetFootballMatchScoreFunc)
	sportRouter.PUT("/updateFootballMatchScore", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.POST("/addFootballGoalByPlayer", footballServer.UpdateFootballMatchScoreFunc)
	sportRouter.GET("/getClub/:id", clubServer.GetClubFunc)
	sportRouter.GET("/getClubs", clubServer.GetClubsFunc)
	sportRouter.GET("/getClubMember", clubServer.GetClubMemberFunc)
	sportRouter.GET("/getAllTournamentMatch", tournamentServer.GetTournamentMatch)

	sportRouter.POST("/addCricketMatchScore", cricketServer.AddCricketMatchScoreFunc)
	//sportRouter.GET("/getCricketTournamentMatches", CricketServer.)
	sportRouter.PUT("/updateCricketMatchRunsScore", cricketServer.UpdateCricketMatchRunsScoreFunc)
	sportRouter.PUT("/updateCricketMatchWicket", cricketServer.UpdateCricketMatchWicketFunc)
	sportRouter.PUT("/updateCricketMatchExtras", cricketServer.UpdateCricketMatchExtrasFunc)
	sportRouter.PUT("/updateCricketMatchInnings", cricketServer.UpdateCricketMatchInningsFunc)
	sportRouter.POST("/addCricketMatchToss", cricketServer.AddCricketMatchTossFunc)
	sportRouter.GET("/getCricketMatchToss", cricketServer.GetCricketMatchTossFunc)
	sportRouter.POST("/addCricketTeamPlayerScore", cricketServer.AddCricketPlayerScoreFunc)
	sportRouter.GET("/getCricketTeamPlayerScore", cricketServer.GetCricketPlayerScoreFunc)
	sportRouter.GET("/getCricketPlayerScore", cricketServer.GetCricketPlayerScoreFunc)
	sportRouter.PUT("/updateCricketMatchScoreBatting", cricketServer.UpdateCricketMatchScoreBattingFunc)
	sportRouter.PUT("/updateCricketMatchScoreBowling", cricketServer.UpdateCricketMatchScoreBowlingFunc)

	//sportRouter.GET("/getClubPlayedTournaments", clubServer.GetClubPlayedTournamentsFunc)
	//sportRouter.GET("/getClubPlayedTournament", clubServer.GetClubPlayedTournamentFunc)
	sportRouter.GET("/getTournamentsByClub", clubServer.GetTournamentsByClubFunc)
	sportRouter.GET("/getMatchByClubName", clubServer.GetMatchByClubNameFunc)
	sportRouter.PUT("/updateTournamentDate", tournamentServer.UpdateTournamentDateFunc)

	sportRouter.POST("/createTournamentStanding", tournamentServer.CreateTournamentStandingFunc)
	sportRouter.POST("/createTournamentGroup", tournamentServer.CreateTournamentGroupFunc)
	sportRouter.GET("/getTournamentGroup", tournamentServer.GetTournamentGroupFunc)
	sportRouter.GET("/getTournamentGroups", tournamentServer.GetTournamentGroupsFunc)
	sportRouter.GET("/getTournamentStanding", tournamentServer.GetTournamentStandingFunc)
	sportRouter.GET("/getClubsBySport", clubServer.GetClubsBySportFunc)
	sportRouter.POST("/addTeam", tournamentServer.AddTeamFunc)
	sportRouter.PUT("/updateTeamJoinedTournament", tournamentServer.UpdateTeamsJoinedFunc)
	sportRouter.GET("/getTeams", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getMatch", tournamentServer.GetTournamentMatch)
	sportRouter.GET("/getTournamentByLevel", tournamentServer.GetTournamentByLevelFunc)
	//sportRouter.GET("/getFootballTournamentMatches", FootballServer.GetFootballMatchScore())
	//sportRouter.GET("/GetMatchByClubFunc", clubServer.GetMatchByClubFunc)
	//sportRouter.GET("/getCricketTournamentMatches", CricketServer.GetCricketTournamentMatchesFunc)
	sportRouter.GET("/getTournamentMatches", tournamentServer.GetTournamentMatch)
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
