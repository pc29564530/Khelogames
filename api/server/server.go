package server

import (
	"khelogames/api/auth"

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
	config               util.Config
	store                *db.Store
	tokenMaker           token.Maker
	logger               *logger.Logger
	router               *gin.Engine
	webSocketHandlerImpl *messenger.WebSocketHandlerImpl
}

func NewServer(config util.Config,
	store *db.Store,
	tokenMaker token.Maker,
	logger *logger.Logger,
	otpServer *auth.OtpServer,
	signupServer *auth.SignupServer,
	loginServer *auth.LoginServer,
	tokenServer *auth.TokenServer,
	sessionServer *auth.DeleteSessionServer,
	threadServer *handlers.ThreadServer,
	profileServer *handlers.ProfileServer,
	likeThread *handlers.LikethreadServer,
	clubServer *handlers.ClubServer,
	userServer *handlers.UserServer,
	followServer *handlers.FollowServer,
	communityServer *handlers.CommunityServer,
	joinCommunityServer *handlers.JoinCommunityServer,
	commentServer *handlers.CommentServer,
	clubMemberServer *handlers.ClubMemberServer,
	groupTeamServer *handlers.GroupTeamServer,
	playerProfileServer *handlers.PlayerProfileServer,
	tournamentGroupServer *tournaments.TournamentGroupServer,
	tournamentMatchServer *tournaments.TournamentMatchServer,
	tournamentStanding *tournaments.TournamentStandingServer,
	footballMatchServer *football.FootballMatchServer,
	cricketMatchServer *cricket.CricketMatchServer,
	tournamentServer *tournaments.TournamentServer,
	cricketMatchTossServer *cricket.CricketMatchTossServer,
	cricketMatchPlayerScoreServer *cricket.CricketPlayerScoreServer,
	clubTournamentServer *handlers.ClubTournamentServer,
	footballUpdateServer *football.FootballServer,
	webSocketHandlerImpl *messenger.WebSocketHandlerImpl,
	messageServer *messenger.MessageSever,
	router *gin.Engine,
	communityMessageServer *messenger.CommunityMessageServer,
	footballServer *football.FootballServer,
	cricketUpdateServer *cricket.CricketUpdateServer,
) (*Server, error) {

	server := &Server{
		config:               config,
		store:                store,
		tokenMaker:           tokenMaker,
		logger:               logger,
		router:               router,
		webSocketHandlerImpl: webSocketHandlerImpl,
	}

	router.Use(corsHandle())
	router.StaticFS("/api/images", http.Dir("/Users/pawan/database/Khelogames/images"))
	router.StaticFS("/api/videos", http.Dir("/Users/pawan/database/Khelogames/videos"))
	public := router.Group("/auth")
	{
		public.POST("/send_otp", otpServer.Otp)
		public.POST("/signup", signupServer.CreateSignupFunc)
		public.POST("/users", userServer.CreateUserFunc)
		public.POST("/login", loginServer.CreateLoginFunc)
		public.DELETE("/removeSession/:username", sessionServer.DeleteSessionFunc)
		public.POST("/tokens/renew_access", tokenServer.RenewAccessTokenFunc)
		public.GET("/user/:username", userServer.GetUsersFunc)
		public.GET("/getProfile/:owner", profileServer.GetProfileFunc)
	}
	authRouter := router.Group("/api").Use(authMiddleware(server.tokenMaker))
	{
		authRouter.GET("/ws", webSocketHandlerImpl.HandleWebSocket)
		authRouter.POST("/addJoinCommunity", joinCommunityServer.AddJoinCommunityFunc)
		authRouter.GET("/getUserByCommunity/:community_name", joinCommunityServer.GetUserByCommunityFunc)
		authRouter.GET("/getCommunityByUser", joinCommunityServer.GetCommunityByUserFunc)
		authRouter.GET("/user_list", userServer.ListUsersFunc)
		authRouter.POST("/communities", communityServer.CreateCommunitesFunc)
		//authRouter.GET("/communities/:id", server.GetCommunitiesFunc)
		authRouter.GET("/community/:id", communityServer.GetCommunityFunc)
		authRouter.GET("/get_all_communities", communityServer.GetAllCommunitiesFunc)
		authRouter.GET("/getCommunityByCommunityName/:communities_name", communityServer.GetCommunityByCommunityNameFunc)
		authRouter.POST("/create_thread", threadServer.CreateThreadFunc)
		authRouter.GET("/getThread/:id", threadServer.GetThreadFunc)
		authRouter.PUT("/update_like", threadServer.UpdateThreadLikeFunc)
		authRouter.GET("/all_threads", threadServer.GetAllThreadsFunc)
		authRouter.GET("/getAllThreadByCommunity/:communities_name", threadServer.GetAllThreadsByCommunitiesFunc)
		authRouter.GET("/get_communities_member/:communities_name", communityServer.GetCommunitiesMemberFunc)
		authRouter.POST("/create_follow/:following_owner", followServer.CreateFollowingFunc)
		authRouter.GET("/getFollower", followServer.GetAllFollowerFunc)
		authRouter.GET("/getFollowing", followServer.GetAllFollowingFunc)
		authRouter.POST("/createComment/:threadId", commentServer.CreateCommentFunc)
		authRouter.GET("/getComment/:thread_id", commentServer.GetAllCommentFunc)
		authRouter.GET("/getCommentByUser/:username", commentServer.GetCommentByUserFunc)
		authRouter.DELETE("/unFollow/:following_owner", followServer.DeleteFollowingFunc)
		authRouter.POST("/createLikeThread/:thread_id", likeThread.CreateLikeFunc)
		authRouter.GET("/countLike/:thread_id", likeThread.CountLikeFunc)
		authRouter.GET("/checkLikeByUser/:thread_id", likeThread.CheckLikeByUserFunc)
		authRouter.POST("/createProfile", profileServer.CreateProfileFunc)
		authRouter.PUT("/editProfile", profileServer.UpdateProfileFunc)
		authRouter.PUT("/updateAvatar", profileServer.UpdateAvatarUrlFunc)
		authRouter.PUT("/updateCover", profileServer.UpdateCoverUrlFunc)
		authRouter.PUT("/updateFullName", profileServer.UpdateFullNameFunc)
		authRouter.PUT("/updateBio", profileServer.UpdateBioFunc)
		authRouter.GET("getThreadByUser/:username", threadServer.GetThreadByUserFunc)
		authRouter.GET("/getMessage/:receiver_username", messageServer.GetMessageByReceiverFunc)
		authRouter.PUT("/updateAvatarUrl", profileServer.UpdateAvatarUrlFunc)
		authRouter.PUT("/updateClubSport", clubServer.UpdateClubSportFunc)
		authRouter.POST("/addClubMember", clubMemberServer.AddClubMemberFunc)
		authRouter.POST("/createTournament", tournamentServer.CreateTournamentFunc)
		authRouter.GET("/getPlayerProfile", playerProfileServer.GetPlayerProfileFunc)
		authRouter.GET("/getAllPlayerProfile", playerProfileServer.GetAllPlayerProfileFunc)
		authRouter.POST("/addPlayerProfile", playerProfileServer.AddPlayerProfileFunc)
		authRouter.PUT("/updatePlayerProfileAvatarUrl", playerProfileServer.UpdatePlayerProfileAvatarFunc)
		authRouter.POST("/addGroupTeam", groupTeamServer.AddGroupTeamFunc)
		//authRouter.POST("/createTournamentOrganization", tournamentOrganizerServer.CreateTournamentOrganizationFunc)
		authRouter.GET("/getMessagedUser", messageServer.GetUserByMessageSendFunc)
		authRouter.POST("/createUploadMedia", communityMessageServer.CreateUploadMediaFunc)
		authRouter.POST("/createMessageMedia", communityMessageServer.CreateMessageMediaFunc)
		authRouter.POST("/createCommunityMessage", communityMessageServer.CreateCommunityMessageFunc)
		authRouter.GET("/getCommunityMessage", communityMessageServer.GetCommunityByMessageFunc)
		authRouter.GET("/getCommunityByMessage", communityMessageServer.GetCommunityByMessageFunc)
		authRouter.POST("/createOrganizer", tournamentServer.CreateOrganizerFunc)
		authRouter.GET("/getOrganizer", tournamentServer.GetOrganizerFunc)
		authRouter.POST("/createClub", clubServer.CreateClubFunc)
		//whole page routes
		authRouter.GET("/GetAllThreadDetailFunc", threadServer.GetAllThreadDetailFunc)
		authRouter.GET("/GetAllThreadsByCommunityDetailsFunc/:communities_name", threadServer.GetAllThreadsByCommunityDetailsFunc)
	}
	sportRouter := router.Group("/api/:sport").Use(authMiddleware(server.tokenMaker))
	sportRouter.POST("/createTournamentMatch", tournamentMatchServer.CreateTournamentMatch)
	sportRouter.GET("/getTeamsByGroup", groupTeamServer.GetTeamsByGroupFunc)
	sportRouter.GET("/getTeams/:tournament_id", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getTeam/:team_id", tournamentServer.GetTeamFunc)
	sportRouter.GET("/getTournamentsBySport", tournamentServer.GetTournamentsBySportFunc)
	sportRouter.GET("/getTournament/:tournament_id", tournamentServer.GetTournamentFunc)
	sportRouter.POST("/addFootballMatchScore", footballServer.AddFootballMatchScoreFunc)
	//sportRouter.GET("/getFootballMatchScore", footballMatchServer.GetFootballMatchScoreFunc)
	sportRouter.PUT("/updateFootballMatchScore", footballUpdateServer.UpdateFootballMatchScoreFunc)
	//sportRouter.POST("/addFootballGoalByPlayer", footballUpdateServer)
	sportRouter.GET("/getClub/:id", clubServer.GetClubFunc)
	sportRouter.GET("/getClubs", clubServer.GetClubsFunc)
	sportRouter.GET("/getClubMember", clubMemberServer.GetClubMemberFunc)
	sportRouter.GET("/getAllTournamentMatch", tournamentMatchServer.GetTournamentMatch)

	sportRouter.POST("/addCricketMatchScore", cricketMatchServer.AddCricketMatchScoreFunc)
	//sportRouter.GET("/getCricketTournamentMatches", cricketMatchServer.)
	sportRouter.PUT("/updateCricketMatchRunsScore", cricketUpdateServer.UpdateCricketMatchRunsScoreFunc)
	sportRouter.PUT("/updateCricketMatchWicket", cricketUpdateServer.UpdateCricketMatchWicketFunc)
	sportRouter.PUT("/updateCricketMatchExtras", cricketUpdateServer.UpdateCricketMatchExtrasFunc)
	sportRouter.PUT("/updateCricketMatchInnings", cricketUpdateServer.UpdateCricketMatchInningsFunc)
	sportRouter.POST("/addCricketMatchToss", cricketMatchTossServer.AddCricketMatchTossFunc)
	sportRouter.GET("/getCricketMatchToss", cricketMatchTossServer.GetCricketMatchTossFunc)
	sportRouter.POST("/addCricketTeamPlayerScore", cricketMatchPlayerScoreServer.AddCricketPlayerScoreFunc)
	sportRouter.GET("/getCricketTeamPlayerScore", cricketMatchPlayerScoreServer.GetCricketPlayerScoreFunc)
	sportRouter.GET("/getCricketPlayerScore", cricketMatchPlayerScoreServer.GetCricketPlayerScoreFunc)
	sportRouter.PUT("/updateCricketMatchScoreBatting", cricketUpdateServer.UpdateCricketMatchScoreBattingFunc)
	sportRouter.PUT("/updateCricketMatchScoreBowling", cricketUpdateServer.UpdateCricketMatchScoreBowlingFunc)

	sportRouter.GET("/getClubPlayedTournaments", clubTournamentServer.GetClubPlayedTournamentsFunc)
	sportRouter.GET("/getClubPlayedTournament", clubTournamentServer.GetClubPlayedTournamentFunc)
	sportRouter.GET("/getTournamentsByClub", clubServer.GetTournamentsByClubFunc)
	sportRouter.GET("/getMatchByClubName", clubServer.GetMatchByClubNameFunc)
	sportRouter.PUT("/updateTournamentDate", tournamentServer.UpdateTournamentDateFunc)

	sportRouter.POST("/createTournamentStanding", tournamentStanding.CreateTournamentStandingFunc)
	sportRouter.POST("/createTournamentGroup", tournamentGroupServer.CreateTournamentGroupFunc)
	sportRouter.GET("/getTournamentGroup", tournamentGroupServer.GetTournamentGroupFunc)
	sportRouter.GET("/getTournamentGroups", tournamentGroupServer.GetTournamentGroupsFunc)
	sportRouter.GET("/getTournamentStanding", tournamentStanding.GetTournamentStandingFunc)
	sportRouter.GET("/getClubsBySport", clubServer.GetClubsBySportFunc)
	sportRouter.POST("/addTeam", tournamentServer.AddTeamFunc)
	//sportRouter.GET("/getTeams", tournamentServer.GetTeamsFunc)
	sportRouter.GET("/getMatch", tournamentMatchServer.GetTournamentMatch)
	sportRouter.GET("/getTournamentByLevel", tournamentServer.GetTournamentByLevelFunc)
	//Football
	//sportRouter.GET("/getFootballTournamentMatches", footballMatchServer.GetFootballMatchScore())
	//Cricket
	//sportRouter.GET("/GetMatchByClubFunc", clubServer.GetMatchByClubFunc)
	//sportRouter.GET("/getCricketTournamentMatches", cricketMatchServer.GetCricketTournamentMatchesFunc)
	sportRouter.GET("/getTournamentMatches", tournamentMatchServer.GetTournamentMatch)
	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	go server.webSocketHandlerImpl.StartWebSocketHub()
	return server.router.Run(address)
}

// func (err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }

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
