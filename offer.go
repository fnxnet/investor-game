package main

import (
    "github.com/google/uuid"
)

type Offer struct {
    Coins   int64       `json:"coins"`
    Price   int64       `json:"price"`
    UUid    uuid.UUID   `json:"id"`
    User    *User        `json:"user"`
    Payload interface{} `json:"payload"`
}

func (o Offer) Total() int64 {
    return o.Coins * o.Price
}

func (o *Offer) HasUuid() bool {
    return o.UUid.String() != (&uuid.UUID{}).String()
}
