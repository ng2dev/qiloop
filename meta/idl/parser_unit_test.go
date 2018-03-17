package idl

import (
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"reflect"
	"testing"
)

func TestParseReturns(t *testing.T) {
	input := "-> int32"
	expected := signature.NewIntType()
	root, _ := returns()(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Fatalf("error parsing returns:\n%s", input)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("cannot parse returns: %v", err)
	}
	if ret, ok := root.(signature.Type); !ok {
		t.Fatalf("return type error: %+v", root)
	} else if ret.Signature() != expected.Signature() {
		t.Fatalf("cannot generate signature: %+v", ret)
	}
}

func helpParseMethod(t *testing.T, label, input string, expected object.MetaMethod) {
	root, _ := method()(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Fatalf("%s: cannot parse input:\n%s", label, input)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("%s: parsing error: %v", label, err)
	}
	if method, ok := root.(*object.MetaMethod); !ok {
		t.Fatalf("%s; type error %+v: %+v", label, reflect.TypeOf(root), root)
	} else if !reflect.DeepEqual(*method, expected) {
		t.Fatalf("%s: expected %#v, got %#v", label, expected, method)
	}
}

func TestParseMethod0(t *testing.T) {
	input := `fn methodName()`
	expected := object.MetaMethod{
		Name:            "methodName",
		ReturnSignature: "v",
	}
	helpParseMethod(t, "TestParseMethod0", input, expected)
}

func TestParseMethod1(t *testing.T) {
	input := `fn methodName() -> int32`
	expected := object.MetaMethod{
		ReturnSignature: "i",
		Name:            "methodName",
	}
	helpParseMethod(t, "TestParseMethod1", input, expected)
}

func TestParseMethod2(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64)`
	expected := object.MetaMethod{
		Name:                "methodName",
		ParametersSignature: "id",
		ReturnSignature:     "v",
		Parameters: []object.MetaMethodParameter{
			object.MetaMethodParameter{
				Name: "param1",
			},
			object.MetaMethodParameter{
				Name: "param2",
			},
		},
	}
	helpParseMethod(t, "TestParseMethod2", input, expected)
}
func TestParseMethod3(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64) -> bool`
	expected := object.MetaMethod{
		ParametersSignature: "id",
		Name:                "methodName",
		ReturnSignature:     "b",
		Parameters: []object.MetaMethodParameter{
			object.MetaMethodParameter{
				Name: "param1",
			},
			object.MetaMethodParameter{
				Name: "param2",
			},
		},
	}
	helpParseMethod(t, "TestParseMethod3", input, expected)
}
