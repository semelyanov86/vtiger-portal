package cache

type Cache interface {
	Set(key string, value any, ttl int64) error
	Get(key string) (any, error)
}
