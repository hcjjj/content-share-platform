-- 验证码在 Redis 上的 key
-- phone_code:login:181xxxxxxxx
local key = KEYS[1]
-- 验证次数，一个验证码最多验证三次，这个记录还可以验证几次
-- phone_code:login:181xxxxxxxx:cnt
local cntKey = key..":cnt"
-- 你的验证码 1234
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))

-- https://redis.io/commands/
-- The command returns -2 if the key does not exist.
-- The command returns -1 if the key exists but has no associated expire.

if ttl == -1 then
    -- key 存在，但是没有过期时间值
    -- 系统错误，手动设置了这个 key
    return -2
    -- 540 = 600s - 60s 九分钟
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    -- 符合预期
    return 0
else
    -- 发送太频繁
    return -1
end

