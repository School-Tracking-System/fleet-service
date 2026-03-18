package grpc

import (
	"context"
	"net"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/fleet/pkg/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(
	lc fx.Lifecycle,
	cfg *env.Config,
	log *zap.Logger,
	vehicleHandler pb.VehicleServiceServer,
	routeHandler pb.RouteServiceServer,
	driverHandler pb.DriverServiceServer,
	studentHandler pb.StudentServiceServer,
	guardianHandler pb.GuardianServiceServer,
	schoolHandler pb.SchoolServiceServer,
) *grpc.Server {

	// In a real approach, you can inject interceptors here (logging, recovery, auth validation)
	grpcServer := grpc.NewServer()

	// Register handlers
	pb.RegisterVehicleServiceServer(grpcServer, vehicleHandler)
	pb.RegisterRouteServiceServer(grpcServer, routeHandler)
	pb.RegisterDriverServiceServer(grpcServer, driverHandler)
	pb.RegisterStudentServiceServer(grpcServer, studentHandler)
	pb.RegisterGuardianServiceServer(grpcServer, guardianHandler)
	pb.RegisterSchoolServiceServer(grpcServer, schoolHandler)

	// Enable reflection for tools like grpcurl/evans
	reflection.Register(grpcServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
			if err != nil {
				log.Fatal("failed to listen for grpc server", zap.Error(err))
			}

			log.Info("gRPC server listening",
				zap.String("addr", lis.Addr().String()),
				zap.String("configured_port", cfg.GRPCPort),
			)

			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					log.Fatal("failed to serve grpc", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping gRPC server gracefully")
			grpcServer.GracefulStop()
			return nil
		},
	})

	return grpcServer
}

var Module = fx.Provide(NewGRPCServer)
