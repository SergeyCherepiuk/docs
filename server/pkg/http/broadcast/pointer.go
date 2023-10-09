package broadcast

import (
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/http/ws"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var pointerRoom = ws.NewRoom()

func Pointer(c echo.Context) error {
	websocket.Server{Handler: websocket.Handler(func(wsc *websocket.Conn) {
		pointerRoom.EnterAndListen(wsc, ws.Listener{
			OnEnter:   func() {},
			OnMessage: func(message any) any { return message },
			OnExit:    func() {},
		})
	})}.ServeHTTP(c.Response(), c.Request())
	return c.NoContent(http.StatusSwitchingProtocols)
}
