package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/radugaf/simplebank/db/sqlc"
	"github.com/radugaf/simplebank/token"
	"github.com/radugaf/simplebank/tools"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store          db.Store
	router         *gin.Engine
	tokenGenerator token.Token
	config         tools.Config
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config tools.Config, store db.Store) (*Server, error) {

	tokenGenerator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token generator: %w", err)
	}

	server := &Server{config: config, store: store, tokenGenerator: tokenGenerator}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenGenerator))

	authRoutes.POST("/bank_accounts", server.createBankAccount)
	authRoutes.GET("/bank_accounts/:id", server.getBankAccount)
	authRoutes.GET("/bank_accounts", server.listBankAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
