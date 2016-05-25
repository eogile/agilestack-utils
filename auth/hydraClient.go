package auth

import (
	"net/http"

	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"bytes"

	"github.com/dgrijalva/jwt-go"
	"github.com/eogile/agilestack-utils/secu"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	policyPath  = "/policies"
	accountPath = "/accounts"
	oauthPath   = "/oauth2"
)

const dummyKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----
`

type HydraClient struct {
	clientCredentialConfig clientcredentials.Config
	oauth2Config           oauth2.Config
	authorizationServer    string
}

func NewClient(authorizationServer, clientID, clientSecret string) *HydraClient {
	clientCredentialConfig := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authorizationServer + "/oauth2/token",
	}
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,

		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizationServer + "/oauth2/auth",
			TokenURL: authorizationServer + "/oauth2/token",
		},

		// RedirectURL is the URL to redirect users going through
		// the OAuth flow, after the resource owner's URLs.
		RedirectURL: "http://localhost:8080/login", //TODO verify

	}
	return &HydraClient{
		clientCredentialConfig: clientCredentialConfig,
		oauth2Config:           oauth2Config,
		authorizationServer:    authorizationServer,
		//clientCredentialConfig.Client(oauth2.NoContext),
	}
}

func (client HydraClient) getHttpClient(tokenInfo *TokenInfo) *http.Client {
	if tokenInfo == nil || tokenInfo.TokenInfo == "" || tokenInfo.TokenInfo == "null" {
		log.Println(" in getHttpClient, tokenInfo == nil || tokenInfo.TokenInfo==\"\" || tokenInfo.TokenInfo == \"null\"")
		return http.DefaultClient

	}

	token, err := DecodeTokenInfo(tokenInfo)

	if err != nil {
		if tokenInfo == nil {
			log.Println("in getHttpClient tokenInfo == nil")
		} else {
			log.Println("Error when decoding : ", tokenInfo.TokenInfo)
		}
		return nil
	}
	if token == nil {
		log.Println(" in getHttpClient, token nil")
		return http.DefaultClient
	}
	oauth2Token := &oauth2.Token{}
	oauth2Token.AccessToken = token.AccessToken
	oauth2Token.TokenType = token.TokenType
	oauth2Token.Expiry = token.Expiry
	oauth2Token.RefreshToken = token.RefreshToken

	return client.oauth2Config.Client(oauth2.NoContext, oauth2Token)
}

func getUserId(tokenInfo *TokenInfo) (string, error) {
	token, err := DecodeTokenInfo(tokenInfo)
	if err != nil {
		log.Printf("error in getUserId>DecodeTokenInfo : %v", err)
		return "", err
	}
	log.Printf("token.AccessToken=%v", token.AccessToken)

	accessToken, err := jwt.Parse(token.AccessToken, func(*jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(dummyKey))
	})
	if err != nil {
		log.Printf("error in getUserId>Unmarshall AccessToken : %v", err)
		return "", err
	}
	log.Printf("accessToken=%v", accessToken)

	return fmt.Sprintf("%v", accessToken.Claims["sub"]), nil
}

func (client HydraClient) GetUser(tokenInfo *TokenInfo) (*secu.User, error, int) {
	userId, err := getUserId(tokenInfo)
	if err != nil {
		return nil, err, 0
	}

	httpClient := client.getHttpClient(tokenInfo)

	account := &account.DefaultAccount{}
	found, err, respCode := client.findElement(&account, accountPath+"/"+userId, httpClient)
	if !found || err != nil {
		log.Printf("in hydraClient.GetUser, err:%v, respCode:%v\n", err, respCode)
		return nil, err, respCode
	}
	user := secu.NewUser(account)
	return user, nil, respCode
}

func (client HydraClient) ListUsers(tokenInfo *TokenInfo) ([]secu.User, error, int) {
	httpClient := client.getHttpClient(tokenInfo)

	accounts := []account.DefaultAccount{}
	found, err, respCode := client.findElement(&accounts, accountPath, httpClient)
	if !found || err != nil {
		log.Printf("in hydraClient.ListUsers, err:%v, respCode:%v\n", err, respCode)
		return nil, err, respCode
	}
	users := make([]secu.User, 0, len(accounts))
	for _, account := range accounts {
		users = append(users, *secu.NewUser(&account))
	}
	return users, nil, respCode
}

func (client HydraClient) FindUser(accountId string, tokenInfo *TokenInfo) (*secu.User, error) {
	httpClient := client.getHttpClient(tokenInfo)
	var account account.DefaultAccount
	if found, err, _ := client.findElement(&account, accountPath+"/"+accountId, httpClient); !found || err != nil {
		return nil, err
	}
	return secu.NewUser(&account), nil
}

func (client HydraClient) CreateUser(user *secu.User, tokenInfo *TokenInfo) (id string, err error) {
	// by default, an user is user active and not blocked
	user.SetActive(true)
	user.SetBlocked(false)

	httpClient := client.getHttpClient(tokenInfo)
	request := user.ToCreateAccountRequest()
	return client.createElement(*request, accountPath, httpClient)
}

func (client HydraClient) DeleteUser(accountId string, tokenInfo *TokenInfo) error {

	httpClient := client.getHttpClient(tokenInfo)
	return client.deleteElement(accountPath+"/"+accountId, httpClient)
}

func (client HydraClient) UpdateUserLogin(userId string, r secu.UpdateLoginRequest) error {
	hydraReq := account.UpdateUsernameRequest{
		Username: r.Login,
		Password: r.Password,
	}
	json, err := json.Marshal(hydraReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", client.authorizationServer+accountPath+"/"+userId+"/username", bytes.NewReader(json))
	if err != nil {
		return err
	}
	//FIXME
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while updating the user", resp.Status))
	}

	return nil
}

func (client HydraClient) UpdateUserPassword(userId string, r secu.UpdatePasswordRequest) error {
	hydraReq := account.UpdatePasswordRequest{
		CurrentPassword: r.CurrentPassword,
		NewPassword:     r.NewPassword,
	}
	json, err := json.Marshal(hydraReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", client.authorizationServer+accountPath+"/"+userId+"/password", bytes.NewReader(json))
	if err != nil {
		return err
	}
	//FIXME
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while updating the user", resp.Status))
	}

	return nil
}

func (client HydraClient) UpdateUserData(userId string, data secu.UserData, tokenInfo *TokenInfo) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	hydraReq := account.UpdateDataRequest{Data: string(jsonData)}
	jsonHydraReq, err := json.Marshal(hydraReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", client.authorizationServer+accountPath+"/"+userId+"/data", bytes.NewReader(jsonHydraReq))
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while updating the user", resp.Status))
	}

	return nil
}

func (client HydraClient) ListProfiles(tokenInfo *TokenInfo) ([]secu.Policy, error) {
	policies := []policy.DefaultPolicy{}
	httpClient := client.getHttpClient(tokenInfo)
	if found, err, _ := client.findElement(&policies, policyPath, httpClient); !found || err != nil {
		return nil, err
	}
	profiles := make([]secu.Policy, 0, len(policies))
	for _, policy := range policies {
		profile := secu.ConvertPolicy(&policy)
		if profile != nil {
			profiles = append(profiles, *profile)
		}
	}
	return profiles, nil
}

//ListPolicies list all the policies
// returns Policy, error, and http respCode
func (client HydraClient) ListPolicies(tokenInfo *TokenInfo) ([]secu.Policy, error, int) {
	defaultPolicies := []policy.DefaultPolicy{}
	httpClient := client.getHttpClient(tokenInfo)
	if found, err, respCode := client.findElement(&defaultPolicies, policyPath, httpClient); !found || err != nil {
		return nil, err, respCode
	}
	policies := make([]secu.Policy, 0, len(defaultPolicies))
	for _, policy := range defaultPolicies {
		profile := secu.ConvertPolicy(&policy)
		if profile != nil {
			policies = append(policies, *profile)
		}
	}
	return policies, nil, http.StatusOK
}

//FindPolicy find a policy by id
func (client HydraClient) FindPolicy(profileId string, tokenInfo *TokenInfo) (*secu.Policy, error) {
	var policy policy.DefaultPolicy
	httpClient := client.getHttpClient(tokenInfo)
	if found, err, _ := client.findElement(&policy, policyPath+"/"+profileId, httpClient); !found || err != nil {
		return nil, err
	}
	return secu.ConvertPolicy(&policy), nil
}

// CreatePolicy creates a new policy
func (client HydraClient) CreatePolicy(policy *secu.Policy, tokenInfo *TokenInfo) (id string, err error) {
	httpClient := client.getHttpClient(tokenInfo)

	hydraPolicy := policy.ToPolicy()
	return client.createElement(hydraPolicy, policyPath, httpClient)
}

func (client HydraClient) FindProfile(profileId string, tokenInfo *TokenInfo) (*secu.Policy, error) {
	var policy policy.DefaultPolicy
	httpClient := client.getHttpClient(tokenInfo)
	if found, err, _ := client.findElement(&policy, policyPath+"/"+profileId, httpClient); !found || err != nil {
		return nil, err
	}
	return secu.ConvertPolicy(&policy), nil
}

func (client HydraClient) CreateProfile(profile *secu.Policy, tokenInfo *TokenInfo) (id string, err error) {
	httpClient := client.getHttpClient(tokenInfo)

	policy := profile.ToPolicy()
	return client.createElement(policy, policyPath, httpClient)
}

func (client HydraClient) DeleteProfile(profileId string, tokenInfo *TokenInfo) error {
	req, err := http.NewRequest("DELETE", client.authorizationServer+policyPath+"/"+profileId, nil)
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while deleting the profile", resp.Status))
	}

	return nil
}

func (client HydraClient) UpdateProfileDescription(profileId string, escapedDescription []byte, tokenInfo *TokenInfo) error {
	req, err := http.NewRequest("PUT", client.authorizationServer+policyPath+"/"+profileId+"/description", bytes.NewReader(escapedDescription))
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while deleting the profile", resp.Status))
	}

	return nil
}

func (client HydraClient) UpdateProfileUsers(profileId string, userIds []string, tokenInfo *TokenInfo) error {
	profile, err := client.FindProfile(profileId, tokenInfo)
	if err != nil {
		return err
	}

	deletedUsers, addedUsers := diff(profile.Subjects, userIds)
	for _, deletedUser := range deletedUsers {
		if err = client.DeleteProfileUser(profileId, deletedUser, tokenInfo); err != nil {
			return err
		}
	}
	for _, addedUser := range addedUsers {
		if err = client.AddProfileUser(profileId, addedUser, tokenInfo); err != nil {
			return err
		}
	}

	return nil
}

func (client HydraClient) AddProfileUser(profileId string, userId string, tokenInfo *TokenInfo) error {
	req, err := http.NewRequest("PUT", client.authorizationServer+policyPath+"/"+profileId+"/subjects/"+userId, nil)
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while adding user to the profile", resp.Status))
	}

	return nil
}

func (client HydraClient) DeleteProfileUser(profileId string, userId string, tokenInfo *TokenInfo) error {
	req, err := http.NewRequest("DELETE", client.authorizationServer+policyPath+"/"+profileId+"/subjects/"+userId, nil)
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while deleting user from the profile", resp.Status))
	}

	return nil
}

func (client HydraClient) UpdateProfileRoles(profileId string, roles []string, tokenInfo *TokenInfo) error {
	//profile, err := client.FindProfile(profileId, tokenInfo)
	//if err != nil {
	//	return err
	//}

	//TODO clean this func, may be not useful anymore
	//deletedRoles, addedRoles := diff(profile.Permissions, roles)
	//for _, deletedRole := range deletedRoles {
	//	if err = client.DeleteProfileRole(profileId, deletedRole, tokenInfo); err != nil {
	//		return err
	//	}
	//}
	//for _, addedRole := range addedRoles {
	//	if err = client.AddProfileRole(profileId, addedRole, tokenInfo); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (client HydraClient) AddProfileRole(profileId string, role string, tokenInfo *TokenInfo) error {
	req, err := http.NewRequest("PUT", client.authorizationServer+policyPath+"/"+profileId+"/permissions/"+role, nil)
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while adding role the profile", resp.Status))
	}

	return nil
}

func (client HydraClient) DeleteProfileRole(profileId string, role string, tokenInfo *TokenInfo) error {
	req, err := http.NewRequest("DELETE", client.authorizationServer+policyPath+"/"+profileId+"/permissions/"+role, nil)
	if err != nil {
		return err
	}
	httpClient := client.getHttpClient(tokenInfo)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Got status '%s' while deleting role from the profile", resp.Status))
	}

	return nil
}

func (client HydraClient) findElement(element interface{}, path string, httpClient *http.Client) (found bool, err error, respCode int) {
	fullPath := client.authorizationServer + path
	resp, err := httpClient.Get(fullPath)
	if err != nil {

		if resp == nil {
			//when resp is nil, it means that the token is empty
			respCode = http.StatusUnauthorized
		} else {
			respCode = resp.StatusCode
		}
		log.Printf("Got error when trying to find %v : %v", fullPath, err)
		return false, errors.New("Error while getting the element: " + path), respCode
	}
	if resp.StatusCode == http.StatusNotFound {
		element = nil
		return false, nil, resp.StatusCode
	}
	if resp.StatusCode == 401 {
		//not authentified
		return false, errors.New("Not authentified"), resp.StatusCode
	}
	if resp.StatusCode == 403 {
		return false, errors.New("Not Authorized"), resp.StatusCode
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, errors.New(fmt.Sprintf("Got status '%s' while getting the element", resp.Status)), resp.StatusCode
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(element)
	if err != nil {
		return false, errors.New("Error while decoding the element: " + err.Error()), resp.StatusCode
	}

	return true, nil, resp.StatusCode
}

func (client HydraClient) createElement(element interface{}, path string, httpClient *http.Client) (id string, err error) {
	jsonElement, err := json.Marshal(element)
	if err != nil {
		return
	}

	resp, err := httpClient.Post(client.authorizationServer+path, "application/json", bytes.NewReader(jsonElement))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusCreated {
		err = errors.New("Got status " + resp.Status)
		return
	}
	location := resp.Header.Get("Location")
	if !strings.HasPrefix(location, path+"/") {
		err = errors.New("Invalid location: " + location)
		return
	}
	id = location[len(path)+1:]

	return
}

func (client HydraClient) deleteElement(path string, httpClient *http.Client) (err error) {

	req, _ := http.NewRequest("DELETE", client.authorizationServer+path, nil)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Got status " + resp.Status)
		return err
	}

	return nil
}

func (client HydraClient) Login(username string, password string) (token *oauth2.Token, err error) {

	return client.oauth2Config.PasswordCredentialsToken(oauth2.NoContext, username, password)
}

func diff(oldIds, newIds []string) (deleted, added []string) {
	deleted = []string{}
	added = []string{}

	for _, oldId := range oldIds {
		found := false
		for _, newId := range newIds {
			if newId == oldId {
				found = true
				break
			}
		}
		if !found {
			deleted = append(deleted, oldId)
		}
	}

	for _, newId := range newIds {
		found := false
		for _, oldId := range oldIds {
			if oldId == newId {
				found = true
				break
			}
		}
		if !found {
			added = append(added, newId)
		}
	}

	return
}
