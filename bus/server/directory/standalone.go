package directory

import (
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
)

func NewServer(addr string, auth server.Authenticator) (*server.Server, error) {

	if auth == nil {
		auth = server.Yes{}
	}

	impl := NewServiceDirectory()
	info := ServiceInfo{
		Name:      "ServiceDirectory",
		ServiceId: 1,
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: []string{addr},
		SessionId: "",
	}
	serviceID, err := impl.RegisterService(info)
	if err != nil {
		return nil, err
	}
	if serviceID != 1 {
		return nil, fmt.Errorf("service directory id: %d", serviceID)
	}
	err = impl.ServiceReady(1)
	if err != nil {
		return nil, err
	}

	service1 := ServiceDirectoryObject(impl)

	listener, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}
	s, err := server.StandAloneServer(listener, auth, impl.Namespace())
	if err != nil {
		return nil, err
	}
	_, err = s.NewService("DirectoryService", service1)
	if err != nil {
		s.Stop()
		return nil, err
	}
	return s, nil
}
