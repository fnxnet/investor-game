package main

import (
    "encoding/json"
    "errors"
    "fmt"
)

type handler interface {
    Handle(message []byte, broadcast chan []byte) error
}

type ErrorMessage struct {
    Action string `json:"action"`
    User   User   `json:"user"`
    Error  string `json:"error"`
}

func newErrorMessage(e string, u User) *ErrorMessage {
    return &ErrorMessage{
        Action: "error",
        User:   u,
        Error:  e,
    }
}

type MessageHandler struct {
    handlers map[string]handler
}

func (mh *MessageHandler) addHandler(name string, h handler) {
    if mh.handlers == nil {
        mh.handlers = make(map[string]handler)
    }
    mh.handlers[name] = h
}

func (mh *MessageHandler) HandleMessage(content []byte, broadcast chan []byte) {
    var e error
    message := &RequestMessage{}

    json.Unmarshal(content, message)

    if userManager.Validate(message.User) {

        switch message.Action {
        case "newOffer", "removeOffer", "acceptOffer":
            e = mh.handlers["offer"].Handle(content, broadcast)
        }

    } else {
        fmt.Println(e.Error())
        e = errors.New("Invalid user token")
    }

    if e != nil {
        em := newErrorMessage(e.Error(), message.User)
        val, _ := json.Marshal(em)

        broadcast <- val
    }
}
