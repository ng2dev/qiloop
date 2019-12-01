// Package services contains a generated proxy
// .

package services

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// ServiceAdded is serializable
type ServiceAdded struct {
	ServiceID uint32
	Name      string
}

// readServiceAdded unmarshalls ServiceAdded
func readServiceAdded(r io.Reader) (s ServiceAdded, err error) {
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceID field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	return s, nil
}

// writeServiceAdded marshalls ServiceAdded
func writeServiceAdded(s ServiceAdded, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("write ServiceID field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	return nil
}

// ServiceRemoved is serializable
type ServiceRemoved struct {
	ServiceID uint32
	Name      string
}

// readServiceRemoved unmarshalls ServiceRemoved
func readServiceRemoved(r io.Reader) (s ServiceRemoved, err error) {
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceID field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	return s, nil
}

// writeServiceRemoved marshalls ServiceRemoved
func writeServiceRemoved(s ServiceRemoved, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("write ServiceID field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	return nil
}

// ServiceDirectory is the abstract interface of the service
type ServiceDirectory interface {
	// Service calls the remote procedure
	Service(name string) (ServiceInfo, error)
	// Services calls the remote procedure
	Services() ([]ServiceInfo, error)
	// RegisterService calls the remote procedure
	RegisterService(info ServiceInfo) (uint32, error)
	// UnregisterService calls the remote procedure
	UnregisterService(serviceID uint32) error
	// ServiceReady calls the remote procedure
	ServiceReady(serviceID uint32) error
	// UpdateServiceInfo calls the remote procedure
	UpdateServiceInfo(info ServiceInfo) error
	// MachineId calls the remote procedure
	MachineId() (string, error)
	// _socketOfService calls the remote procedure
	_socketOfService(serviceID uint32) (object.ObjectReference, error)
	// SubscribeServiceAdded subscribe to a remote signal
	SubscribeServiceAdded() (unsubscribe func(), updates chan ServiceAdded, err error)
	// SubscribeServiceRemoved subscribe to a remote signal
	SubscribeServiceRemoved() (unsubscribe func(), updates chan ServiceRemoved, err error)
}

// ServiceDirectoryProxy represents a proxy object to the service
type ServiceDirectoryProxy interface {
	object.Object
	bus.Proxy
	ServiceDirectory
}

// proxyServiceDirectory implements ServiceDirectoryProxy
type proxyServiceDirectory struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeServiceDirectory returns a specialized proxy.
func MakeServiceDirectory(sess bus.Session, proxy bus.Proxy) ServiceDirectoryProxy {
	return &proxyServiceDirectory{bus.MakeObject(proxy), sess}
}

// ServiceDirectory returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) ServiceDirectory(closer func(error)) (ServiceDirectoryProxy, error) {
	proxy, err := c.session.Proxy("ServiceDirectory", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeServiceDirectory(c.session, proxy), nil
}

// Service calls the remote procedure
func (p *proxyServiceDirectory) Service(name string) (ServiceInfo, error) {
	var err error
	var ret ServiceInfo
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Call("service", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call service failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = readServiceInfo(resp)
	if err != nil {
		return ret, fmt.Errorf("parse service response: %s", err)
	}
	return ret, nil
}

// Services calls the remote procedure
func (p *proxyServiceDirectory) Services() ([]ServiceInfo, error) {
	var err error
	var ret []ServiceInfo
	var buf bytes.Buffer
	response, err := p.Call("services", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call services failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []ServiceInfo, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]ServiceInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readServiceInfo(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse services response: %s", err)
	}
	return ret, nil
}

// RegisterService calls the remote procedure
func (p *proxyServiceDirectory) RegisterService(info ServiceInfo) (uint32, error) {
	var err error
	var ret uint32
	var buf bytes.Buffer
	if err = writeServiceInfo(info, &buf); err != nil {
		return ret, fmt.Errorf("serialize info: %s", err)
	}
	response, err := p.Call("registerService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerService failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadUint32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse registerService response: %s", err)
	}
	return ret, nil
}

// UnregisterService calls the remote procedure
func (p *proxyServiceDirectory) UnregisterService(serviceID uint32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("serialize serviceID: %s", err)
	}
	_, err = p.Call("unregisterService", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterService failed: %s", err)
	}
	return nil
}

// ServiceReady calls the remote procedure
func (p *proxyServiceDirectory) ServiceReady(serviceID uint32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("serialize serviceID: %s", err)
	}
	_, err = p.Call("serviceReady", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call serviceReady failed: %s", err)
	}
	return nil
}

// UpdateServiceInfo calls the remote procedure
func (p *proxyServiceDirectory) UpdateServiceInfo(info ServiceInfo) error {
	var err error
	var buf bytes.Buffer
	if err = writeServiceInfo(info, &buf); err != nil {
		return fmt.Errorf("serialize info: %s", err)
	}
	_, err = p.Call("updateServiceInfo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call updateServiceInfo failed: %s", err)
	}
	return nil
}

// MachineId calls the remote procedure
func (p *proxyServiceDirectory) MachineId() (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	response, err := p.Call("machineId", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call machineId failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse machineId response: %s", err)
	}
	return ret, nil
}

