package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type NavType int

const (
	Select NavType = iota
	Create
	Join
	Exit
)

var actionsMap = map[string]NavType{
	"select": Select,
	"create": Create,
	"join":   Join,
	"exit":   Exit,
}

type Message struct {
	Name string    `json:"name"`
	Time time.Time `json:"created_at"`
	Data string    `json:"data"`
}

type NavigationDto struct {
	selectedCh string
	list       []string
	name       string
}

var navDTO *NavigationDto = &NavigationDto{}
var channelsInst = NewChannels()

func mapActions(action NavType) {
	switch action {
	case Join:
		handleJoinChannelAction()
	case Select:
		handleSelectChannelAction()
	case Create:
		handleCreateChannelAction()
	case Exit:
	}
}

func handleSelectChannelAction() {
	list, err := channelsInst.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oops, something went frong on channel selection: %s", err)
	}
	if len(list) == 0 {
		fmt.Println("No created channels yet")
		return
	}

	navDTO.list = list
	navDTO.selectedCh = getSelectedChannel(list)
	mapActions(Join)
}

func handleJoinChannelAction() {
	navDTO.name = getUserName()
	u := url.URL{
		Scheme: "ws",
		Host:   *addr, Path: "/connect",
		RawQuery: fmt.Sprintf("channel=%s&name=%s", navDTO.selectedCh, navDTO.name)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}

	go tt(c)
	for {

		ms := &Message{}
		//var raw map[string]interface{}

		err := c.ReadJSON(ms)

		fmt.Println(ms)
		//_, msg, err := c.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message.Data)
	}
	navDTO.selectedCh = ""
}

func tt(c *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		fmt.Println("client", input)
		message := &Message{
			Data: input,
			Name: navDTO.name,
		}
		c.WriteJSON(message)
	}
}

func handleCreateChannelAction() {
	cn := getChannelName()
	err := channelsInst.Create(cn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oops, something went frong on channel creation: %d", err)
	}
	navDTO.selectedCh = cn
	mapActions(Join)
}

func GetTopLevelActions() NavType {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("do you want to \"select\" or \"create\" channel")
	for {
		scanner.Scan()
		input := scanner.Text()

		if val, ok := actionsMap[input]; ok {
			return val
		} else {
			fmt.Println("Unknown parameter: Please choose one of avaliable option")
		}
	}
}
func getUserName() string {
	scanner := bufio.NewScanner(os.Stdin)
	message := "Please enter username"
	fmt.Println(message)
	for {
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			return input
		} else {
			fmt.Println(message)
		}
	}
}

func getChannelName() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("please provide the name for the channel")
	for {
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			return input
		} else {
			fmt.Println("Please enter channel name")
		}
	}
}

func getSelectedChannel(chs []string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Please select channel from the list: %s", strings.Join(chs, ", "))
	for {
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			return input
		} else {
			fmt.Println("Please select the channel")
		}
	}
}
