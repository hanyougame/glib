package metadata

import (
	"context"
	"testing"
)

func TestGetMetadata(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetadata(ctx, CtxJWTUserId, 9527)
	t.Log(GetMetadata[int](ctx, CtxJWTUserId))
}

func TestGetUidFromCtx(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetadata(ctx, CtxJWTUserId, 9527)

	t.Log(GetUidFromCtx(ctx))
}

func TestGetUsernameFromCtx(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetadata(ctx, CtxJWTUsername, "tom")
	t.Log(GetUsernameFromCtx(ctx))
}
