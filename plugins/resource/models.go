package resource

//A Resource typically represents the security information for a plugin
//
// Resource is a struct containing the resources provided by running plugins to be authorized
//contains 2 keys, one for translations, and one for security

type Keys struct {
	Lang     string `json:"lang"`     // example : accounts
	Security string `json:"security"` // example : rn:hydra:accounts
}
type Resource struct {
	Key         string `json:"key"`
	SecurityKey string `json:"security_key"`
	//Keys        Keys     `json:"keys"`
	Permissions []string `json:"permissions"`
}
