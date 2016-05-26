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
		Login:    "user1@eogile.com",
		UserData: secu.UserData{
			FirstName: "First name 1",
			LastName:  "Last name 1",
		},
	}, tokenInfo)

	require.Nil(t, err)
	require.NotNil(t, id)
	require.Regexp(t, "^[0-9a-f\\-]+", id)

	user, err := client.FindUser(id, tokenInfo)

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

func TestUpdateUser(t *testing.T) {
	client := auth.NewClient("http://localhost:9090", "superapp2", "supersecret2")
	token, err := client.Login("superadmin@eogile.com", "supersecret")

	tokenInfo, err := auth.EncodeTokenInfo(token)
	require.Nil(t, err)

	id, err := client.CreateUser(&secu.User{
		Password: "1234",
		Login:    "user2@eogile.com",
		UserData: secu.UserData{
			FirstName: "First name 2",
			LastName:  "Last name 2",
		},
	}, tokenInfo)

	require.Nil(t, err)
	require.NotNil(t, id)

	/*
		Updating the client.
	*/
	newUserData := secu.UserData{
		FirstName: "First name 2 updated",
		LastName:  "Last name 2 updated",
		Active:    false,
		Blocked:   true,
	}
	require.Nil(t, client.UpdateUserData(id, newUserData, tokenInfo))

	/*
		Checking the client
	*/
	user, err := client.FindUser(id, tokenInfo)

	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, id, user.Id)
	require.Equal(t, "user2@eogile.com", user.Login)
	require.Equal(t, "", user.Password)
	require.Equal(t, newUserData.FirstName, user.FirstName)
	require.Equal(t, newUserData.LastName, user.LastName)
	require.Equal(t, newUserData.Active, user.IsActive())
	require.Equal(t, newUserData.Blocked, user.IsBlocked())
}

// Tests that the default policy allows an user to access its own data
// but not other users data.
func TestDefaultPolicy(t *testing.T) {
	client := auth.NewClient("http://localhost:9090", "superapp2", "supersecret2")
	token, err := client.Login("superadmin@eogile.com", "supersecret")
	require.Nil(t, err)
	require.NotNil(t, token)

	tokenInfo, err := auth.EncodeTokenInfo(token)
	require.Nil(t, err)

	_, err = client.CreateDefaultPolicy(tokenInfo)
	require.Nil(t, err)

	// Creating user3
	id3, err := client.CreateUser(&secu.User{
		Password: "1234",
		Login:    "user3@eogile.com",
		UserData: secu.UserData{
			FirstName: "First name 3",
			LastName:  "Last name 3",
		},
	}, tokenInfo)
	require.Nil(t, err)
	require.NotNil(t, id3)

	// Creating user4
	id4, err := client.CreateUser(&secu.User{
		Password: "1234",
		Login:    "user4@eogile.com",
		UserData: secu.UserData{
			FirstName: "First name 4",
			LastName:  "Last name 4",
		},
	}, tokenInfo)
	require.Nil(t, err)
	require.NotNil(t, id4)

	// Authenticate with user3
	user3Token, err := client.Login("user3@eogile.com", "1234")
	require.Nil(t, err)
	tokenInfo3, err := auth.EncodeTokenInfo(user3Token)
	require.Nil(t, err)

	// user3 can access its data
	user3, err := client.FindUser(id3, tokenInfo3)
	require.Nil(t, err)
	require.Equal(t, "user3@eogile.com", user3.Login)

	// user3 cannot access user4 data
	user4, err := client.FindUser(id4, tokenInfo3)
	require.NotNil(t, err)
	require.Nil(t, user4)

}
