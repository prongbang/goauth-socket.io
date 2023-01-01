package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/labstack/echo/v4"
)

const (
	deviceNamespace = "/device"
)

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

type Device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	// Socket ------------------------------------------------------------------------------
	io := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	io.OnConnect("/", AuthorizeConnection(func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	}))

	// Device namespace
	// https://socket.io/docs/v4/namespaces
	io.OnEvent(deviceNamespace, "subscribe", func(s socketio.Conn, payload any) any {

		id := payload.(string)
		room := fmt.Sprintf("device/%s", id)
		joined := io.JoinRoom(deviceNamespace, room, s)
		fmt.Println("room:", room, "joined:", joined)

		return payload.(string)
	})

	io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := io.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer io.Close()

	// Mqtt --------------------------------------------------------------------------------
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://broker.emqx.io:1883").
		SetClientID("emqx_test_client").
		SetUsername("emqx_test").
		SetPassword("emqx_test")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))

		// Parse
		d := Device{}
		json.Unmarshal(msg.Payload(), &d)

		// Broadcast
		room := fmt.Sprintf("device/%s", d.ID)
		event := room
		send := io.BroadcastToRoom(deviceNamespace, room, event, d)
		fmt.Println("send:", send)
	})
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("Mqtt connected")
	}

	// Subscribe to a topic
	if token := c.Subscribe("device", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	// Echo --------------------------------------------------------------------------------
	e := echo.New()
	e.Static("/", "./public")

	e.GET("/device/publish", func(c echo.Context) error {

		id := "1e4832e7-1ffa-4cf4-b9d9-0b8eff286c52"
		room := fmt.Sprintf("device/%s", id)
		event := room

		data := echo.Map{
			"id":   id,
			"name": "Temp",
		}

		send := io.BroadcastToRoom(deviceNamespace, room, event, data)
		fmt.Println("send:", send)

		return c.JSON(http.StatusOK, data)
	})

	e.Any("/socket.io/", func(context echo.Context) error {
		io.ServeHTTP(context.Response(), context.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":3000"))
}
