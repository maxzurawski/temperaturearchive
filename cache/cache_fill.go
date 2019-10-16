package cache

func (c Cache) FillCache(sensors []Sensor) {
	lock.Lock()
	defer lock.Unlock()
	for _, item := range sensors {
		c[item.Uuid] = item
	}
}
