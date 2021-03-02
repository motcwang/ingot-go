package provider

import (
	"github.com/google/wire"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/container"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/container/provider/preset"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/endpoint"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token/granter"
)

// AuthorizationServerContainer 授权服务器容器
var AuthorizationServerContainer = wire.NewSet(wire.Struct(new(container.AuthorizationServerContainer), "*"))

// AuthorizationServerContainerFields 授权服务器容器所有字段
var AuthorizationServerContainerFields = wire.NewSet(
	AuthorizationAuthenticationManager,
	AuthorizationServerTokenServices,
	ConsumerTokenServices,
	TokenEndpoint,
	TokenEnhancer,
	TokenEnhancers,
	TokenGranters,
	TokenGranter,
	PasswordTokenGranter,
)

// AuthorizationAuthenticationManager 授权服务器中的认证管理器
func AuthorizationAuthenticationManager(securityContainer *container.SecurityContainer, injector container.SecurityInjector) authentication.AuthorizationManager {
	if injector.GetAuthorizationAuthenticationManager() != nil {
		return injector.GetAuthorizationAuthenticationManager()
	}
	return preset.AuthorizationAuthenticationManager(securityContainer)
}

// AuthorizationServerTokenServices 授权服务器 token 服务
func AuthorizationServerTokenServices(oauth2Container *container.OAuth2Container, securityContainer *container.SecurityContainer, enhancer token.Enhancer, manager authentication.AuthorizationManager, injector container.SecurityInjector) token.AuthorizationServerTokenServices {
	if injector.GetAuthorizationServerTokenServices() != nil {
		return injector.GetAuthorizationServerTokenServices()
	}
	return preset.AuthorizationServerTokenServices(oauth2Container, securityContainer, enhancer, manager)
}

// ConsumerTokenServices 令牌撤销
func ConsumerTokenServices(oauth2Container *container.OAuth2Container, injector container.SecurityInjector) token.ConsumerTokenServices {
	if injector.GetConsumerTokenServices() != nil {
		return injector.GetConsumerTokenServices()
	}
	return preset.ConsumerTokenServices(oauth2Container)
}

// TokenEndpoint 端点
func TokenEndpoint(granter token.Granter, securityContainer *container.SecurityContainer, injector container.SecurityInjector) *endpoint.TokenEndpoint {
	if injector.GetTokenEndpoint() != nil {
		return injector.GetTokenEndpoint()
	}
	return preset.TokenEndpoint(granter, securityContainer)
}

// TokenEnhancer token增强，默认使用增强链
func TokenEnhancer(enhancers token.Enhancers, oauth2Container *container.OAuth2Container, injector container.SecurityInjector) token.Enhancer {
	if injector.GetTokenEnhancer() != nil {
		return injector.GetTokenEnhancer()
	}
	return preset.TokenEnhancer(enhancers, oauth2Container)
}

// TokenEnhancers 自定义增强
func TokenEnhancers(injector container.SecurityInjector) token.Enhancers {
	if len(injector.GetTokenEnhancers()) != 0 {
		return injector.GetTokenEnhancers()
	}
	return preset.TokenEnhancers()
}

// TokenGranters 自定义授权
func TokenGranters(injector container.SecurityInjector) token.Granters {
	if injector.GetTokenGranters() != nil {
		return injector.GetTokenGranters()
	}
	return preset.TokenGranters()
}

// TokenGranter token 授权
func TokenGranter(granters token.Granters, password *granter.PasswordTokenGranter, injector container.SecurityInjector) token.Granter {
	if injector.GetTokenGranter() != nil {
		return injector.GetTokenGranter()
	}
	return preset.TokenGranter(granters, password)
}

// PasswordTokenGranter 密码模式授权
func PasswordTokenGranter(tokenServices token.AuthorizationServerTokenServices, manager authentication.AuthorizationManager, injector container.SecurityInjector) *granter.PasswordTokenGranter {
	if injector.GetPasswordTokenGranter() != nil {
		return injector.GetPasswordTokenGranter()
	}
	return preset.PasswordTokenGranter(tokenServices, manager)
}
