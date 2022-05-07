package main

import (
	"flag"
	"log"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	mapActions(GetTopLevelActions())

	// u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	// log.Printf("Connecting to %s", u.String())

	// c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	// if err != nil {
	// 	log.Fatal("dial:", err)
	// }

	// c.WriteMessage(websocket.TextMessage, []byte("Hi"))
	//communicationChannels := make(map[string](chan string))

	// for {
	// 	_, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// }
}
