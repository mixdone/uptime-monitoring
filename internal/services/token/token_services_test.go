package token_test

import (
	"testing"

	"github.com/mixdone/uptime-monitoring/internal/services/constants"
	"github.com/mixdone/uptime-monitoring/internal/services/token"
	"github.com/stretchr/testify/assert"
)

const userID = int64(123456789)

func TestGenerate_NoErr(t *testing.T) {
	srv := token.NewTokenService("my-very-secret-access-key",
		"my-very-secret-refresh-key", constants.AccessTokenTTL, constants.RefreshTokenTTL)

	_, _, err := srv.Generate(userID)

	assert.NoError(t, err)
}

func TestGenerateValidate_Access(t *testing.T) {
	srv := token.NewTokenService("my-very-secret-access-key",
		"my-very-secret-refresh-key", constants.AccessTokenTTL, constants.RefreshTokenTTL)

	aT, _, err := srv.Generate(userID)

	assert.NoError(t, err)

	id, err := srv.ValidateAccess(aT)

	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}

func TestGenerateValidate_Refresh(t *testing.T) {
	srv := token.NewTokenService("my-very-secret-access-key",
		"my-very-secret-refresh-key", constants.AccessTokenTTL, constants.RefreshTokenTTL)

	_, rT, err := srv.Generate(userID)

	assert.NoError(t, err)

	id, err := srv.ValidateRefresh(rT)

	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}
