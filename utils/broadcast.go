package utils

import (
	"fmt"
	"sync"
)

type Broadcaster struct {
	// When events are recieved on this channel, broadcast them to all
	// listeners.
	receiver chan interface{}

	listeners []chan interface{}

	sync.RWMutex
	isLive bool
}

func NewBroadcaster() *Broadcaster {
	broadcaster := &Broadcaster{
		receiver: make(chan interface{}),
		isLive:   false,
	}
	return broadcaster
}

func (broadcast *Broadcaster) AddListener(channel chan interface{}) {
	broadcast.RWMutex.Lock()
	broadcast.listeners = append(broadcast.listeners, channel)
	broadcast.RWMutex.Unlock()
}

func (broadcast *Broadcaster) RemoveListener(channel chan interface{}) {
	newListeners := []chan interface{}{}

	broadcast.RWMutex.Lock()
	fmt.Println("before remove", broadcast.listeners, channel)
	for _, listener := range broadcast.listeners {
		if listener != channel {
			newListeners = append(newListeners, listener)
		}
	}
	broadcast.listeners = newListeners
	fmt.Println("removed", broadcast.listeners)
	broadcast.RWMutex.Unlock()
}

func (broadcast *Broadcaster) Broadcast(thing interface{}) {
	if len(broadcast.listeners) == 0 {
		return
	}

	broadcast.receiver <- thing
}

func (broadcast *Broadcaster) StartBroadcasting() {
	if broadcast.isLive {
		// Already broadcasting
		return
	}

	broadcast.isLive = true
	for event := range broadcast.receiver {
		broadcast.RWMutex.RLock()
		fmt.Println(broadcast.listeners)
		for _, listener := range broadcast.listeners {
			listener <- event
		}
		broadcast.RWMutex.RUnlock()
	}
}
