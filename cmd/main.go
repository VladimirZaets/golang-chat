package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/VladimirZaets/vz-chat/routes"
	"github.com/VladimirZaets/vz-chat/services/channels"
	"github.com/VladimirZaets/vz-chat/services/client"
	"github.com/VladimirZaets/vz-chat/subscription"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var ch = channels.NewChannels()
var sm = subscription.NewSubscriptionManager(ch)

func main() {
	flag.Parse()
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/channels", channelsRoutes)
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/c", stop)
	log.Println("Server running")
	http.ListenAndServe(*addr, nil)
}

func stop(w http.ResponseWriter, r *http.Request) {
	c := ch.Get("vova")
	c.Close <- true
}

func channelsRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		routes.ChannelsGet(w, r, ch)
	case http.MethodPost:
		routes.ChannelsPost(w, r, ch)
	case http.MethodDelete:
		routes.ChannelsDelete(w, r)
	case http.MethodPut:
		routes.ChannelsPut(w, r)
	}
}

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	cn := r.URL.Query().Get("channel")
	un := r.URL.Query().Get("name")
	cl := client.NewClient(un, c)
	sm.Subscribe(cn, cl)
	fmt.Println("ENd")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	qs := r.URL.Query()
	fmt.Println("sukabl9t", qs["channel"])
	//communicationChannels := make(map[string](chan string))
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
