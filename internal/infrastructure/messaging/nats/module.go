package nats

import (
	"context"
	"fmt"

	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/resources"
	"github.com/fercho/school-tracking/services/fleet/pkg/env"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewConnection establishes a NATS connection using the URL from config.
func NewConnection(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.NatsURL,
		nats.Name("fleet-service"),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS at %s: %w", cfg.NatsURL, err)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Connected to NATS", zap.String("url", cfg.NatsURL))
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Closing NATS connection")
			nc.Drain()
			return nil
		},
	})

	return nc, nil
}

// NewEventPublisher wraps the NATS publisher as the EventPublisher port.
func NewEventPublisher(nc *nats.Conn) (resources.EventPublisher, error) {
	return NewPublisher(nc)
}

// Module provides NATS connection and EventPublisher to the fx dependency graph.
var Module = fx.Module("messaging.nats",
	fx.Provide(NewConnection),
	fx.Provide(NewEventPublisher),
)