// _socketOfService calls the remote procedure
func (p *proxyServiceDirectory) _socketOfService(serviceID uint32) (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf bytes.Buffer
	if err = basic.WriteUint32(serviceID, &buf); err != nil {
		return ret, fmt.Errorf("serialize serviceID: %s", err)
	}
	response, err := p.Call("_socketOfService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _socketOfService failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(resp)
	if err != nil {
		return ret, fmt.Errorf("parse _socketOfService response: %s", err)
	}
	return ret, nil
}

// SubscribeServiceAdded subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceAdded() (func(), chan ServiceAdded, error) {
	propertyID, err := p.SignalID("serviceAdded")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceAdded", err)
	}
	ch := make(chan ServiceAdded)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := readServiceAdded(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// SubscribeServiceRemoved subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceRemoved() (func(), chan ServiceRemoved, error) {
	propertyID, err := p.SignalID("serviceRemoved")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceRemoved", err)
	}
	ch := make(chan ServiceRemoved)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := readServiceRemoved(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// ServiceInfo is serializable
type ServiceInfo struct {
	Name      string
	ServiceId uint32
	MachineId string
	ProcessId uint32
	Endpoints []string
	SessionId string
}

// readServiceInfo unmarshalls ServiceInfo
func readServiceInfo(r io.Reader) (s ServiceInfo, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.ServiceId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceId field: %s", err)
	}
	if s.MachineId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read MachineId field: %s", err)
	}
	if s.ProcessId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ProcessId field: %s", err)
	}
	if s.Endpoints, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(r)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("read Endpoints field: %s", err)
	}
	if s.SessionId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read SessionId field: %s", err)
	}
	return s, nil
}

// writeServiceInfo marshalls ServiceInfo
func writeServiceInfo(s ServiceInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteUint32(s.ServiceId, w); err != nil {
		return fmt.Errorf("write ServiceId field: %s", err)
	}
	if err := basic.WriteString(s.MachineId, w); err != nil {
		return fmt.Errorf("write MachineId field: %s", err)
	}
	if err := basic.WriteUint32(s.ProcessId, w); err != nil {
		return fmt.Errorf("write ProcessId field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Endpoints)), w)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range s.Endpoints {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Endpoints field: %s", err)
	}
	if err := basic.WriteString(s.SessionId, w); err != nil {
		return fmt.Errorf("write SessionId field: %s", err)
	}
	return nil
}

// LogLevel is serializable
type LogLevel struct {
	Level int32
}

// readLogLevel unmarshalls LogLevel
func readLogLevel(r io.Reader) (s LogLevel, err error) {
	if s.Level, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("read Level field: %s", err)
	}
	return s, nil
}

// writeLogLevel marshalls LogLevel
func writeLogLevel(s LogLevel, w io.Writer) (err error) {
	if err := basic.WriteInt32(s.Level, w); err != nil {
		return fmt.Errorf("write Level field: %s", err)
	}
	return nil
}

// TimePoint is serializable
type TimePoint struct {
	Ns uint64
}

// readTimePoint unmarshalls TimePoint
func readTimePoint(r io.Reader) (s TimePoint, err error) {
	if s.Ns, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("read Ns field: %s", err)
	}
	return s, nil
}

// writeTimePoint marshalls TimePoint
func writeTimePoint(s TimePoint, w io.Writer) (err error) {
	if err := basic.WriteUint64(s.Ns, w); err != nil {
		return fmt.Errorf("write Ns field: %s", err)
	}
	return nil
}

// LogMessage is serializable
type LogMessage struct {
	Source     string
	Level      LogLevel
	Category   string
	Location   string
	Message    string
	Id         uint32
	Date       TimePoint
	SystemDate TimePoint
}

// readLogMessage unmarshalls LogMessage
func readLogMessage(r io.Reader) (s LogMessage, err error) {
	if s.Source, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Source field: %s", err)
	}
	if s.Level, err = readLogLevel(r); err != nil {
		return s, fmt.Errorf("read Level field: %s", err)
	}
	if s.Category, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Category field: %s", err)
	}
	if s.Location, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Location field: %s", err)
	}
	if s.Message, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Message field: %s", err)
	}
	if s.Id, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Id field: %s", err)
	}
	if s.Date, err = readTimePoint(r); err != nil {
		return s, fmt.Errorf("read Date field: %s", err)
	}
	if s.SystemDate, err = readTimePoint(r); err != nil {
		return s, fmt.Errorf("read SystemDate field: %s", err)
	}
	return s, nil
}

// writeLogMessage marshalls LogMessage
func writeLogMessage(s LogMessage, w io.Writer) (err error) {
	if err := basic.WriteString(s.Source, w); err != nil {
		return fmt.Errorf("write Source field: %s", err)
	}
	if err := writeLogLevel(s.Level, w); err != nil {
		return fmt.Errorf("write Level field: %s", err)
	}
	if err := basic.WriteString(s.Category, w); err != nil {
		return fmt.Errorf("write Category field: %s", err)
	}
	if err := basic.WriteString(s.Location, w); err != nil {
		return fmt.Errorf("write Location field: %s", err)
	}
	if err := basic.WriteString(s.Message, w); err != nil {
		return fmt.Errorf("write Message field: %s", err)
	}
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("write Id field: %s", err)
	}
	if err := writeTimePoint(s.Date, w); err != nil {
		return fmt.Errorf("write Date field: %s", err)
	}
	if err := writeTimePoint(s.SystemDate, w); err != nil {
		return fmt.Errorf("write SystemDate field: %s", err)
	}
	return nil
}

