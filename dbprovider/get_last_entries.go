package dbprovider

import "github.com/maxzurawski/temperaturearchive/model"

func (mgr *manager) GetLastEntries(amount int, uuid string) ([]model.Measurement, error) {
	var results []model.Measurement
	error := mgr.GetDb().Limit(amount).Order("reported_at desc").Find(results).Error
	return results, error
}
