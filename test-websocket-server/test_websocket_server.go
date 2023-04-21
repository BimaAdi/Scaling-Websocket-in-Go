package testwebsocketserver

type LogMessageStruct struct {
	Action   string
	SocketId string
	Message  string
}

type TestWebsocketServer struct {
	Log []LogMessageStruct
}

func InitTestWebsocketServer() TestWebsocketServer {
	return TestWebsocketServer{
		Log: []LogMessageStruct{},
	}
}

func (socket *TestWebsocketServer) Emit(message string) {
	socket.Log = append(socket.Log, LogMessageStruct{
		Action:   "Emit",
		SocketId: "",
		Message:  message,
	})
}

func (socket *TestWebsocketServer) Broadcast(message string) {
	socket.Log = append(socket.Log, LogMessageStruct{
		Action:   "Broadcast",
		SocketId: "",
		Message:  message,
	})
}

func (socket *TestWebsocketServer) ToSocketId(socketId string, message string) {
	socket.Log = append(socket.Log, LogMessageStruct{
		Action:   "ToSocketId",
		SocketId: socketId,
		Message:  message,
	})
}
