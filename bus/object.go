package bus

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
)

// ErrWrongObjectID is returned when a method argument is given the
// wrong object ID.
var ErrWrongObjectID = errors.New("Wrong object ID")

// ErrNotYetImplemented is returned when a feature is not yet
// implemented.
var ErrNotYetImplemented = errors.New("Not supported")

func (s *stubObject) UpdateSignal(signal uint32, data []byte) error {
	return s.signal.UpdateSignal(signal, data)
}

func (s *stubObject) UpdateProperty(id uint32, sig string, data []byte) error {
	objImpl, ok := (s.impl).(*objectImpl)
	if !ok {
		return fmt.Errorf("unexpected implementation")
	}
	prop, ok := objImpl.meta.Properties[id]
	if !ok {
		return fmt.Errorf("missing property (%d), %#v", id,
			objImpl.meta)
	}
	err := objImpl.onPropertyChange(prop.Name, data)
	if err != nil {
		return err
	}
	newValue := value.Opaque(sig, data)
	err = objImpl.saveProperty(prop.Name, newValue)
	if err != nil {
		return err
	}
	return s.signal.UpdateProperty(id, sig, data)
}

type objectImpl struct {
	signalHandler    *signalHandler
	meta             object.MetaObject
	onPropertyChange func(string, []byte) error
	objectID         uint32
	signal           ObjectSignalHelper
	properties       map[string]value.Value
	propertiesMutex  sync.RWMutex
	terminate        func()
	stats            map[uint32]MethodStatistics
	statsMutex       sync.RWMutex
	statsEnabled     bool
	traceEnabled     bool
	nextTrace        uint32
}

// NewBasicObject returns an BasicObject which implements Actor. It
// handles all the generic methods and signals common to all objects.
// onPropertyChange is called each time a property is udpated.
func NewBasicObject(obj Actor, meta object.MetaObject,
	onPropertyChange func(string, []byte) error) BasicObject {

	impl := &objectImpl{
		meta:             object.FullMetaObject(meta),
		onPropertyChange: onPropertyChange,
		signalHandler:    newSignalHandler(),
		properties:       make(map[string]value.Value),
		stats:            make(map[uint32]MethodStatistics),
	}

	for uid, _ := range impl.meta.Methods {
		var m MethodStatistics
		impl.stats[uid] = m
	}

	return &stubObject{
		impl:   impl,
		obj:    obj,
		signal: impl.signalHandler,
	}
}

func (o *objectImpl) Activate(activation Activation,
	signal ObjectSignalHelper) error {

	o.signal = signal
	o.objectID = activation.ObjectID
	o.terminate = activation.Terminate

	return nil
}

func (o *objectImpl) OnTerminate() {
	o.signalHandler.OnTerminate()
}

func (o *objectImpl) RegisterEvent(msg *net.Message, from Channel) error {

	buf := bytes.NewBuffer(msg.Payload)
	_, err := basic.ReadUint32(buf)
	if err == nil {
		signalID, err := basic.ReadUint32(buf)
		if err == nil && signalID == 0x56 {
			o.EnableTrace(true)
		}
	}

	return o.signalHandler.RegisterEvent(msg, from)
}

func (o *objectImpl) UnregisterEvent(msg *net.Message, from Channel) error {
	return o.signalHandler.UnregisterEvent(msg, from)
}

func (o *objectImpl) MetaObject(objectID uint32) (object.MetaObject, error) {
	// remote objects don't know their real object id.
	if objectID != 0 && o.objectID < (1<<31) && objectID != o.objectID {
		return o.meta, ErrWrongObjectID
	}
	return o.meta, nil
}

func (o *objectImpl) Terminate(objectID uint32) error {
	// remote objects don't know their real object id.
	if objectID != 0 && o.objectID < (1<<31) && objectID != o.objectID {
		return ErrWrongObjectID
	}
	o.terminate()
	return nil
}

func (o *objectImpl) Property(name value.Value) (value.Value, error) {
	stringValue, ok := name.(value.StringValue)
	if !ok {
		return nil, fmt.Errorf("property name must be a string value")
	}
	nameStr := stringValue.Value()
	o.propertiesMutex.RLock()
	defer o.propertiesMutex.RUnlock()
	val, ok := o.properties[nameStr]
	if !ok {
		return nil, fmt.Errorf("property unknown: %s, %#v", nameStr,
			o.properties)
	}
	return val, nil
}

