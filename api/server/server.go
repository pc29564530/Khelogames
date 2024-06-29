package server

import (
	"fmt"
	"khelogames/api/auth"
	"khelogames/api/cricket"
	"khelogames/api/football"
	"khelogames/api/handlers"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	util "khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

// type ServerInterface interface {
// 	Start(address string) error
// 	HandleWebSocket(ctx *gin.Context)
// 	GetLogger() *logger.Logger
// 	GetStore() *db.Store
// }

// type Server struct {
// 	config     util.Config
// 	store      *db.Store
// 	tokenMaker token.Maker
// 	router     *gin.Engine
// 	upgrader   websocket.Upgrader
// 	clients    map[*websocket.Conn]bool
// 	broadcast  chan []byte
// 	rabbitConn *ampq.Connection
// 	rabbitChan *ampq.Channel
// 	mutex      sync.Mutex
// 	logger     *logger.Logger
// }

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	logger     *logger.Logger
	upgrader   websocket.Upgrader
	router     *gin.Engine
	rabbitConn *ampq.Connection
	rabbitChan *ampq.Channel
}

// func (server *Server) GetLogger() *logger.Logger {
// 	return server.logger
// }

// func (server *Server) GetStore() *db.Store {
// 	return server.store
// }

// func (server *Server) TokenMaker() token.Maker {
// 	return server.tokenMaker
// }

// func (server *Server) Users() *Users {
// 	return server.Usersse
// }

// func startWebSocketHub() {
// 	for {
// 		select {
// 		case message := <-server.broadcast:
// 			server.mutex.Lock()
// 			for client := range server.clients {
// 				err := client.WriteMessage(websocket.TextMessage, message)
// 				if err != nil {
// 					delete(server.clients, client)
// 					client.Close()
// 				}
// 			}
// 			server.mutex.Unlock()
// 		}
// 	}
// }

// func (server *Server) HandleWebSocket(ctx *gin.Context) {
// 	authHeader := ctx.GetHeader("Authorization")
// 	auth := strings.Split(authHeader, " ")

// 	if len(auth) == 0 {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
// 		return
// 	}

// 	_, err := server.tokenMaker.VerifyToken(auth[1])
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	conn, err := server.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		return
// 	}

// 	defer conn.Close()

// 	server.clients[conn] = true

// 	for {
// 		_, msg, err := conn.ReadMessage()

// 		if err != nil {
// 			delete(server.clients, conn)
// 			break
// 		}

// 		var message map[string]string
// 		err = json.Unmarshal(msg, &message)
// 		if err != nil {
// 			fmt.Print("unable to unmarshal msg ", err)
// 			return
// 		}

// 		err = server.rabbitChan.Publish(
// 			"",
// 			"message",
// 			false,
// 			false,
// 			ampq.Publishing{
// 				ContentType: "application/json",
// 				Body:        msg,
// 			},
// 		)

// 		authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 		b64data := message["media_url"][strings.IndexByte(message["media_url"], ',')+1:]
// 		data, err := base64.StdEncoding.DecodeString(b64data)
// 		if err != nil {
// 			fmt.Print("unable to decode :", err)
// 			return
// 		}
// 		mediaType := "image"
// 		path, err := util.SaveImageToFile(data, mediaType)
// 		if err != nil {
// 			fmt.Print("unable to create a file")
// 			return
// 		}

// 		arg := db.CreateNewMessageParams{
// 			Content:          message["content"],
// 			IsSeen:           false,
// 			SenderUsername:   authToken.Username,
// 			ReceiverUsername: message["receiver_username"],
// 			MediaUrl:         path,
// 			MediaType:        message["media_type"],
// 		}

// 		_, err = server.store.CreateNewMessage(ctx, arg)
// 		if err != nil {
// 			fmt.Print("unable to store new message: ", err)
// 			return
// 		}

// 		server.broadcast <- msg
// 	}
// }

