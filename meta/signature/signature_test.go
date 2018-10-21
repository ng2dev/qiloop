package signature

import (
	"bytes"
	"strings"
	"testing"
)

func testUtil(t *testing.T, input string, expected Type) {
	result, err := Parse(input)
	if err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("wrong return")
	} else if strings.ToLower(result.Signature()) != strings.ToLower(expected.Signature()) {
		buf := bytes.NewBufferString("")
		result.TypeName().Render(buf)
		t.Errorf("invalid type: %s (%s)", buf.String(), result.Signature())
	}
}

func testSignature(t *testing.T, signature string) {
	result, err := Parse(signature)
	if err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("wrong return")
	} else if result.Signature() != signature {
		t.Error("invalid signature: " + result.Signature())
	}
}

func TestParseBasics(t *testing.T) {
	testUtil(t, "i", NewIntType())
	testUtil(t, "I", NewUIntType())
	testUtil(t, "s", NewStringType())
	testUtil(t, "L", NewULongType())
	testUtil(t, "l", NewLongType())
	testUtil(t, "b", NewBoolType())
	testUtil(t, "f", NewFloatType())
	testUtil(t, "d", NewDoubleType())
	testUtil(t, "m", NewValueType())
	testUtil(t, "X", NewUnknownType())
}

func TestParseMultipleString(t *testing.T) {
	testUtil(t, "ss", NewStringType())
}

func TestParseEmpty(t *testing.T) {
	t.SkipNow()
	testUtil(t, "", nil)
}

func TestParseMap(t *testing.T) {
	testUtil(t, "{ss}", NewMapType(NewStringType(), NewStringType()))
	testUtil(t, "{sI}", NewMapType(NewStringType(), NewUIntType()))
	testUtil(t, "{is}", NewMapType(NewIntType(), NewStringType()))
	testUtil(t, "{iI}", NewMapType(NewIntType(), NewUIntType()))
	testUtil(t, "{Li}", NewMapType(NewULongType(), NewIntType()))
	testUtil(t, "{sl}", NewMapType(NewStringType(), NewLongType()))
}

func TestParseList(t *testing.T) {
	testUtil(t, "[s]", NewListType(NewStringType()))
	testUtil(t, "[i]", NewListType(NewIntType()))
	testUtil(t, "[b]", NewListType(NewBoolType()))
	testUtil(t, "[{bi}]", NewListType(NewMapType(NewBoolType(), NewIntType())))
	testUtil(t, "{b[i]}", NewMapType(NewBoolType(), NewListType(NewIntType())))
}

func TestParseTuple(t *testing.T) {
	testUtil(t, "(s)", NewTupleType([]Type{NewStringType()}))
	testUtil(t, "(i)", NewTupleType([]Type{NewIntType()}))
	testUtil(t, "(ii)", NewTupleType([]Type{NewIntType(), NewIntType()}))
	testUtil(t, "(fbd)", NewTupleType([]Type{NewFloatType(), NewBoolType(), NewDoubleType()}))
}

func TestParseDefinition(t *testing.T) {
	testUtil(t, "()<test>", NewStructType("test", []MemberType{}))
	testUtil(t, "(s)<test,a>", NewStructType("test", []MemberType{NewMemberType("a", NewStringType())}))
	testUtil(t, "(ss)<test,a,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewStringType()),
		NewMemberType("a", NewStringType()),
	}))
	testUtil(t, "(sss)<test,a,a,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewStringType()),
		NewMemberType("a", NewStringType()),
		NewMemberType("a", NewStringType()),
	}))
}

func TestParseEmbeddedDefinition(t *testing.T) {
	testUtil(t, "([s])<test,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewListType(NewStringType()))}))
	testUtil(t, "({si})<test,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewMapType(NewStringType(), NewIntType()))}))
}

func TestParseMapMap(t *testing.T) {
	testSignature(t, "{{ii}i}")
	testSignature(t, "{i{ii}}")
	testSignature(t, "{{ss}{ii}}")
	testSignature(t, "{{{si}s}{ii}}")
}

func TestParseDefinitionSignature(t *testing.T) {
	testSignature(t, "(s)<test,a>")
	testSignature(t, "(si)<test,a,b>")
	testSignature(t, "(iii)<test,a,b,c>")
	testSignature(t, "(s{ii})<test,a,b>")
	testSignature(t, "({ss})<test,a>")
}

func TestParseEmbeddedDefinitionSignature(t *testing.T) {
	testSignature(t, "([(s)<test2,b>])<test,a>")
	testSignature(t, "(s[(s)<test2,b>])<test,a,b>")
	testSignature(t, "([(s)<test2,b>]s)<test,a,b>")
}

func TestParseMetaSignal(t *testing.T) {
	testSignature(t, "(Iss)<MetaSignal,uid,name,signature>")
}
func TestParseMetaProperty(t *testing.T) {
	testSignature(t, "(Iss)<MetaProperty,uid,name,signature>")
}
func TestParseMetaMethodParameter(t *testing.T) {
	testSignature(t, "(ss)<MetaMethodParameter,name,description>")
}
func TestParseMetaMethod(t *testing.T) {
	testSignature(t, "(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>")
}
func TestParseMetaSignalMap2(t *testing.T) {
	testSignature(t, "{(Iss)<MetaSignal,uid,name,signature>I}")
}
func TestParseMetaSignalMap(t *testing.T) {
	testSignature(t, "{I(Iss)<MetaSignal,uid,name,signature>}")
}
func TestParseMetaPropertyMap(t *testing.T) {
	testSignature(t, "{I(Iss)<MetaProperty,uid,name,signature>}")
}
func TestParseMetaMethodMap(t *testing.T) {
	testSignature(t, "{I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}")
}
func TestParseMetaObject(t *testing.T) {
	testSignature(t, "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>")
}
func TestParseServiceInfo(t *testing.T) {
	testSignature(t, "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]")
}
func TestParseTextProcessingContext(t *testing.T) {
	testSignature(t, "{sc}")
	testSignature(t, "({sc}fs)<AgentGrd,concepts,confidence,userId>")
	testSignature(t, "(s(({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}fs)<AgentGrd,concepts,confidence,userId>i)<TextProcessingContext,author,recever,language>)")
	testSignature(t, "((iis({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}f({ii})<Duration,timeInfos>)<TimeGrd,concepts,confidence,reference>(im)<Expression,type,exp>)<Knowledge,fromSource,fromLanguage,fromText,author,time,exp>(({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}fs)<AgentGrd,concepts,confidence,userId>i)<TextProcessingContext,author,recever,language>m)")
}
func TestParseKnowledge(t *testing.T) {
	testSignature(t, "((iis({sc}fs)<AgentGrd,concepts,confidence,userId>({sc}f({ii})<Duration,timeInfos>)<TimeGrd,concepts,confidence,reference>(im)<Expression,type,exp>)<Knowledge,fromSource,fromLanguage,fromText,author,time,exp>)")
}
func TestParsePersonState(t *testing.T) {
	testSignature(t, "((ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>((ff)<BodyLanguageEase,level,confidence>)<BodyLanguageState,ease>(ff)<Smile,value,confidence>((ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>)<Expressions,calm,anger,joy,sorrow,laughter,excitement,surprise>)<PersonState,valence,attention,bodyLanguageState,smile,expressions>")
}
func TestParseRobotFullState(t *testing.T) {
	testSignature(t, "((ff)<ValueConfidence<float>,value,confidence>(ff)<ValueConfidence<float>,value,confidence>)<RobotFullState,pleasure,excitement>")
}
