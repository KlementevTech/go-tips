package errlog

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KlementevTech/gotips/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var internalStatus = status.New(codes.Internal, "internal server error, see details log")

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		return resp, handleErr(ctx, err, internalStatus)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, stream)
		return handleErr(stream.Context(), err, internalStatus)
	}
}

func handleErr(ctx context.Context, err error, internalStatus *status.Status) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	code := toGrpcCode(err)

	if code == codes.Internal {
		slog.Default().ErrorContext(ctx, "internal error", slog.String("error", err.Error()))
		return internalStatus.Err()
	}

	return status.Error(code, err.Error())
}

func toGrpcCode(err error) codes.Code {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return codes.NotFound
	case errors.Is(err, domain.ErrPreconditionFailed):
		return codes.FailedPrecondition
	case errors.Is(err, domain.ErrAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}
