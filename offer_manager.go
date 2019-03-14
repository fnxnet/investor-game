package main

import (
    "errors"
    "strconv"
    "strings"

    "github.com/go-redis/redis"
    "github.com/google/uuid"
)

type OfferManager struct {
    rdb *redis.Client
}

func NewOfferManager(rdb *redis.Client) *OfferManager {
    return &OfferManager{
        rdb: rdb,
    }
}

func (om *OfferManager) Remove(offer Offer) bool {
    key, _ := om.key(offer)

    om.rdb.Del(key)

    return true
}

func (om *OfferManager) Find(id uuid.UUID) (offer *Offer, e error) {
    return om.FindByKey(om.generateKey(id))
}
func (om *OfferManager) FindByKey(key string) (offer *Offer, e error) {
    data := om.rdb.HGetAll(key).Val()

    if len(data) == 0 {
        return
    }

    coins, _ := strconv.ParseInt(data["coins"], 0, 64);
    price, _ := strconv.ParseInt(data["price"], 0, 64);
    userId, _ := uuid.Parse(data["user"]);

    o := &Offer{
        Coins: int64(coins),
        Price: int64(price),
        User:  &User{Uuid: userId},
    }
    o.UUid, _ = uuid.Parse(key[2:])

    return o, nil
}

func (om *OfferManager) All() (data []Offer) {
    keys := om.rdb.Keys("o_*").Val()

    for _, key := range keys {
        offer, _ := om.FindByKey(key)
        data = append(data, *offer)
    }

    if len(data)==0{
        return make([]Offer,0)
    }

    return
}
func (om *OfferManager) Save(offer Offer) (e error) {
    key, _ := om.key(offer)

    data := make(map[string]interface{})
    data["coins"] = offer.Coins
    data["price"] = offer.Price
    data["user"] = offer.User.Uuid.String()

    om.rdb.HMSet(key, data)

    return
}

func (om *OfferManager) Exists(offer Offer) bool {
    if !offer.HasUuid() {
        return false
    }

    key, _ := om.key(offer)

    return om.rdb.HExists(key, "price").Val()
}

func (om *OfferManager) generateKey(id uuid.UUID) (key string) {
    var b = strings.Builder{}
    b.WriteString("o_")
    b.WriteString(id.String())

    return b.String()
}

func (om *OfferManager) key(offer Offer) (key string, e error) {
    if !offer.HasUuid() {
        return "", errors.New("No id found")
    }

    return om.generateKey(offer.UUid), nil
}
