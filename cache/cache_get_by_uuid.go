package cache

func (c Cache) GetByUuid(uuid string) *Sensor {
	lock.Lock()
	defer lock.Unlock()

	if sensor, ok := c[uuid]; !ok {
		return nil
	} else {
		return &sensor
	}
}
