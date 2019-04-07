package main

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "strings"
)

type SimpleMessage struct {
    Password string `json="password"`
}

func (sm *SimpleMessage) GetPassword() string {
    return sm.Password
}

type AdminMessage interface {
    GetPassword() string
}

type ShareMessage struct {
    SimpleMessage
    Prev    float64 `json:"prev"`
    Current float64 `json:"current"`
    Diff    float64 `json:"diff"`
}

func lock(hub *Hub, w http.ResponseWriter, r *http.Request) {
    msg := &SimpleMessage{}
    if e := validateAdminMessage(w, r, msg); e != nil {
        http.Error(w, e.Error(), 500)
        return
    }

    hub.lock()
}

func unlock(hub *Hub, w http.ResponseWriter, r *http.Request) {
    msg := &SimpleMessage{}
    if e := validateAdminMessage(w, r, msg); e != nil {
        http.Error(w, e.Error(), 500)
        return
    }

    hub.unlock()
}

func best(hub *Hub, w http.ResponseWriter, r *http.Request) {
    msg := &SimpleMessage{}
    if e := validateAdminMessage(w, r, msg); e != nil {
        http.Error(w, e.Error(), 403)
        return
    }

    m, _ := json.Marshal(userManager.MostCash())
    w.Write(m)
}

func clearAll(w http.ResponseWriter, r *http.Request) {
    msg := &SimpleMessage{}
    if e := validateAdminMessage(w, r, msg); e != nil {
        http.Error(w, e.Error(), 500)
        return
    }

    rdb.FlushAll()
}

func registerAdminUser(w http.ResponseWriter, r *http.Request) {
    msg := &SimpleMessage{}
    if e := validateAdminMessage(w, r, msg); e != nil {
        http.Error(w, e.Error(), 500)
        return
    }

    registerAdmin(w, r)
}

func shareIncome(hub *Hub, w http.ResponseWriter, r *http.Request) {
    msg := &ShareMessage{}

    if e := validateAdminMessage(w, r, msg); e != nil {
        http.Error(w, e.Error(), 500)
        return
    }

    users := userManager.All()

    for _, user := range users {
        if user.Coins > 0 {
            user.Cash += msg.Diff / 1000 * float64(user.Coins)
            userManager.Save(user)

            data := make(map[string]interface{})
            data["action"] = "incomeReceived"
            data["user"] = user

            m, _ := json.Marshal(data)

            hub.send(m)
        }
    }
}

var AdminPassword = ""

func init() {
    if dat, e := ioutil.ReadFile("./.passwd"); e == nil {
        AdminPassword = strings.Replace(string(dat), "\n", "", -1)
    }
}

func validateAdminMessage(w http.ResponseWriter, r *http.Request, message AdminMessage) (e error) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return errors.New("Invalid JSON")
    }

    json.Unmarshal(body, message)

    if message.GetPassword() != AdminPassword && len(AdminPassword) > 0 {
        return errors.New("Invalid Admin credentials")
    }

    return
}
