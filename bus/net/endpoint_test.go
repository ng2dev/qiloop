package net_test

import (
	"io"
	"io/ioutil"
	gonet "net"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
)

func TestConnectUnix(t *testing.T) {
	f, err := ioutil.TempFile("", "go-net-test")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	os.Remove(name)

	listener, err := gonet.Listen("unix", name)
	if err != nil {
		t.Fatal(err)
	}
	_, err = net.DialEndPoint("unix://" + name)
	if err != nil {
		t.Fatal(err)
	}
	listener.Close()
}

func TestSend(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	defer b.Close()

	e := net.ConnEndPoint(a)

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	m := net.NewMessage(h, []byte{0xab, 0xcd})

	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		err := e.Send(m)
		if err != nil {
			panic(err)
		}
		wait.Done()
	}()

	buf := make([]byte, net.HeaderSize+2)
	count, err := b.Read(buf)
	if err != nil {
		panic(err)
	}
	if count != len(buf) {
		t.Errorf("read less than expected: %d", count)
	}
	wait.Wait()
}

func TestReceiveOne(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	defer b.Close()

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	m := net.NewMessage(h, []byte{0xab, 0xcd})

	var wait sync.WaitGroup

	e := net.ConnEndPoint(b)
	msgChan, err := e.ReceiveAny()
	if err != nil {
		panic(err)
	}
	err = m.Write(a)
	if err != nil {
		panic(err)
	}

	msg, ok := <-msgChan
	if !ok {
		panic("receive")
	}

	if len(msg.Payload) != 2 {
		t.Errorf("read less than expected: %d", len(msg.Payload))
	}
	wait.Wait()
}

func TestPingPong(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var wait sync.WaitGroup

	srvChan, err := server.ReceiveAny()
	// reply to a single message
	wait.Add(1)
	go func() {
		m, ok := <-srvChan
		if !ok {
			t.Fatalf("receive meesage")
		}
		err = server.Send(*m)
		if err != nil {
			t.Errorf("send meesage: %s", err)
		}
		wait.Done()
	}()

	var mRcv *net.Message
	var mSent net.Message

	// server replied
	cltChan, err := client.ReceiveAny()
	if err != nil {
		t.Fatalf("receive net.")
	}

	// client send a message
	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	mSent = net.NewMessage(h, []byte{0xab, 0xcd})
	if err := client.Send(mSent); err != nil {
		t.Errorf("send paquet: %s", err)
	}

	mRcv, ok := <-cltChan
	if !ok {
		panic("chanel closed")
	}
	wait.Wait()

	// check packet integrity
	if !reflect.DeepEqual(mSent, *mRcv) {
		t.Errorf("expected %#v, go %#v", mSent, *mRcv)
	}
}

func TestEndPointFinalizer(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	go func() {
		msg := net.NewMessage(net.NewHeader(net.Call, 1, 1, 1, 1), make([]byte, 0))
		msg.Write(a)
	}()
	filter := func(hrd *net.Header) (bool, bool) {
		return true, true
	}
	wait := make(chan *net.Message)
	consumer := func(msg *net.Message) error {
		wait <- msg
		return nil
	}
	closer := func(err error) {
	}
	finalizer := func(e net.EndPoint) {
		e.AddHandler(filter, consumer, closer)
	}
	endpoint := net.EndPointFinalizer(net.ConnStream(b), finalizer)
	defer endpoint.Close()
	msg, ok := <-wait
	if !ok {
		panic("expecting a message")
	}
	if msg == nil {
		panic("not expected nil")
	}
}

func TestEndpointShallAcceptMultipleHandlers(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()

	filter := func(hrd *net.Header) (bool, bool) {
		return true, true
	}
	var wait sync.WaitGroup

	consumer := func(msg *net.Message) error {
		wait.Done()
		return nil
	}
	closer := func(err error) {
	}

	endpoint := net.ConnEndPoint(b)
	defer endpoint.Close()

	ids := make([]int, 20)
	for i := range ids {
		wait.Add(1)
		ids[i] = endpoint.AddHandler(filter, consumer, closer)
	}
	msg := net.NewMessage(net.NewHeader(net.Call, 1, 1, 1, 1), make([]byte, 0))
	msg.Write(a)

	wait.Wait()
	for _, id := range ids {
		err := endpoint.RemoveHandler(id)
		if err != nil {
			panic(err)
		}
	}
	err := endpoint.RemoveHandler(ids[0])
	if err == nil {
		panic("shall fail")
	}
}

