--验证码在 Redis 上的 key
-- phone_code:login:152xxxxxxxx
local key = KEYS[1]
-- 验证次数，一个验证码，最多重复三次，这个记录还可以验证几次
-- phone_code:login:152xxxxxxxx:cnt
local cntKey = key..":cnt"
-- 验证码 123456
local val= ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    --    key 存在，但是没有过期时间
    -- 系统错误，手动设置了这个 key，但是没给过期时间
    return -2
    --如果验证码的剩余有效期小于 540 秒，或者验证码不存在
    -- 1 min 只能发一条 验证码的有效期是 10 min
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    -- 完美，符合预期
    return 0
else
    -- 发送太频繁
    return -1
end

