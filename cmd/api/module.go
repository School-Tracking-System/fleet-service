package main

import (
	"github.com/fercho/school-tracking/services/fleet/internal/core/fleet"
	internalGrpc "github.com/fercho/school-tracking/services/fleet/internal/infrastructure/grpc"
	"github.com/fercho/school-tracking/services/fleet/internal/infrastructure/grpc/handlers"
	"github.com/fercho/school-tracking/services/fleet/internal/infrastructure/persistence/postgres"
	"github.com/fercho/school-tracking/services/fleet/pkg/env"
	"github.com/fercho/school-tracking/services/fleet/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func AppModule() fx.Option {
	return fx.Options(
		env.Module,
		logger.Module,
		postgres.Module,
		fleet.Module,
		handlers.Module,
		internalGrpc.Module, // Provides the gRPC server
		fx.Invoke(func(*grpc.Server) {}), // Forces fx to instantiate and start it
	)
}
