package object

import (
	"github.com/lugu/qiloop/bus"
)

func MakeObject(proxy bus.Proxy) ObjectObject {
	return &proxyObject{proxy}
}
