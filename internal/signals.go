package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignals(ctx context.Context) os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		return sig
	case <-ctx.Done():
		return nil
	}
}
