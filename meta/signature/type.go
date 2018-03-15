package signature

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"strings"
)

// TypeSet is a container which contains exactly one instance of each
// Type currently generated. It is used to generate the
// type declaration only once.
type TypeSet struct {
	Signatures map[string]string // maps type name with signatures
	Types      []Type
}

// Declare writes all the registered Type into the jen.File.
func (s *TypeSet) Declare(f *jen.File) {
	for _, v := range s.Types {
		v.typeDeclaration(f)
	}
}

// NewTypeSet construct a new TypeSet.
func NewTypeSet() *TypeSet {
	sig := make(map[string]string)
	typ := make([]Type, 0)
	return &TypeSet{sig, typ}
}

// Type represents a type of a signature or a type
// embedded inside a signature. Type represents types for
// primitive types (int, long, float, string), vectors of a type,
// associative maps and structures.
type Type interface {
	Signature() string
	TypeName() *Statement
	typeDeclaration(*jen.File)
	RegisterTo(s *TypeSet)
	Marshal(id string, writer string) *Statement // returns an error
	Unmarshal(reader string) *Statement          // returns (type, err)
}

// Print render the type into a string. It is only used for testing.
func Print(v Type) string {
	buf := bytes.NewBufferString("")
	v.TypeName().Render(buf)
	return buf.String()
}

// NewLongType is a contructor for the representation of a uint64.
func NewLongType() LongType {
	return LongType{}
}

// NewULongType is a contructor for the representation of a uint64.
func NewULongType() ULongType {
	return ULongType{}
}

// NewFloatType is a contructor for the representation of a float32.
func NewFloatType() FloatType {
	return FloatType{}
}

// NewDoubleType is a contructor for the representation of a float32.
func NewDoubleType() DoubleType {
	return DoubleType{}
}

// NewIntType is a contructor for the representation of an int32.
func NewIntType() IntType {
	return IntType{}
}

// NewUIntType is a contructor for the representation of an uint32.
func NewUIntType() UIntType {
	return UIntType{}
}

// NewStringType is a contructor for the representation of a string.
func NewStringType() StringType {
	return StringType{}
}

// NewVoidType is a contructor for the representation of the
// absence of a return type. Only used in the context of a returned
// type.
func NewVoidType() VoidType {
	return VoidType{}
}

// NewValueType is a contructor for the representation of a Value.
func NewValueType() ValueType {
	return ValueType{}
}

// NewBoolType is a contructor for the representation of a bool.
func NewBoolType() BoolType {
	return BoolType{}
}

// NewListType is a contructor for the representation of a slice.
func NewListType(value Type) *ListType {
	return &ListType{value}
}

// NewMapType is a contructor for the representation of a map.
func NewMapType(key, value Type) *MapType {
	return &MapType{key, value}
}

// NewMemberType is a contructor for the representation of a field in
// a struct.
func NewMemberType(name string, value Type) MemberType {
	return MemberType{name, value}
}

// NewStrucType is a contructor for the representation of a struct.
func NewStrucType(name string, members []MemberType) *StructType {
	return &StructType{name, members}
}

// NewTupleType is a contructor for the representation of a series of
// types. Used to describe a method parameters list.
func NewTupleType(values []Type) *TupleType {
	return &TupleType{values}
}

// NewMetaObjectType is a contructor for the representation of an
// object.
func NewMetaObjectType() MetaObjectType {
	return MetaObjectType{}
}

// NewObjectType is a contructor for the representation of a Value.
func NewObjectType() ObjectType {
	return ObjectType{}
}

// NewUnknownType is a contructor for an unkown type.
func NewUnknownType() UnknownType {
	return UnknownType{}
}

// IntType represents an integer.
type IntType struct {
}

