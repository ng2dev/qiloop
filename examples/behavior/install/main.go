// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"flag"
	"fmt"

	"bitbucket.org/swoldt/pkg/xerrors/iferr"
	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/bus/services"
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

	ok, err := bhm.InstallBehavior("examples/behavior/install/tai_chi_chuan.crg")
	iferr.Exit(err)
	fmt.Println("bahavior installed -> ", ok)

	// subscribe to the signal "serviceAdded" of the service directory.
	behaviors, err := bhm.GetInstalledBehaviors()
	iferr.Exit(err)

	fmt.Println("installed bahaviors -> ", len(behaviors))

	// for _, behavior := range behaviors {
	// 	fmt.Println(behavior)
	// }
}
