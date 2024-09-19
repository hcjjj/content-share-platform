--利用 Redis 的有序集合（Sorted Set）数据结构来记录请求的时间戳，并通过计算当前窗口内的请求数来进行限流


-- 1, 2, 3, 4, 5, 6, 7 这是的元素
-- ZREMRANGEBYSCORE key1 0 6
-- 7 执行完之后

-- 限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber( ARGV[2])
local now = tonumber(ARGV[3])
-- 窗口的起始时间
local min = now - window

--删除集合中时间早于 min 的记录，确保有序集合只保留当前窗口内的请求
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)
--计当前窗口内的请求数
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
if cnt >= threshold then
    -- 执行限流
    return "true"
else
    -- 把 score 和 member 都设置成 now
    --分数用于排序，成员用于唯一性标识
    --向有序集合中添加当前请求的时间戳
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window)
    return "false"
end