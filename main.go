package main

import (
    "encoding/json"
    "flag"
    "log"
    "net/http"

    "github.com/go-redis/redis"
)

const (
    def_coins       = 0
    def_coinsLocked = 0
    def_cash        = 1000
)

var rdb = redis.NewClient(&redis.Options{
    Addr: "redis:6379",
})

var userManager = NewUserManager(rdb)
var offerManager = NewOfferManager(rdb)

func offers(hub *Hub, w http.ResponseWriter, r *http.Request) {
    m, _ := json.Marshal(offerManager.All())
    w.Write(m)
}

func enableCORS(w http.ResponseWriter, r *http.Request, callback func()) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

    if r.Method == http.MethodOptions {
        return
    }

    callback()
}

var frontServer = flag.Bool("front", false, "Use front server")

func main() {

    flag.Parse()

    hub := newHub()

    if *frontServer {
        http.Handle("/", http.FileServer(http.Dir("./public"))) //handles static html / css etc. under ./public
    }

    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            registerUser(hub, w, r)
        })
    })

    http.HandleFunc("/register-admin", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            registerAdminUser(w, r)
        })
    })

    http.HandleFunc("/admin/best", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            best(hub, w, r)
        })
    })

    http.HandleFunc("/admin/share", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            shareIncome(w, r)
        })
    })

    http.HandleFunc("/admin/clear", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            clearAll(w, r)
        })
    })

    http.HandleFunc("/admin/lock", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            lock(hub, w, r)
        })
    })

    http.HandleFunc("/admin/unlock", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            unlock(hub, w, r)
        })
    })

    http.HandleFunc("/offers", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r, func() {
            offers(hub, w, r)
        })
    })

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    go hub.run()

    err := http.ListenAndServe("0.0.0.0:3001", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
