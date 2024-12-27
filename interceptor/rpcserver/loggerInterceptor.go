package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	errors2 "github.com/zeromicro/x/errors"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err) // err类型
		var e *errors2.CodeMsg
		logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
		if errors.As(causeErr, &e) { //自定义错误类型
			err = status.Error(codes.Code(e.Code), e.Msg)
		}
	}

	return resp, err
}
