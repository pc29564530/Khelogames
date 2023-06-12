package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"khelogames/util"
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
	router.POST("/users", server.createUser)
	router.POST("/signup", server.createSignup)
	router.POST("/login", server.createLogin)

	router.POST("/tokens/renew_access", server.renewAccessToken)
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.GET("/user_list", server.listUsers)
	router.POST("/blogs", server.createBlog)
	authRouter.GET("/blogs/:id", server.getBlog)
	authRouter.POST("/communities", server.createCommunites)
	authRouter.GET("/communities/:id", server.getCommunity)
	authRouter.POST("/friend_request", server.getRecieverUsername)
	authRouter.GET("/friend_request", server.ListConnections)
	authRouter.POST("/accept_friend/:id", server.acceptFriend)
	authRouter.GET("/get_all_friends", server.getAllFriends)
	//router.POST("/connections", server)

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
