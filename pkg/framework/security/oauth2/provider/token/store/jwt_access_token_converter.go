package store

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/authentication"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/oauth2/provider/token"
)

// JwtAccessTokenConverter jwt和oauth2身份验证信息转换器
type JwtAccessTokenConverter struct {
	tokenConverter token.AccessTokenConverter
	signingMethod  jwt.SigningMethod
}

// NewJwtAccessTokenConverter 实例化
func NewJwtAccessTokenConverter(method jwt.SigningMethod) *JwtAccessTokenConverter {
	return &JwtAccessTokenConverter{
		tokenConverter: token.NewDefaultAccessTokenConverter(),
		signingMethod:  method,
	}
}

// ConvertAccessToken 返回访问令牌映射内容
func (c *JwtAccessTokenConverter) ConvertAccessToken(accessToken token.OAuth2AccessToken, authentication *authentication.OAuth2Authentication) (map[string]interface{}, error) {
	return c.tokenConverter.ConvertAccessToken(accessToken, authentication)
}

// ExtractAccessToken 根据token value和映射内容提取访问令牌
func (c *JwtAccessTokenConverter) ExtractAccessToken(accessToken string, mapInfo map[string]interface{}) (token.OAuth2AccessToken, error) {
	return c.tokenConverter.ExtractAccessToken(accessToken, mapInfo)
}

// ExtractAuthentication 根据token映射信息提取身份验证信息
func (c *JwtAccessTokenConverter) ExtractAuthentication(mapInfo map[string]interface{}) (*authentication.OAuth2Authentication, error) {
	return c.tokenConverter.ExtractAuthentication(mapInfo)
}

// Enhance 增强
func (c *JwtAccessTokenConverter) Enhance(accessToken token.OAuth2AccessToken, authentication *authentication.OAuth2Authentication) token.OAuth2AccessToken {
	// todo
	return nil
}

// SetAccessTokenConverter 设置访问令牌转换器
func (c *JwtAccessTokenConverter) SetAccessTokenConverter(tokenConverter token.AccessTokenConverter) {
	c.tokenConverter = tokenConverter
}

// GetAccessTokenConverter 获取访问令牌转换器
func (c *JwtAccessTokenConverter) GetAccessTokenConverter() token.AccessTokenConverter {
	if c.tokenConverter == nil {
		c.tokenConverter = &token.DefaultAccessTokenConverter{}
	}
	return c.tokenConverter
}

// Encode 编码
func (c *JwtAccessTokenConverter) Encode(accessToken token.OAuth2AccessToken, auth *authentication.OAuth2Authentication) (string, error) {
	token, err := c.GetAccessTokenConverter().ConvertAccessToken(accessToken, auth)
	if err != nil {
		return "", nil
	}
	// todo jwt 转换 token
	return "", nil
}

// Decode 解码
func (c *JwtAccessTokenConverter) Decode(token string) {

}
