package main

import (
	"github.com/patrickmn/go-cache"
)

func main() {
	InitCache()
	ConnectToDiscord()
}
func InitCache() {
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
}
