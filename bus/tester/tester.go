package tester

import (
	"fmt"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/type/object"
)

var (
	// Hook let the test program record what is happening
	Hook = make(chan string)
)

// NewSpacecraftObject creates a new server side Spacecraft object.
func NewSpacecraftObject() bus.Actor {
	return SpacecraftObject(&spacecraftImpl{})
}

type spacecraftImpl struct {
	session   bus.Session
	terminate func()
	service   bus.Service
	ammo      BombProxy
}

func (f *spacecraftImpl) Activate(activation bus.Activation,
	helper SpacecraftSignalHelper) error {
	f.session = activation.Session
	f.terminate = activation.Terminate
	f.service = activation.Service
	ammo, err := CreateBomb(f.session, f.service)
	f.ammo = ammo
	return err
}

func (f *spacecraftImpl) OnTerminate() {
	select {
	case Hook <- "SpaceCraft.OnTerminate()":
	default:
	}
}

func (f *spacecraftImpl) Shoot() (BombProxy, error) {
	return f.ammo, nil
}

func (f *spacecraftImpl) Ammo(b BombProxy) error {
	f.ammo = b
	return nil
}

type bombImpl struct{}

func (f *bombImpl) Activate(activation bus.Activation,
	helper BombSignalHelper) error {

	err := helper.UpdateDelay(10)
	if err != nil {
		return err
	}
	return nil
}

func (f *bombImpl) OnTerminate() {
	select {
	case Hook <- "Bomb.OnTerminate()":
	default:
	}

}

func (f *bombImpl) OnDelayChange(duration int32) error {
	if duration < 0 {
		return fmt.Errorf("duration cannot be negative (%d)", duration)
	}
	return nil
}

// NewBombObject returns the server side implementation of a Bomb
// object.
func NewBombObject() bus.Actor {
	return BombObject(&bombImpl{})
}

// CreateBomb returns a new Bomb object.
//
// Not entirely satisfying: need to allow for client side object
// generation... Here comes the ObjectID question..
func CreateBomb(session bus.Session, service bus.Service) (BombProxy, error) {

	var stb stubBomb
	stb.impl = &bombImpl{}
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj

	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}

	meta := object.FullMetaObject(stb.metaObject())

	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)

	return MakeBomb(session, proxy), nil
}
