package redis

import (
	"fmt"
	"strings"
	"time"
)

// PrefixedRedisClient struct
type PrefixedRedisClient struct {
	prefix string
	client Client
}

// withPrefix adds a prefix to the key if the prefix supplied has a length greater than 0
func (p *PrefixedRedisClient) withPrefix(key string) string {
	if len(p.prefix) > 0 {
		return fmt.Sprintf("%s.%s", p.prefix, key)
	}
	return key
}

// withoutPrefix removes the prefix from a key if the prefix has a length greater than 0
func (p *PrefixedRedisClient) withoutPrefix(key string) string {
	if len(p.prefix) > 0 {
		return strings.Replace(key, fmt.Sprintf("%s.", p.prefix), "", 1)
	}
	return key
}

// Get wraps around redis get method by adding prefix and returning string and error directly
func (p *PrefixedRedisClient) Get(key string) (string, error) {
	return p.client.Get(p.withPrefix(key)).ResultString()
}

// Set wraps around redis get method by adding prefix and returning error directly
func (p *PrefixedRedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return p.client.Set(p.withPrefix(key), value, expiration).Err()
}

// Keys wraps around redis keys method by adding prefix and returning []string and error directly
func (p *PrefixedRedisClient) Keys(pattern string) ([]string, error) {
	keys, err := p.client.Keys(p.withPrefix(pattern)).Multi()
	if err != nil {
		return nil, err
	}

	woPrefix := make([]string, len(keys))
	for index, key := range keys {
		woPrefix[index] = p.withoutPrefix(key)
	}
	return woPrefix, nil

}

// Del wraps around redis del method by adding prefix and returning int64 and error directly
func (p *PrefixedRedisClient) Del(keys ...string) (int64, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = p.withPrefix(k)
	}
	return p.client.Del(prefixedKeys...).Result()
}

// SMembers returns a slice with all the members of a set
func (p *PrefixedRedisClient) SMembers(key string) ([]string, error) {
	return p.client.SMembers(p.withPrefix(key)).Multi()
}

// SAdd adds new members to a set
func (p *PrefixedRedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	return p.client.SAdd(p.withPrefix(key), members...).Result()
}

// SRem removes members from a set
func (p *PrefixedRedisClient) SRem(key string, members ...string) (int64, error) {
	return p.client.SRem(p.withPrefix(key), members).Result()
}

// Exists returns true if a key exists in redis
func (p *PrefixedRedisClient) Exists(keys ...string) (int64, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = p.withPrefix(k)
	}
	val, err := p.client.Exists(prefixedKeys...).Result()
	return val, err
}

// Incr increments a key. Sets it in one if it doesn't exist
func (p *PrefixedRedisClient) Incr(key string) error {
	return p.client.Incr(p.withPrefix(key)).Err()
}

// RPush insert all the specified values at the tail of the list stored at key
func (p *PrefixedRedisClient) RPush(key string, values ...interface{}) (int64, error) {
	return p.client.RPush(p.withPrefix(key), values...).Result()
}

// LRange Returns the specified elements of the list stored at key
func (p *PrefixedRedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return p.client.LRange(p.withPrefix(key), start, stop).Multi()
}

// LTrim Trim an existing list so that it will contain only the specified range of elements specified
func (p *PrefixedRedisClient) LTrim(key string, start, stop int64) error {
	return p.client.LTrim(p.withPrefix(key), start, stop).Err()
}

// LLen Returns the length of the list stored at key
func (p *PrefixedRedisClient) LLen(key string) int64 {
	return p.client.LLen(p.withPrefix(key)).Int()
}

// Expire set expiration time for particular key
func (p *PrefixedRedisClient) Expire(key string, value time.Duration) bool {
	return p.client.Expire(p.withPrefix(key), value).Bool()
}

// TTL for particular key
func (p *PrefixedRedisClient) TTL(key string) time.Duration {
	return p.client.TTL(p.withPrefix(key)).Duration()
}

// MGet fetchs multiple results
func (p *PrefixedRedisClient) MGet(keys []string) ([]interface{}, error) {
	keysWithPrefix := make([]string, 0)
	for _, key := range keys {
		keysWithPrefix = append(keysWithPrefix, p.withPrefix(key))
	}
	return p.client.MGet(keysWithPrefix).MultiInterface()
}

// NewPrefixedRedisClient returns a new Prefixed Redis Client
func NewPrefixedRedisClient(redisClient Client, prefix string) (*PrefixedRedisClient, error) {
	return &PrefixedRedisClient{
		client: redisClient,
		prefix: prefix,
	}, nil
}
