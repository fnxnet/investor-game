package main

type RequestMessage struct {
    Callback string      `json:"callback"`
    Action   string      `json:"action"`
    User     User        `json:"user"`
    Payload  interface{} `json:"payload"`
}

type ResponseMessage struct {
    Action  string      `json:"action"`
    Payload interface{} `json:"payload"`
}
