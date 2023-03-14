package cache

import (
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
)

type PreemptibleNodeDetails struct {
	MarkedToBeDeleted bool
	DeletionSuccess   bool
	TimeSinceCordon   time.Time
}

var NodeCache *ttlcache.Cache[string, PreemptibleNodeDetails]

func Cache() *ttlcache.Cache[string, PreemptibleNodeDetails] {
	return ttlcache.New[string, PreemptibleNodeDetails](
		ttlcache.WithTTL[string, PreemptibleNodeDetails](24 * time.Hour),
	)
}
