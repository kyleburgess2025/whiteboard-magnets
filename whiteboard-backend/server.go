package main

import "encoding/json"

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan MessageWithClient
	emit 	   chan []byte
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan MessageWithClient),
		emit: 	 make(chan []byte),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.emit:
			server.emitToClients(message)
		
		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) emitToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) broadcastToClients(message MessageWithClient) {
    // Send the message to all clients except the sender
	jsonMessage, _ := json.Marshal(message.Message)
	for client := range server.clients {
		if client != message.Client {
			client.send <- jsonMessage
		}
	}
}