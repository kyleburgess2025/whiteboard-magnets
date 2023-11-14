package main

import (
	"log"
	"net/http"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
		},
}

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms   map[*Room]bool
	Name  string
}

type Word struct {
	Word string `json:"word"`
	XValue int `json:"xValue"`
	YValue int `json:"yValue"`
	DeltaX int `json:"deltaX"`
	DeltaY int `json:"deltaY"`
	Id string `json:"id"`
}

type MessageWithClient struct {
	Message Message `json:"message"`
	Client *Client `json:"client"`
}

type Message struct {
	Word Word `json:"word,omitempty"`
	ClientId string `json:"clientId,omitempty"`
	Type string `json:"type"`
	Target string `json:"target,omitempty"`
}

type WordListMessage struct {
	Type string `json:"type"`
	Words []Word `json:"words"`
}

var WordMap map[string]Word = make(map[string]Word)

func newClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		Name: name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}

func (client *Client) getName() string {
	return client.Name
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		var message Message
		err = json.Unmarshal(jsonMessage, &message)
		if err != nil {
			log.Printf("error unmarshalling json: %v", err)
			break
		}

		if message.Type == "add" || message.Type == "move" {
			roomName := message.Target
			if room := client.wsServer.findRoomByName(roomName); room != nil {
				room.WordMap[message.Word.Id] = message.Word
				room.broadcast <- MessageWithClient{message, client}
			}
		} else if message.Type == "delete" {
			roomName := message.Target
			if room := client.wsServer.findRoomByName(roomName); room != nil {
				delete(room.WordMap, message.Word.Id)
				room.broadcast <- MessageWithClient{message, client}
			}
		} else if message.Type == "get" {
			roomName := message.Target
			if room := client.wsServer.findRoomByName(roomName); room != nil {
				wordList := WordListMessage{}
				wordList.Type = "get"
				for _, word := range WordMap {
					wordList.Words = append(wordList.Words, word)
				}
				jsonString, err := json.Marshal(wordList)
				if err != nil {
					log.Printf("error marshalling json: %v", err)
					break
				}
				client.send <- jsonString
				continue
			}
		} else if message.Type == "join" {
			client.handleJoinRoomMessage(message)
		}
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, name[0])

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByName(message.Message)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}