// LogProvider is the abstract interface of the service
type LogProvider interface {
	// SetVerbosity calls the remote procedure
	SetVerbosity(level LogLevel) error
	// SetCategory calls the remote procedure
	SetCategory(category string, level LogLevel) error
	// ClearAndSet calls the remote procedure
	ClearAndSet(filters map[string]int32) error
}

// LogProviderProxy represents a proxy object to the service
type LogProviderProxy interface {
	object.Object
	bus.Proxy
	LogProvider
}

// proxyLogProvider implements LogProviderProxy
type proxyLogProvider struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogProvider returns a specialized proxy.
func MakeLogProvider(sess bus.Session, proxy bus.Proxy) LogProviderProxy {
	return &proxyLogProvider{bus.MakeObject(proxy), sess}
}

// LogProvider returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) LogProvider(closer func(error)) (LogProviderProxy, error) {
	proxy, err := c.session.Proxy("LogProvider", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeLogProvider(c.session, proxy), nil
}

// SetVerbosity calls the remote procedure
func (p *proxyLogProvider) SetVerbosity(level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Call("setVerbosity", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setVerbosity failed: %s", err)
	}
	return nil
}

// SetCategory calls the remote procedure
func (p *proxyLogProvider) SetCategory(category string, level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(category, &buf); err != nil {
		return fmt.Errorf("serialize category: %s", err)
	}
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Call("setCategory", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearAndSet calls the remote procedure
func (p *proxyLogProvider) ClearAndSet(filters map[string]int32) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(filters)), &buf)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range filters {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize filters: %s", err)
	}
	_, err = p.Call("clearAndSet", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearAndSet failed: %s", err)
	}
	return nil
}

// LogListener is the abstract interface of the service
type LogListener interface {
	// SetCategory calls the remote procedure
	SetCategory(category string, level LogLevel) error
	// ClearFilters calls the remote procedure
	ClearFilters() error
	// SubscribeOnLogMessage subscribe to a remote signal
	SubscribeOnLogMessage() (unsubscribe func(), updates chan LogMessage, err error)
	// GetVerbosity returns the property value
	GetVerbosity() (LogLevel, error)
	// SetVerbosity sets the property value
	SetVerbosity(LogLevel) error
	// SubscribeVerbosity regusters to a property
	SubscribeVerbosity() (unsubscribe func(), updates chan LogLevel, err error)
	// GetFilters returns the property value
	GetFilters() (map[string]int32, error)
	// SetFilters sets the property value
	SetFilters(map[string]int32) error
	// SubscribeFilters regusters to a property
	SubscribeFilters() (unsubscribe func(), updates chan map[string]int32, err error)
}

// LogListenerProxy represents a proxy object to the service
type LogListenerProxy interface {
	object.Object
	bus.Proxy
	LogListener
}

// proxyLogListener implements LogListenerProxy
type proxyLogListener struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogListener returns a specialized proxy.
func MakeLogListener(sess bus.Session, proxy bus.Proxy) LogListenerProxy {
	return &proxyLogListener{bus.MakeObject(proxy), sess}
}

// LogListener returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) LogListener(closer func(error)) (LogListenerProxy, error) {
	proxy, err := c.session.Proxy("LogListener", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeLogListener(c.session, proxy), nil
}

