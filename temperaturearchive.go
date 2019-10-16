package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/xdevices/temperaturearchive/cache"
	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/temperaturearchive/dbprovider"
	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/handlers"
	"github.com/xdevices/temperaturearchive/observer"
	"github.com/xdevices/temperaturearchive/processors"
	"github.com/xdevices/temperaturearchive/publishers"
	"github.com/xdevices/temperaturearchive/service"
)

func main() {

	go observer.ObserveSensorChanges()
	go observer.TemperatureObserver([]func(dto.MeasurementDTO){processors.NotifierProcessor, processors.ArchiveProcessor})

	e := echo.New()
	e.GET("/", handlers.HandleFind)
	e.GET("/ws", handlers.HandleWebsocket)
	e.Logger.Fatal(e.Start(config.TemperaturearchiveConfig().Address()))
}

func init() {
	manager := config.EurekaManagerInit()
	manager.SendRegistrationOrFail()
	manager.ScheduleHeartBeat(config.TemperaturearchiveConfig().ServiceName(), 10)

	dbprovider.InitDbManager()
	service.Init()

	publishers.InitLogger()
	publishers.InitNotifier()

	_ = cache.InitSensorsCache(uuid.New().String())

}
