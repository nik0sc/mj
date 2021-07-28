package cache

type Cache interface {
	Add(key, value interface{}) (evicted bool)
	Get(key interface{}) (value interface{}, ok bool)
}
