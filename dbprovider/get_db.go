package dbprovider

import "github.com/jinzhu/gorm"

func (mgr *manager) GetDb() *gorm.DB {
	return mgr.db
}
