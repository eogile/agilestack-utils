package auth_test

import (
	"log"
	"testing"

	"github.com/eogile/agilestack-utils/auth"
	"github.com/eogile/agilestack-utils/secu"
	"github.com/ory-am/osin-storage/Godeps/_workspace/src/github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	client := auth.NewClient("http://localhost:9090", "superapp2", "supersecret2")
	token, err := client.Login("superadmin@eogile.com", "supersecret")
	require.Nil(t, err)
	require.NotNil(t, token)

	log.Println("TOKEN:", token)
	require.Equal(t, "Bearer", token.TokenType)
	require.NotEqual(t, "^\\s*$", token.AccessToken)
	require.NotEqual(t, "^\\s*$", token.RefreshToken)
}

func TestCreateUser(t *testing.T) {
	client := auth.NewClient("http://localhost:9090", "superapp2", "supersecret2")
	token, err := client.Login("superadmin@eogile.com", "supersecret")
	require.Nil(t, err)
	require.NotNil(t, token)

	tokenInfo, err := auth.EncodeTokenInfo(token)
	require.Nil(t, err)

	id, err := client.CreateUser(&secu.User{
		Password: "1234",
		Login:       "user1@eogile.com",
		UserData: secu.UserData{
			FirstName: "First name 1",
			LastName:  "Last name 1",
		},
	}, tokenInfo)

	require.Nil(t, err)
	require.NotNil(t, id)
	require.Regexp(t, "^[0-9a-f\\-]+", id)

	user, err:= client.FindUser(id, tokenInfo)

	log.Println("Romain2", user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, id, user.Id)
	require.Equal(t, "user1@eogile.com", user.Login)
	require.Equal(t, "", user.Password)
	require.Equal(t, "First name 1", user.FirstName)
	require.Equal(t, "Last name 1", user.LastName)
	require.Equal(t, true, user.IsActive())
	require.Equal(t, false, user.IsBlocked())
}
