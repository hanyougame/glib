package userx

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// 定义常量
const (
	// KeyPrefix Redis键前缀
	KeyPrefix = "{cache}:user:exception:%d"
)

type UserX struct {
	rdb redis.UniversalClient //redis
}

func NewUserX(rdb redis.UniversalClient) *UserX {
	return &UserX{rdb: rdb}
}

// StoreUserExceptionStatus 存储用户异常状态到Redis
// userId: 用户ID
// exceptionValue: 异常状态值
// expireSeconds: 过期时间（秒），如果为0则使用默认过期时间
// 返回: 存储是否成功，错误信息
func (u *UserX) StoreUserExceptionStatus(ctx context.Context, userId, exceptionValue int64) error {
	if userId == 0 || exceptionValue == 0 {
		return fmt.Errorf("用户ID和异常状态值不能为空")
	}
	return u.rdb.Set(ctx, fmt.Sprintf(KeyPrefix, userId), exceptionValue, 0).Err()
}

// GetUserExceptionStatus 从Redis获取用户异常状态
// userId: 用户ID
// 返回: 异常状态值，错误信息
func (u *UserX) GetUserExceptionStatus(ctx context.Context, userId int64) (val int64, err error) {
	if userId == 0 {
		return 0, fmt.Errorf("用户ID不能为空")
	}
	// 获取数据
	value, err := u.rdb.Get(ctx, fmt.Sprintf(KeyPrefix, userId)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 键不存在的情况
			return 1, nil
		}
		// 其他错误
		return
	}
	return strconv.ParseInt(value, 10, 64)
}

// DelUserExceptionStatus 从Redis删除用户异常状态
// userId: 用户ID
// 返回: 异常状态值，错误信息
func (u *UserX) DelUserExceptionStatus(ctx context.Context, userId int64) error {
	if userId == 0 {
		return fmt.Errorf("用户ID不能为空")
	}
	return u.rdb.Del(ctx, fmt.Sprintf(KeyPrefix, userId)).Err()
}