// SetCategory calls the remote procedure
func (p *proxyLogListener) SetCategory(category string, level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(category, &buf); err != nil {
		return fmt.Errorf("serialize category: %s", err)
	}
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Call("setCategory", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearFilters calls the remote procedure
func (p *proxyLogListener) ClearFilters() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("clearFilters", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearFilters failed: %s", err)
	}
	return nil
}

// SubscribeOnLogMessage subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessage() (func(), chan LogMessage, error) {
	propertyID, err := p.SignalID("onLogMessage")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "onLogMessage", err)
	}
	ch := make(chan LogMessage)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := readLogMessage(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// GetVerbosity updates the property value
func (p *proxyLogListener) GetVerbosity() (ret LogLevel, err error) {
	name := value.String("verbosity")
	value, err := p.Property(name)
	if err != nil {
		return ret, fmt.Errorf("get property: %s", err)
	}
	var buf bytes.Buffer
	err = value.Write(&buf)
	if err != nil {
		return ret, fmt.Errorf("read response: %s", err)
	}
	s, err := basic.ReadString(&buf)
	if err != nil {
		return ret, fmt.Errorf("read signature: %s", err)
	}
	// check the signature
	sig := "(i)<LogLevel,level>"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = readLogLevel(&buf)
	return ret, err
}

// SetVerbosity updates the property value
func (p *proxyLogListener) SetVerbosity(update LogLevel) error {
	name := value.String("verbosity")
	var buf bytes.Buffer
	err := writeLogLevel(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("(i)<LogLevel,level>", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeVerbosity subscribe to a remote property
func (p *proxyLogListener) SubscribeVerbosity() (func(), chan LogLevel, error) {
	propertyID, err := p.PropertyID("verbosity")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "verbosity", err)
	}
	ch := make(chan LogLevel)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := readLogLevel(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// GetFilters updates the property value
func (p *proxyLogListener) GetFilters() (ret map[string]int32, err error) {
	name := value.String("filters")
	value, err := p.Property(name)
	if err != nil {
		return ret, fmt.Errorf("get property: %s", err)
	}
	var buf bytes.Buffer
	err = value.Write(&buf)
	if err != nil {
		return ret, fmt.Errorf("read response: %s", err)
	}
	s, err := basic.ReadString(&buf)
	if err != nil {
		return ret, fmt.Errorf("read signature: %s", err)
	}
	// check the signature
	sig := "{si}"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = func() (m map[string]int32, err error) {
		size, err := basic.ReadUint32(&buf)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[string]int32, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(&buf)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := basic.ReadInt32(&buf)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}()
	return ret, err
}

// SetFilters updates the property value
func (p *proxyLogListener) SetFilters(update map[string]int32) error {
	name := value.String("filters")
	var buf bytes.Buffer
	err := func() error {
		err := basic.WriteUint32(uint32(len(update)), &buf)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range update {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("{si}", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeFilters subscribe to a remote property
func (p *proxyLogListener) SubscribeFilters() (func(), chan map[string]int32, error) {
	propertyID, err := p.PropertyID("filters")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "filters", err)
	}
	ch := make(chan map[string]int32)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (m map[string]int32, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return m, fmt.Errorf("read map size: %s", err)
				}
				m = make(map[string]int32, size)
				for i := 0; i < int(size); i++ {
					k, err := basic.ReadString(buf)
					if err != nil {
						return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
					}
					v, err := basic.ReadInt32(buf)
					if err != nil {
						return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
					}
					m[k] = v
				}
				return m, nil
			}()
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// LogManager is the abstract interface of the service
type LogManager interface {
	// Log calls the remote procedure
	Log(messages []LogMessage) error
	// CreateListener calls the remote procedure
	CreateListener() (LogListenerProxy, error)
	// GetListener calls the remote procedure
	GetListener() (LogListenerProxy, error)
	// AddProvider calls the remote procedure
	AddProvider(source LogProviderProxy) (int32, error)
	// RemoveProvider calls the remote procedure
	RemoveProvider(providerID int32) error
}

// LogManagerProxy represents a proxy object to the service
type LogManagerProxy interface {
	object.Object
	bus.Proxy
	LogManager
}

// proxyLogManager implements LogManagerProxy
type proxyLogManager struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogManager returns a specialized proxy.
func MakeLogManager(sess bus.Session, proxy bus.Proxy) LogManagerProxy {
	return &proxyLogManager{bus.MakeObject(proxy), sess}
}

// LogManager returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) LogManager(closer func(error)) (LogManagerProxy, error) {
	proxy, err := c.session.Proxy("LogManager", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeLogManager(c.session, proxy), nil
}

// Log calls the remote procedure
func (p *proxyLogManager) Log(messages []LogMessage) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(messages)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range messages {
			err = writeLogMessage(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize messages: %s", err)
	}
	_, err = p.Call("log", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call log failed: %s", err)
	}
	return nil
}

// CreateListener calls the remote procedure
func (p *proxyLogManager) CreateListener() (LogListenerProxy, error) {
	var err error
	var ret LogListenerProxy
	var buf bytes.Buffer
	response, err := p.Call("createListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call createListener failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (LogListenerProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("get proxy: %s", err)
		}
		return MakeLogListener(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse createListener response: %s", err)
	}
	return ret, nil
}

// GetListener calls the remote procedure
func (p *proxyLogManager) GetListener() (LogListenerProxy, error) {
	var err error
	var ret LogListenerProxy
	var buf bytes.Buffer
	response, err := p.Call("getListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getListener failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (LogListenerProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("get proxy: %s", err)
		}
		return MakeLogListener(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getListener response: %s", err)
	}
	return ret, nil
}

// AddProvider calls the remote procedure
func (p *proxyLogManager) AddProvider(source LogProviderProxy) (int32, error) {
	var err error
	var ret int32
	var buf bytes.Buffer
	if err = func() error {
		meta, err := source.MetaObject(source.ObjectID())
		if err != nil {
			return fmt.Errorf("get meta: %s", err)
		}
		ref := object.ObjectReference{
			MetaObject: meta,
			ServiceID:  source.ServiceID(),
			ObjectID:   source.ObjectID(),
		}
		return object.WriteObjectReference(ref, &buf)
	}(); err != nil {
		return ret, fmt.Errorf("serialize source: %s", err)
	}
	response, err := p.Call("addProvider", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call addProvider failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse addProvider response: %s", err)
	}
	return ret, nil
}

// RemoveProvider calls the remote procedure
func (p *proxyLogManager) RemoveProvider(providerID int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteInt32(providerID, &buf); err != nil {
		return fmt.Errorf("serialize providerID: %s", err)
	}
	_, err = p.Call("removeProvider", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
}

// ALTextToSpeech is the abstract interface of the service
type ALTextToSpeech interface {
	// Say calls the remote procedure
	Say(stringToSay string) error
}

// ALTextToSpeechProxy represents a proxy object to the service
type ALTextToSpeechProxy interface {
	object.Object
	bus.Proxy
	ALTextToSpeech
}

// proxyALTextToSpeech implements ALTextToSpeechProxy
type proxyALTextToSpeech struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALTextToSpeech returns a specialized proxy.
func MakeALTextToSpeech(sess bus.Session, proxy bus.Proxy) ALTextToSpeechProxy {
	return &proxyALTextToSpeech{bus.MakeObject(proxy), sess}
}

// ALTextToSpeech returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) ALTextToSpeech(closer func(error)) (ALTextToSpeechProxy, error) {
	proxy, err := c.session.Proxy("ALTextToSpeech", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeALTextToSpeech(c.session, proxy), nil
}

// Say calls the remote procedure
func (p *proxyALTextToSpeech) Say(stringToSay string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(stringToSay, &buf); err != nil {
		return fmt.Errorf("serialize stringToSay: %s", err)
	}
	_, err = p.Call("say", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call say failed: %s", err)
	}
	return nil
}

// ALAnimatedSpeech is the abstract interface of the service
type ALAnimatedSpeech interface {
	// Say calls the remote procedure
	Say(text string) error
	// IsBodyTalkEnabled calls the remote procedure
	IsBodyTalkEnabled() (bool, error)
	// IsBodyLanguageEnabled calls the remote procedure
	IsBodyLanguageEnabled() (bool, error)
	// SetBodyTalkEnabled calls the remote procedure
	SetBodyTalkEnabled(enable bool) error
	// SetBodyLanguageEnabled calls the remote procedure
	SetBodyLanguageEnabled(enable bool) error
}

// ALAnimatedSpeechProxy represents a proxy object to the service
type ALAnimatedSpeechProxy interface {
	object.Object
	bus.Proxy
	ALAnimatedSpeech
}

// proxyALAnimatedSpeech implements ALAnimatedSpeechProxy
type proxyALAnimatedSpeech struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALAnimatedSpeech returns a specialized proxy.
func MakeALAnimatedSpeech(sess bus.Session, proxy bus.Proxy) ALAnimatedSpeechProxy {
	return &proxyALAnimatedSpeech{bus.MakeObject(proxy), sess}
}

// ALAnimatedSpeech returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) ALAnimatedSpeech(closer func(error)) (ALAnimatedSpeechProxy, error) {
	proxy, err := c.session.Proxy("ALAnimatedSpeech", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeALAnimatedSpeech(c.session, proxy), nil
}

// Say calls the remote procedure
func (p *proxyALAnimatedSpeech) Say(text string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(text, &buf); err != nil {
		return fmt.Errorf("serialize text: %s", err)
	}
	_, err = p.Call("say", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call say failed: %s", err)
	}
	return nil
}

// IsBodyTalkEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) IsBodyTalkEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("isBodyTalkEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBodyTalkEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBodyTalkEnabled response: %s", err)
	}
	return ret, nil
}

// IsBodyLanguageEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) IsBodyLanguageEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("isBodyLanguageEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBodyLanguageEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBodyLanguageEnabled response: %s", err)
	}
	return ret, nil
}

// SetBodyTalkEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) SetBodyTalkEnabled(enable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enable, &buf); err != nil {
		return fmt.Errorf("serialize enable: %s", err)
	}
	_, err = p.Call("setBodyTalkEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setBodyTalkEnabled failed: %s", err)
	}
	return nil
}

// SetBodyLanguageEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) SetBodyLanguageEnabled(enable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enable, &buf); err != nil {
		return fmt.Errorf("serialize enable: %s", err)
	}
	_, err = p.Call("setBodyLanguageEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setBodyLanguageEnabled failed: %s", err)
	}
	return nil
}

// BehaviorFailed is serializable
type BehaviorFailed struct {
	P0 string
	P1 string
	P2 string
}

// readBehaviorFailed unmarshalls BehaviorFailed
func readBehaviorFailed(r io.Reader) (s BehaviorFailed, err error) {
	if s.P0, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read P0 field: %s", err)
	}
	if s.P1, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read P1 field: %s", err)
	}
	if s.P2, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read P2 field: %s", err)
	}
	return s, nil
}

// writeBehaviorFailed marshalls BehaviorFailed
func writeBehaviorFailed(s BehaviorFailed, w io.Writer) (err error) {
	if err := basic.WriteString(s.P0, w); err != nil {
		return fmt.Errorf("write P0 field: %s", err)
	}
	if err := basic.WriteString(s.P1, w); err != nil {
		return fmt.Errorf("write P1 field: %s", err)
	}
	if err := basic.WriteString(s.P2, w); err != nil {
		return fmt.Errorf("write P2 field: %s", err)
	}
	return nil
}

