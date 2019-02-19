package net

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	gonet "net"
	"net/url"
	"strings"
	"sync"
)

// Filter returns true if given message shall be processed by a
// Consumer. Returns two values:
// - matched: true if the message should be processed by the Consumer.
// - keep: true if the handler shall be kept in the dispatcher.
type Filter func(hdr *Header) (matched bool, keep bool)

// Consumer process a message which has been selected by a filter.
type Consumer func(msg *Message) error

// Closer informs the handler about a disconnection
type Closer func(err error)

// EndPoint reprensents a network socket capable of sending and
// receiving messages.
type EndPoint interface {

	// Send pushes the message into the network.
	Send(m Message) error

	// ReceiveAny returns a chanel to receive a single message.
	ReceiveAny() (chan *Message, error)

	// AddHandler registers the associated Filter and Consumer to the
	// EndPoint. Do not attempt to add another handler from within a
	// Filter.
	AddHandler(f Filter, c Consumer, cl Closer) int

	// RemoveHandler removes the associated Filter and Consumer.
	// RemoveHandler must not be called from within the Filter: use
	// the Filter returned value keep for this purpose.
	RemoveHandler(id int) error

	// Close close the underlying connection
	Close() error

	String() string
}

// Handler represents a client of a incomming message stream. It
// contains a filter used to decided if the message shall be sent to
// the handler, a queue and a consumuer which process the messages and
// a closer called when the connection is closed.
type Handler struct {
	filter   Filter
	consumer Consumer
	closer   Closer
	queue    chan *Message
	cancel   chan struct{}
	err      error
}

// NewHandler returns an Handler: f is call on each incomming message,
// c is called if f returns true. cl is always called when the handler
// is effectively closed.
func NewHandler(f Filter, c Consumer, cl Closer) *Handler {
	h := &Handler{
		filter:   f,
		consumer: c,
		closer:   cl,
		queue:    make(chan *Message, 10),
		cancel:   make(chan struct{}),
	}
	go h.run()
	return h
}

// Stop stops the handler main loop. If immediatly, the pending
// messages in the queue will be dropped, else the queue will be
// processed before terminating the handler.
func (h *Handler) Stop(immediatly bool) {
	if immediatly {
		close(h.cancel)
	} else {
		close(h.queue)
	}
}

func (h *Handler) run() {
loop:
	for {
		select {
		case msg, ok := <-h.queue:
			if !ok {
				break loop
			}
			err := h.consumer(msg)
			if err != nil {
				log.Printf("failed to consume message: %s", err)
			}
		case <-h.cancel:
			break loop
		}
	}
	h.closer(h.err)
}

type endPoint struct {
	conn          gonet.Conn
	handlers      []*Handler
	handlersMutex sync.Mutex
}

// EndPointFinalizer creates a new EndPoint and let you process it
// before it start handling messages. This allows you to add handler
// or/and avoid data races.
func EndPointFinalizer(conn gonet.Conn, finalizer func(EndPoint)) EndPoint {
	e := &endPoint{
		conn:     conn,
		handlers: make([]*Handler, 10),
	}
	finalizer(e)
	go e.process()
	return e
}

// NewEndPoint returns an EndPoint which already process incomming
// messages. Since no handler have been register at the time of the
// creation of the EndPoint, any message receive will be droped until
// and Handler is registered. Prefer EndPointFinalizer for a safe way
// to construct EndPoint.
func NewEndPoint(conn gonet.Conn) EndPoint {
	e := &endPoint{
		conn:     conn,
		handlers: make([]*Handler, 10),
	}
	go e.process()
	return e
}

func dialUNIX(name string) (EndPoint, error) {
	conn, err := gonet.Dial("unix", name)
	if err != nil {
		return nil, fmt.Errorf(`failed to connect unix socket "%s": %s`,
			name, err)
	}
	return NewEndPoint(conn), nil
}

func dialTCP(addr string) (EndPoint, error) {
	conn, err := gonet.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return NewEndPoint(conn), nil
}

// dialTLS connects regardless of the certificate.
func dialTLS(addr string) (EndPoint, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return NewEndPoint(conn), nil
}

