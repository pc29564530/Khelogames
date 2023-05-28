package api

import (
	db "khelogames/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	router := gin.Default()
	router.POST("/blogs", server.createBlog)
	router.GET("/blogs/:id", server.getBlog)
	router.POST("/users", server.createUser)
	router.POST("/signup", server.createSignup)
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
