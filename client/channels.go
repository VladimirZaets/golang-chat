package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type channels struct{}

func NewChannels() *channels {
	return &channels{}
}

func (ch *channels) Create(cn string) error {
	postBody, _ := json.Marshal(map[string]string{"name": cn})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8080/channels", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return err
	}
	defer resp.Body.Close()
	responseBodyReader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return err
	}

	if resp.StatusCode == 400 {
		responseBody := make(map[string]string)
		json.Unmarshal(responseBodyReader, &responseBody)
		log.Fatalf("An Error Occured %v", err)
		//implement http error type
		return errors.New(responseBody["error"])
	}
	return nil
}

func (ch *channels) List() ([]string, error) {
	type responseBody struct {
		Data []string `json: data`
	}
	resp, err := http.Get("http://localhost:8080/channels")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rb := new(responseBody)
	err = json.NewDecoder(resp.Body).Decode(&rb)
	if err != nil {
		return nil, err
	}
	return rb.Data, nil
}
