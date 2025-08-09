package gapi

import (
	"testing"
	"time"

	db "github.com/Samudra-G/simplebank/db/sqlc"
	"github.com/Samudra-G/simplebank/util"
	"github.com/Samudra-G/simplebank/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}