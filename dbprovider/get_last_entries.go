package dbprovider

import "github.com/maxzurawski/temperaturearchive/model"

func (mgr *manager) GetLastEntries(amount int, uuid, processid string) ([]model.Measurement, error) {
	var results []model.Measurement
	error := mgr.GetDb().Where("uuid = ?", uuid).
		Where("process_id <> ?", processid).
		Limit(amount).Order("reported_at desc").Find(results).Error
	return results, error
}
