package grpc_api

import (
	"fmt"

	db "github.com/radugaf/simplebank/db/sqlc"
	"github.com/radugaf/simplebank/pb"
	"github.com/radugaf/simplebank/token"
	"github.com/radugaf/simplebank/tools"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	store          db.Store
	tokenGenerator token.Token
	config         tools.Config
}

// NewServer creates a new gRPC server.
func NewServer(config tools.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token generator: %w", err)
	}

	server := &Server{
		config:         config,
		store:          store,
		tokenGenerator: tokenGenerator,
	}

	return server, nil
}
