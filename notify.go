package kafka_playground

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	return ctx, cancel
}
