package main

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
)

type UserResponse struct {
    Payload User   `json:"payload"`
    Error   string `json:"error"`
}

func registerUser(hub *Hub, w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        http.Error(w, "", 403)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid JSON", 500)
        return
    }

    user := &User{}
    json.Unmarshal(body, user)

    if userManager.Exists(*user) {
        sendConflictResponse(w, r, *user)
        return
    }

    user = NewUser(user.Login)
    if user.Login == "admin" {
        sendConflictResponse(w, r, *user)
        return
    }

    result := userManager.Save(*user)

    if !result {
        sendConflictResponse(w, r, *user)
        return
    }

    d, _ := json.Marshal(user)
    w.Write(d)

    data := make(map[string]interface{})
    data["action"] = "userRegistered"

    m, _ := json.Marshal(data)

    hub.send(m)
}

func registerAdmin(w http.ResponseWriter, r *http.Request) {
    admin := NewUser("admin")
    admin.Coins = 1000

    if !userManager.LoginExists("admin") {
        if ok := userManager.Save(*admin); ok {
            m, _ := json.Marshal(admin)
            w.Write(m)
            return
        }
    }

    sendConflictResponse(w, r, *admin)
}

func sendConflictResponse(w http.ResponseWriter, r *http.Request, user User) {
    d, _ := json.Marshal(&UserResponse{
        user,
        "User exists",
    })

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(409)
    w.Write(d)
}