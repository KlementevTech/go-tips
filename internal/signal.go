package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func waitForSignal(ctx context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	sigChan := make(chan os.Signal, len(sigs))
	signal.Notify(sigChan, sigs...)

	newCtx, cancel := context.WithCancel(ctx)

	go func() {
		defer signal.Stop(sigChan)

		select {
		case sig := <-sigChan:
			slog.Default().InfoContext(newCtx, "received signal", slog.Any("signal", sig.String()))
			cancel()
		case <-newCtx.Done():
		}
	}()

	return newCtx, cancel
}
