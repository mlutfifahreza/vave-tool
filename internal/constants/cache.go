package constants

import "time"

const (
	ProductCacheKeyPrefix = "product:%s"
	ProductCacheTTL       = 15 * time.Minute
)
