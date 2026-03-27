package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"

	transport "github.com/KlementevTech/gotips/internal/transport/http"
)

type Config struct {
	Address string `mapstructure:"address"`
	Enable  bool   `mapstructure:"enable"`
}

func RunServer(ctx context.Context, cfg Config) error {
	mux := http.NewServeMux()
	registerRoutes(mux)
	return transport.RunServer(ctx, transport.NewDefaultConfig(cfg.Address, "pprof"), mux)
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
