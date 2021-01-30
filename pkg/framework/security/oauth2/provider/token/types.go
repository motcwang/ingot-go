package token

import (
	"time"

	"github.com/ingot-cloud/ingot-go/pkg/framework/core/model/enums"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/core"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/authentication"
)

// OAuth2AccessToken OAuth2 访问令牌
type OAuth2AccessToken interface {
	// 获取额外信息
	GetAdditionalInformation() map[string]interface{}
	// 获取令牌访问范围
	GetScope() []string
	// 获取刷新令牌
	GetRefreshToken() OAuth2RefreshToken
	// 获取令牌类型
	GetTokenType() enums.TokenType
	// 令牌是否过期
	IsExpired() bool
	// 令牌到期时间
	GetExpiration() time.Time
	// 令牌有效期，单位秒
	GetExpiresIn() int
	// 获取令牌值
	GetValue() string
}

// OAuth2RefreshToken OAuth2 刷新令牌
type OAuth2RefreshToken interface {
	GetRefreshTokenValue() string
}

// ResourceServerTokenServices 资源服务器 token 服务
type ResourceServerTokenServices interface {
	// 通过token加载身份验证信息
	LoadAuthentication(string) (*authentication.OAuth2Authentication, error)
	// 读取指定token详细信息
	ReadAccessToken(string) OAuth2AccessToken
}

// UserAuthenticationConverter 用户map信息和身份验证信息互相转换接口
type UserAuthenticationConverter interface {
	// 在身份验证信息中提取访问令牌使用的信息
	ConvertUserAuthentication(core.Authentication) (map[string]interface{}, error)
	// 从map中提取身份验证信息
	ExtractAuthentication(map[string]interface{}) (core.Authentication, error)
}
