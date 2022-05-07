package channels

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VladimirZaets/vz-chat/subscription"
	"github.com/gorilla/websocket"
)

type Channels struct {
	channels map[string](*subscription.CommunicationChannels)
}

func NewChannels() *Channels {
	return &Channels{
		channels: make(map[string](*subscription.CommunicationChannels)),
	}
}

func (ch *Channels) List() []string {
	keys := make([]string, len(ch.channels))
	i := 0
	for k := range ch.channels {
		keys[i] = k
		i++
	}
	return keys
}

func (ch *Channels) Get(name string) *subscription.CommunicationChannels {
	return ch.channels[name]
}

func (ch *Channels) Add(name string) error {
	if ch.channels[name] != nil {
		return fmt.Errorf("the channel with the name %s is already created", name)
	}
	ch.channels[name] = &subscription.CommunicationChannels{
		Messages: make(chan subscription.MessageInterface),
		Close:    make(chan bool),
		Clients:  make([]subscription.ClientInterface, 0),
	}
	return nil
}

func (ch *Channels) Delete(name string) {

}

func (ch *Channels) Listen(name string, clients *[]subscription.ClientInterface) {
	broadcast := ch.channels[name]
	for {
		select {
		case <-broadcast.Close:
			for _, subsc := range *clients {
				subsc.Conn().WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "The channel is deleted"))
			}
			return
		case msg := <-broadcast.Messages:
			msgJson, err := json.Marshal(msg)
			fmt.Println(string(msgJson))
			if err != nil {
				log.Fatalf("json marsharing error: %v", err)
			}
			for _, subsc := range *clients {
				subsc.Conn().WriteJSON(msg)
			}
		}
	}
}