// ALBehaviorManager is the abstract interface of the service
type ALBehaviorManager interface {
	// InstallBehavior calls the remote procedure
	InstallBehavior(localPath string) (bool, error)
	// PreloadBehavior calls the remote procedure
	PreloadBehavior(behavior string) (bool, error)
	// StartBehavior calls the remote procedure
	StartBehavior(behavior string) error
	// RunBehavior calls the remote procedure
	RunBehavior(behavior string) error
	// StopBehavior calls the remote procedure
	StopBehavior(behavior string) error
	// StopAllBehaviors calls the remote procedure
	StopAllBehaviors() error
	// RemoveBehavior calls the remote procedure
	RemoveBehavior(behavior string) (bool, error)
	// IsBehaviorInstalled calls the remote procedure
	IsBehaviorInstalled(name string) (bool, error)
	// IsBehaviorPresent calls the remote procedure
	IsBehaviorPresent(prefixedBehavior string) (bool, error)
	// GetBehaviorNames calls the remote procedure
	GetBehaviorNames() ([]string, error)
	// GetUserBehaviorNames calls the remote procedure
	GetUserBehaviorNames() ([]string, error)
	// GetSystemBehaviorNames calls the remote procedure
	GetSystemBehaviorNames() ([]string, error)
	// GetInstalledBehaviors calls the remote procedure
	GetInstalledBehaviors() ([]string, error)
	// GetBehaviorsByTag calls the remote procedure
	GetBehaviorsByTag(tag string) ([]string, error)
	// IsBehaviorRunning calls the remote procedure
	IsBehaviorRunning(behavior string) (bool, error)
	// IsBehaviorLoaded calls the remote procedure
	IsBehaviorLoaded(behavior string) (bool, error)
	// GetRunningBehaviors calls the remote procedure
	GetRunningBehaviors() ([]string, error)
	// GetLoadedBehaviors calls the remote procedure
	GetLoadedBehaviors() ([]string, error)
	// GetTagList calls the remote procedure
	GetTagList() ([]string, error)
	// GetBehaviorTags calls the remote procedure
	GetBehaviorTags(behavior string) ([]string, error)
	// GetBehaviorNature calls the remote procedure
	GetBehaviorNature(behavior string) (string, error)
	// AddDefaultBehavior calls the remote procedure
	AddDefaultBehavior(behavior string) error
	// RemoveDefaultBehavior calls the remote procedure
	RemoveDefaultBehavior(behavior string) error
	// GetDefaultBehaviors calls the remote procedure
	GetDefaultBehaviors() ([]string, error)
	// PlayDefaultProject calls the remote procedure
	PlayDefaultProject() error
	// SubscribeBehaviorsAdded subscribe to a remote signal
	SubscribeBehaviorsAdded() (unsubscribe func(), updates chan []string, err error)
	// SubscribeBehaviorsRemoved subscribe to a remote signal
	SubscribeBehaviorsRemoved() (unsubscribe func(), updates chan []string, err error)
	// SubscribeBehaviorLoaded subscribe to a remote signal
	SubscribeBehaviorLoaded() (unsubscribe func(), updates chan string, err error)
	// SubscribeBehaviorStarted subscribe to a remote signal
	SubscribeBehaviorStarted() (unsubscribe func(), updates chan string, err error)
	// SubscribeBehaviorStopped subscribe to a remote signal
	SubscribeBehaviorStopped() (unsubscribe func(), updates chan string, err error)
	// SubscribeBehaviorFailed subscribe to a remote signal
	SubscribeBehaviorFailed() (unsubscribe func(), updates chan BehaviorFailed, err error)
}

// ALBehaviorManagerProxy represents a proxy object to the service
type ALBehaviorManagerProxy interface {
	object.Object
	bus.Proxy
	ALBehaviorManager
}

// proxyALBehaviorManager implements ALBehaviorManagerProxy
type proxyALBehaviorManager struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALBehaviorManager returns a specialized proxy.
func MakeALBehaviorManager(sess bus.Session, proxy bus.Proxy) ALBehaviorManagerProxy {
	return &proxyALBehaviorManager{bus.MakeObject(proxy), sess}
}

// ALBehaviorManager returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) ALBehaviorManager(closer func(error)) (ALBehaviorManagerProxy, error) {
	proxy, err := c.session.Proxy("ALBehaviorManager", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeALBehaviorManager(c.session, proxy), nil
}

// InstallBehavior calls the remote procedure
func (p *proxyALBehaviorManager) InstallBehavior(localPath string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(localPath, &buf); err != nil {
		return ret, fmt.Errorf("serialize localPath: %s", err)
	}
	response, err := p.Call("installBehavior", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call installBehavior failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse installBehavior response: %s", err)
	}
	return ret, nil
}

// PreloadBehavior calls the remote procedure
func (p *proxyALBehaviorManager) PreloadBehavior(behavior string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize behavior: %s", err)
	}
	response, err := p.Call("preloadBehavior", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call preloadBehavior failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse preloadBehavior response: %s", err)
	}
	return ret, nil
}