// DialEndPoint construct an endpoint by contacting a given address.
func DialEndPoint(addr string) (EndPoint, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("dial: invalid address: %s", err)
	}
	switch u.Scheme {
	case "tcp":
		return dialTCP(u.Host)
	case "tcps":
		return dialTLS(u.Host)
	case "unix":
		return dialUNIX(strings.TrimPrefix(addr, "unix://"))
	default:
		return nil, fmt.Errorf("unknown URL scheme: %s", addr)
	}
}

// Send post a message to the other side of the endpoint.
func (e *endPoint) Send(m Message) error {
	return m.Write(e.conn)
}

// closeWith close all handler
func (e *endPoint) closeWith(err error) error {

	ret := e.conn.Close()

	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()

	for id, handler := range e.handlers {
		if handler != nil {
			handler.err = err
			handler.Stop(false)
			e.handlers[id] = nil
		}
	}
	return ret
}

// Close wait for a message to be received.
func (e *endPoint) Close() error {
	return e.closeWith(nil)
}

// RemoveHandler unregister the associated Filter and Consumer.
// WARNING: RemoveHandler must not be called from within the Filter or
// the Consumer.
func (e *endPoint) RemoveHandler(id int) error {
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	if id >= 0 && id < len(e.handlers) && e.handlers[id] != nil {
		e.handlers[id].Stop(true)
		e.handlers[id] = nil
		return nil
	}
	return fmt.Errorf("invalid handler id: %d", id)
}

// AddHandler register the associated Filter and Consumer to the
// EndPoint.
func (e *endPoint) AddHandler(f Filter, c Consumer, cl Closer) int {
	newHandler := NewHandler(f, c, cl)
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	for i, handler := range e.handlers {
		if handler == nil {
			e.handlers[i] = newHandler
			return i
		}
	}
	e.handlers = append(e.handlers, newHandler)
	return len(e.handlers) - 1
}

// ErrNoMatch is returned when the message did not match any handler
var ErrNoMatch = errors.New("message dropped: no handler match")

// ErrNoHandler is returned when there is no handler registered.
var ErrNoHandler = errors.New("message dropped: no handler registered")

// dispatch requests each handler if it match the message, if so it
// sends it to the handler queue. If no handler is registered, it
// returns ErrNoHandler and if no handler match the message it returns
// ErrNoMatch.
func (e *endPoint) dispatch(msg *Message) error {
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	if len(e.handlers) == 0 {
		return ErrNoHandler
	}
	ret := ErrNoMatch
	for i, h := range e.handlers {
		if h == nil {
			continue
		}
		matched, keep := h.filter(&msg.Header)
		if matched {
			h.queue <- msg
			ret = nil
		}
		if !keep {
			h.Stop(false)
			e.handlers[i] = nil
		}
	}
	return ret
}

// process read all messages from the end point and dispatch them one
// by one.
func (e *endPoint) process() {
	var err error

	for {
		msg := new(Message)
		err = msg.Read(e.conn)
		if err != nil {
			e.closeWith(err)
			return
		}
		err = e.dispatch(msg)
		if err != nil {
			log.Printf("%s: %#v", err, msg.Header)
		}
	}
}

// ReceiveAny returns a chanel to receive one message. If the
// connection close, the chanel is closed.
func (e *endPoint) ReceiveAny() (chan *Message, error) {
	found := make(chan *Message, 1)
	filter := func(hdr *Header) (matched bool, keep bool) {
		return true, false
	}
	consumer := func(msg *Message) error {
		found <- msg
		return nil
	}
	closer := func(err error) {
		close(found)
	}
	_ = e.AddHandler(filter, consumer, closer)
	return found, nil
}

func (e *endPoint) String() string {
	return e.conn.RemoteAddr().Network() + "://" +
		e.conn.RemoteAddr().String()
}

// Pipe returns a set of EndPoint connected to each other.
func Pipe() (EndPoint, EndPoint) {
	a, b := gonet.Pipe()
	return NewEndPoint(a), NewEndPoint(b)
}
