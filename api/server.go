package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()

	router.Use(corsHandle())
	router.StaticFS("/images", http.Dir("/Users/pawan/project/Khelogames/images"))

	router.POST("/send_otp", server.Otp)
	router.POST("/signup", server.createSignup)
	router.POST("/users", server.createUser)
	router.POST("/login", server.createLogin)
	router.DELETE("/removeSession/:username", server.deleteSession)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/user/:username", server.getUsers)
	router.GET("/getProfile/:owner", server.getProfile)
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.POST("/joinUserCommunity/:community_name", server.addJoinCommunity)
	authRouter.GET("/getUserByCommunity/:community_name", server.getUserByCommunity)
	authRouter.GET("/getCommunityByUser", server.getCommunityByUser)
	authRouter.GET("/user_list", server.listUsers)
	authRouter.POST("/communities", server.createCommunites)
	authRouter.GET("/communities/:id", server.getCommunity)
	authRouter.GET("/community/:id", server.getCommunity)
	authRouter.GET("/get_all_communities/:owner", server.getAllCommunities)
	authRouter.POST("/create_thread", server.createThread)
	authRouter.GET("/getThread/:id", server.getThread)
	authRouter.PUT("/update_like", server.updateThreadLike)
	authRouter.GET("/all_threads", server.getAllThreads)
	authRouter.GET("/get_all_communities_by_owner", server.getAllThreadsByCommunities)
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

	server.router = router
	return server, nil
}

// Start run the HTTP server as specific address.
func (server *Server) Start(address string) error {
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
