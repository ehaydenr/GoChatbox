package chat

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Id           int64
	Name         string
	Picture      string
	Conn         *websocket.Conn
	Server       *Server
	OutgoingChan chan *OutgoingMessage
}

// Client Constructor
func NewClient(ws *websocket.Conn, server *Server) *Client {
	outgoingChan := make(chan *OutgoingMessage)
	return &Client{
		Id:           time.Now().UnixNano(),
		Name:         "",
		Picture:      "",
		Conn:         ws,
		Server:       server,
		OutgoingChan: outgoingChan}
}

// Initiate Listening to Client incoming
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead() // Not go routine to block Listen from returning
}

// Write to client
func (c *Client) Write(msg *OutgoingMessage) {
	c.OutgoingChan <- msg
}

// Destroy Client
func (c *Client) Destroy() {
	c.Server.RemoveClient(c)
	c.Conn.Close()
	log.Printf("Destroyed: %v", c)
}

// Listen for Messages to send out to client
func (c *Client) listenWrite() {
	for {
		msg := <-c.OutgoingChan
		websocket.JSON.Send(c.Conn, msg)
	}
}

func (c *Client) digestToken(token string) bool {
	res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token)
	if err != nil {
		c.Destroy()
		return false
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Destroy()
		return false
	}

	var m map[string]string
	json.Unmarshal(body, &m)
	c.Name = m["name"]
	c.Picture = m["picture"]
	return true
}

func interpretIncomingChat(c *Client, msg Message) {
	orig := &Origin{
		Name:    c.Name,
		Picture: c.Picture,
	}
	out := &OutgoingMessage{
		Origin: orig,
	}
	out.Message.Content = msg.Content
	out.Message.Operation = msg.Operation
	c.Server.Broadcast(out)
}

func interpretIncomingToken(c *Client, token string) {
	res := c.digestToken(token)
	var content string

	if res {
		content = "true"
	} else {
		content = "false"
	}

	m := &Message{
		Content:   content,
		Operation: "token_verified",
	}
	if err := websocket.JSON.Send(c.Conn, m); err != nil {
		c.Destroy()
	}
}

// List for Messages coming from the client
func (c *Client) listenRead() {
	for {
		var msg Message
		err := websocket.JSON.Receive(c.Conn, &msg)
		if err != nil {
			c.Destroy()
			return
		}

		if msg.Operation == "chat" {
			interpretIncomingChat(c, msg)
		} else if msg.Operation == "giveToken" {
			interpretIncomingToken(c, msg.Content)
		}
	}
}
