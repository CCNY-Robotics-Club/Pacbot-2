package webserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Protect global websocket values from race conditions with a
// mutex - it will lock the resources until they are successfully written to
var muWS sync.Mutex

/*
Websockets are the way that the server will communicate with
the browser clients over the LAN. When connecting to a web-socket,
we connect normally over HTTP and then upgrade the
connection upon agreement between the server and client.
*/
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
This handler makes sure that once we connect to the websocket,
all communication goes smoothly.
*/
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrades the connection, and quits if it didn't work out.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a websocket session object
	ws := newWebSession(conn)

	// Register the websocket session
	ws.register()

	// Close the connection at the end of the function, or if
	// something goes wrong.
	defer func() {
		ws.unregister()
	}()

	// go ws.sendLoop()

	// Goroutine: handles reads for the websocket session
	go ws.readLoop()

	// Goroutine: handles writes for the websocket session
	go ws.writeLoop()

	// Wait until the quit channel has items before stopping the program
	<-ws.quitCh
}
