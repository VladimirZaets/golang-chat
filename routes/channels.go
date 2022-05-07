package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/VladimirZaets/vz-chat/services/channels"
)

type createChannelBody struct {
	Name string
}

func ChannelsPost(w http.ResponseWriter, r *http.Request, ch *channels.Channels) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	channelBody := new(createChannelBody)
	err = json.Unmarshal(body, &channelBody)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = ch.Add(channelBody.Name)
	response := make(map[string]string)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		response["error"] = err.Error()
		responseJson, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJson)
		return
	}
	response["message"] = "ok"
	fmt.Println(ch.List())
	response["channels"] = strings.Join(ch.List(), ",")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	responseJson, _ := json.Marshal(response)
	w.Write(responseJson)
}

func ChannelsGet(w http.ResponseWriter, r *http.Request, ch *channels.Channels) {
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string][]string)
	response["data"] = ch.List()
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
	}
	w.Write(responseJson)
}

func ChannelsPut(w http.ResponseWriter, r *http.Request) {

}

func ChannelsDelete(w http.ResponseWriter, r *http.Request) {

}
