package config

import (
	"github.com/ingot-cloud/ingot-go/pkg/framework/security"
	"github.com/ingot-cloud/ingot-go/pkg/framework/security/web/builders"
)

// WebSecurityConfigurerAdapter 安全配置适配器
type WebSecurityConfigurerAdapter struct {
	AdditionalHTTPSecurityConfigurer security.HTTPSecurityConfigurer
}

// NewWebSecurityConfigurerAdapter 实例化
func NewWebSecurityConfigurerAdapter(httpSecurity security.HTTPSecurityConfigurer) *WebSecurityConfigurerAdapter {
	return &WebSecurityConfigurerAdapter{
		AdditionalHTTPSecurityConfigurer: httpSecurity,
	}
}

// WebConfigure Web安全配置
func (adapter *WebSecurityConfigurerAdapter) WebConfigure(web security.WebSecurityBuilder) error {
	http, err := adapter.getHTTP()
	if err != nil {
		return err
	}

	web.AddSecurityFilterChainBuilder(http)

	return nil
}

func (adapter *WebSecurityConfigurerAdapter) getHTTP() (security.HTTPSecurityBuilder, error) {

	http := builders.NewHTTPSecurity()

	err := adapter.applyDefaultConfiguration(http)
	if err != nil {
		return nil, err
	}

	if adapter.AdditionalHTTPSecurityConfigurer != nil {
		err = adapter.AdditionalHTTPSecurityConfigurer.HTTPConfigure(http)
	}

	return http, err
}

func (adapter *WebSecurityConfigurerAdapter) applyDefaultConfiguration(http *builders.HTTPSecurity) error {
	// 应用默认配置
	// http.A()  http.B()

	return nil
}
