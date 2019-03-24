package tester_test

import (
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/tester"
	"github.com/lugu/qiloop/bus/util"
	"sync"
	"testing"
)

func TestAddRemoveObject(t *testing.T) {

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := server.PrivateNamespace()
	srv, err := server.StandAloneServer(listener, server.Yes{}, ns)
	if err != nil {
		t.Error(err)
	}

	obj := tester.NewSpacecraftObject()
	service, err := srv.NewService("Spacecraft", obj)
	if err != nil {
		t.Error(err)
	}

	session := srv.Session()
	proxies := tester.Services(session)

	spacecraft, err := proxies.Spacecraft()
	if err != nil {
		t.Fatal(err)
	}

	bomb, err := spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	// initial delay is 10 seconds
	delay, err := bomb.GetDelay()
	if err != nil {
		t.Error(err)
	} else if delay != 10 {
		t.Errorf("unexpected delay: %d", delay)
	}

	err = bomb.SetDelay(12)
	if err != nil {
		t.Error(err)
	}

	delay, err = bomb.GetDelay()
	if err != nil {
		t.Error(err)
	} else if delay != 12 {
		t.Errorf("unexpected delay: %d", delay)
	}

	err = bomb.SetDelay(-1)
	if err == nil {
		t.Error("error expected")
	}

	delay, err = bomb.GetDelay()
	if err != nil {
		t.Error(err)
	} else if delay != 12 {
		t.Errorf("unexpected delay: %d", delay)
	}

	err = spacecraft.Ammo(bomb)
	if err != nil {
		t.Error(err)
	}

	ammo := tester.NewBombObject()
	id, err := service.Add(ammo)
	if err != nil {
		t.Error(err)
	}

	proxy, err := session.Proxy("Spacecraft", id)
	if err != nil {
		t.Error(err)
	}

	ammoProxy := tester.MakeBomb(session, proxy)
	err = spacecraft.Ammo(ammoProxy)
	if err != nil {
		t.Error(err)
	}

	service.Terminate()
	srv.Terminate()

}

func TestOnTerminate(t *testing.T) {
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := server.PrivateNamespace()
	srv, err := server.StandAloneServer(listener, server.Yes{}, ns)
	if err != nil {
		t.Error(err)
	}

	obj := tester.NewSpacecraftObject()
	service, err := srv.NewService("Spacecraft", obj)
	if err != nil {
		t.Error(err)
	}
	defer service.Terminate()

	session := srv.Session()
	proxies := tester.Services(session)

	spacecraft, err := proxies.Spacecraft()
	if err != nil {
		t.Fatal(err)
	}

	bomb, err := spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	var wait sync.WaitGroup
	wait.Add(1)
	tester.Hook = func(event string) {
		if event == "Bomb.OnTerminate()" {
			wait.Done()
		}
	}

	err = bomb.Terminate(bomb.ObjectID())
	if err != nil {
		t.Fatal(err)
	}
	wait.Wait()
	tester.Hook = func(event string) {}
}