func (o *objectImpl) SetProperty(name value.Value, newValue value.Value) error {
	var nameStr string
	stringValue, ok := name.(value.StringValue)
	if ok {
		nameStr = stringValue.Value()
	} else {
		idValue, ok := name.(value.UintValue)
		if !ok {
			return fmt.Errorf("incorrect name type")
		}
		property, ok := o.meta.Properties[idValue.Value()]
		if !ok {
			return fmt.Errorf(
				"incorrect property id value, got %d",
				idValue.Value())
		}
		nameStr = property.Name
	}
	var buf bytes.Buffer
	err := newValue.Write(&buf)
	if err != nil {
		return fmt.Errorf("cannot write value: %s", err)
	}
	sig, err := basic.ReadString(&buf)
	if err != nil {
		return fmt.Errorf("invalid signature: %s", err)
	}
	data := buf.Bytes()
	err = o.onPropertyChange(nameStr, data)
	if err != nil {
		return err
	}
	err = o.saveProperty(nameStr, newValue)
	if err != nil {
		return err
	}
	id, err := o.meta.PropertyID(nameStr)
	if err != nil {
		return fmt.Errorf("cannot set property: %s", err)
	}
	return o.signalHandler.UpdateProperty(id, sig, data)
}

func (o *objectImpl) saveProperty(name string, newValue value.Value) error {
	o.propertiesMutex.Lock()
	defer o.propertiesMutex.Unlock()
	o.properties[name] = newValue
	return nil
}

func (o *objectImpl) Properties() ([]string, error) {
	properties := make([]string, 0)
	o.propertiesMutex.RLock()
	defer o.propertiesMutex.RUnlock()
	for property := range o.properties {
		properties = append(properties, property)
	}
	return properties, nil
}

func (o *objectImpl) RegisterEventWithSignature(objectID uint32,
	actionID uint32, handler uint64, P3 string) (uint64, error) {
	return 0, fmt.Errorf("Not yet implemented")
}

func (o *objectImpl) IsStatsEnabled() (bool, error) {
	return o.statsEnabled, nil
}

func (m MethodStatistics) updateWith(t time.Duration) MethodStatistics {
	m.Count++
	duration := float32(t.Seconds())
	m.Wall.CumulatedValue += duration
	if m.Wall.MinValue == 0 || duration < m.Wall.MinValue {
		m.Wall.MinValue = duration
	}
	if duration > m.Wall.MaxValue {
		m.Wall.MaxValue = duration
	}
	return m
}

func (o *objectImpl) EnableStats(enabled bool) error {
	o.statsEnabled = enabled
	return nil
}

func (o *objectImpl) Stats() (map[uint32]MethodStatistics, error) {
	o.statsMutex.Lock()
	defer o.statsMutex.Unlock()
	stats := make(map[uint32]MethodStatistics)
	for id, stat := range o.stats {
		stats[id] = stat
	}
	return stats, nil
}

func (o *objectImpl) ClearStats() error {
	o.statsMutex.Lock()
	defer o.statsMutex.Unlock()
	o.stats = make(map[uint32]MethodStatistics)
	for uid, _ := range o.meta.Methods {
		var m MethodStatistics
		o.stats[uid] = m
	}
	return nil
}

func (o *objectImpl) IsTraceEnabled() (bool, error) {
	return o.traceEnabled, nil
}

func (o *objectImpl) EnableTrace(enable bool) error {
	o.traceEnabled = enable
	return nil
}

// Tracer records the arrival or departure of the message.
type Tracer interface {
	Trace(msg *net.Message, id uint32)
}

func signature(msg *net.Message, meta *object.MetaObject) string {
	if msg.Header.Type == net.Call ||
		msg.Header.Type == net.Post ||
		msg.Header.Type == net.Reply {
		m, ok := meta.Methods[msg.Header.Action]
		if !ok {
			return "X"
		}
		if msg.Header.Type == net.Reply {
			return m.ReturnSignature
		}
		return m.ParametersSignature
	} else if msg.Header.Type == net.Event {
		s, ok := meta.Signals[msg.Header.Action]
		if ok {
			return s.Signature
		}
		p, ok := meta.Properties[msg.Header.Action]
		if ok {
			return p.Signature
		}
	}
	return "X"
}

