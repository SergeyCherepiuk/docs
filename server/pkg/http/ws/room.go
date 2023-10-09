package ws

import (
	"golang.org/x/exp/maps"
	"golang.org/x/net/websocket"
)

type room struct {
	conns map[*websocket.Conn]struct{}
}

type Listener struct {
	OnEnter   func()
	OnMessage func(message any) any
	OnExit    func()
}

func NewRoom() *room {
	return &room{
		conns: make(map[*websocket.Conn]struct{}),
	}
}

func (r room) EnterAndListen(conn *websocket.Conn, listener Listener) {
	r.conns[conn] = struct{}{}
	defer delete(r.conns, conn)
	defer conn.Close()

	listener.OnEnter()

	for {
		var message string
		if err := websocket.Message.Receive(conn, &message); err != nil {
			break
		}

		transformedMessage := listener.OnMessage(message)

		for _, c := range r.GetConns() {
			if c != conn {
				websocket.Message.Send(c, transformedMessage)
			}
		}
	}

	listener.OnExit()
}

func (r room) GetConns() []*websocket.Conn {
	return maps.Keys(r.conns)
}
