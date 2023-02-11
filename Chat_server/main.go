package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var Conn *websocket.Conn
var clients []websocket.Conn

func main() {

	router := gin.Default()

	router.Static("/css", "Frontend/css")
	router.LoadHTMLGlob("Frontend/*.html")

	router.GET("/echo", Chat)

	router.GET("/", MainForm) //Главная страница 

	router.Run("localhost:8080")
}

func Chat (c *gin.Context) {
	Conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clients = append(clients, *Conn)

	for {
		// Read message from browser
		msgType, msg, err := Conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", Conn.RemoteAddr(), string(msg))

		for _, client := range clients {
			// Write message back to browser
			if err = client.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
		
	}
}

func MainForm(c *gin.Context) {
	c.HTML(200, "index", nil)
}