package requesttesting

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type PerformRequestResult struct {
	Req     *http.Request
	ReqBody []byte
	Resp    []byte
}

// PerformRequest performs the HTTP request in 'request' against a http.Server and returns the http.Request that is seen by a http.Handler and the response that the server generates as a []byte.
func PerformRequest(ctx context.Context, request []byte) (PerformRequestResult, error) {
	handler := &saveRequestHandler{}

	srv := http.Server{Handler: handler}
	listener := newInMemoryListener()
	defer listener.Close()

	go srv.Serve(&listener)
	defer srv.Close()

	result := PerformRequestResult{}
	if err := listener.SendRequest(request); err != nil {
		return result, err
	}

	var err error
	result.Resp, err = listener.ReadResponse()
	if err != nil {
		return result, err
	}
	result.Req = handler.LastRequest
	result.ReqBody = handler.LastRequestBody

	return result, srv.Shutdown(ctx)
}

// saveRequestHandler puts the most recent request it has received in LastRequest.
// LastRequestBody needs to be saved separately since the Server closes LastRequest.Body after it has passed the handler.
type saveRequestHandler struct {
	LastRequest     *http.Request
	LastRequestBody []byte
}

func (h *saveRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LastRequest = r
	h.LastRequestBody, _ = ioutil.ReadAll(r.Body)
	io.WriteString(w, "Hello World!")
}

type inMemoryListener struct {
	s2c         io.Closer
	c2s         net.Conn
	connChannel chan net.Conn
	closeOnce   sync.Once
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
func (l *inMemoryListener) SendRequest(request []byte) error {
	_, err := l.c2s.Write(request)
	return err
}

// ReadResponse reads the response from the c2s connection which is sent by the listening server.
// Blocks until the server has sent it's response or times out.
func (l *inMemoryListener) ReadResponse() ([]byte, error) {
	bytes := make([]byte, 4096)
	n, err := l.c2s.Read(bytes)
	if n == 4096 {
		return nil, errors.New("response larger than 4096 bytes")
	}
	return bytes[:n], err
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
	l.closeOnce.Do(func() {
		l.s2c.Close()
		l.c2s.Close()
		close(l.connChannel)
	})
	return nil
}

// Addr returns the listener's network address.
func (l *inMemoryListener) Addr() net.Addr {
	return l.c2s.LocalAddr()
}
