package cache

func (c *Cache) Reset() {
	*c = make(map[string]Sensor)
}
