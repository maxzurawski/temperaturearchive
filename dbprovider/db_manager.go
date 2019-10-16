package dbprovider

import (
	"fmt"
	"log"

	"github.com/xdevices/temperaturearchive/config"

	"github.com/jinzhu/gorm"
	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/model"
	"github.com/xdevices/utilities/db"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBManager interface {
	Find(searchDTO dto.SearchDTO) ([]model.Measurement, error)
	Save(measurement dto.MeasurementDTO) (*model.Measurement, error)
	GetDb() *gorm.DB

	// mappers
	MapToEntity(dto dto.MeasurementDTO) (measurement *model.Measurement)
}

var Mgr DBManager

type manager struct {
	db *gorm.DB
}

func InitDbManager() {
	dbPath := config.TemperaturearchiveConfig().DBPath()

	if path, exists := db.AdjustDBPath(dbPath); !exists {
		dbPath = path
	}

	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to init db[%s]:", dbPath)
		log.Fatal(errorMsg, err)
	}

	db.SingularTable(true)
	db.AutoMigrate(&model.Measurement{})
	Mgr = &manager{db: db}
}
