package main

import (
	"sync"
)

var values sync.Map

func cacheUserText(userId, channel, text string) {
	values.Store(userId+channel, text)
}

func removeUserText(userId, channel string) {
	values.Delete(userId + channel)
}

func getUserText(userId, channel string) (string, bool) {
	value, ok := values.Load(userId + channel)
	if !ok {
		return "", false
	}
	return value.(string), true
}
