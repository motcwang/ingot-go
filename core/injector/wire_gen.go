// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package injector

import (
	"ingot/api"
	"ingot/core/container"
	"ingot/core/provider"
	"ingot/model/dao"
	"ingot/router"
	"ingot/service"
)

// Injectors from wire.go:

func BuildContainer() (*container.Container, func(), error) {
	authentication, cleanup, err := provider.AuthenticationProvider()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := provider.GormProvider()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	user := &dao.User{
		DB: db,
	}
	test := &service.Test{
		UserDao: user,
	}
	apiTest := &api.Test{
		Test: test,
	}
	routerRouter := &router.Router{
		Auth: authentication,
		Test: apiTest,
	}
	engine := provider.HTTPHandlerProvider(routerRouter)
	containerContainer := &container.Container{
		Engine: engine,
	}
	return containerContainer, func() {
		cleanup2()
		cleanup()
	}, nil
}
