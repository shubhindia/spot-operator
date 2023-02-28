package cache

import (
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
)

func Cache() *ttlcache.Cache[string, bool] {
	return ttlcache.New[string, bool](
		ttlcache.WithTTL[string, bool](24 * time.Hour),
	)
}
