// +build wireinject
// The build tag makes sure the stub is not built in the final build

package injector

import (
	"github.com/google/wire"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/container"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/container/provider"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/core/userdetails"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/config"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/clientdetails"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token"
)

func BuildSecurityContainer(userDetailsService userdetails.Service, clientDetailsService clientdetails.Service) (*container.SecurityContainer, error) {
	wire.Build(
		provider.Providers,
		provider.PasswordEncoder,
		provider.UserCache,
		provider.PreChecker,
		provider.PostChecker,
		provider.DaoAuthenticationProviderSet,
		provider.SecurityContainerSet,
	)
	return nil, nil
}

func BuildOAuth2Container(oauth2Config config.OAuth2) (*container.OAuth2Container, error) {
	wire.Build(
		provider.DefaultTokenServices,
		provider.TokenStore,
		provider.JwtAccessTokenConverter,
		provider.AccessTokenConverter,
		provider.UserAuthenticationConverter,
		provider.OAuth2ContainerSet,
	)
	return nil, nil
}

func BuildResourceServerContainer(oauth2Container *container.OAuth2Container) (*container.ResourceServerContainer, error) {
	wire.Build(
		provider.ResourceServerTokenServices,
		provider.OAuth2SecurityConfigurer,
		provider.TokenExtractor,
		provider.ResourceAuthenticationManager,
		provider.ResourceServerContainerSet,
	)
	return nil, nil
}

func BuildAuthorizationServerContainer(oauth2Container *container.OAuth2Container, securityContainer *container.SecurityContainer, enhancers token.Enhancers) (*container.AuthorizationServerContainer, error) {
	wire.Build(
		provider.AuthorizationServerTokenServices,
		provider.ConsumerTokenServices,
		provider.TokenEnhancer,
		provider.AuthorizationAuthenticationManager,
		provider.AuthorizationServerContainerSet,
	)
	return nil, nil
}
