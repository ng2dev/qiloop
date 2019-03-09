// Package tester contains a generated stub
// File generated. DO NOT EDIT.
package tester

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	generic "github.com/lugu/qiloop/bus/server/generic"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
)

// BombImplementor interface of the service implementation
type BombImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper BombSignalHelper) error
	OnTerminate()
}

// BombSignalHelper provided to Bomb a companion object
type BombSignalHelper interface {
	SignalBoom(energy int32) error
}

// stubBomb implements server.ServerObject.
type stubBomb struct {
	obj     generic.Object
	impl    BombImplementor
	session bus.Session
}

// BombObject returns an object using BombImplementor
func BombObject(impl BombImplementor) server.ServerObject {
	var stb stubBomb
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	return &stb
}

func (p *stubBomb) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubBomb) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubBomb) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubBomb) SignalBoom(energy int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(energy, &buf); err != nil {
		return fmt.Errorf("failed to serialize energy: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x64), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalBoom: %s", err)
	}
	return nil
}
func (p *stubBomb) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Bomb",
		Methods:     map[uint32]object.MetaMethod{},
		Signals: map[uint32]object.MetaSignal{uint32(0x64): {
			Name:      "boom",
			Signature: "(i)",
			Uid:       uint32(0x64),
		}},
	}
}

// SpacecraftImplementor interface of the service implementation
type SpacecraftImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper SpacecraftSignalHelper) error
	OnTerminate()
	Shoot() (BombProxy, error)
	Ammo(ammo BombProxy) error
}

// SpacecraftSignalHelper provided to Spacecraft a companion object
type SpacecraftSignalHelper interface{}

// stubSpacecraft implements server.ServerObject.
type stubSpacecraft struct {
	obj     generic.Object
	impl    SpacecraftImplementor
	session bus.Session
}

// SpacecraftObject returns an object using SpacecraftImplementor
func SpacecraftObject(impl SpacecraftImplementor) server.ServerObject {
	var stb stubSpacecraft
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x64), stb.Shoot)
	stb.obj.Wrap(uint32(0x65), stb.Ammo)
	return &stb
}
func (p *stubSpacecraft) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubSpacecraft) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubSpacecraft) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubSpacecraft) Shoot(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.Shoot()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := func() error {
		meta, err := ret.MetaObject(ret.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s (%d)", err,
				ret.ObjectID())
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			ret.ServiceID(),
			ret.ObjectID(),
		}
		return object.WriteObjectReference(ref, &out)
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubSpacecraft) Ammo(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	ammo, err := func() (BombProxy, error) {
		ref, err := object.ReadObjectReference(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get proxy: %s", err)
		}
		return MakeBomb(p.session, proxy), nil
	}()
	if err != nil {
		return nil, fmt.Errorf("cannot read ammo: %s", err)
	}
	callErr := p.impl.Ammo(ammo)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubSpacecraft) metaObject() object.MetaObject {
	meta := object.MetaObject{
		Description: "Spacecraft",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "shoot",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "ammo",
				ParametersSignature: "(o)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x65),
			},
		},
		Signals: map[uint32]object.MetaSignal{},
	}
	return object.FullMetaObject(meta)
}
