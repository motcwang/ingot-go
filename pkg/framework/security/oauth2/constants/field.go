package constants

// OAuth2 请求映射参数常量
const (
	ClientID          = "client_id"
	State             = "state"
	Scope             = "scope"
	RedirectURI       = "redirect_uri"
	ResponseType      = "response_type"
	UserOAuthApproval = "user_oauth_approval"
	ScopePrefix       = "scope."
	GrantType         = "grant_type"
	Password          = "password"
	ClientSecret      = "client_secret"
)

// TokenPayloadKey Token载体key
type TokenPayloadKey string

// Token载体key
const (
	TokenScope       TokenPayloadKey = "scope"
	TokenAuthorities TokenPayloadKey = "authorities"
	TokenClientID    TokenPayloadKey = "client_id"
	TokenUser        TokenPayloadKey = "user"
	TokenJti         TokenPayloadKey = "jti"
)
