package dbprovider

import (
	"fmt"
	"log"

	"github.com/maxzurawski/temperaturearchive/config"

	"github.com/jinzhu/gorm"
	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/temperaturearchive/model"
	"github.com/maxzurawski/utilities/db"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBManager interface {
	Find(searchDTO dto.SearchDTO) ([]model.Measurement, error)
	Save(measurement dto.MeasurementDTO) (*model.Measurement, error)
	GetDb() *gorm.DB

	// mappers
	MapToEntity(dto dto.MeasurementDTO) (measurement *model.Measurement)
	MapToDto(measurement *model.Measurement) dto.MeasurementDTO
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
