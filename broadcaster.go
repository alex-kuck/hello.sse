package main

type Broadcaster[T any] struct {
	source  <-chan T
	clients []chan T
}

func NewBroadcaster[T any](source <-chan T) *Broadcaster[T] {
	b := &Broadcaster[T]{source, nil}
	go func() {
		for v := range source {
			for _, ch := range b.clients {
				ch <- v
			}
		}
	}()
	return b
}

func (b *Broadcaster[T]) Subscribe() <-chan T {
	ch := make(chan T)
	b.clients = append(b.clients, ch)
	return ch
}

func (b *Broadcaster[T]) Unsubscribe(c <-chan T) {
	for i, ch := range b.clients {
		if ch == c {
			close(ch)
			b.clients[i] = b.clients[len(b.clients)-1] // replace with last element
			b.clients = b.clients[:len(b.clients)-1]   // drop duplicate last element
			break
		}
	}
}
