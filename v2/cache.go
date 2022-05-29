package ggm

type Cache interface {
	Get(key string) (string, error)
	Set(key string, data string) error
}

var cache Cache

func SetCache(c Cache) {
	cache = c
}