func (o *objectImpl) Trace(msg *net.Message, id uint32) {

	// do not trace traceObject signal
	if msg.Header.Action == 0x56 {
		return
	}

	now := time.Now()
	arguments := value.Opaque(signature(msg, &o.meta), msg.Payload)
	timeval := Timeval{
		Tv_sec:  int64(now.Second()),
		Tv_usec: int64(now.Nanosecond() / 1000),
	}
	event := EventTrace{
		Id:        o.nextTrace,
		Kind:      int32(msg.Header.Type),
		SlotId:    msg.Header.Action,
		Arguments: arguments,
		Timestamp: timeval,
	}
	err := o.signal.SignalTraceObject(event)
	if err != nil {
		log.Printf("trace error: %s", err)
	}
}

func (o *objectImpl) updateMethodStatistics(uid uint32, d time.Duration) {
	o.statsMutex.Lock()
	defer o.statsMutex.Unlock()
	stat, ok := o.stats[uid]
	if ok {
		o.stats[uid] = stat.updateWith(d)
	}
}

func (o *objectImpl) Tracer(msg *net.Message, from Channel) Channel {

	if o.statsEnabled {
		from = &statChannel{from, time.Now(), o}
	}

	if !o.traceEnabled {
		return from
	}

	traceID := o.nextTrace
	o.nextTrace++
	o.Trace(msg, traceID)

	return &tracedChannel{from, o, traceID}
}

// clientObject implements Actor. It is used to forward incomming
// messages to a remote client object.
type clientObject struct {
	serviceID uint32
	remoteID  uint32
	channel   Channel
	client    Client
}

// NewClientObject returns an Actor which forwards messages to a
// remote object.
//
// TODO: Why this can't work correctly:
// - Pb 1: the client side does not know its public object id. this
// means for methods like MetaObject(id uint32), it can't correctly
// compare its object id with the one embedded in the payload. in
// other words: rewritting the header is not enougth for actions which
// includes the object id (ex: metaObject, terminate, register event,
// ...). Work around: do not check object id in such cases.
// - Pb 2: the client side does not know its public object id. this
// means it cannot share this id with different services. in other
// words: each time the object is shared with a service it must
// register again to this service using the 2^31 tricks.
//
// The obvious solution is to inform the remote object of its true
// identity. This can be done using a new method LendObjectID() uint32
// to the services: those services would dedicate an id on demand and
// route the traffic at this stage. The process would go like:
// 1. client side request the service to lend her an object id.
// 2. service allocate an id for the client and setup a route.
// 3. client side can share its "official" reference to anyone.
func NewClientObject(remoteID uint32, from Channel) Actor {
	return &clientObject{
		remoteID: remoteID,
		channel:  from,
		client:   NewClient(from.EndPoint()),
	}
}

func (c *clientObject) handleRegister(msg *net.Message, from Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return from.SendError(msg, err)
	}
	if objectID != c.remoteID {
		err = fmt.Errorf("invalid object ID: %d instead of %d",
			objectID, c.remoteID)
		return from.SendError(msg, err)
	}
	signalID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read signal uid: %s", err)
		return from.SendError(msg, err)
	}
	// FIXME: hook to the unregister message and call cancel.
	_, events, err := c.client.Subscribe(msg.Header.Service,
		c.remoteID, signalID)
	if err != nil {
		return from.SendError(msg, err)
	}
	for event := range events {
		m := net.NewMessage(msg.Header, event)
		m.Header.Type = net.Event
		return from.Send(&m)
	}
	return nil
}

func (c *clientObject) handleCall(msg *net.Message, from Channel) error {
	resp, err := c.client.Call(msg.Header.Service, c.remoteID,
		msg.Header.Action, msg.Payload)
	if err != nil {
		return from.SendError(msg, err)
	}
	return from.SendReply(msg, resp)
}

func (c *clientObject) Receive(msg *net.Message, from Channel) error {
	// call to RegisterEvent
	if msg.Header.Type == net.Call && msg.Header.Action == 0x0 {
		go c.handleRegister(msg, from)
		return nil
	} else if msg.Header.Type == net.Call {
		go c.handleCall(msg, from)
		return nil
	} else if msg.Header.Type == net.Post {
		msg.Header.Object = c.remoteID
		return c.channel.Send(msg)
	}
	return from.SendError(msg, fmt.Errorf("unexpected message type: %#v", msg))
}
func (c *clientObject) Activate(activation Activation) error {
	c.serviceID = activation.ServiceID
	return nil
}

func (c *clientObject) OnTerminate() {
}
