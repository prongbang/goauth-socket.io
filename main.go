package main

import (
	"fmt"
	"log"
	"net/http"

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

func main() {

	e := echo.New()
	e.Static("/", "./public")

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

	io.OnEvent("/chat", "message", func(s socketio.Conn, msg any) any {

		s.Emit("reply", msg)

		s.SetContext(msg)

		return "recv " + msg.(string)
	})

	io.OnEvent("/device", "subscribe", func(s socketio.Conn, msg any) any {

		id := msg.(string)
		room := fmt.Sprintf("device/%s", id)
		joined := io.JoinRoom(deviceNamespace, room, s)
		fmt.Println("room:", room, "joined:", joined)

		return msg.(string)
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

	e.GET("/device/payload", func(c echo.Context) error {

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
