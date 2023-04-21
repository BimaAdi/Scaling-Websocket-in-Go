package memoryadapter

import (
	"strings"

	wsc "github.com/BimaAdi/use-websocket-scaler/websocketscaler"
)

type MemoryAdapter struct {
	channel         chan string
	WebsocketServer wsc.WebsocketServer
}

func NewMemoryAdapter() *MemoryAdapter {
	channel := make(chan string)
	return &MemoryAdapter{
		channel: channel,
	}
}

func (ma *MemoryAdapter) Prelude() {
	ma.channel = make(chan string)
	go ma.Subscribe()
}

func (ma *MemoryAdapter) AddWebsokcetServer(ws wsc.WebsocketServer) {
	ma.WebsocketServer = ws
}

func (ma MemoryAdapter) Subscribe() {
	isChannelOpen := true
	var v string
	for isChannelOpen {
		v, isChannelOpen = <-ma.channel
		action, socketId, message := messageParser(v)
		if action == "toSocketId" {
			ma.WebsocketServer.ToSocketId(socketId, message)
		} else if action == "Broadcast" {
			ma.WebsocketServer.Broadcast(message)
		}
	}
}

// return action, socketId, message
func messageParser(message string) (string, string, string) {
	words := strings.Split(message, `|`)
	if words[0] == "toSocketId" {
		if len(words) > 3 {
			sanitize_message := ""
			for ix, x := range words {
				if ix != 0 && ix != 1 {
					sanitize_message = sanitize_message + x
				}
			}
			words[2] = sanitize_message
		}
		return words[0], words[1], words[2]
	} else if words[0] == "Broadcast" {
		if len(words) > 2 {
			sanitize_message := ""
			for ix, x := range words {
				if ix != 0 {
					sanitize_message = sanitize_message + x
				}
			}
			words[1] = sanitize_message
		}
		return words[0], "", words[1]
	}
	return "", "", ""
}

// message format toSocketId|{socketId}|{message}
func (ma MemoryAdapter) ToSocketId(socketId string, message string) {
	ma.channel <- "toSocketId|" + socketId + "|" + message
}

// message format Broadcast|{message}
func (ma MemoryAdapter) Broadcast(message string) {
	ma.channel <- "Broadcast|" + message
}
