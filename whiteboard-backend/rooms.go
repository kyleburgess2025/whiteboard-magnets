package main

import (
	"encoding/json"
)

type Room struct {
	name 	   string
	words 	   map[string]Word
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan MessageWithClient
	emit 	   chan []byte
}

func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		words:      make(WordMap),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan MessageWithClient),
		emit: 	 	make(chan []byte),
	}
}

func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message)

		case message := <-room.emit:
			room.emitToClientsInRoom(message)
		}

	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message MessageWithClient) {
	jsonMessage, _ := json.Marshal(message.Message)
	for client := range room.clients {
		if client != message.Client {
			client.send <- jsonMessage
		}
	}
}

func (room *Room) emitToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := MessageWithClient{}
	message.Message.Type = "join"
	message.Message.ClientId = client.getName()
	message.Client = client
	room.broadcastToClientsInRoom(message)
}

