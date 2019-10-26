package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/maxzurawski/temperaturearchive/publishers"

	"github.com/maxzurawski/utilities/rabbit/crosscutting"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/maxzurawski/temperaturearchive/config"
	"github.com/maxzurawski/utilities/rabbit/observer"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartWsTemperatureObserver() *observer.Observer {
	observer := config.TemperaturearchiveConfig().InitObserver()
	observer.DeclareTopicExchange(crosscutting.TopicMeasurements.String())
	observer.BindQueue(observer.Queue, crosscutting.RoutingKeyTemperatureMeasurement.String(), crosscutting.TopicMeasurements.String())
	return observer
}

func HandleWebsocket(c echo.Context) error {
	log.Info("#### ENTERING ####")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		publishers.Logger().Error(uuid.New().String(), "", "websocket upgrade failed", err.Error())
		log.Error(err.Error())
	}
	defer conn.Close()

	observer := StartWsTemperatureObserver()
	defer observer.Channel.Close()

	states := make(chan SocketState)
	go readMessagesFromClient(conn, states)
	go communicateWithClient(conn, observer, states)

	for {
		data := <-states
		// NOTE: until go-channel does not receive CLOSED state -> wait for values on the channel
		if data == CLOSED {
			break
		}
	}

	log.Info("#### LEAVING WEBSOCKET ####")
	return nil
}

func readMessagesFromClient(conn *websocket.Conn, c chan SocketState) {

	// NOTE: until client does not send WebSocketEvent.Close -> assume websocket should stay opened
	for {
		var websocketEvent WebSocketEvent
		_, msg, err := conn.ReadMessage()
		err = json.Unmarshal(msg, &websocketEvent)

		if err != nil {
			websocketEvent.Close = true
		}
		if websocketEvent.Close {
			c <- CLOSED
			return
		}
		c <- RECEIVING
	}
}

func communicateWithClient(ws *websocket.Conn, observer *observer.Observer, c chan SocketState) {

	deliveries := observer.Observe()

	// NOTE: Until there are messages in temperature measurement/queue is not closed -> observe queue
	// and send RefreshStateEvent.Source = 'temperature' to client
	for range deliveries {
		event := RefreshStateEvent{
			Source: Temperature.String(),
		}
		_, _ = json.Marshal(event)
		err := ws.WriteJSON(event)
		if err != nil {
			c <- CLOSED
			return
		}
	}

}
