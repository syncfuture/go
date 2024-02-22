package sredis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/syncfuture/go/serr"
)

var (
	ErrorUnlockFailed = errors.New("unlock failed or lock already released")
)

// StrLock 结构体
type StrLock struct {
	client    redis.UniversalClient
	lockKey   string
	lockValue string
	timeout   time.Duration
}

// NewStrLock 创建一个新的 RedisLock 实例
func NewStrLock(client redis.UniversalClient, lockKey, lockValue string, timeout time.Duration) *StrLock {
	return &StrLock{
		client:    client,
		lockKey:   lockKey,
		lockValue: lockValue,
		timeout:   timeout,
	}
}

// Lock 尝试获取锁
func (lock *StrLock) Lock() (bool, error) {
	// 使用 SET 命令的 NX 选项尝试获取锁，并设置超时时间以自动释放锁
	result, err := lock.client.SetNX(context.Background(), lock.lockKey, lock.lockValue, lock.timeout).Result()
	if err != nil {
		return false, serr.WithStack(err)
	}
	return result, nil
}

// Unlock 释放锁
func (lock *StrLock) Unlock() error {
	// 使用 Lua 脚本来原子地释放锁
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	result, err := lock.client.Eval(context.Background(), script, []string{lock.lockKey}, lock.lockValue).Result()
	if err != nil {
		return serr.WithStack(err)
	}

	// 如果结果不是 1，说明锁已经被其他进程释放或过期
	if result.(int64) != 1 {
		return ErrorUnlockFailed
	}
	return nil
}

// Locked 判断锁是否上锁
func (lock *StrLock) Locked() (bool, error) {
	result, err := lock.client.Exists(context.Background(), lock.lockKey).Result()
	if err != nil {
		return false, serr.WithStack(err)
	}
	return result > 0, nil
}
