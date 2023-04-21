package websocketscaler

import (
	"encoding/json"
	"math/rand"
)

type Scaler struct {
	Namespaces      map[string]Namespace
	Adapter         Adapter
	WebsocketServer WebsocketServer
}

type WebsocketServer interface {
	AddAdapter(adapter Adapter)
	ToSocketId(socketId string, message string)
	// Emit(message string)
	Broadcast(message string)
}

type Adapter interface {
	Prelude()
	AddWebsokcetServer(ws WebsocketServer)
	ToSocketId(socketId string, message string)
	Broadcast(message string)
}

func NewScaler(wsserver WebsocketServer, adapter Adapter) Scaler {
	wsserver.AddAdapter(adapter)
	adapter.AddWebsokcetServer(wsserver)
	adapter.Prelude()
	return Scaler{
		Namespaces:      map[string]Namespace{},
		Adapter:         adapter,
		WebsocketServer: wsserver,
	}
}

func (scaler *Scaler) Of(namespace_name string, namespace Namespace) {
	scaler.Namespaces[namespace_name] = namespace
}

type Namespace struct {
	Event map[string]func(Adapter, string)
}

func InitNamespace() Namespace {
	return Namespace{
		Event: map[string]func(Adapter, string){},
	}
}

func (namespace *Namespace) On(event_name string, method func(Adapter, string)) {
	namespace.Event[event_name] = method
}

func (scaler Scaler) Run(namespace_name string, event_name string, message string, socket_id string) {
	namespace, is_exists := scaler.Namespaces[namespace_name]
	if is_exists {
		event, is_exists := namespace.Event[event_name]
		if is_exists {
			event(scaler.Adapter, message)
		}
	}
}

type MessageParserFormat struct {
	Namespace string
	Event     string
	Message   string
}

func (scaler Scaler) MessageParser(message string) (MessageParserFormat, error) {
	data := MessageParserFormat{}
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		return MessageParserFormat{}, err
	}
	return data, nil
}

const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateSocketId() string {
	b := make([]byte, 30)
	for i := range b {
		b[i] = LetterBytes[rand.Int63()%int64(len(LetterBytes))]
	}
	return string(b)
}
