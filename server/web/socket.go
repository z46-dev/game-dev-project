package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/z46-dev/golog"
)

func GetUpgrader() *websocket.Upgrader {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &upgrader
}

var (
	socketId int                 = 0
	upgrader *websocket.Upgrader = GetUpgrader()
)

// Socket structure

type Socket struct {
	connection        *websocket.Conn
	Open, onCloseDone bool
	mu                sync.Mutex
	ID                int
	OnClose           func()
	Logger            *golog.Logger
}

func NewSocket(connection *websocket.Conn) (socket *Socket) {
	// Prepare
	socket = new(Socket)
	var closeHandler func(int, string) error = connection.CloseHandler()

	connection.SetCloseHandler(func(code int, text string) error {
		var err error = nil

		if closeHandler != nil {
			err = closeHandler(code, text)
		}

		socket.Open = false

		if socket.OnClose != nil && !socket.onCloseDone {
			socket.onCloseDone = true
			socket.OnClose()
		}

		return err
	})

	// Create socket
	socket.connection = connection
	socket.Open = true
	socket.mu = sync.Mutex{}
	socket.ID = socketId
	socketId++

	return
}

func (socket *Socket) ReadAndValidate() (message []byte, err error) {
	var messageType int
	if messageType, message, err = socket.connection.ReadMessage(); err == nil && messageType != websocket.BinaryMessage {
		err = fmt.Errorf("invalid message type")
	}

	return
}

func (socket *Socket) handleClose() {
	socket.Open = false
	if socket.OnClose != nil && !socket.onCloseDone {
		socket.onCloseDone = true
		socket.OnClose()
	}
}

func (socket *Socket) InitiateUpdateLoop(callback func(message []byte)) {
	var (
		err     error
		message []byte
	)

	for {
		if message, err = socket.ReadAndValidate(); err != nil || !socket.Open {
			if err != nil {
				socket.handleClose()
			}
			break
		}

		callback(message)
	}
}

func (socket *Socket) CreateLoop(callback func(), t time.Duration) {
	var ticker *time.Ticker = time.NewTicker(t)

	defer ticker.Stop()

	for {
		<-ticker.C

		if !socket.Open {
			break
		}

		callback()
	}
}

func (socket *Socket) Write(message []byte) error {
	socket.mu.Lock()
	defer socket.mu.Unlock()
	return socket.connection.WriteMessage(websocket.BinaryMessage, message)
}

func (socket *Socket) Close() error {
	socket.handleClose()

	return socket.connection.Close()
}

// Upgrading and other stuff
func Upgrade(w http.ResponseWriter, r *http.Request) (socket *Socket, err error) {
	var connection *websocket.Conn
	if connection, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	socket = NewSocket(connection)
	return
}
