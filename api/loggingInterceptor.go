package api

import (
	"context"
	"fmt"
	logr "github.com/adam-bunce/grpc-todo/helpers"
	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logr.Info(fmt.Sprintf(info.FullMethod + " called"))
	resp, err := handler(ctx, req)
	if err != nil {
		logr.Error(err.Error())
	}
	return resp, err
}
