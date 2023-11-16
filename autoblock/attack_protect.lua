-- 攻击检测key
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
return { tonumber(val), blocked }