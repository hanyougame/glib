package xhttp

import (
	"context"
	"github.com/hanyougame/glib/utils/xerr"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	// BusinessCodeOK represents the business code for success.
	BusinessCodeOK = 0
	// BusinessMsgOk represents the business message for success.
	BusinessMsgOk = "ok"
)

// BaseResponse is the base response struct.
type BaseResponse[T any] struct {
	// Code represents the business code, not the http xstatus code.
	Code int ` json:"code" xml:"code"`
	// Msg represents the business message, if Code = BusinessCodeOK,
	// and Msg is empty, then the Msg will be set to BusinessMsgOk.
	Message string `json:"message" xml:"message"`
	// Data represents the business data.
	Data T `json:"data,omitempty" xml:"data,omitempty"`
}

// JsonBaseResponseCtx writes v into w with http.StatusOK.
func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v any) {
	var (
		traceId = trace.TraceIDFromContext(ctx)
		spanID  = trace.SpanIDFromContext(ctx)
	)

	if len(spanID) > 0 {
		w.Header().Set("X-Span-Id", spanID)
		w.Header().Set("X-Trace-Id", traceId)
	}
	httpx.OkJsonCtx(ctx, w, wrapBaseResponse(v))
}

func wrapBaseResponse(v any) BaseResponse[any] {
	var resp BaseResponse[any]

	switch data := v.(type) {
	case *errors.CodeMsg:
		resp.Code = data.Code
		resp.Message = data.Msg
	case errors.CodeMsg:
		resp.Code = data.Code
		resp.Message = data.Msg
	case *status.Status:
		resp.Code = int(data.Code())
		resp.Message = data.Message()
	case error:
		resp.Code = int(xerr.ServerError)
		resp.Message = data.Error()
	default:
		resp.Code = BusinessCodeOK
		resp.Message = BusinessMsgOk
		resp.Data = v
	}
	return resp
}
