// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"flag"
	"fmt"

	"bitbucket.org/swoldt/pkg/xerrors/iferr"
	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/bus/services"
)

const (
	HEY_ANIMATION_1 = "animations/Stand/Gestures/Hey_1"
)

func main() {
	flag.Parse()
	// session represents a connection to the service directory.
	session, err := app.SessionFromFlag()
	iferr.Exit(err)
	defer session.Terminate()

	// Access the specialized proxy constructor.
	proxy := services.Services(session)

	// Obtain a proxy to the service
	bhm, err := proxy.ALBehaviorManager(nil)
	iferr.Exit(err)

	err = bhm.StopAllBehaviors()
	iferr.Exit(err)

	var unsubscribe func()
	var channel chan string
	// subscribe to the signal "behaviorStopped" of the behavior manager.
	unsubscribe, channel, err = bhm.SubscribeBehaviorStopped()
	iferr.Exit(err)

	err = bhm.RunBehavior(HEY_ANIMATION_1)
	iferr.Exit(err)

	// wait until behavior stop signal is received.
	event := <-channel
	fmt.Println("finished bahavior -> ", event)
	unsubscribe()

}
