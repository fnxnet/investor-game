package main

import (
    "encoding/json"
    "errors"
    "reflect"
    "strings"

    "github.com/google/uuid"
)

type OfferMessage struct {
    RequestMessage
    Payload Offer `json:"payload"`
}

type OfferHandler struct{}

// Decorator for all methods
func (oh *OfferHandler) Handle(message []byte, broadcast chan []byte) (e error) {
    om := &OfferMessage{}
    if e = json.Unmarshal(message, om); e != nil {
        return e
    }

    inputs := make([]reflect.Value, 2)
    inputs[0] = reflect.ValueOf(om)
    inputs[1] = reflect.ValueOf(broadcast)

    result := reflect.ValueOf(oh).MethodByName(strings.Title(om.Action)).Call(inputs)[0].Interface()

    if result != nil {
        return result.(error)
    }
    return
}

func (oh *OfferHandler) NewOffer(om *OfferMessage, channel chan []byte) (e error) {
    offer := om.Payload
    if !offer.HasUuid() {
        offer.UUid = uuid.New()
    } else if offerManager.Exists(offer) {
        e = errors.New("Offer exists")
        return
    }

    user, _ := userManager.Find(om.User.Uuid)

    if user == nil {
        return
    }

    offer.User = user

    user.CoinsLocked += offer.Coins
    if user.CoinsLocked > user.Coins {
        return errors.New("You have insufficient amount of Coins")
    }

    offerManager.Save(offer)
    userManager.Save(*user)

    rm := &ResponseMessage{
        Action:  om.Callback,
        Payload: offer,
    }

    val, _ := json.Marshal(rm)

    channel <- val

    return
}

func (oh *OfferHandler) RemoveOffer(om *OfferMessage, channel chan []byte) (e error) {

    offer := om.Payload
    dbOffer, e := offerManager.Find(offer.UUid);
    if e != nil{
        return
    }

    if dbOffer == nil{
        oh.sendResponse(om.Callback, offer, channel)
        return
    }

    if dbOffer != nil {
        user, _ := userManager.Find(dbOffer.User.Uuid)
        dbOffer.User = user

        user.CoinsLocked -= dbOffer.Coins

        offerManager.Remove(*dbOffer)
        userManager.Save(*user)

        offer.User = user
    }

    oh.sendResponse(om.Callback, offer, channel)

    return
}

func (oh *OfferHandler) AcceptOffer(om *OfferMessage, channel chan []byte) (e error) {

    offer := om.Payload
    dbOffer, _ := offerManager.Find(offer.UUid)


    if dbOffer == nil {
        return errors.New("offer not found")
    }

    owner, _ := userManager.Find(dbOffer.User.Uuid)
    buyer, _ := userManager.Find(om.User.Uuid)

    if owner.Same(*buyer) {
        return errors.New("Can not buy your own shares")
    }

    total := offer.Total()
    if !buyer.HasMoney(float64(total)) {
        return errors.New("You have no funds on the account")
    }

    owner.Cash += float64(total)
    buyer.Cash -= float64(total)

    owner.Coins -= dbOffer.Coins
    buyer.Coins += dbOffer.Coins

    dbOffer.Coins -= offer.Coins
    owner.CoinsLocked -= offer.Coins

    userManager.Save(*owner)
    userManager.Save(*buyer)

    if dbOffer.Coins <= 0 {
        offerManager.Remove(*dbOffer)
    } else {
        offerManager.Save(*dbOffer)
    }

    chunkStats.RecordOffer(offer)
    mainStats.RecordOffer(offer)

    offer.Payload = map[string]User{"buyer": *buyer, "owner": *owner}

    oh.sendResponse(om.Callback, offer, channel)

    return
}

func (oh *OfferHandler) sendResponse(callback string, offer Offer, channel chan []byte) {
    rm := &ResponseMessage{
        Action:  callback,
        Payload: offer,
    }

    val, _ := json.Marshal(rm)

    channel <- val
}
