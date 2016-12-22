package utils

import (
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

func NewBroadcasterFromChan(channel chan interface{}) *Broadcaster {
	broadcaster := &Broadcaster{
		receiver: channel,
		isLive:   false,
	}

	return broadcaster
}

func (broadcast *Broadcaster) AddListener(channel chan interface{}) {
	broadcast.RWMutex.Lock()
	broadcast.listeners = append(broadcast.listeners, channel)
	broadcast.RWMutex.Unlock()
}

func (broadcast *Broadcaster) NewListener() chan interface{} {
	channel := make(chan interface{})
	broadcast.AddListener(channel)
	return channel
}

func (broadcast *Broadcaster) RemoveListener(channel chan interface{}) {
	newListeners := []chan interface{}{}

	broadcast.RWMutex.Lock()
	for _, listener := range broadcast.listeners {
		if listener != channel {
			newListeners = append(newListeners, listener)
		}
	}
	broadcast.listeners = newListeners
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
		for _, listener := range broadcast.listeners {
			listener <- event
		}
		broadcast.RWMutex.RUnlock()
	}
}

func (broadcast *Broadcaster) Destroy() {
	broadcast.RWMutex.Lock()
	close(broadcast.receiver)
	broadcast.RWMutex.Unlock()
}