func startRabbitMQ(config util.Config) (*ampq.Connection, *ampq.Channel, error) {
	rabbitConn, err := ampq.Dial(config.RabbitSource)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to RabbitMQ :%w", err)
	}
	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open RabbitMQ channel :%w", err)
	}
	return rabbitConn, rabbitChan, nil
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
	tournamentGroupServer *handlers.TournamentGroupServer,
	tournamentOrganizerServer *handlers.TournamentOrganizerServer,
	tournamentMatchServer *handlers.TournamentMatchServer,
	tournamentStanding *handlers.TournamentStandingServer,
	footballMatchServer *football.FootballMatches,
	cricketMatchServer *cricket.CricketMatchServer,
	tournamentServer *handlers.TournamentServer,
	cricketMatchTossServer *cricket.CricketMatchTossServer,
	cricketMatchPlayerScoreServer *cricket.CricketMatchScoreServer,
	clubTournamentServer *handlers.ClubTournamentServer,
) (*Server, error) {
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create token maker: %w", err)
	// }

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	rabbitConn, rabbitChan, err := startRabbitMQ(config)
	if err != nil {
		return nil, fmt.Errorf("cannot run the rabbit mq :%w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		upgrader:   upgrader,
		rabbitConn: rabbitConn,
		rabbitChan: rabbitChan,
		logger:     logger,
	}

	router := gin.Default()

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
		//authRouter.GET("/ws", wsHandler.HandleWebSocket)
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
		//authRouter.GET("/getMessage/:receiver_username", threadServer.GetMessageByReceiverFunc)
		authRouter.PUT("/updateAvatarUrl", profileServer.UpdateAvatarUrlFunc)
		authRouter.PUT("/updateClubSport", clubServer.UpdateClubSportFunc)
		authRouter.POST("/addClubMember", clubMemberServer.AddClubMemberFunc)
		authRouter.POST("/createTournament", tournamentServer.CreateTournamentFunc)
		authRouter.GET("/getPlayerProfile", playerProfileServer.GetPlayerProfileFunc)
		authRouter.GET("/getAllPlayerProfile", playerProfileServer.GetAllPlayerProfileFunc)
		authRouter.POST("/addPlayerProfile", playerProfileServer.AddPlayerProfileFunc)
		authRouter.PUT("/updatePlayerProfileAvatarUrl", playerProfileServer.UpdatePlayerProfileAvatarFunc)
		authRouter.POST("/addGroupTeam", groupTeamServer.AddGroupTeamFunc)
		authRouter.POST("/createTournamentOrganization", tournamentOrganizerServer.CreateTournamentOrganizationFunc)
		//authRouter.GET("/getMessagedUser", server.GetUserByMessageSendFunc)
		// authRouter.POST("/createUploadMedia", server.CreateUploadMediaFunc)
		// authRouter.POST("/createMessageMedia", server.CreateMessageMediaFunc)
		// authRouter.POST("/createCommunityMessage", server.CreateCommunityMessageFunc)
		// authRouter.GET("/getCommunityMessage", server.GetCommunityMessageFunc)
		// authRouter.GET("/getCommunityByMessage", server.GetCommunityByMessageFunc)
		authRouter.POST("/createOrganizer", tournamentServer.CreateOrganizerFunc)
		authRouter.GET("/getOrganizer", tournamentServer.GetOrganizerFunc)
		authRouter.POST("/createClub", clubServer.CreateClubFunc)

	}
	sportRouter := router.Group("/api/:sport").Use(authMiddleware(server.tokenMaker))
	sportRouter.POST("/createTournamentMatch", tournamentMatchServer.CreateTournamentMatchFunc)
	sportRouter.GET("/getTeamsByGroup", groupTeamServer.GetTeamsByGroupFunc)
	sportRouter.GET("/getTeams/:tournament_id", tournamentServer.GetTeamFunc)
	sportRouter.GET("/getTournamentsBySport", tournamentServer.GetTournamentsBySportFunc)
	sportRouter.GET("/getTournament/:tournament_id", tournamentServer.GetTournamentFunc)
	//sportRouter.POST("/addFootballMatchScore", footballMatchServer.AddFootballMatchScoreFunc)
	//sportRouter.GET("/getFootballMatchScore", footballMatchServer.GetFootballMatchScoreFunc)
	//sportRouter.PUT("/updateFootballMatchScore", football.UpdateFootballMatchScoreFunc)
	//sportRouter.POST("/addFootballGoalByPlayer", server.AddFootballGoalByPlayerFunc)
	sportRouter.GET("/getClub/:id", clubServer.GetClubFunc)
	sportRouter.GET("/getClubs", clubServer.GetClubsFunc)
	sportRouter.GET("/getClubMember", clubMemberServer.GetClubMemberFunc)
	sportRouter.GET("/getAllTournamentMatch", tournamentMatchServer.GetAllTournamentMatchFunc)

	sportRouter.POST("/addCricketMatchScore", cricketMatchPlayerScoreServer.AddCricketMatchScoreFunc)
	//sportRouter.GET("/getCricketMatchScore", cricketMatchServer.GetCricketTournamentMatchesFunc)
	sportRouter.PUT("/updateCricketMatchRunsScore", cricketMatchPlayerScoreServer.UpdateCricketMatchRunsScoreFunc)
	sportRouter.PUT("/updateCricketMatchWicket", cricketMatchPlayerScoreServer.UpdateCricketMatchWicketFunc)
	sportRouter.PUT("/updateCricketMatchExtras", cricketMatchPlayerScoreServer.UpdateCricketMatchExtrasFunc)
	sportRouter.PUT("/updateCricketMatchInnings", cricketMatchPlayerScoreServer.UpdateCricketMatchInningsFunc)
	sportRouter.POST("/addCricketMatchToss", cricketMatchTossServer.AddCricketMatchTossFunc)
	sportRouter.GET("/getCricketMatchToss", cricketMatchTossServer.GetCricketMatchTossFunc)
	sportRouter.POST("/addCricketTeamPlayerScore", cricketMatchPlayerScoreServer.AddCricketPlayerScoreFunc)
	sportRouter.GET("/getCricketTeamPlayerScore", cricketMatchPlayerScoreServer.GetCricketPlayerScoreFunc)
	sportRouter.GET("/getCricketPlayerScore", cricketMatchPlayerScoreServer.GetCricketPlayerScoreFunc)
	sportRouter.PUT("/updateCricketMatchScoreBatting", cricketMatchPlayerScoreServer.UpdateCricketMatchScoreBattingFunc)
	sportRouter.PUT("/updateCricketMatchScoreBowling", cricketMatchPlayerScoreServer.UpdateCricketMatchScoreBowlingFunc)

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
	sportRouter.GET("/getMatch", tournamentMatchServer.GetMatchFunc)
	sportRouter.GET("/getTournamentByLevel", tournamentServer.GetTournamentByLevelFunc)
	// fm := football.NewFootballMatches(server)
	sportRouter.GET("/getFootballTournamentMatches", footballMatchServer.GetFootballTournamentMatchesFunc)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	// go server.wsHandler.startWebSocketHub()
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
