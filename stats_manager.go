package main

import (
    "github.com/go-redis/redis"
)

type Stats struct {
    Name      string
    redis     redis.Client
    CoinsSold int64
    TotalCash float64
    BestOffer float64
}

func NewStats(name string, redis *redis.Client) *Stats {
    return &Stats{
        name,
        *redis,
        0,
        0,
        0,
    }
}

func (s *Stats) RecordOffer(offer Offer) {
    s.TotalCash += float64(offer.Total())
    s.CoinsSold += offer.Coins

    if s.BestOffer < float64(offer.Price) {
        s.BestOffer = float64(offer.Price)
    }
}

func (s *Stats) Reset() {
    s.CoinsSold = 0
    s.TotalCash = 0
}

func (s *Stats) Average() float64 {
    if s.CoinsSold == 0 || s.TotalCash == 0 {
        return 0
    }
    return s.TotalCash / float64(s.CoinsSold)
}
