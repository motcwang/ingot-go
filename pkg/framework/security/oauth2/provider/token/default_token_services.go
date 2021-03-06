package token

import (
	"strings"
	"time"

	"github.com/ingot-cloud/ingot-go/pkg/framework/core/utils/uuid"
	securityAuthentication "github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication/preauth"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/core"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/constants"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/errors"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/authentication"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/clientdetails"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/request"
)

// DefaultTokenServices 默认token服务
// 实现 AuthorizationServerTokenServices, ResourceServerTokenServices 和 ConsumerTokenServices
type DefaultTokenServices struct {
	// 默认RefreshToken有效时间，单位秒
	RefreshTokenValiditySeconds int
	// 默认AccessToken有效时间，单位秒
	AccessTokenValiditySeconds int
	// 是否支持RefreshToken
	SupportRefreshToken bool
	// 是否重复使用RefreshToken
	ReuseRefreshToken     bool
	TokenStore            Store
	ClientDetailsService  clientdetails.Service
	TokenEnhancer         Enhancer
	AuthenticationManager securityAuthentication.Manager
}

// NewDefaultTokenServices 实例化默认 TokenServices
func NewDefaultTokenServices(tokenStore Store) *DefaultTokenServices {
	return &DefaultTokenServices{
		RefreshTokenValiditySeconds: 60 * 60 * 24 * 30, // default 30 days.
		AccessTokenValiditySeconds:  60 * 60 * 2,       // default 2 hours.
		SupportRefreshToken:         false,
		ReuseRefreshToken:           true,
		TokenStore:                  tokenStore,
	}
}

// AuthorizationServerTokenServices

