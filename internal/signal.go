package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func waitForSignal(ctx context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	toStrings := func(sigs ...os.Signal) (res []string) {
		for _, sig := range sigs {
			res = append(res, sig.String())
		}
		return res
	}

	slog.Default().InfoContext(ctx, "waiting for signal", slog.Any("signal", toStrings(sigs...)))
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
