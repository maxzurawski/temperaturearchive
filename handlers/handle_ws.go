package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/xdevices/utilities/rabbit/crosscutting"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/utilities/rabbit/observer"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartWsTemperatureObserver() *observer.Observer {
	observer := config.TemperaturearchiveConfig().InitObserver()
	observer.DeclareTopicExchange(crosscutting.TopicMeasurements.String())
	observer.BindQueue(observer.Queue, "#", crosscutting.TopicMeasurements.String())
	return observer
}

func HandleWebsocket(c echo.Context) error {
	log.Info("#### ENTERING ####")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		log.Error(err.Error())
	}
	defer conn.Close()

	observer := StartWsTemperatureObserver()
	defer observer.Channel.Close()

	states := make(chan SocketState)

	go readMessagesFromClient(conn, states)
	go writeMessagesToClient(conn, observer, states)

	for {
		data := <-states
		if data == CLOSED {
			break
		}
	}

	log.Info("#### LEAVING WEBSOCKET ####")
	return nil
}

func readMessagesFromClient(conn *websocket.Conn, c chan SocketState) {
	for {
		var websocketEvent WebSocketEvent
		_, msg, err := conn.ReadMessage()
		err = json.Unmarshal(msg, &websocketEvent)
		if err != nil || websocketEvent.Close {
			c <- CLOSED
			return
		}
		c <- RECEIVING
	}
}

func writeMessagesToClient(ws *websocket.Conn, observer *observer.Observer, c chan SocketState) {

	deliveries := observer.Observe()

	for range deliveries {
		event := RefreshStateEvent{
			Source: Temperature.String(),
		}
		_, _ = json.Marshal(event)
		err := ws.WriteJSON(event)
		if err != nil {
			c <- CLOSED
			return
		} else {
			c <- RECEIVING
		}
	}

}
