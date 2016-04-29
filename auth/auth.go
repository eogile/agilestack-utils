package auth

import (
	"encoding/json"
	"golang.org/x/oauth2"
)

type TokenInfo struct {
	TokenInfo string `json:"tokenInfo"`
}

func EncodeTokenInfo(token *oauth2.Token) (*TokenInfo, error) {
	if token == nil {
		return nil, nil
	}
	tokenBytes, err := json.Marshal(*token)

	return &TokenInfo{
		TokenInfo: string(tokenBytes),
	}, err

}

func DecodeTokenInfo(tokenInfo *TokenInfo) (*oauth2.Token, error) {
	if tokenInfo == nil {
		return nil, nil
	}

	token := &oauth2.Token{}
	err := json.Unmarshal([]byte(tokenInfo.TokenInfo), token)
	if err != nil {
		return nil, err
	}
	return token, err
}
