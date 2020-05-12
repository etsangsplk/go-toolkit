package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// ===== Command output / return value types

// Result generic interface
type Result interface {
	Int() int64
	String() string
	Bool() bool
	Duration() time.Duration
	Result() (int64, error)
	ResultString() (string, error)
	Multi() ([]string, error)
	MultiInterface() ([]interface{}, error)
	Err() error
}

// ResultImpl generic interface
type ResultImpl struct {
	value          int64
	valueString    string
	valueBool      bool
	valueDuration  time.Duration
	err            error
	multi          []string
	multiInterface []interface{}
}

// Int implementation
func (r *ResultImpl) Int() int64 {
	return r.value
}

// String implementation
func (r *ResultImpl) String() string {
	return r.valueString
}

// Bool implementation
func (r *ResultImpl) Bool() bool {
	return r.valueBool
}

// Duration implementation
func (r *ResultImpl) Duration() time.Duration {
	return r.valueDuration
}

// Err implementation
func (r *ResultImpl) Err() error {
	return r.err
}

// Result implementation
func (r *ResultImpl) Result() (int64, error) {
	return r.value, r.err
}

// ResultString implementation
func (r *ResultImpl) ResultString() (string, error) {
	return r.valueString, r.err
}

// Multi implementation
func (r *ResultImpl) Multi() ([]string, error) {
	return r.multi, r.err
}

// MultiInterface implementation
func (r *ResultImpl) MultiInterface() ([]interface{}, error) {
	return r.multiInterface, r.err
}

// ====== Client

// Client interface which specifies the currently used subset of redis operations
type Client interface {
	Del(keys ...string) Result
	Exists(keys ...string) Result
	Get(key string) Result
	Set(key string, value interface{}, expiration time.Duration) Result
	Ping() Result
	Keys(pattern string) Result
	SMembers(key string) Result
	SAdd(key string, members ...interface{}) Result
	SRem(key string, members ...interface{}) Result
	Incr(key string) Result
	RPush(key string, values ...interface{}) Result
	LRange(key string, start, stop int64) Result
	LTrim(key string, start, stop int64) Result
	LLen(key string) Result
	Expire(key string, value time.Duration) Result
	TTL(key string) Result
	MGet(keys []string) Result
}

// ClientImpl wrapps redis client
type ClientImpl struct {
	wrapped redis.UniversalClient
}

func (c *ClientImpl) wrapResult(result interface{}) Result {
	if result == nil {
		return nil
	}
	switch v := result.(type) {
	case *redis.StatusCmd:
		return &ResultImpl{
			valueString: v.Val(),
			err:         v.Err(),
		}
	case *redis.IntCmd:
		return &ResultImpl{
			value:       v.Val(),
			valueString: v.String(),
			err:         v.Err(),
		}
	case *redis.StringCmd:
		return &ResultImpl{
			valueString: v.Val(),
			err:         v.Err(),
		}
	case *redis.StringSliceCmd:
		return &ResultImpl{
			err:   v.Err(),
			multi: v.Val(),
		}
	case *redis.BoolCmd:
		return &ResultImpl{
			valueBool: v.Val(),
			err:       v.Err(),
		}
	case *redis.DurationCmd:
		return &ResultImpl{
			valueDuration: v.Val(),
			err:           v.Err(),
		}
	case *redis.SliceCmd:
		return &ResultImpl{
			err:            v.Err(),
			multiInterface: v.Val(),
		}
	default:
		return nil
	}
}

// Del implements Del wrapper for redis
func (c *ClientImpl) Del(keys ...string) Result {
	res := c.wrapped.Del(keys...)
	return c.wrapResult(res)
}

// Exists implements Exists wrapper for redis
func (c *ClientImpl) Exists(keys ...string) Result {
	res := c.wrapped.Exists(keys...)
	return c.wrapResult(res)
}

// Get implements Get wrapper for redis
func (c *ClientImpl) Get(key string) Result {
	res := c.wrapped.Get(key)
	return c.wrapResult(res)
}

// Set implements Set wrapper for redis
func (c *ClientImpl) Set(key string, value interface{}, expiration time.Duration) Result {
	res := c.wrapped.Set(key, value, expiration)
	return c.wrapResult(res)
}

// Ping implements Ping wrapper for redis
func (c *ClientImpl) Ping() Result {
	res := c.wrapped.Ping()
	return c.wrapResult(res)
}

// Keys implements Keys wrapper for redis
func (c *ClientImpl) Keys(pattern string) Result {
	res := c.wrapped.Keys(pattern)
	return c.wrapResult(res)
}

// SMembers implements SMembers wrapper for redis
func (c *ClientImpl) SMembers(key string) Result {
	res := c.wrapped.SMembers(key)
	return c.wrapResult(res)
}

// SAdd implements SAdd wrapper for redis
func (c *ClientImpl) SAdd(key string, members ...interface{}) Result {
	res := c.wrapped.SAdd(key, members...)
	return c.wrapResult(res)
}

// SRem implements SRem wrapper for redis
func (c *ClientImpl) SRem(key string, members ...interface{}) Result {
	res := c.wrapped.SRem(key, members...)
	return c.wrapResult(res)
}

// Incr implements Incr wrapper for redis
func (c *ClientImpl) Incr(key string) Result {
	res := c.wrapped.Incr(key)
	return c.wrapResult(res)
}

// RPush implements RPush wrapper for redis
func (c *ClientImpl) RPush(key string, values ...interface{}) Result {
	res := c.wrapped.RPush(key, values...)
	return c.wrapResult(res)
}

// LRange implements LRange wrapper for redis
func (c *ClientImpl) LRange(key string, start, stop int64) Result {
	res := c.wrapped.LRange(key, start, stop)
	return c.wrapResult(res)
}

// LTrim implements LTrim wrapper for redis
func (c *ClientImpl) LTrim(key string, start, stop int64) Result {
	res := c.wrapped.LTrim(key, start, stop)
	return c.wrapResult(res)
}

// LLen implements LLen wrapper for redis
func (c *ClientImpl) LLen(key string) Result {
	res := c.wrapped.LLen(key)
	return c.wrapResult(res)
}

// Expire implements Expire wrapper for redis
func (c *ClientImpl) Expire(key string, value time.Duration) Result {
	res := c.wrapped.Expire(key, value)
	return c.wrapResult(res)
}

// TTL implements TTL wrapper for redis
func (c *ClientImpl) TTL(key string) Result {
	res := c.wrapped.TTL(key)
	return c.wrapResult(res)
}

// MGet implements MGet wrapper for redis
func (c *ClientImpl) MGet(keys []string) Result {
	res := c.wrapped.MGet(keys...)
	return c.wrapResult(res)
}

// NewClient returns new client implementation
func NewClient(options *UniversalOptions) (Client, error) {
	return &ClientImpl{
		wrapped: redis.NewUniversalClient(options),
	}, nil
}
