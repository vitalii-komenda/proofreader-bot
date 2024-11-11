package slashcommands

import (
	"sync"
)

var values sync.Map

func CacheUserText(userId, channel, role, text string) {
	values.Store(userId+channel+role, text)
}

func removeUserText(userId, channel, role string) {
	values.Delete(userId + channel + role)
}

func GetUserText(userId, channel, role string) (string, bool) {
	value, ok := values.Load(userId + channel + role)
	if !ok {
		return "", false
	}
	return value.(string), true
}
