package errors

import (
	"net/http"

	"github.com/ingot-cloud/ingot-go/pkg/framework/core/errors"
	"github.com/ingot-cloud/ingot-go/pkg/framework/core/utils"
)

// InvalidToken 无效的Token，自定义提示信息
func InvalidToken(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusUnauthorized, InvalidTokenCode, message)
}

// OAuth2AccessDenied error
func OAuth2AccessDenied(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusForbidden, AccessDeniedCode, message)
}

// InvalidGrant 无效的授权
func InvalidGrant(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusBadRequest, InvalidGrantCode, message)
}

// InvalidScope 无效的Scope
func InvalidScope(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusBadRequest, InvalidScopeCode, message)
}

// InvalidClient 无效的Client
func InvalidClient(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusUnauthorized, InvalidClientCode, message)
}

// InvalidRequest 无效的情况
func InvalidRequest(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusBadRequest, InvalidRequestCode, message)
}

// InsufficientAuthentication 不充足的认证
func InsufficientAuthentication(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.Unauthorized(message)
}

// UnsupportedGrantType 不支持的 grant type
func UnsupportedGrantType(args ...string) error {
	message := utils.StringCombine(args...)
	return errors.New(http.StatusBadRequest, UnsupportedGrantTypeCode, message)
}