// StartBehavior calls the remote procedure
func (p *proxyALBehaviorManager) StartBehavior(behavior string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return fmt.Errorf("serialize behavior: %s", err)
	}
	_, err = p.Call("startBehavior", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call startBehavior failed: %s", err)
	}
	return nil
}

// RunBehavior calls the remote procedure
func (p *proxyALBehaviorManager) RunBehavior(behavior string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return fmt.Errorf("serialize behavior: %s", err)
	}
	_, err = p.Call("runBehavior", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call runBehavior failed: %s", err)
	}
	return nil
}

// StopBehavior calls the remote procedure
func (p *proxyALBehaviorManager) StopBehavior(behavior string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return fmt.Errorf("serialize behavior: %s", err)
	}
	_, err = p.Call("stopBehavior", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stopBehavior failed: %s", err)
	}
	return nil
}

// StopAllBehaviors calls the remote procedure
func (p *proxyALBehaviorManager) StopAllBehaviors() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("stopAllBehaviors", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stopAllBehaviors failed: %s", err)
	}
	return nil
}

// RemoveBehavior calls the remote procedure
func (p *proxyALBehaviorManager) RemoveBehavior(behavior string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize behavior: %s", err)
	}
	response, err := p.Call("removeBehavior", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call removeBehavior failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse removeBehavior response: %s", err)
	}
	return ret, nil
}

// IsBehaviorInstalled calls the remote procedure
func (p *proxyALBehaviorManager) IsBehaviorInstalled(name string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Call("isBehaviorInstalled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBehaviorInstalled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBehaviorInstalled response: %s", err)
	}
	return ret, nil
}

// IsBehaviorPresent calls the remote procedure
func (p *proxyALBehaviorManager) IsBehaviorPresent(prefixedBehavior string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(prefixedBehavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize prefixedBehavior: %s", err)
	}
	response, err := p.Call("isBehaviorPresent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBehaviorPresent failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBehaviorPresent response: %s", err)
	}
	return ret, nil
}

// GetBehaviorNames calls the remote procedure
func (p *proxyALBehaviorManager) GetBehaviorNames() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getBehaviorNames", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBehaviorNames failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getBehaviorNames response: %s", err)
	}
	return ret, nil
}

// GetUserBehaviorNames calls the remote procedure
func (p *proxyALBehaviorManager) GetUserBehaviorNames() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getUserBehaviorNames", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getUserBehaviorNames failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getUserBehaviorNames response: %s", err)
	}
	return ret, nil
}

// GetSystemBehaviorNames calls the remote procedure
func (p *proxyALBehaviorManager) GetSystemBehaviorNames() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getSystemBehaviorNames", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getSystemBehaviorNames failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getSystemBehaviorNames response: %s", err)
	}
	return ret, nil
}

// GetInstalledBehaviors calls the remote procedure
func (p *proxyALBehaviorManager) GetInstalledBehaviors() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getInstalledBehaviors", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getInstalledBehaviors failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getInstalledBehaviors response: %s", err)
	}
	return ret, nil
}

// GetBehaviorsByTag calls the remote procedure
func (p *proxyALBehaviorManager) GetBehaviorsByTag(tag string) ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	if err = basic.WriteString(tag, &buf); err != nil {
		return ret, fmt.Errorf("serialize tag: %s", err)
	}
	response, err := p.Call("getBehaviorsByTag", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBehaviorsByTag failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getBehaviorsByTag response: %s", err)
	}
	return ret, nil
}

// IsBehaviorRunning calls the remote procedure
func (p *proxyALBehaviorManager) IsBehaviorRunning(behavior string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize behavior: %s", err)
	}
	response, err := p.Call("isBehaviorRunning", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBehaviorRunning failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBehaviorRunning response: %s", err)
	}
	return ret, nil
}

// IsBehaviorLoaded calls the remote procedure
func (p *proxyALBehaviorManager) IsBehaviorLoaded(behavior string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize behavior: %s", err)
	}
	response, err := p.Call("isBehaviorLoaded", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBehaviorLoaded failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBehaviorLoaded response: %s", err)
	}
	return ret, nil
}

// GetRunningBehaviors calls the remote procedure
func (p *proxyALBehaviorManager) GetRunningBehaviors() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getRunningBehaviors", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getRunningBehaviors failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getRunningBehaviors response: %s", err)
	}
	return ret, nil
}

// GetLoadedBehaviors calls the remote procedure
func (p *proxyALBehaviorManager) GetLoadedBehaviors() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getLoadedBehaviors", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getLoadedBehaviors failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getLoadedBehaviors response: %s", err)
	}
	return ret, nil
}

// GetTagList calls the remote procedure
func (p *proxyALBehaviorManager) GetTagList() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getTagList", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getTagList failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getTagList response: %s", err)
	}
	return ret, nil
}

// GetBehaviorTags calls the remote procedure
func (p *proxyALBehaviorManager) GetBehaviorTags(behavior string) ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize behavior: %s", err)
	}
	response, err := p.Call("getBehaviorTags", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBehaviorTags failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getBehaviorTags response: %s", err)
	}
	return ret, nil
}

