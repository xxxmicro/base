package cache

import(
	"github.com/xxxmicro/base/store"
	"github.com/xxxmicro/base/store/redis"
	"github.com/micro/go-micro/v2/config"
)

func NewCacheProvider(config config.Config) Cache {
	t := config.Get("cache", "type").String("redis")
	switch t {
		case "redis":
			return newRedisCache(config)
	}
	
	return newRedisCache(config)
}

func newRedisCache(config config.Config) Cache {
	addrs := config.Get("redis", "addrs").StringSlice(nil)

	options := make([]store.Option, 0)
	options = append(options, store.Addrs(addrs...))

	store := NewStore(redis.NewStore(options...))
	store.Init()

	return store.(Cache)
}