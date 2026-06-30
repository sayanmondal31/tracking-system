local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4] or 1)

-- Retrieve the current token count and last updated timestamp
local data = redis.call('HMGET', key, 'tokens', 'last_updated')
local tokens = tonumber(data[1])
local last_updated = tonumber(data[2])

if not tokens then
    -- Initialize the bucket if it does not exist
    tokens = capacity
    last_updated = now
else
    -- Calculate how many tokens were generated since the last request
    local elapsed = math.max(0, now - last_updated)
    tokens = math.min(capacity, tokens + (elapsed * refill_rate))
end

-- Allow request and deduct a token if possible
if tokens >= requested then
    tokens = tokens - requested
    redis.call('HMSET', key, 'tokens', tokens, 'last_updated', now)
    redis.call('EXPIRE', key, 86400) -- Set TTL of 24h to clean up stale client IPs
    return 1 -- Allowed
else
    return 0 -- Rate Limited
end
