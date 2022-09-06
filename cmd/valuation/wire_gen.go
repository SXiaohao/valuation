// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"valuation/internal/biz"
	"valuation/internal/conf"
	"valuation/internal/data"
	"valuation/internal/server"
	"valuation/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	grpcServer := server.NewGRPCServer(confServer, logger)

	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	goodRepo := data.NewGoodRepo(dataData, logger)
	goodUsecase := biz.NewGoodUsecase(goodRepo, logger)
	goodService := service.NewGoodService(goodUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, goodService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
