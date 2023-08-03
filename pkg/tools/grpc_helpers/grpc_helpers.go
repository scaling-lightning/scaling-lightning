package grpc_helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const UnaryTimeout = 2 * time.Minute

// 150mb
const RecvMsgSize = 150 << (10 * 2)

// InterceptorLogger adapts zerolog logger to interceptor logger.
func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l := l.With().Fields(fields).Logger()

			switch lvl {
			case logging.LevelDebug:
				l.Debug().Msg(msg)
			case logging.LevelInfo:
				l.Info().Msg(msg)
			case logging.LevelWarn:
				l.Warn().Msg(msg)
			case logging.LevelError:
				l.Error().Msg(msg)
			default:
				panic(fmt.Sprintf("unknown level %v", lvl))
			}
		},
	)
}

func GetLoggingOptions() []logging.Option {
	return []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
}

func ClientInterceptor(name string) func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		headerMap := make(map[string]string)
		headerMap["X-Node"] = name
		md := metadata.New(headerMap)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
