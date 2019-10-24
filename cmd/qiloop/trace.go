package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

var (
	infos = make([]services.ServiceInfo, 0)
	metas = make([]object.MetaObject, 0)
)

func getObject(sess bus.Session, info services.ServiceInfo,
	objectID uint32) (bus.ObjectProxy, error) {

	proxy, err := sess.Proxy(info.Name, objectID)
	if err != nil {
		return nil, fmt.Errorf("connect service (%s): %s", info.Name, err)
	}
	return bus.MakeObject(proxy), nil
}

func printEvent(e bus.EventTrace, info *services.ServiceInfo,
	meta *object.MetaObject) {

	var typ = "unknown"
	switch e.Kind {
	case int32(net.Call):
		typ = "call "
	case int32(net.Reply):
		typ = "reply"
	case int32(net.Error):
		typ = "error"
	case int32(net.Post):
		typ = "post "
	case int32(net.Event):
		typ = "event"
	case int32(net.Capability):
		typ = "capability"
	case int32(net.Cancel):
		typ = "cancel"
	case int32(net.Cancelled):
		typ = "cancelled"
	}
	action, err := meta.ActionName(e.SlotId)
	if err != nil {
		action = fmt.Sprintf("unknown (%d)", e.SlotId)
	}
	var size = -1
	var sig = "unknown"
	var data = []byte{}
	var buf bytes.Buffer
	err = e.Arguments.Write(&buf)
	if err == nil {
		sig, err = basic.ReadString(&buf)
		if err == nil {
			data = buf.Bytes()
			size = len(data)
		}
	}

	fmt.Printf("[%s %4d bytes] %s.%s: %s\n", typ, size, info.Name,
		action, sig)
}

func trace(serverURL, serviceName string, objectID uint32) {

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("%s: %s", serverURL, err)
	}

	proxies := services.Services(sess)

	directory, err := proxies.ServiceDirectory()
	if err != nil {
		log.Fatalf("directory: %s", err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		log.Fatalf("services: %s", err)
	}

	stop := make(chan struct{})

	serviceCount := 0

	serviceID, err := strconv.Atoi(serviceName)
	if err != nil {
		serviceID = -1
	} else {
		serviceName = ""
	}

	for _, info := range serviceList {

		if serviceID != -1 && uint32(serviceID) != info.ServiceId {
			continue
		} else if serviceName != "" && serviceName != info.Name {
			continue
		}

		obj, err := getObject(sess, info, objectID)
		if err != nil {
			log.Printf("cannot trace %s: %s", info.Name, err)
			continue
		}

		serviceCount++
		go func(info services.ServiceInfo, obj bus.ObjectProxy) {

			err = obj.EnableTrace(true)
			if err != nil {
				log.Fatalf("Failed to start traces: %s.", err)
			}
			defer obj.EnableTrace(false)

			cancel, trace, err := obj.SubscribeTraceObject()
			if err != nil {
				log.Fatalf("Failed to stop stats: %s.", err)
			}
			defer cancel()

			meta, err := obj.MetaObject(objectID)
			if err != nil {
				log.Fatalf("%s: MetaObject: %s.", info.Name, err)
			}

			for {
				select {
				case event, ok := <-trace:
					if !ok {
						return
					}
					printEvent(event, &info, &meta)
				case <-stop:
					return
				}
			}
		}(info, obj)
	}

	if serviceCount == 0 {
		log.Fatalf("Service not found")
		return
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	<-signalChannel
	close(stop)
}
