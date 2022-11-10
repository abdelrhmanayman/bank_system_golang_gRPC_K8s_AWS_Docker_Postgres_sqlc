package api

import (
	db "banksystem/db/sqlc"
	"banksystem/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, store db.Store) (*Server, string) {
	symmetricKey := util.GenerateRandomString(32)
	config := util.Config{
		SymmetricKey:  symmetricKey,
		TokenDuration: time.Minute,
	}

	server, err := SetupRoutes(config, store)
	assert.NoError(t, err)

	return server, symmetricKey
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
