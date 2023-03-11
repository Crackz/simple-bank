package api

import (
	"fmt"

	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/crackz/simple-bank/token"
	"github.com/crackz/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     *util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("couldn't create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// User Endpoints
	router.POST("/users/register", server.registerUser)
	router.POST("/users/login", server.loginUser)

	// Accounts Endpoints
	router.GET("/accounts", server.getAccounts)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:accountID", server.getAccount)
	router.PATCH("/accounts/:accountID", server.updateAccount)
	router.DELETE("/accounts/:accountID", server.deleteAccount)

	// Transfer Endpoints
	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}
