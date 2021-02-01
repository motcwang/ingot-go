package authentication

import (
	ginwrapper "github.com/ingot-cloud/ingot-go/pkg/framework/core/wrapper/gin"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/authentication/preauth"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/core"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/core/ingot"
)

// BearerTokenExtractor bearer token 提取器
type BearerTokenExtractor struct {
}

// Extract 提取token
func (e *BearerTokenExtractor) Extract(ctx *ingot.Context) core.Authentication {
	token := ginwrapper.GetBearerToken(ctx.Context)
	if token != "" {
		// 返回 preauth token
		return preauth.NewAuthenticationToken(token, "", nil)
	}
	return nil
}