// GetBehaviorNature calls the remote procedure
func (p *proxyALBehaviorManager) GetBehaviorNature(behavior string) (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return ret, fmt.Errorf("serialize behavior: %s", err)
	}
	response, err := p.Call("getBehaviorNature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBehaviorNature failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getBehaviorNature response: %s", err)
	}
	return ret, nil
}

// AddDefaultBehavior calls the remote procedure
func (p *proxyALBehaviorManager) AddDefaultBehavior(behavior string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return fmt.Errorf("serialize behavior: %s", err)
	}
	_, err = p.Call("addDefaultBehavior", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call addDefaultBehavior failed: %s", err)
	}
	return nil
}

// RemoveDefaultBehavior calls the remote procedure
func (p *proxyALBehaviorManager) RemoveDefaultBehavior(behavior string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(behavior, &buf); err != nil {
		return fmt.Errorf("serialize behavior: %s", err)
	}
	_, err = p.Call("removeDefaultBehavior", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeDefaultBehavior failed: %s", err)
	}
	return nil
}

// GetDefaultBehaviors calls the remote procedure
func (p *proxyALBehaviorManager) GetDefaultBehaviors() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getDefaultBehaviors", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getDefaultBehaviors failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getDefaultBehaviors response: %s", err)
	}
	return ret, nil
}

// PlayDefaultProject calls the remote procedure
func (p *proxyALBehaviorManager) PlayDefaultProject() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("playDefaultProject", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call playDefaultProject failed: %s", err)
	}
	return nil
}

// SubscribeBehaviorsAdded subscribe to a remote property
func (p *proxyALBehaviorManager) SubscribeBehaviorsAdded() (func(), chan []string, error) {
	propertyID, err := p.SignalID("behaviorsAdded")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "behaviorsAdded", err)
	}
	handlerID := rand.Uint64()

	_, err = p.RegisterEvent(p.ObjectID(), propertyID, handlerID)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "behaviorsAdded", err)
	}
	ch := make(chan []string)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (b []string, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return b, fmt.Errorf("read slice size: %s", err)
				}
				b = make([]string, size)
				for i := 0; i < int(size); i++ {
					b[i], err = basic.ReadString(buf)
					if err != nil {
						return b, fmt.Errorf("read slice value: %s", err)
					}
				}
				return b, nil
			}()
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}

// SubscribeBehaviorsRemoved subscribe to a remote property
func (p *proxyALBehaviorManager) SubscribeBehaviorsRemoved() (func(), chan []string, error) {
	propertyID, err := p.SignalID("behaviorsRemoved")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "behaviorsRemoved", err)
	}
	handlerID := rand.Uint64()

	_, err = p.RegisterEvent(p.ObjectID(), propertyID, handlerID)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "behaviorsRemoved", err)
	}
	ch := make(chan []string)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (b []string, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return b, fmt.Errorf("read slice size: %s", err)
				}
				b = make([]string, size)
				for i := 0; i < int(size); i++ {
					b[i], err = basic.ReadString(buf)
					if err != nil {
						return b, fmt.Errorf("read slice value: %s", err)
					}
				}
				return b, nil
			}()
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}

// SubscribeBehaviorLoaded subscribe to a remote property
func (p *proxyALBehaviorManager) SubscribeBehaviorLoaded() (func(), chan string, error) {
	propertyID, err := p.SignalID("behaviorLoaded")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "behaviorLoaded", err)
	}
	handlerID := rand.Uint64()

	_, err = p.RegisterEvent(p.ObjectID(), propertyID, handlerID)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "behaviorLoaded", err)
	}
	ch := make(chan string)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := basic.ReadString(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}

// SubscribeBehaviorStarted subscribe to a remote property
func (p *proxyALBehaviorManager) SubscribeBehaviorStarted() (func(), chan string, error) {
	propertyID, err := p.SignalID("behaviorStarted")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "behaviorStarted", err)
	}
	handlerID := rand.Uint64()

	_, err = p.RegisterEvent(p.ObjectID(), propertyID, handlerID)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "behaviorStarted", err)
	}
	ch := make(chan string)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := basic.ReadString(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}

// SubscribeBehaviorStopped subscribe to a remote property
func (p *proxyALBehaviorManager) SubscribeBehaviorStopped() (func(), chan string, error) {
	propertyID, err := p.SignalID("behaviorStopped")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "behaviorStopped", err)
	}
	handlerID := rand.Uint64()

	_, err = p.RegisterEvent(p.ObjectID(), propertyID, handlerID)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "behaviorStopped", err)
	}
	ch := make(chan string)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := basic.ReadString(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}

// SubscribeBehaviorFailed subscribe to a remote property
func (p *proxyALBehaviorManager) SubscribeBehaviorFailed() (func(), chan BehaviorFailed, error) {
	propertyID, err := p.SignalID("behaviorFailed")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "behaviorFailed", err)
	}
	handlerID := rand.Uint64()

	_, err = p.RegisterEvent(p.ObjectID(), propertyID, handlerID)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "behaviorFailed", err)
	}
	ch := make(chan BehaviorFailed)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := readBehaviorFailed(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}
