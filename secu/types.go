package secu

import (
	"encoding/json"

	"github.com/ory-am/hydra/account"
	"github.com/ory-am/ladon/policy"
)

type Role struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type UserData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Active    bool   `json:"active"`
	Blocked   bool   `json:"blocked"`
}

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	UserData
}

type Policy struct {
	Id          string   `json:"id"`
	Description string   `json:"description"`
	Subjects    []string `json:"subjects"`
	Permissions []string `json:"permissions"`
	Resource    string   `json:"resource"` // only one resource per policy
}

type UpdateLoginRequest struct {
	Password string `json:"password" valid:"required"`
	Login    string `json:"login" valid:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" valid:"required"`
	NewPassword     string `json:"newPassword" valid:"required"`
}

func NewUser(account *account.DefaultAccount) *User {
	if account == nil {
		return nil
	}
	var userData UserData
	json.Unmarshal([]byte(account.Data), &userData)
	return &User{
		Id:       account.ID,
		Login:    account.Username,
		UserData: userData,
	}
}

func (user *User) ToCreateAccountRequest() *account.CreateAccountRequest {
	if user == nil {
		return nil
	}
	jsonData, _ := json.Marshal(user.UserData)
	return &account.CreateAccountRequest{
		Username: user.Login,
		Password: user.Password,
		Data:     string(jsonData),
	}
}

func (user *User) IsActive() bool {
	return user.Active
}

func (user *User) IsBlocked() bool {
	return user.Blocked
}

func (user *User) SetActive(active bool) {
	user.Active = active
}

func (user *User) SetBlocked(blocked bool) {
	user.Blocked = blocked
}

func ConvertPolicy(policy *policy.DefaultPolicy) *Policy {
	if policy == nil || policy.Effect != "allow" || len(policy.Resources) > 1 || policy.Resources[0] == "<.*>" {
		//if policy == nil || policy.Effect != "allow" || len(policy.Resources) != 1 {
		return nil
	}

	return &Policy{
		Id:          policy.ID,
		Description: policy.Description,
		Subjects:    policy.Subjects,
		Permissions: policy.Permissions,
		Resource:    policy.Resources[0],
	}
}

//func convertPermissions(permissionsString []string) []Permission {
//	//TODO get label, value of permissions from consul kv store
//	permissions := make([]Permission, 0, len(permissionsString))
//	for _, perm := range permissionsString {
//		permissions = append(permissions, Permission{perm, perm})
//	}
//	return permissions
//}
//
//func getPermissionValues(permisssions []Permission) ([]string){
//	perms := make([]string,0, len(permisssions))
//	for _, perm := range permisssions {
//		perms = append(perms, perm.value)
//	}
//	return perms
//}

func (p *Policy) ToPolicy() *policy.DefaultPolicy {
	if p == nil {
		return nil
	}
	return &policy.DefaultPolicy{
		ID:          p.Id,
		Description: p.Description,
		Subjects:    p.Subjects,
		Effect:      "allow",
		Resources:   []string{p.Resource},
		Permissions: p.Permissions,
	}
}
