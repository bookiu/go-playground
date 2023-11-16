package autoblock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"

	"github.com/bookiu/go-playground/utils/randutil"
)

const (
	redisAddr = "127.0.0.1:6379"
	redisPass = "123456"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
}

// AutoBlock 自动封禁
// 客户端 clientIP 在 duration 秒内访问次数超过 threshold 次时，自动封禁客户端 blockTTL 秒
func AutoBlock(ctx context.Context, clientIP string, duration, threshold, blockTTL int) bool {
	checkKey := "attack_check#ip:" + clientIP
	blockKey := "attack_block#ip:" + clientIP
	timestamp := time.Now().Unix()
	windowSize := duration
	uniqueID := randutil.RandomString(10)

	script := getScript()
	ret, err := rdb.
		Eval(ctx, script, []string{checkKey, blockKey}, timestamp, windowSize, uniqueID, threshold, blockTTL).
		Int64Slice()
	if err != nil {
		// TODO: 记录日志
		return false
	}
	blocked := ret[1]
	// reqCount := ret[0]
	// TODO: 将 reqCount 记录到日志中
	if blocked == 1 {
		return true
	}
	return false
}

func getScript() string {
	return `-- 攻击检测key
local key = KEYS[1]
-- 封禁key
local block_key = KEYS[2]
-- 当前时间戳
local timestamp = tonumber(ARGV[1])
-- 滑动窗口大小，单位秒
local win_size = tonumber(ARGV[2])
-- 请求ID（唯一）
local rid = ARGV[3]
-- 封禁阈值
local threshold = tonumber(ARGV[4])
-- 封禁时长
local block_ttl = ARGV[5]
local start_time = timestamp - win_size

--
redis.call("ZADD", key, timestamp, rid)
local val = redis.call("ZCOUNT", key, start_time, timestamp)
redis.call("ZREMRANGEBYSCORE", key, 0, start_time)
-- 是否封禁
local blocked = 0
if val and tonumber(val) >= threshold then
    -- 请求量超过阈值，封禁
    redis.call("SET", block_key, "", "NX", "EX", block_ttl)
    blocked = 1
else
    -- 判断封禁key是否存在，存在表示封禁中
    local ret = redis.call("EXISTS", block_key)
    if ret == 1 then
        blocked = 1
    end
end
-- 返回滑动窗口内请求量和是否封禁了
return { tonumber(val), blocked }`
}
