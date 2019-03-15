package main

import (
    "errors"
    "strconv"
    "strings"

    "github.com/go-redis/redis"
    "github.com/google/uuid"
)

type UserManager struct {
    rdb *redis.Client
}

func NewUserManager(c *redis.Client) *UserManager {
    return &UserManager{
        c,
    }
}

func (um *UserManager) Save(user User) (bool) {
    key, _ := um.key(user)

    data := make(map[string]interface{})
    data["login"] = user.Login
    data["cash"] = user.Cash
    data["coins"] = user.Coins
    data["coins_locked"] = user.CoinsLocked

    um.rdb.HMSet(key, data)
    um.rdb.SAdd("users", user.Login)

    return true
}

func (um *UserManager) Update(user User) (bool) {
    key, _ := um.key(user)

    data := make(map[string]interface{})
    data["login"] = user.Login
    data["cash"] = user.Cash
    data["coins"] = user.Coins
    data["coins_locked"] = user.CoinsLocked

    um.rdb.HMSet(key, data)
    um.rdb.SAdd("users", user.Login)

    return true
}

func (um *UserManager) Find(id uuid.UUID) (u *User, e error) {
    return um.FindByKey(um.uuidKey(id))
}

func (um *UserManager) All() (data []User) {
    keys := um.rdb.Keys("u_*").Val()

    for _, key := range keys {
        user, _ := userManager.FindByKey(key)
        if user != nil && user.Login != "admin" {
            data = append(data, *user)
        }
    }

    return
}

func (um *UserManager) MostCash() (users []User) {

    for _, user := range um.All() {

        switch(true) {
        case user.Coins == 0:
            break
        case len(users) == 0:
            users = append(users, user)
            break
        case user.Cash == users[0].Cash:
            users = append(users, user)
            break
        case user.Cash > users[0].Cash:
            users = []User{user}
            break
        }
    }

    return
}

func (um *UserManager) FindByKey(key string) (u *User, e error) {
    data := um.rdb.HGetAll(key).Val()

    if len(data) == 0 {
        return
    }

    cash, _ := strconv.ParseFloat(data["cash"], 64)
    coins, _ := strconv.ParseInt(data["coins"], 0, 64)
    coinsLocked, _ := strconv.ParseInt(data["coins_locked"], 0, 64)

    u = &User{}
    u.Uuid, _ = uuid.Parse(key[2:])
    u.Login = data["login"]
    u.Cash = cash
    u.Coins = coins
    u.CoinsLocked = coinsLocked

    return u, nil
}

func (um *UserManager) Validate(user User) bool {

    dbUser, _ := um.Find(user.Uuid)

    if dbUser == nil {
        return false
    }

    return dbUser.Same(user)
}

func (um *UserManager) Exists(user User) bool {

    return um.rdb.SIsMember("users", user.Login).Val()
}

func (um *UserManager) LoginExists(login string) bool {

    return um.rdb.SIsMember("users", login).Val()
}

func (um *UserManager) uuidKey(id uuid.UUID) (key string) {
    var b = strings.Builder{}
    b.WriteString("u_")
    b.WriteString(id.String())

    return b.String()
}

func (om *UserManager) key(user User) (key string, e error) {
    if !user.HasUuid() {
        return "", errors.New("No id found")
    }

    return om.uuidKey(user.Uuid), nil
}
