package fiberwebsocketserver

import (
	wsc "github.com/BimaAdi/use-websocket-scaler/websocketscaler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ConnectionFormat struct {
	Room       string
	Namespace  string
	SocketId   string
	Connection *websocket.Conn
}

type FiberWebsocketServer struct {
	connections []ConnectionFormat
	scaler      *wsc.Scaler
	Adapter     wsc.Adapter
}

func NewFiberWebsocketServer() *FiberWebsocketServer {
	return &FiberWebsocketServer{
		connections: []ConnectionFormat{},
		scaler:      nil,
	}
}

func (fwc *FiberWebsocketServer) AddAdapter(adapter wsc.Adapter) {
	fwc.Adapter = adapter
}

func (fwc *FiberWebsocketServer) AddConnection(conn *websocket.Conn, room string, namespace string) {
	fwc.RemoveDisconnectedConnection()
	socket_id := wsc.GenerateSocketId()
	fwc.connections = append(fwc.connections, ConnectionFormat{
		Room:       room,
		Namespace:  namespace,
		SocketId:   socket_id,
		Connection: conn,
	})
}

func (fwc *FiberWebsocketServer) RemoveDisconnectedConnection() {
	newFwcConnections := []ConnectionFormat{}
	for _, conn := range fwc.connections {
		if conn.Connection.Conn != nil {
			newFwcConnections = append(newFwcConnections, conn)
		}
	}
	fwc.connections = newFwcConnections
}

func (fwc FiberWebsocketServer) ToSocketId(socketid string, message string) {
	selected_conn := []ConnectionFormat{}
	for _, conn := range fwc.connections {
		if conn.SocketId == socketid {
			selected_conn = append(selected_conn, conn)
		}
	}

	for _, conn := range selected_conn {
		conn.Connection.WriteMessage(1, []byte(message))
	}
}

func (fwc FiberWebsocketServer) Broadcast(message string) {
	for _, conn := range fwc.connections {
		conn.Connection.WriteMessage(1, []byte(message))
	}
}

func (fwc FiberWebsocketServer) GetSocketId(c *websocket.Conn) *ConnectionFormat {
	if c == nil {
		return nil
	}

	var selected_conn *ConnectionFormat = nil
	for _, conn := range fwc.connections {
		if c == conn.Connection {
			selected_conn = &conn
		}
	}
	return selected_conn
}

func (fwc *FiberWebsocketServer) Serve(scaler wsc.Scaler) func(c *fiber.Ctx) error {
	fwc.scaler = &scaler
	return websocket.New(func(c *websocket.Conn) {
		var (
			// mt  int
			msg []byte
			err error
		)
		namespace := c.Params("namespace")
		fwc.AddConnection(c, "", namespace)

		for {
			_, msg, err = c.ReadMessage()
			if err != nil {
				break
			}
			fwc.RemoveDisconnectedConnection()
			x := fwc.GetSocketId(c)

			mpf, err := scaler.MessageParser(string(msg))
			if err == nil {
				scaler.Run(mpf.Namespace, mpf.Event, mpf.Message, x.SocketId)
			}

			// if err = c.WriteMessage(mt, msg); err != nil {
			// 	log.Println("write:", err)
			// 	break
			// }
		}
	})
}
