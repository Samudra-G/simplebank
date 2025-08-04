package gapi

import (
	"fmt"

	db "github.com/Samudra-G/simplebank/db/sqlc"
	"github.com/Samudra-G/simplebank/pb"
	"github.com/Samudra-G/simplebank/token"
	"github.com/Samudra-G/simplebank/util"
	"github.com/Samudra-G/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

//NewServer creates a new gRPC Server and setup routing
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}