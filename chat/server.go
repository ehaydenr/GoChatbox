package chat

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Server struct {
	route            string
	clients          map[int64]*Client
	addChannel       chan *Client
	removeChannel    chan *Client
	broadcastChannel chan *OutgoingMessage
	doneChannel      chan bool
}

func NewServer(route string) *Server {
	clients := make(map[int64]*Client)
	addChannel := make(chan *Client)
	removeChannel := make(chan *Client)
	broadcastChannel := make(chan *OutgoingMessage)
	doneChannel := make(chan bool)

	return &Server{
		route,
		clients,
		addChannel,
		removeChannel,
		broadcastChannel,
		doneChannel,
	}
}

func (s *Server) AddClient(c *Client) {
	s.addChannel <- c
}
func (s *Server) RemoveClient(c *Client) {
	s.removeChannel <- c
}
func (s *Server) Broadcast(m *OutgoingMessage) {
	s.broadcastChannel <- m
}

func (s *Server) Listen() {
	onConnection := func(ws *websocket.Conn) {
		defer ws.Close()
		client := NewClient(ws, s)
		s.AddClient(client)
		client.Listen()
	}
	http.Handle(s.route, websocket.Handler(onConnection))

	for {
		select {
		case client := <-s.addChannel:
			s.clients[client.Id] = client
		case client := <-s.removeChannel:
			delete(s.clients, client.Id)
			log.Printf("Number of Clients: %d", len(s.clients))
		case msg := <-s.broadcastChannel:
			for _, c := range s.clients {
				c.Write(msg)
			}
		}
	}
}
