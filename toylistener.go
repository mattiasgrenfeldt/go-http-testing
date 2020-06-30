package main

import (
	"errors"
	"net"
	"sync"
)

// ToyListener implements net.Listener
type ToyListener struct {
	s2c          net.Conn
	c2s          net.Conn
	wg           *sync.WaitGroup
	acceptedOnce bool
	request      string
}

// NewToyListener creates a new ToyListener that is only in memory.
func NewToyListener(request string) ToyListener {
	s2c, c2s := net.Pipe()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	return ToyListener{
		s2c:          s2c,
		c2s:          c2s,
		wg:           wg,
		acceptedOnce: false,
		request:      request,
	}
}

// Accept waits for and returns the next connection to the listener.
func (l *ToyListener) Accept() (net.Conn, error) {
	if l.acceptedOnce {
		l.wg.Wait()
		return l.s2c, errors.New("listener is closed")
	}

	go func() {
		l.c2s.Write([]byte(l.request))
	}()
	l.acceptedOnce = true
	return l.s2c, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *ToyListener) Close() error {
	l.s2c.Close()
	l.c2s.Close()
	l.wg.Done()
	return nil
}

// Addr returns the listener's network address.
func (l *ToyListener) Addr() net.Addr {
	return l.c2s.LocalAddr()
}
