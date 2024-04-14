package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	rabbitConn *ampq.Connection
	rabbitChan *ampq.Channel
	mutex      sync.Mutex
}

func (server *Server) startWebSocketHub() {
	for {
		select {
		case message := <-server.broadcast:
			server.mutex.Lock()
			for client := range server.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					delete(server.clients, client)
					client.Close()
				}
			}
			server.mutex.Unlock()
		}
	}
}

func (server *Server) handleWebSocket(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	auth := strings.Split(authHeader, " ")

	if len(auth) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	_, err := server.tokenMaker.VerifyToken(auth[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := server.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	fmt.Println("Error: ", err)
	if err != nil {
		return
	}

	defer conn.Close()

	server.clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			delete(server.clients, conn)
			break
		}

		var message map[string]string
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Errorf("unable to unmarshal msg ", err)
			return
		}

		err = server.rabbitChan.Publish(
			"",
			"message",
			false,
			false,
			ampq.Publishing{
				ContentType: "application/json",
				Body:        msg,
			},
		)

		authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		b64data := message["media_url"][strings.IndexByte(message["media_url"], ',')+1:]
		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Println("unable to decode :", err)
			return
		}
		mediaType := "image"
		path, err := saveImageToFile(data, mediaType)
		if err != nil {
			fmt.Println("unable to create a file")
			return
		}

		arg := db.CreateNewMessageParams{
			Content:          message["content"],
			IsSeen:           false,
			SenderUsername:   authToken.Username,
			ReceiverUsername: message["receiver_username"],
			MediaUrl:         path,
			MediaType:        message["media_type"],
		}

		_, err = server.store.CreateNewMessage(ctx, arg)
		if err != nil {
			fmt.Errorf("unable to store new message: ", err)
			return
		}

		server.broadcast <- msg
		fmt.Println("Server Message: ", server.broadcast)
	}
}

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

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

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
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		rabbitConn: rabbitConn,
		rabbitChan: rabbitChan,
		mutex:      sync.Mutex{},
	}

	router := gin.Default()

	router.Use(corsHandle())
	router.StaticFS("/images", http.Dir("/Users/pawan/database/Khelogames/images"))
	router.StaticFS("/videos", http.Dir("/Users/pawan/database/Khelogames/videos"))
	router.POST("/send_otp", server.Otp)
	router.POST("/signup", server.createSignup)
	router.POST("/users", server.createUser)
	router.POST("/login", server.createLogin)
	router.DELETE("/removeSession/:username", server.deleteSession)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/user/:username", server.getUsers)
	router.GET("/getProfile/:owner", server.getProfile)
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.GET("/ws", server.handleWebSocket)
	authRouter.POST("/joinUserCommunity/:community_name", server.addJoinCommunity)
	authRouter.GET("/getUserByCommunity/:community_name", server.getUserByCommunity)
	authRouter.GET("/getCommunityByUser", server.getCommunityByUser)
	authRouter.GET("/user_list", server.listUsers)
	authRouter.POST("/communities", server.createCommunites)
	authRouter.GET("/communities/:id", server.getCommunity)
	authRouter.GET("/community/:id", server.getCommunity)
	authRouter.GET("/get_all_communities", server.getAllCommunities)
	authRouter.GET("/getCommunityByCommunityName/:communities_name", server.getCommunityByCommunityName)
	authRouter.POST("/create_thread", server.createThread)
	authRouter.GET("/getThread/:id", server.getThread)
	authRouter.PUT("/update_like", server.updateThreadLike)
	authRouter.GET("/all_threads", server.getAllThreads)
	authRouter.GET("/getAllThreadByCommunity/:communities_name", server.getAllThreadsByCommunities)
	authRouter.GET("/get_communities_member/:communities_name", server.getCommunitiesMember)
	authRouter.POST("/create_follow/:following_owner", server.createFollowing)
	authRouter.GET("/getFollower", server.getAllFollower)
	authRouter.GET("/getFollowing", server.getAllFollowing)
	authRouter.POST("/createComment/:threadId", server.createComment)
	authRouter.GET("/getComment/:thread_id", server.getAllComment)
	authRouter.DELETE("/unFollow/:following_owner", server.deleteFollowing)
	authRouter.POST("/createLikeThread/:thread_id", server.createLike)
	authRouter.GET("/countLike/:thread_id", server.countLike)
	authRouter.GET("/checkLikeByUser/:thread_id", server.checkLikeByUser)
	authRouter.POST("/createProfile", server.createProfile)
	authRouter.PUT("/editProfile", server.updateProfile)
	authRouter.PUT("/updateAvatar", server.updateAvatarUrl)
	authRouter.PUT("/updateCover", server.updateCoverUrl)
	authRouter.PUT("/updateFullName", server.updateFullName)
	authRouter.PUT("/updateBio", server.updateBio)
	authRouter.GET("getThreadByUser/:username", server.getThreadByUser)
	authRouter.GET("/getMessage/:receiver_username", server.getMessageByReceiver)
	authRouter.POST("/createClub", server.createClub)
	authRouter.GET("/getClub/:id", server.getClub)
	authRouter.GET("/getClubs", server.getClubs)
	authRouter.PUT("/updateAvatarUrl", server.updateAvatarUrl)
	authRouter.PUT("/updateClubSport", server.updateClubSport)
	authRouter.POST("/addClubMember", server.addClubMember)
	authRouter.GET("/getClubMember/:club_name", server.getClubMember)
	authRouter.POST("/createTournament", server.createTournament)
	authRouter.GET("/getTournament/:tournament_id", server.getTournament)
	authRouter.GET("/getTournaments", server.getTournaments)
	authRouter.POST("/createOrganizer", server.createOrganizer)
	authRouter.GET("/getOrganizer/:tournament_id", server.getOrganizer)
	authRouter.GET("/getTeam/:team_id", server.getTeam)
	authRouter.POST("/addTeam", server.addTeam)
	authRouter.GET("/getTeams/:tournament_id", server.getTeams)
	authRouter.GET("/getTournamentTeamCount/:tournament_id", server.getTournamentTeamCount)
	authRouter.PUT("/updateTeamsJoined", server.updateTeamsJoined)
	authRouter.POST("/searchTeam", server.searchTeam)
	authRouter.POST("/createTournamentMatch", server.createTournamentMatch)
	authRouter.GET("/getAllTournamentMatch", server.getAllTournamentMatch)
	authRouter.GET("/getClubsBySport/:sport", server.getClubsBySport)
	authRouter.POST("/createTournamentOrganization", server.createTournamentOrganization)
	authRouter.GET("/getTournamentOrganization", server.getTournamentOrganization)
	authRouter.POST("/createTournamentStanding", server.createTournamentStanding)
	authRouter.POST("/createTournamentGroup", server.createTournamentGroup)
	authRouter.GET("/getTournamentGroup", server.getTournamentGroup)
	authRouter.GET("/getTournamentGroups", server.getTournamentGroups)
	authRouter.GET("/getTournamentStanding", server.getTournamentStanding)
	authRouter.POST("/addGroupTeam", server.addGroupTeam)
	authRouter.GET("/getTeamsByGroup", server.getTeamsByGroup)
	authRouter.GET("/getCommentByUser/:owner", server.getCommentByUser)
	authRouter.GET("/getThreadByThreadID/:thread_id", server.getThreadByThreadID)
	authRouter.GET("/getMessagedUser", server.getUserByMessageSend)
	authRouter.POST("/createUploadMedia", server.createUploadMedia)
	authRouter.POST("/createMessageMedia", server.createMessageMedia)
	authRouter.POST("/createCommunityMessage", server.createCommunityMessage)
	authRouter.GET("/getCommunityMessage", server.getCommuntiyMessage)
	authRouter.GET("/getCommunityByMessage", server.getCommunityByMessage)
	authRouter.GET("/getClubPlayedTournaments", server.getClubPlayedTournaments)
	authRouter.GET("/getClubPlayedTournament", server.getClubPlayedTournament)
	authRouter.GET("/getTournamentsByClub", server.getTournamentsByClub)
	authRouter.GET("/getMatchByClubName", server.getMatchByClubName)
	authRouter.PUT("/updateTournamentDate", server.updateTournamentDate)
	authRouter.GET("/getTournamentsBySport", server.getTournamentsBySport)

	authRouter.POST("/createClub", server.createClub)
	authRouter.GET("/getClub/:id", server.getClub)
	authRouter.GET("/getClubs", server.getClubs)
	authRouter.PUT("/updateAvatarUrl", server.updateAvatarUrl)
	authRouter.PUT("/updateClubSport", server.updateClubSport)
	authRouter.POST("/addClubMember", server.addClubMember)
	authRouter.GET("/getClubMember/:club_name", server.getClubMember)
	authRouter.POST("/createTournament", server.createTournament)
	authRouter.GET("/getTournament/:tournament_id", server.getTournament)
	authRouter.GET("/getTournaments", server.getTournaments)
	authRouter.POST("/createOrganizer", server.createOrganizer)
	authRouter.GET("/getOrganizer/:tournament_id", server.getOrganizer)
	authRouter.GET("/getTeam/:team_id", server.getTeam)
	authRouter.POST("/addTeam", server.addTeam)
	authRouter.GET("/getTeams/:tournament_id", server.getTeams)
	authRouter.GET("/getTournamentTeamCount/:tournament_id", server.getTournamentTeamCount)
	authRouter.PUT("/updateTeamsJoined", server.updateTeamsJoined)
	authRouter.POST("/searchTeam", server.searchTeam)
	authRouter.POST("/createTournamentMatch", server.createTournamentMatch)
	authRouter.GET("/getAllTournamentMatch", server.getAllTournamentMatch)
	authRouter.GET("/getClubsBySport/:sport", server.getClubsBySport)
	authRouter.POST("/createTournamentOrganization", server.createTournamentOrganization)
	authRouter.GET("/getTournamentOrganization", server.getTournamentOrganization)
	authRouter.POST("/createTournamentStanding", server.createTournamentStanding)
	authRouter.POST("/createTournamentGroup", server.createTournamentGroup)
	authRouter.GET("/getTournamentGroup", server.getTournamentGroup)
	authRouter.GET("/getTournamentGroups", server.getTournamentGroups)
	authRouter.GET("/getTournamentStanding", server.getTournamentStanding)
	authRouter.POST("/addGroupTeam", server.addGroupTeam)
	authRouter.GET("/getTeamsByGroup", server.getTeamsByGroup)

	server.router = router
	return server, nil
}

// Start run the HTTP server as specific address.
func (server *Server) Start(address string) error {
	go server.startWebSocketHub()
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
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
