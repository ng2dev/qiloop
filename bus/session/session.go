package session

import (
	"fmt"
	"log"
	"sync"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/type/object"
)

// Session implements the Session interface. It is an
// implementation of Session. It does not update the list of services
// and returns clients.
type Session struct {
	serviceList      []services.ServiceInfo
	serviceListMutex sync.Mutex
	Directory        services.ServiceDirectoryProxy
	cancel           func()
	added            chan services.ServiceAdded
	removed          chan services.ServiceRemoved
	userName         string
	userToken        string
	poll             map[string]bus.Client
	pollMutex        sync.RWMutex
}

func (s *Session) newObject(info services.ServiceInfo, ref object.ObjectReference) (bus.ObjectProxy, error) {
	c, err := s.client(info)
	if err != nil {
		return nil, fmt.Errorf("object connection error (%s): %s",
			info.Name, err)
	}
	proxy := bus.NewProxy(c, ref.MetaObject, ref.ServiceID, ref.ObjectID)
	return bus.MakeObject(proxy), nil
}

func (s *Session) newService(info services.ServiceInfo, objectID uint32) (p bus.Proxy, err error) {
	c, err := s.client(info)
	if err != nil {
		return nil, fmt.Errorf("service connection error (%s): %s", info.Name, err)
	}
	proxy, err := metaProxy(c, info.ServiceId, objectID)
	if err != nil {
		return nil, fmt.Errorf("get service meta object (%s): %s", info.Name, err)
	}
	return proxy, nil
}

// endpoint returns an net.EndPoint matching the info description. If
// an existing connection exists, it reuse the connection, otherwise
// it establish a new connection.
func (s *Session) client(info services.ServiceInfo) (bus.Client, error) {
	if len(info.Endpoints) == 0 {
		return nil, fmt.Errorf("empty address list")
	}
	s.pollMutex.RLock()
	for _, addr := range info.Endpoints {
		c, ok := s.poll[addr]
		if ok {
			s.pollMutex.RUnlock()
			return c, nil
		}
	}
	s.pollMutex.RUnlock()
	addr, endpoint, err := bus.SelectEndPoint(info.Endpoints, s.userName, s.userToken)
	if err != nil {
		return nil, fmt.Errorf("service connection error (%s): %s", info.Name, err)
	}
	filter := func(hdr *net.Header) (matched bool, keep bool) { return false, true }
	consumer := func(msg *net.Message) error { panic("unexpected") }
	closer := func(err error) {
		s.pollMutex.Lock()
		delete(s.poll, addr)
		s.pollMutex.Unlock()
	}
	s.pollMutex.Lock()
	c, ok := s.poll[addr]
	if ok {
		s.pollMutex.RUnlock()
		endpoint.Close()
		return c, nil
	}
	c = bus.NewClient(endpoint)
	s.poll[addr] = c
	s.pollMutex.Unlock()
	endpoint.AddHandler(filter, consumer, closer)
	return c, nil
}

func (s *Session) findServiceName(name string) (i services.ServiceInfo, err error) {
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	for _, service := range s.serviceList {
		if service.Name == name {
			return service, nil
		}
	}
	return i, fmt.Errorf("Service not found: %s", name)
}

func (s *Session) findServiceID(uid uint32) (i services.ServiceInfo, err error) {
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	for _, service := range s.serviceList {
		if service.ServiceId == uid {
			return service, nil
		}
	}
	return i, fmt.Errorf("Service ID not found: %d", uid)
}

// Proxy resolve the service name and returns a proxy to it.
func (s *Session) Proxy(name string, objectID uint32) (p bus.Proxy, err error) {
	info, err := s.findServiceName(name)
	if err != nil {
		return p, err
	}
	return s.newService(info, objectID)
}

// Object returns a reference to ref.
// TODO: cache the returned objects in order to benefit from the
// signal registration caching.
func (s *Session) Object(ref object.ObjectReference) (o bus.Proxy, err error) {
	info, err := s.findServiceID(ref.ServiceID)
	if err != nil {
		return o, err
	}
	return s.newObject(info, ref)
}

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(c bus.Client, serviceID, objectID uint32) (p bus.Proxy, err error) {
	meta, err := bus.GetMetaObject(c, serviceID, objectID)
	if err != nil {
		return p, fmt.Errorf("Can not reach metaObject: %s", err)
	}
	return bus.NewProxy(c, meta, serviceID, objectID), nil
}

// NewAuthSession connects an address and return a new session.
func NewAuthSession(addr, user, token string) (bus.Session, error) {

	s := new(Session)
	s.userName = user
	s.userToken = token
	s.poll = map[string]bus.Client{}
	// Manually create a serviceList with just the ServiceInfo
	// needed to contact ServiceDirectory.
	s.serviceList = []services.ServiceInfo{
		services.ServiceInfo{
			Name:      "ServiceDirectory",
			ServiceId: 1,
			Endpoints: []string{
				addr,
			},
		},
	}
	var err error
	s.Directory, err = services.Services(s).ServiceDirectory(nil)
	if err != nil {
		return nil, fmt.Errorf("contact server: %s", err)
	}

	s.serviceList, err = s.Directory.Services()
	if err != nil {
		return nil, fmt.Errorf("list services: %s", err)
	}
	var cancelRemoved, cancelAdded func()
	cancelRemoved, s.removed, err = s.Directory.SubscribeServiceRemoved()
	if err != nil {
		return nil, fmt.Errorf("subscribe remove signal: %s", err)
	}
	cancelAdded, s.added, err = s.Directory.SubscribeServiceAdded()
	if err != nil {
		return nil, fmt.Errorf("subscribe added signal: %s", err)
	}
	s.cancel = func() {
		cancelRemoved()
		cancelAdded()
	}
	go s.updateLoop()
	return s, nil
}

// NewSession connects an address and return a new session.
func NewSession(addr string) (bus.Session, error) {
	return NewAuthSession(addr, "", "")
}

func (s *Session) updateServiceList() {
	services, err := s.Directory.Services()
	if err != nil {
		log.Printf("error: failed to update service directory list: %s", err)
		log.Printf("error: closing session.")
		if err := s.Terminate(); err != nil {
			log.Printf("error: session destruction: %s", err)
		}
	}
	s.serviceListMutex.Lock()
	s.serviceList = services
	s.serviceListMutex.Unlock()
}

// Terminate close the session.
func (s *Session) Terminate() error {
	s.cancel()
	return fmt.Errorf("Session.Terminate: Not yet implemented")
}

func (s *Session) updateLoop() {
	for {
		select {
		case _, ok := <-s.removed:
			if !ok {
				return
			}
			s.updateServiceList()
		case _, ok := <-s.added:
			if !ok {
				return
			}
			s.updateServiceList()
		}
	}
}