// CreateAccessToken 通过身份验证信息创建访问令牌
func (service *DefaultTokenServices) CreateAccessToken(auth *authentication.OAuth2Authentication) (OAuth2AccessToken, error) {
	existingAccessToken, err := service.TokenStore.GetAccessToken(auth)
	if err != nil {
		return nil, err
	}

	var refreshToken OAuth2RefreshToken

	if existingAccessToken != nil {
		// 如果已经过去，则移除该访问令牌以及对应的刷新令牌
		if existingAccessToken.IsExpired() {
			if existingAccessToken.GetRefreshToken() != nil {
				refreshToken = existingAccessToken.GetRefreshToken()
				service.TokenStore.RemoveRefreshToken(refreshToken)
			}
			service.TokenStore.RemoveAccessToken(existingAccessToken)
		} else {
			service.TokenStore.StoreAccessToken(existingAccessToken, auth)
			return existingAccessToken, nil
		}
	}

	// 为什么要判断nil？
	// 1. 如果过期的 AccessToken 存在相应的 RefreshToken，
	// 	  那么客户端可能持有该 RefreshToken，所以我们需要重用该 RefreshToken
	// 2. 如果不存在 RefreshToken，直接实例新的 RefreshToken
	// 非nil，如果 RefreshToken 已经过期，那么重新创建
	if refreshToken == nil || service.isExpired(refreshToken) {
		refreshToken, err = service.createRefreshToken(auth)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := service.createAccessToken(auth, refreshToken)
	if err != nil {
		return nil, err
	}
	service.TokenStore.StoreAccessToken(accessToken, auth)

	refreshToken = accessToken.GetRefreshToken()
	if refreshToken != nil {
		service.TokenStore.StoreRefreshToken(refreshToken, auth)
	}

	return accessToken, nil
}

// RefreshAccessToken 通过refresh token和请求信息刷新token
func (service *DefaultTokenServices) RefreshAccessToken(refreshTokenValue string, tokenRequest *request.TokenRequest) (OAuth2AccessToken, error) {
	if !service.SupportRefreshToken {
		return nil, errors.InvalidGrant("Invalid refresh token: ", refreshTokenValue)
	}

	refreshToken, err := service.TokenStore.ReadRefreshToken(refreshTokenValue)
	if err != nil {
		return nil, err
	}
	if refreshToken == nil {
		return nil, errors.InvalidGrant("Invalid refresh token: ", refreshTokenValue)
	}

	auth, err := service.TokenStore.ReadAuthenticationForRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if service.AuthenticationManager != nil && !auth.IsClientOnly() {
		var user core.Authentication
		user = preauth.NewAuthenticationToken(auth.UserAuthentication, "", auth.GetAuthorities())
		user, err := service.AuthenticationManager.Authenticate(user)
		if err != nil {
			return nil, err
		}
		details := auth.GetDetails()
		auth = authentication.NewOAuth2Authentication(auth.GetOAuth2Request(), user)
		auth.SetDetails(details)
	}
	clientID := auth.GetOAuth2Request().GetClientID()
	if clientID == "" || clientID != tokenRequest.GetClientID() {
		return nil, errors.InvalidGrant("Wrong client for this refresh token: ", refreshTokenValue)
	}

	// 清楚当前存储的访问令牌
	service.TokenStore.RemoveAccessTokenUsingRefreshToken(refreshToken)

	if service.isExpired(refreshToken) {
		service.TokenStore.RemoveRefreshToken(refreshToken)
		return nil, errors.InvalidGrant("Invalid refresh token (expired): ", refreshTokenValue)
	}

	auth, err = service.createRefreshedAuthentication(auth, tokenRequest)

	// 判断是否再次使用当前 RefreshToken，如果不再使用当前RefreshToken，那么创建一个新的
	if !service.ReuseRefreshToken {
		service.TokenStore.RemoveRefreshToken(refreshToken)
		refreshToken, err = service.createRefreshToken(auth)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := service.createAccessToken(auth, refreshToken)
	if err != nil {
		return nil, err
	}
	service.TokenStore.StoreAccessToken(accessToken, auth)
	// 如果不在使用当前RefreshToken，那么保存新的
	if !service.ReuseRefreshToken {
		service.TokenStore.StoreRefreshToken(accessToken.GetRefreshToken(), auth)
	}

	return accessToken, nil
}

// GetAccessToken 根据身份验证信息获取访问令牌
func (service *DefaultTokenServices) GetAccessToken(auth *authentication.OAuth2Authentication) (OAuth2AccessToken, error) {
	return service.TokenStore.GetAccessToken(auth)
}

// ResourceServerTokenServices

// LoadAuthentication 通过access token加载身份验证信息
func (service *DefaultTokenServices) LoadAuthentication(accessTokenValue string) (*authentication.OAuth2Authentication, error) {
	accessToken, err := service.TokenStore.ReadAccessToken(accessTokenValue)
	if err != nil {
		return nil, err
	}
	if accessToken == nil {
		return nil, errors.InvalidToken("Invalid access token: ", accessTokenValue)
	} else if accessToken.IsExpired() {
		service.TokenStore.RemoveAccessToken(accessToken)
		return nil, errors.InvalidToken("Access token expired: ", accessTokenValue)
	}

	result, err := service.TokenStore.ReadAuthentication(accessToken)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.InvalidToken("Invalid access token: ", accessTokenValue)
	}

	if service.ClientDetailsService != nil {
		clientID := result.GetOAuth2Request().GetClientID()

		_, err := service.ClientDetailsService.LoadClientByClientID(clientID)
		if err != nil {
			return nil, errors.InvalidToken("Client not valid: ", clientID, ", original error = ", err.Error())
		}
	}

	return result, nil
}

// ReadAccessToken 读取指定access token详细信息
func (service *DefaultTokenServices) ReadAccessToken(accessToken string) (OAuth2AccessToken, error) {
	return service.TokenStore.ReadAccessToken(accessToken)
}

// ConsumerTokenServices

// RevokeToken 撤销令牌
func (service *DefaultTokenServices) RevokeToken(tokenValue string) bool {
	accessToken, err := service.TokenStore.ReadAccessToken(tokenValue)
	if err != nil || accessToken == nil {
		return false
	}

	if accessToken.GetRefreshToken() != nil {
		service.TokenStore.RemoveRefreshToken(accessToken.GetRefreshToken())
	}
	service.TokenStore.RemoveAccessToken(accessToken)
	return true
}

func (service *DefaultTokenServices) createRefreshToken(auth *authentication.OAuth2Authentication) (OAuth2RefreshToken, error) {
	support, err := service.isSupportRefreshToken(auth.GetOAuth2Request())
	if err != nil {
		return nil, err
	}
	if !support {
		return nil, nil
	}
	validitySeconds, err := service.getRefreshTokenValiditySeconds(auth.GetOAuth2Request())
	if err != nil {
		return nil, err
	}
	value := uuid.MustString()
	if validitySeconds > 0 {
		return NewDefaultExpiringOAuth2RefreshToken(value, time.Now().Add(time.Duration(validitySeconds)*time.Second)), nil
	}

	return NewDefaultOAuth2RefreshToken(value), nil
}

func (service *DefaultTokenServices) createAccessToken(auth *authentication.OAuth2Authentication, refreshToken OAuth2RefreshToken) (OAuth2AccessToken, error) {
	token := NewDefaultOAuth2AccessToken(uuid.MustString())
	validitySeconds, err := service.getAccessTokenValiditySeconds(auth.GetOAuth2Request())
	if err != nil {
		return nil, err
	}
	if validitySeconds > 0 {
		token.Expiration = time.Now().Add(time.Duration(validitySeconds) * time.Second)
	}
	token.RefreshToken = refreshToken
	token.Scope = auth.GetOAuth2Request().GetScope()
	if service.TokenEnhancer != nil {
		return service.TokenEnhancer.Enhance(token, auth)
	}

	return token, nil
}

func (service *DefaultTokenServices) createRefreshedAuthentication(auth *authentication.OAuth2Authentication, tokenRequest *request.TokenRequest) (*authentication.OAuth2Authentication, error) {
	narrowed := auth
	scope := tokenRequest.GetScope()
	clientAuth := auth.GetOAuth2Request().UpdateRefresh(tokenRequest)
	if len(scope) != 0 {
		originalScope := clientAuth.GetScope()
		originalScopeMap := make(map[string]string)
		// 判断原始scope中是否包含新的scope
		for _, item := range originalScope {
			originalScopeMap[item] = item
		}
		for _, item := range scope {
			if _, ok := originalScopeMap[item]; !ok {
				return nil, errors.InvalidToken("Unable to narrow the scope of the client authentication to ", strings.Join(scope, ","), ".", strings.Join(originalScope, ","))
			}
		}

		clientAuth = clientAuth.NarrowScope(scope)
	}

	narrowed = authentication.NewOAuth2Authentication(clientAuth, auth.UserAuthentication)
	return narrowed, nil
}

// 判断 OAuth2RefreshToken 是否过期
func (service *DefaultTokenServices) isExpired(refreshToken OAuth2RefreshToken) bool {
	if expiring, ok := refreshToken.(ExpiringOAuth2RefreshToken); ok {
		if expiring.GetExpiration().Before(time.Now()) {
			return true
		}
	}
	return false
}

func (service *DefaultTokenServices) getAccessTokenValiditySeconds(clientAuth *request.OAuth2Request) (int, error) {
	client, err := service.getClientDetails(clientAuth.ClientID)
	if err != nil {
		return 0, err
	}
	if client != nil {
		validity := client.GetAccessTokenValiditySeconds()
		if validity != 0 {
			return validity, nil
		}
	}

	return service.AccessTokenValiditySeconds, nil
}

func (service *DefaultTokenServices) getRefreshTokenValiditySeconds(clientAuth *request.OAuth2Request) (int, error) {
	client, err := service.getClientDetails(clientAuth.ClientID)
	if err != nil {
		return 0, err
	}
	if client != nil {
		validity := client.GetRefreshTokenValiditySeconds()
		if validity != 0 {
			return validity, nil
		}
	}

	return service.RefreshTokenValiditySeconds, nil
}

func (service *DefaultTokenServices) isSupportRefreshToken(clientAuth *request.OAuth2Request) (bool, error) {
	client, err := service.getClientDetails(clientAuth.ClientID)
	if err != nil {
		return false, err
	}
	if client != nil {
		grantTypes := client.GetAuthorizedGrantTypes()
		for _, grant := range grantTypes {
			if grant == constants.GrantTypeRefresh {
				return true, nil
			}
		}
	}
	return service.SupportRefreshToken, nil
}

func (service *DefaultTokenServices) getClientDetails(clientID string) (clientdetails.ClientDetails, error) {
	if service.ClientDetailsService != nil {
		return service.ClientDetailsService.LoadClientByClientID(clientID)
	}
	return nil, nil
}
