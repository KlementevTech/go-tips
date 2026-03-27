package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func waitForSignal(ctx context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	slog.Default().InfoContext(ctx, "waiting for signal", slog.Any("signal", sigsToStrings(sigs...)))
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

func sigsToStrings(sigs ...os.Signal) []string {
	if len(sigs) == 0 {
		return nil
	}
	res := make([]string, len(sigs))
	for i, sig := range sigs {
		res[i] = sig.String()
	}
	return res
}