func TestEndPoint_ShallDropMessages(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	msg := net.NewMessage(net.NewHeader(net.Call, 1, 1, 1, 1), make([]byte, 0))
	endpoint := net.ConnEndPoint(b)
	if endpoint.String() == "" {
		panic("empty name")
	}
	defer endpoint.Close()
	msg.Write(a)
}

func TestEndPoint_DialTCP(t *testing.T) {
	addr := "tcp://localhost:23456"
	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	defer endpoint.Close()

	addr = "tcp://localhost"
	_, err = net.Listen(addr)
	if err == nil {
		panic("shall fail")
	}

	_, err = net.DialEndPoint(addr)
	if err == nil {
		panic("shall fail")
	}
}

func TestEndPointInvalidAddress(t *testing.T) {
	test := func(addr, msg string) {
		_, err := net.Listen(addr)
		if err == nil {
			panic("listen: " + msg)
		}
		_, err = net.DialEndPoint(addr)
		if err == nil {
			panic("dial: " + msg)
		}
	}
	test("localhost:32432", "missing scheme")
	test("udp://localhost:32432", "invalid scheme")
	test("tcp://localhost", "missing port")
	test("tcp://_localhost", "invalid host")
	test("$#(*W@LKASDL(1", "random")
}

func TestEndPoint_DialTLS(t *testing.T) {
	addr := "tcps://localhost:54321"
	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 10)
		conn.Read(buf)
		conn.Close()
	}()
	defer listener.Close()
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	defer endpoint.Close()

	addr = "tcps://localhost"
	_, err = net.Listen(addr)
	if err == nil {
		panic("shall fail")
	}

	_, err = net.DialEndPoint(addr)
	if err == nil {
		panic("shall fail")
	}
}

func TestEndPoint_DialQUIC(t *testing.T) {
	addr := "quic://localhost:54322"
	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go func() {
		conn, err := listener.Accept()
		if err == io.EOF {
			panic("unexpected closed listener")
		}
		if err != nil {
			panic(err)
		}

		finalizer := func(e net.EndPoint) {
			filter := func(hrd *net.Header) (bool, bool) {
				return true, false
			}
			consumer := func(msg *net.Message) error {
				err = e.Send(*msg)
				if err != nil {
					panic(err)
				}
				return nil
			}
			closer := func(err error) {
			}
			e.AddHandler(filter, consumer, closer)
		}

		net.EndPointFinalizer(conn, finalizer)

	}()
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	defer endpoint.Close()

	response, err := endpoint.ReceiveAny()
	if err != nil {
		panic(err)
	}

	msg := net.NewMessage(net.NewHeader(net.Call, 1, 1, 1, 1), make([]byte, 0))
	endpoint.Send(msg)

	msg2, ok := <-response
	if !ok {
		panic("receive message")
	}

	if msg.Header != msg2.Header {
		panic("different messages")
	}

	addr = "quic://localhost"
	_, err = net.Listen(addr)
	if err == nil {
		panic("shall fail")
	}

	_, err = net.DialEndPoint(addr)
	if err == nil {
		panic("shall fail")
	}
}

func TestEndPoint_DialUnix(t *testing.T) {
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	defer endpoint.Close()

	addr = "unix:///test"
	_, err = net.Listen(addr)
	if err == nil {
		panic("shall fail")
	}

	_, err = net.DialEndPoint(addr)
	if err == nil {
		panic("shall fail")
	}
}

func TestEndPoint_DialPipe(t *testing.T) {
	addr := util.NewPipeAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		conn.Close()
	}()
	defer listener.Close()
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	defer endpoint.Close()

	addr = "pipe:///test"
	_, err = net.Listen(addr)
	if err == nil {
		panic("shall fail")
	}

	_, err = net.DialEndPoint(addr)
	if err == nil {
		panic("shall fail")
	}
}

func TestLoadCertificateError(t *testing.T) {
	os.Setenv("QILOOP_CERT_CONF", "incorrect")
	addr := "tcps://localhost:23432"
	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
}
