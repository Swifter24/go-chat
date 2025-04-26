package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go_chat/internal/config"
	"go_chat/pkg/zlog"
	"log"
	"time"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.Host
	port := conf.RedisConfig.Port
	password := conf.RedisConfig.Password
	db := conf.RedisConfig.Db
	addr := fmt.Sprintf("%s:%d", host, port)
	//配置redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// SetKeyEx 设置 Ex 形数据
func SetKeyEx(key string, value string, expiration time.Duration) error {
	err := redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetKey 查询该 key 可以返回空
func GetKey(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info("改key不存在")
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// GetKeyNilIsErr 查询该 key 不可返回空
func GetKeyNilIsErr(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetKeyWithPrefixNilIsErr	获取有prefix前缀的key，如果不存在，err ≠ nil
func GetKeyWithPrefixNilIsErr(prefix string) (string, error) {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return "", err
		}
		if len(keys) == 0 {
			zlog.Info("没有找到相关前缀的key")
			return "", nil
		}
		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关前缀的key", keys))
			return keys[0], nil
		} else {
			zlog.Error("找到了数量大于1的key，查找异常")
			return "", errors.New("找到了数量大于1的key，查找异常")
		}
	}
}

// GetKeyWithSuffixNilIsErr 获取有suffix后缀的key，如果不存在，err ≠ nil
func GetKeyWithSuffixNilIsErr(suffix string) (string, error) {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
		if err != nil {
			return "", err
		}
		if len(keys) == 0 {
			zlog.Info("没有找到相关后缀的key")
			return "", nil
		}
		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关后缀的key", keys))
			return keys[0], nil
		} else {
			zlog.Error("找到了数量大于1的key，查找异常")
			return "", errors.New("找到了数量大于1的key，查找异常")
		}
	}
}

// DelKeyIfExists 如果key存在，删除该key
func DelKeyIfExists(key string) error {
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		delErr := redisClient.Del(ctx, key).Err()
		if delErr != nil {
			return delErr
		}
	}
	return nil
}

// DelKeysWithPrefix 如果有该前缀的key，则删除
func DelKeysWithPrefix(prefix string) error {
	var keys []string
	var err error

	for {
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			log.Println("没有找到对应的key")
			break
		}

		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相应key", keys)
		}
	}
	return nil
}

// DelKeysWithSuffix 如果有该前缀的key，则删除
func DelKeysWithSuffix(suffix string) error {
	var keys []string
	var err error

	for {
		keys, err = redisClient.Keys(ctx, suffix+"*").Result()
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			log.Println("没有找到对应的key")
			break
		}

		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相应key", keys)
		}
	}
	return nil
}

// DeleteAllRedisKeys 删除所有的key
func DeleteAllRedisKeys() error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}
		cursor = nextCursor

		if len(keys) > 0 {
			_, err := redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

func DelKeysWithPattern(pattern string) error {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			log.Println("没有找到对应key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相关对应key", keys)
		}
	}

	return nil
}
