package main

import (
	"github.com/labstack/echo"
	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/temperaturearchive/dbprovider"
	"github.com/xdevices/temperaturearchive/handlers"
	"github.com/xdevices/temperaturearchive/service"
)

func main() {

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

}
