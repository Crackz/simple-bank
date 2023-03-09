package api

import (
	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	// User Endpoints
	router.POST("/users", server.createUser)

	// Accounts Endpoints
	router.GET("/accounts", server.getAccounts)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:accountID", server.getAccount)
	router.PATCH("/accounts/:accountID", server.updateAccount)
	router.DELETE("/accounts/:accountID", server.deleteAccount)

	// Transfer Endpoints
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}
