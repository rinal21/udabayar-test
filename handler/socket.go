package handler

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/googollee/go-socket.io"
)

type SocketServerHandler struct {
	Server *socketio.Server
	Socket socketio.Socket
}

func (h *SocketServerHandler) SocketServer(c *gin.Context) {
	h.Server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("transaction")

		h.Socket = so

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	h.Server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := c.Request.Header.Get("Origin")
	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	h.Server.ServeHTTP(c.Writer, c.Request)
}

func NewSocketServerHandler() *SocketServerHandler {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	return &SocketServerHandler{
		Server: server,
	}
}
