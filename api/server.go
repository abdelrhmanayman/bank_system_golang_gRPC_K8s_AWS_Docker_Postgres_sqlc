package api

import (
	db "banksystem/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func SetupRoutes(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", CurrencyValidator)
	}

	// Accounts Routers
	router.POST("/accounts", server.CreateAccountController)
	router.GET("/accounts/:ID", server.GetAccountController)
	router.GET("/accounts", server.ListAccountsController)

	// Transfers Routers
	router.POST("/transfers", server.CreateTransferController)

	server.router = router
	return server

}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
