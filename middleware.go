package main

import (
	"errors"
	"fmt"
	"net/url"

	socketio "github.com/googollee/go-socket.io"
)

var ErrUnauthorization = errors.New("Unauthorization")

func AuthorizeConnection(next func(socketio.Conn) error) func(socketio.Conn) error {
	return func(conn socketio.Conn) error {
		q, _ := url.ParseQuery(conn.URL().RawQuery)
		token := q.Get("token")

		payload, err := VerifyToken(token)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return ErrUnauthorization
		}
		fmt.Println("payload:", payload)

		return next(conn)
	}
}
