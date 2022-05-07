package subscription

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type MessageInterface interface {
	Message() string
}

type CommunicationChannels struct {
	Messages chan MessageInterface
	Close    chan bool
	Clients  []ClientInterface
}

type SourceInterface interface {
	Get(name string) *CommunicationChannels
	Listen(name string, clients *[]ClientInterface)
}

type ClientInterface interface {
	Name() string
	Listen(*CommunicationChannels)
	Conn() *websocket.Conn
}

type SubscriptionManager struct {
	source  SourceInterface
	clients map[string][]ClientInterface
}

func NewSubscriptionManager(source SourceInterface) *SubscriptionManager {
	return &SubscriptionManager{
		source:  source,
		clients: make(map[string][]ClientInterface),
	}
}

func (sm *SubscriptionManager) Subscribe(name string, client ClientInterface) {
	channel := sm.source.Get(name)
	fmt.Println("Add client")
	if len(channel.Clients) == 0 {
		fmt.Println("0 clients, add one")
		go sm.source.Listen(name, &channel.Clients)
	}
	channel.Clients = append(channel.Clients, client)
	client.Listen(channel) //Blocking function. Execution will continue when client disconected
	sm.Unsubscribe(name, client)
}

func (sm *SubscriptionManager) Unsubscribe(name string, client ClientInterface) {
	channel := sm.source.Get(name)
	chanClients := channel.Clients
	client.Conn().Close()
	fmt.Println("Removing clients", chanClients)
	for i, v := range chanClients {
		if v.Name() == client.Name() {
			chanClients[i] = chanClients[len(chanClients)-1]
			channel.Clients = chanClients[:len(chanClients)-1]
		}
	}
	fmt.Println("Removed clients", channel.Clients)

	if len(channel.Clients) == 0 {
		fmt.Println("0 clients, close channel")
		channel.Close <- true
	}
}
