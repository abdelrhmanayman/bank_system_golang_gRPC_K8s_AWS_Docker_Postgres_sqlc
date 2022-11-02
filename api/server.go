package api

import (
	db "banksystem/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func SetupRoutes(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Accounts Routers
	router.POST("/accounts", server.CreateAccountController)
	router.GET("/accounts/:ID", server.GetAccountController)
	router.GET("/accounts", server.ListAccountsController)

	server.router = router
	return server

}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
