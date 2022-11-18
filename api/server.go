package api

import (
	db "banksystem/db/sqlc"
	"banksystem/token"
	"banksystem/util"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	appConfig  util.Config
}

func SetupRoutes(config util.Config, store db.Store) (*Server, error) {
	server := &Server{store: store}
	router := gin.Default()

	pMaker, err := token.CreateNewPasetoMaker([]byte(config.SymmetricKey))

	if err != nil {
		err = fmt.Errorf("can't start token maker %v", err)
		return nil, err
	}

	server.tokenMaker = pMaker
	server.appConfig = config

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", CurrencyValidator)
	}
	// Auth Routers
	router.POST("/auth/login", server.LoginController)

	// Users Routers
	router.POST("/users", server.CreateUserController)

	router.Use(AuthMiddleware(pMaker))

	// Accounts Routers
	router.POST("/accounts", server.CreateAccountController)
	router.GET("/accounts/:ID", server.GetAccountController)
	router.GET("/accounts", server.ListAccountsController)

	// Transfers Routers
	router.POST("/transfers", server.CreateTransferController)

	server.router = router
	return server, nil

}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
