package requesttesting

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
)

// PerformRequest performs the HTTP request in 'request' against a http.Server and returns the http.Request that is seen by a http.Handler and the response that the server generates.
func PerformRequest(ctx context.Context, request string) (*http.Request, *http.Response, error) {
	handler := &saveRequestHandler{LastRequest: nil}

	srv := http.Server{Handler: handler}
	listener := newInMemoryListener()

	go func() {
		srv.Serve(&listener)
	}()
	listener.SendRequest(request)
	resp, err := listener.ReadResponse()
	srv.Shutdown(ctx)

	return handler.LastRequest, resp, err
}

// saveRequestHandler puts the most recent request it has received in LastRequest
type saveRequestHandler struct {
	LastRequest *http.Request
}

func (h *saveRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LastRequest = r
	io.WriteString(w, "Hello World!")
}

type inMemoryListener struct {
	s2c         net.Conn
	c2s         net.Conn
	connChannel chan net.Conn
}

func newInMemoryListener() inMemoryListener {
	s2c, c2s := net.Pipe()
	connChannel := make(chan net.Conn, 1)
	connChannel <- s2c

	return inMemoryListener{
		s2c:         s2c,
		c2s:         c2s,
		connChannel: connChannel,
	}
}

// SendRequest writes 'request' to the c2s connection which will send the request to the server listening on this listener.
// Blocks until the server has read the message.
func (l *inMemoryListener) SendRequest(request string) {
	l.c2s.Write([]byte(request))
}

// ReadResponse reads the response from the c2s connection which is sent by the listening server.
// Blocks until the server has sent it's response or times out.
func (l *inMemoryListener) ReadResponse() (*http.Response, error) {
	reader := bufio.NewReader(l.c2s)
	return http.ReadResponse(reader, nil)
}

// Accept waits for and returns the next connection to the listener.
func (l *inMemoryListener) Accept() (net.Conn, error) {
	s2c, ok := <-l.connChannel
	if !ok {
		return s2c, errors.New("listener is closed")
	}
	return s2c, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *inMemoryListener) Close() error {
	l.s2c.Close()
	l.c2s.Close()
	close(l.connChannel)
	return nil
}

// Addr returns the listener's network address.
func (l *inMemoryListener) Addr() net.Addr {
	return l.c2s.LocalAddr()
}
