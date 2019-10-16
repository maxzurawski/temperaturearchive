package cache

import "sync"

type Cache map[string]Sensor

var (
	lock sync.Mutex
)

var SensorsCache *Cache
