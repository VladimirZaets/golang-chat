package client

import (
	"fmt"
	"log"

	"github.com/VladimirZaets/vz-chat/services/message"
	"github.com/VladimirZaets/vz-chat/subscription"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientInterface interface {
	GetName() string
	GetSessionId() *uuid.UUID
	GetConnection() *websocket.Conn
}

type Client struct {
	name      string
	sessionid *uuid.UUID
	conn      *websocket.Conn
}

func NewClient(name string, conn *websocket.Conn) *Client {
	sessionid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("Ann error on UIID creattion: %w", err)
	}

	return &Client{
		name:      name,
		conn:      conn,
		sessionid: &sessionid,
	}
}

func (cl *Client) Name() string {
	return cl.name
}

func (cl *Client) GetSessionId() *uuid.UUID {
	return cl.sessionid
}

func (cl *Client) Conn() *websocket.Conn {
	return cl.conn
}

func (cl *Client) Listen(ch *subscription.CommunicationChannels) {
	for {
		message := &message.Message{
			Name: cl.name,
		}
		err := cl.conn.ReadJSON(message)
		fmt.Println("Receive client message", message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		fmt.Println("Input msg", message.Data)
		ch.Messages <- message
	}
}
