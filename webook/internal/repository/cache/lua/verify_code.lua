local key = KEYS[1]
-- 用户输入的 code
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 说明试错机会没了
    return -1
elseif expectedCode == code then
    -- 输入正确
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户手抖 输错了
    -- 可验证次数 -1
    redis.call("decr", cntKey, -1)
    return -2
end