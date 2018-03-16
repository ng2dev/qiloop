package idl

import (
	"github.com/lugu/qiloop/type/object"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServiceServer(t *testing.T) {
	var w strings.Builder
	if err := GenerateIDL(&w, "Server", object.MetaService0); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface Server
	fn authenticate(P0: Map<str,any>) -> Map<str,any>
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}

func TestObject(t *testing.T) {
	var w strings.Builder
	if err := GenerateIDL(&w, "Object", object.ObjectMetaObject); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface Object
	fn registerEvent(P0: uint32, P1: uint32, P2: uint64) -> uint64
	fn unregisterEvent(P0: uint32, P1: uint32, P2: uint64)
	fn metaObject(P0: uint32) -> MetaObject
	fn terminate(P0: uint32)
	fn property(P0: any) -> any
	fn setProperty(P0: any, P1: any)
	fn properties() -> Vec<str>
	fn registerEventWithSignature(P0: uint32, P1: uint32, P2: uint64, P3: str) -> uint64
end
struct MetaMethodParameter
	name: str
	description: str
end
struct MetaMethod
	uid: uint32
	returnSignature: str
	name: str
	parametersSignature: str
	description: str
	parameters: Vec<MetaMethodParameter>
	returnDescription: str
end
struct MetaSignal
	uid: uint32
	name: str
	signature: str
end
struct MetaProperty
	uid: uint32
	name: str
	signature: str
end
struct MetaObject
	methods: Map<uint32,MetaMethod>
	signals: Map<uint32,MetaSignal>
	properties: Map<uint32,MetaProperty>
	description: str
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}

func TestServiceDirectory(t *testing.T) {
	var w strings.Builder
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	if err := GenerateIDL(&w, "ServiceDirectory", metaObj); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface ServiceDirectory
	fn registerEvent(P0: uint32, P1: uint32, P2: uint64) -> uint64
	fn unregisterEvent(P0: uint32, P1: uint32, P2: uint64)
	fn metaObject(P0: uint32) -> MetaObject
	fn terminate(P0: uint32)
	fn property(P0: any) -> any
	fn setProperty(P0: any, P1: any)
	fn properties() -> Vec<str>
	fn registerEventWithSignature(P0: uint32, P1: uint32, P2: uint64, P3: str) -> uint64
	fn isStatsEnabled() -> bool
	fn enableStats(P0: bool)
	fn stats() -> Map<uint32,MethodStatistics>
	fn clearStats()
	fn isTraceEnabled() -> bool
	fn enableTrace(P0: bool)
	fn service(P0: str) -> ServiceInfo
	fn services() -> Vec<ServiceInfo>
	fn registerService(P0: ServiceInfo) -> uint32
	fn unregisterService(P0: uint32)
	fn serviceReady(P0: uint32)
	fn updateServiceInfo(P0: ServiceInfo)
	fn machineId() -> str
	fn _socketOfService(P0: uint32) -> obj
	sig traceObject(P0: EventTrace)
	sig serviceAdded(P0: uint32, P1: str)
	sig serviceRemoved(P0: uint32, P1: str)
end
struct MetaMethodParameter
	name: str
	description: str
end
struct MetaMethod
	uid: uint32
	returnSignature: str
	name: str
	parametersSignature: str
	description: str
	parameters: Vec<MetaMethodParameter>
	returnDescription: str
end
struct MetaSignal
	uid: uint32
	name: str
	signature: str
end
struct MetaProperty
	uid: uint32
	name: str
	signature: str
end
struct MetaObject
	methods: Map<uint32,MetaMethod>
	signals: Map<uint32,MetaSignal>
	properties: Map<uint32,MetaProperty>
	description: str
end
struct MinMaxSum
	minValue: float32
	maxValue: float32
	cumulatedValue: float32
end
struct MethodStatistics
	count: uint32
	wall: MinMaxSum
	user: MinMaxSum
	system: MinMaxSum
end
struct ServiceInfo
	name: str
	serviceId: uint32
	machineId: str
	processId: uint32
	endpoints: Vec<str>
	sessionId: str
end
struct timeval
	tv_sec: int64
	tv_usec: int64
end
struct EventTrace
	id: uint32
	kind: int32
	slotId: uint32
	arguments: any
	timestamp: timeval
	userUsTime: int64
	systemUsTime: int64
	callerContext: uint32
	calleeContext: uint32
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}
