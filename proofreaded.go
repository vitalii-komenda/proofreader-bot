package main

import (
	"sync"
)

var values sync.Map

func addProofreaded(userId, channel, text string) {
	values.Store(userId+channel, text)
}

func removeProofreaded(userId, channel string) {
	values.Delete(userId + channel)
}

func getProofreaded(userId, channel string) (string, bool) {
	value, ok := values.Load(userId + channel)
	if !ok {
		return "", false
	}
	return value.(string), true
}
