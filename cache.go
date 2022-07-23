package fly

type Cache interface {
	Get(key string) (string, error)
	Del(key string) error
	Set(key string, data string) error
}

var cache Cache

func SetCache(c Cache) {
	cache = c
}
