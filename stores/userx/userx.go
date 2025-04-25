package userx

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
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
func (u *UserX) StoreUserExceptionStatus(ctx context.Context, userId int64, exceptionValue int) error {
	if userId == 0 || exceptionValue == 0 {
		return fmt.Errorf("用户ID和异常状态值不能为空")
	}
	return u.rdb.Set(ctx, fmt.Sprintf(KeyPrefix, userId), exceptionValue, 0).Err()
}

// GetUserExceptionStatus 从Redis获取用户异常状态
// userId: 用户ID
// 返回: 异常状态值，错误信息
func (u *UserX) GetUserExceptionStatus(ctx context.Context, userId int64) (val int, err error) {
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
	return cast.ToInt(value), nil
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

// BatchStoreUserExceptionStatus 批量存储多个用户相同的异常状态到Redis
// userIds: 用户ID列表
// exceptionValue: 统一的异常状态值
// expireSeconds: 过期时间（秒），如果为0则使用默认过期时间
// 返回: 存储是否成功，错误信息
func (u *UserX) BatchStoreUserExceptionStatus(ctx context.Context, userIds []int64, exceptionValue int) error {
	if len(userIds) == 0 {
		return fmt.Errorf("用户ID列表不能为空")
	}
	if exceptionValue == 0 {
		return fmt.Errorf("异常状态值不能为空")
	}

	pipe := u.rdb.Pipeline()
	for _, userId := range userIds {
		if userId == 0 {
			return fmt.Errorf("用户ID不能为空")
		}
		pipe.Set(ctx, fmt.Sprintf(KeyPrefix, userId), exceptionValue, 0)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// BatchGetUserExceptionStatus 批量获取用户异常状态
// userIds: 用户ID列表
// 返回: 用户ID与异常状态值的映射，错误信息
func (u *UserX) BatchGetUserExceptionStatus(ctx context.Context, userIds []int64) (map[int64]int, error) {
	if len(userIds) == 0 {
		return nil, fmt.Errorf("用户ID列表不能为空")
	}

	keys := make([]string, len(userIds))
	for i, userId := range userIds {
		if userId == 0 {
			return nil, fmt.Errorf("用户ID不能为空")
		}
		keys[i] = fmt.Sprintf(KeyPrefix, userId)
	}

	values, err := u.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[int64]int, len(userIds))
	for i, val := range values {
		if val == nil {
			result[userIds[i]] = 1 // 键不存在，返回默认值
		} else {
			result[userIds[i]] = cast.ToInt(val)
		}
	}

	return result, nil
}

// BatchDelUserExceptionStatus 批量删除用户异常状态
// userIds: 用户ID列表
// 返回: 是否删除成功，错误信息
func (u *UserX) BatchDelUserExceptionStatus(ctx context.Context, userIds []int64) error {
	if len(userIds) == 0 {
		return fmt.Errorf("用户ID列表不能为空")
	}

	keys := make([]string, len(userIds))
	for i, userId := range userIds {
		if userId == 0 {
			return fmt.Errorf("用户ID不能为空")
		}
		keys[i] = fmt.Sprintf(KeyPrefix, userId)
	}

	return u.rdb.Del(ctx, keys...).Err()
}
