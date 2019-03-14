package main

import (
    "github.com/google/uuid"
)

type User struct {
    Login       string    `json:"login"`
    Uuid        uuid.UUID `json:"token"`
    Coins       int64     `json:"coins"`
    CoinsLocked int64     `json:"coins_locked"`
    Cash        float64   `json:"cash"`
}

func (u *User) Same(user User) bool {
    return u.Uuid.String() == user.Uuid.String()
}

func (u *User) HasUuid() bool {
    return u.Uuid.String() != (&uuid.UUID{}).String()
}

func (u *User) HasMoney(needed float64) bool {
    return u.Cash >= needed
}

func NewUser(login string) *User {
    return &User{
        Login:       login,
        Uuid:        uuid.New(),
        Coins:       def_coins,
        CoinsLocked: def_coinsLocked,
        Cash:        def_cash,
    }
}
