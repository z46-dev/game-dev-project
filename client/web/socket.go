package web

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/z46-dev/golog"
)

type Socket struct {
	connection        *websocket.Conn
	Open, onCloseDone bool
	mu                sync.Mutex
	OnClose           func()
	Logger            *golog.Logger
}

func Connect(url string) (socket *Socket, err error) {
	var conn *websocket.Conn
	if conn, _, err = websocket.DefaultDialer.Dial(url, nil); err != nil {
		return nil, err
	}

	socket = NewSocket(conn)
	socket.Logger.Info("Connected to server")
	return socket, nil
}

func NewSocket(connection *websocket.Conn) (socket *Socket) {
	socket = &Socket{
		connection: connection,
		Open:       true,
		Logger:     golog.New().Prefix("[SOCKET]", golog.BoldGreen).Timestamp(),
	}

	var closeHandler func(int, string) error = connection.CloseHandler()
	connection.SetCloseHandler(func(code int, text string) error {
		var err error
		if closeHandler != nil {
			err = closeHandler(code, text)
		}

		socket.handleClose()
		return err
	})

	return socket
}

func (socket *Socket) ReadAndValidate() (message []byte, err error) {
	var messageType int
	if messageType, message, err = socket.connection.ReadMessage(); err == nil && messageType != websocket.BinaryMessage {
		err = fmt.Errorf("invalid message type")
	}

	return
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

func (socket *Socket) Write(message []byte) error {
	socket.mu.Lock()
	defer socket.mu.Unlock()
	return socket.connection.WriteMessage(websocket.BinaryMessage, message)
}

func (socket *Socket) Close() error {
	socket.handleClose()
	return socket.connection.Close()
}

func (socket *Socket) handleClose() {
	socket.Open = false
	if socket.OnClose != nil && !socket.onCloseDone {
		socket.onCloseDone = true
		socket.OnClose()
	}
}