// Signature returns "i".
func (i IntType) Signature() string {
	return "i"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (i IntType) TypeName() *Statement {
	return jen.Int32()
}

// RegisterTo adds the type to the TypeSet.
func (i IntType) RegisterTo(s *TypeSet) {
	return
}

func (i IntType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (i IntType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteInt32").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (i IntType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadInt32").Call(jen.Id(reader))
}

// UIntType represents an integer.
type UIntType struct {
}

// Signature returns "I".
func (i UIntType) Signature() string {
	return "I"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (i UIntType) TypeName() *Statement {
	return jen.Uint32()
}

// RegisterTo adds the type to the TypeSet.
func (i UIntType) RegisterTo(s *TypeSet) {
	return
}

func (i UIntType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (i UIntType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteUint32").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (i UIntType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadUint32").Call(jen.Id(reader))
}

// LongType represents a long.
type LongType struct {
}

// Signature returns "l".
func (i LongType) Signature() string {
	return "l"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (i LongType) TypeName() *Statement {
	return jen.Int64()
}

// RegisterTo adds the type to the TypeSet.
func (i LongType) RegisterTo(s *TypeSet) {
	return
}

func (i LongType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (i LongType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteInt64").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (i LongType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadInt64").Call(jen.Id(reader))
}

// ULongType represents a long.
type ULongType struct {
}

// Signature returns "L".
func (i ULongType) Signature() string {
	return "L"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (i ULongType) TypeName() *Statement {
	return jen.Uint64()
}

// RegisterTo adds the type to the TypeSet.
func (i ULongType) RegisterTo(s *TypeSet) {
	return
}

func (i ULongType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (i ULongType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteUint64").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (i ULongType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadUint64").Call(jen.Id(reader))
}

// FloatType represents a float.
type FloatType struct {
}

// Signature returns "f".
func (f FloatType) Signature() string {
	return "f"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (f FloatType) TypeName() *Statement {
	return jen.Float32()
}

// RegisterTo adds the type to the TypeSet.
func (f FloatType) RegisterTo(s *TypeSet) {
	return
}

func (f FloatType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (f FloatType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteFloat32").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (f FloatType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadFloat32").Call(jen.Id(reader))
}

// DoubleType represents a float.
type DoubleType struct {
}

// Signature returns "f".
func (d DoubleType) Signature() string {
	return "f"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (d DoubleType) TypeName() *Statement {
	return jen.Float64()
}

// RegisterTo adds the type to the TypeSet.
func (d DoubleType) RegisterTo(s *TypeSet) {
	return
}

func (d DoubleType) typeDeclaration(dile *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (d DoubleType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteFloat64").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (d DoubleType) Unmarshal(reader string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "ReadFloat64").Call(jen.Id(reader))
}

// BoolType represents a bool.
type BoolType struct {
}

// Signature returns "b".
func (b BoolType) Signature() string {
	return "b"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (b BoolType) TypeName() *Statement {
	return jen.Bool()
}

// RegisterTo adds the type to the TypeSet.
func (b BoolType) RegisterTo(s *TypeSet) {
	return
}

func (b BoolType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (b BoolType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteBool").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (b BoolType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadBool").Call(jen.Id(reader))
}

// ValueType represents a Value.
type ValueType struct {
}

// Signature returns "m".
func (b ValueType) Signature() string {
	return "m"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (b ValueType) TypeName() *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/value", "Value")
}

// RegisterTo adds the type to the TypeSet.
func (b ValueType) RegisterTo(s *TypeSet) {
	return
}

func (b ValueType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (b ValueType) Marshal(id string, writer string) *Statement {
	return jen.Id(id).Dot("Write").Call(jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (b ValueType) Unmarshal(reader string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/value", "NewValue").Call(jen.Id(reader))
}

// VoidType represents the return type of a method.
type VoidType struct {
}

// Signature returns "v".
func (v VoidType) Signature() string {
	return "v"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (v VoidType) TypeName() *Statement {
	return jen.Empty()
}

// RegisterTo adds the type to the TypeSet.
func (v VoidType) RegisterTo(s *TypeSet) {
	return
}

func (v VoidType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (v VoidType) Marshal(id string, writer string) *Statement {
	return jen.Nil()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (v VoidType) Unmarshal(reader string) *Statement {
	return jen.Empty()
}

// StringType represents a string.
type StringType struct {
}

// Signature returns "s".
func (s StringType) Signature() string {
	return "s"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s StringType) TypeName() *Statement {
	return jen.String()
}

// RegisterTo adds the type to the TypeSet.
func (s StringType) RegisterTo(set *TypeSet) {
	return
}

func (s StringType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s StringType) Marshal(id string, writer string) *Statement {
	return jen.Id("basic.WriteString").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s StringType) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadString").Call(jen.Id(reader))
}

// ListType represents a slice.
type ListType struct {
	value Type
}

// Signature returns "[<signature>]" where <signature> is the
// signature of the type of the list.
func (l *ListType) Signature() string {
	return fmt.Sprintf("[%s]", l.value.Signature())
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (l *ListType) TypeName() *Statement {
	return jen.Index().Add(l.value.TypeName())
}

// RegisterTo adds the type to the TypeSet.
func (l *ListType) RegisterTo(s *TypeSet) {
	l.value.RegisterTo(s)
	return
}

func (l *ListType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (l *ListType) Marshal(listID string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
		jen.Err().Op(":=").Qual("github.com/lugu/qiloop/type/basic", "WriteUint32").Call(jen.Id("uint32").Call(
			jen.Id("len").Call(jen.Id(listID))),
			jen.Id(writer)),
		jen.Id(`if (err != nil) {
            return fmt.Errorf("failed to write slice size: %s", err)
        }`),
		jen.For(
			jen.Id("_, v := range "+listID),
		).Block(
			jen.Err().Op("=").Add(l.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write slice value: %s", err)
            }`),
		),
		jen.Return(jen.Nil()),
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (l *ListType) Unmarshal(reader string) *Statement {
	return jen.Func().Params().Params(
		jen.Id("b").Index().Add(l.value.TypeName()),
		jen.Err().Error(),
	).Block(
		jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
		jen.If(jen.Id("err != nil")).Block(
			jen.Return(jen.Id("b"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"failed to read slice size: %s", err`)))),
		jen.Id("b").Op("=").Id("make").Call(l.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("b[i], err =").Add(l.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return b, fmt.Errorf("failed to read slice value: %s", err)
            }`),
		),
		jen.Return(jen.Id("b"), jen.Nil()),
	).Call()
}

type MetaObjectType struct {
}

func (m MetaObjectType) Signature() string {
	return MetaObjectSignature
}

func (m MetaObjectType) TypeName() *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/object", "MetaObject")
}

func (m MetaObjectType) typeDeclaration(*jen.File) {
	return
}

func (m MetaObjectType) RegisterTo(s *TypeSet) {
	return
}

func (m MetaObjectType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/object", "WriteMetaObject").Call(jen.Id(id), jen.Id(writer))
}

func (m MetaObjectType) Unmarshal(reader string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/object", "ReadMetaObject").Call(jen.Id(reader))
}

type ObjectType struct {
}

func (o ObjectType) Signature() string {
	return "o"
}

func (o ObjectType) TypeName() *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/object", "ObjectReference")
}

func (o ObjectType) typeDeclaration(*jen.File) {
	return
}

func (o ObjectType) RegisterTo(s *TypeSet) {
	return
}

func (o ObjectType) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/object", "WriteObjectReference").Call(jen.Id(id), jen.Id(writer))
}

func (o ObjectType) Unmarshal(reader string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/type/object", "ReadObjectReference").Call(jen.Id(reader))
}

type UnknownType struct {
}

func (u UnknownType) Signature() string {
	return "X"
}

func (u UnknownType) TypeName() *Statement {
	return jen.Id("interface{}")
}

func (u UnknownType) typeDeclaration(*jen.File) {
	return
}

func (u UnknownType) RegisterTo(s *TypeSet) {
	return
}

func (u UnknownType) Marshal(id string, writer string) *Statement {
	return jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown type serialization not supported: %v"), jen.Id(id))
}

func (u UnknownType) Unmarshal(reader string) *Statement {
	return jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown type deserialization not supported")))
}

// MapType represents a map.
type MapType struct {
	key   Type
	value Type
}

// Signature returns "{<signature key><signature value>}" where
// <signature key> is the signature of the key and <signature value>
// the signature of the value.
func (m *MapType) Signature() string {
	return fmt.Sprintf("{%s%s}", m.key.Signature(), m.value.Signature())
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (m *MapType) TypeName() *Statement {
	return jen.Map(m.key.TypeName()).Add(m.value.TypeName())
}

// RegisterTo adds the type to the TypeSet.
func (m *MapType) RegisterTo(s *TypeSet) {
	m.key.RegisterTo(s)
	m.value.RegisterTo(s)
	return
}

func (m *MapType) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (m *MapType) Marshal(mapID string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
		jen.Err().Op(":=").Qual("github.com/lugu/qiloop/type/basic", "WriteUint32").Call(jen.Id("uint32").Call(
			jen.Id("len").Call(jen.Id(mapID))),
			jen.Id(writer)),
		jen.Id(`if (err != nil) {
            return fmt.Errorf("failed to write map size: %s", err)
        }`),
		jen.For(
			jen.Id("k, v := range "+mapID),
		).Block(
			jen.Err().Op("=").Add(m.key.Marshal("k", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write map key: %s", err)
            }`),
			jen.Err().Op("=").Add(m.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write map value: %s", err)
            }`),
		),
		jen.Return(jen.Nil()),
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (m *MapType) Unmarshal(reader string) *Statement {
	return jen.Func().Params().Params(
		jen.Id("m").Map(m.key.TypeName()).Add(m.value.TypeName()),
		jen.Err().Error(),
	).Block(
		jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
		jen.If(jen.Id("err != nil")).Block(
			jen.Return(jen.Id("m"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"failed to read map size: %s", err`)))),
		jen.Id("m").Op("=").Id("make").Call(m.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("k, err :=").Add(m.key.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("failed to read map key: %s", err)
            }`),
			jen.Id("v, err :=").Add(m.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("failed to read map value: %s", err)
            }`),
			jen.Id("m[k] = v"),
		),
		jen.Return(jen.Id("m"), jen.Nil()),
	).Call()
}

// MemberType a field in a struct.
type MemberType struct {
	Name  string
	Value Type
}

// Title is the public name of the field.
func (m MemberType) Title() string {
	return strings.Title(m.Name)
}

// TupleType a list of a parameter of a method.
type TupleType struct {
	values []Type
}

// Signature returns "(<signature 1><signature 2>...)" where
// <signature X> is the signature of the elements.
func (t *TupleType) Signature() string {
	sig := "("
	for _, v := range t.values {
		sig += v.Signature()
	}
	sig += ")"
	return sig
}

// Members returns the list of the types composing the TupleType.
func (t *TupleType) Members() []MemberType {
	members := make([]MemberType, len(t.values))
	for i, v := range t.values {
		members[i] = MemberType{fmt.Sprintf("P%d", i), v}
	}
	return members
}

// Params returns a statement representing the list of parameter of
// a method.
func (t *TupleType) Params() *Statement {
	arguments := make([]jen.Code, len(t.values))
	for i, v := range t.values {
		arguments[i] = jen.Id(fmt.Sprintf("P%d", i)).Add(v.TypeName())
	}
	return jen.Params(arguments...)
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (t *TupleType) TypeName() *Statement {
	params := make([]jen.Code, 0)
	for _, typ := range t.Members() {
		params = append(params, jen.Id(typ.Name).Add(typ.Value.TypeName()))
	}
	return jen.Struct(params...)
}

// RegisterTo adds the type to the TypeSet.
func (t *TupleType) RegisterTo(s *TypeSet) {
	for _, v := range t.values {
		v.RegisterTo(s)
	}
	return
}

func (t *TupleType) typeDeclaration(*jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (t *TupleType) Marshal(tupleID string, writer string) *Statement {
	statements := make([]jen.Code, 0)
	for _, typ := range t.Members() {
		s1 := jen.Err().Op("=").Add(typ.Value.Marshal(tupleID+"."+typ.Name, writer))
		s2 := jen.Id(`if (err != nil) {
			return fmt.Errorf("failed to write tuple member: %s", err)
		}`)
		statements = append(statements, s1)
		statements = append(statements, s2)
	}
	statements = append(statements, jen.Return(jen.Nil()))
	return jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown type serialization not implemented: %v"), jen.Id(tupleID))
	return jen.Func().Params().Params(jen.Error()).Block(
		statements...,
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (t *TupleType) Unmarshal(reader string) *Statement {
	statements := make([]jen.Code, 0)
	for _, typ := range t.Members() {
		s1 := jen.List(jen.Id("s."+typ.Name), jen.Err()).Op("=").Add(typ.Value.Unmarshal(reader))
		s2 := jen.Id(`if (err != nil) {
			return s, fmt.Errorf("failed to read tuple member: %s", err)
		}`)
		statements = append(statements, s1)
		statements = append(statements, s2)
	}
	statements = append(statements, jen.Return(jen.Id("s"), jen.Nil()))
	return jen.Func().Params().Params(
		jen.Id("s").Add(t.TypeName()),
		jen.Err().Error(),
	).Block(
		statements...,
	).Call()
}

// ConvertMetaObjects replace any element type which has the same
// signature as MetaObject with an element of the type
// object.MetaObject. This is required to generate proxy services
// which implements the object.Object interface and avoid a circular
// dependancy.
func (t *TupleType) ConvertMetaObjects() {
	for i, member := range t.values {
		if member.Signature() == MetaObjectSignature {
			t.values[i] = NewMetaObjectType()
		}
	}
}

// StructType represents a struct.
type StructType struct {
	name    string
	members []MemberType
}

// Signature returns the signature of the struct.
func (s *StructType) Signature() string {
	types := ""
	names := make([]string, 0, len(s.members))
	for _, v := range s.members {
		names = append(names, v.Name)
		if s, ok := v.Value.(*StructType); ok {
			types += "[" + s.Signature() + "]"
		} else {
			types += v.Value.Signature()
		}
	}
	return fmt.Sprintf("(%s)<%s,%s>", types,
		s.name, strings.Join(names, ","))
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s *StructType) TypeName() *Statement {
	return jen.Id(s.name)
}

// RegisterTo adds the type to the TypeSet.
func (s *StructType) RegisterTo(set *TypeSet) {
	for _, v := range s.members {
		v.Value.RegisterTo(set)
	}

	// register the name of the struct. if the name is used by a
	// struct of a different signature, change the name.
	for i := 0; i < 100; i++ {
		if sgn, ok := set.Signatures[s.name]; !ok {
			// name not yet used
			set.Signatures[s.name] = s.Signature()
			set.Types = append(set.Types, s)
			break
		} else if sgn == s.Signature() {
			// same name and same signatures
			break
		} else {
			s.name = fmt.Sprintf("%s_%d", s.name, i)
		}
	}
	return
}

func (s *StructType) typeDeclaration(file *jen.File) {
	fields := make([]jen.Code, len(s.members))
	for i, v := range s.members {
		fields[i] = jen.Id(v.Title()).Add(v.Value.TypeName())
	}
	file.Type().Id(s.name).Struct(fields...)

	readFields := make([]jen.Code, len(s.members)+1)
	writeFields := make([]jen.Code, len(s.members)+1)
	for i, v := range s.members {
		readFields[i] = jen.If(
			jen.Id("s."+v.Title()+", err =").Add(v.Value.Unmarshal("r")),
			jen.Id("err != nil")).Block(
			jen.Id(`return s, fmt.Errorf("failed to read ` + v.Title() + ` field: %s", err)`),
		)
		writeFields[i] = jen.If(
			jen.Id("err :=").Add(v.Value.Marshal("s."+v.Title(), "w")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id(`return fmt.Errorf("failed to write ` + v.Title() + ` field: %s", err)`),
		)
	}
	readFields[len(s.members)] = jen.Return(jen.Id("s"), jen.Nil())
	writeFields[len(s.members)] = jen.Return(jen.Nil())

	file.Func().Id("Read"+s.name).Params(
		jen.Id("r").Id("io.Reader"),
	).Params(
		jen.Id("s").Id(s.name), jen.Err().Error(),
	).Block(readFields...)
	file.Func().Id("Write"+s.name).Params(
		jen.Id("s").Id(s.name),
		jen.Id("w").Qual("io", "Writer"),
	).Params(jen.Err().Error()).Block(writeFields...)
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s *StructType) Marshal(structID string, writer string) *Statement {
	return jen.Id("Write"+s.name).Call(jen.Id(structID), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s *StructType) Unmarshal(reader string) *Statement {
	return jen.Id("Read" + s.name).Call(jen.Id(reader))
}
