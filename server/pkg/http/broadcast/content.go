package broadcast

import (
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/http/ws"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var contentRoom = ws.NewRoom()
var lastContentMessage any = ""

func Content(c echo.Context) error {
	websocket.Server{Handler: func(wsc *websocket.Conn) {
		contentRoom.EnterAndListen(wsc, ws.Listener{
			OnEnter: func() { websocket.Message.Send(wsc, lastContentMessage) },
			OnMessage: func(message any) any {
				lastContentMessage = message
				return message
			},
			OnExit: func() {},
		})
	}}.ServeHTTP(c.Response(), c.Request())
	return c.NoContent(http.StatusSwitchingProtocols)
}
