package container

import (
	"github.com/ingot-cloud/ingot-go/pkg/framework/security"
	coreAuth "github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication/provider/basic"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication/provider/dao"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/core/userdetails"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/crypto/password"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/authentication"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/config"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/clientdetails"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/endpoint"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token/granter"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token/store"
)

// CommonContainer 容器
type CommonContainer struct {
	WebSecurityConfigurers security.WebSecurityConfigurers
	PasswordEncoder        password.Encoder
	UserCache              userdetails.UserCache
	PreChecker             userdetails.PreChecker
	PostChecker            userdetails.PostChecker
	UserDetailsService     userdetails.Service
	ClientDetailsService   clientdetails.Service
}

// OAuth2Container OAuth2 容器
type OAuth2Container struct {
	OAuth2Config                config.OAuth2
	TokenStore                  token.Store
	JwtAccessTokenConverter     *store.JwtAccessTokenConverter
	AccessTokenConverter        token.AccessTokenConverter
	UserAuthenticationConverter token.UserAuthenticationConverter
}

// ResourceServerContainer 资源服务器容器
type ResourceServerContainer struct {
	AuthenticationManager       coreAuth.ResourceManager
	ResourceServerConfigurer    security.ResourceServerConfigurer
	ResourceServerTokenServices token.ResourceServerTokenServices
	TokenExtractor              authentication.TokenExtractor
}

// AuthorizationServerContainer 授权服务器容器
type AuthorizationServerContainer struct {
	AuthenticationManager            coreAuth.AuthorizationManager
	AuthorizationServerConfigurer    security.AuthorizationServerConfigurer
	AuthorizationServerTokenServices token.AuthorizationServerTokenServices
	ConsumerTokenServices            token.ConsumerTokenServices
	TokenEndpoint                    *endpoint.TokenEndpoint
	TokenEndpointHTTPConfigurer      endpoint.OAuth2HTTPConfigurer
	TokenEnhancer                    token.Enhancer
	TokenGranter                     token.Granter
	PasswordTokenGranter             *granter.PasswordTokenGranter
}

// AuthProvidersContainer 认证提供者容器
type AuthProvidersContainer struct {
	Providers coreAuth.Providers
	Basic     *basic.AuthenticationProvider
	Dao       *dao.AuthenticationProvider
}
