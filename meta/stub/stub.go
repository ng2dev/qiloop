package stub

import (
	"fmt"
	"io"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
)

func stubName(name string) string {
	return "stub" + name
}

// GeneratePackage generate the stub for the package declaration.
func GeneratePackage(w io.Writer, packagePath string,
	pkg *idl.PackageDeclaration) error {

	if pkg.Name == "" {
		if strings.Contains(packagePath, "/") {
			pkg.Name = packagePath[strings.LastIndex(packagePath, "/")+1:]
		} else {
			return fmt.Errorf("empty package name")
		}
	}
	file := jen.NewFilePathName(packagePath, pkg.Name)

	set := signature.NewTypeSet()
	for _, typ := range pkg.Types {
		typ.RegisterTo(set)

		idl.InterfaceTypeForStub = true
		itf, ok := typ.(*idl.InterfaceType)
		if ok {
			err := generateInterface(file, set, itf)
			if err != nil {
				return err
			}
		}
		idl.InterfaceTypeForStub = false
	}

	proxy.GenerateNewServices(file)

	set.Declare(file)
	return file.Render(w)
}

func generateInterface(f *jen.File, set *signature.TypeSet, itf *idl.InterfaceType) error {
	if err := generateObjectInterface(f, set, itf); err != nil {
		return err
	}
	if err := generateStub(f, itf); err != nil {
		return err
	}
	return nil
}

func implName(name string) string {
	return name + "Implementor"
}

func generateStub(f *jen.File, itf *idl.InterfaceType) error {
	if err := generateStubType(f, itf); err != nil {
		return err
	}
	if err := generateStubConstructor(f, itf); err != nil {
		return err
	}
	if err := generateStubObject(f, itf); err != nil {
		return err
	}
	if err := generateStubMethods(f, itf); err != nil {
		return err
	}
	if err := generateStubMetaObject(f, itf); err != nil {
		return err
	}
	return nil
}

func generateMethodDef(itf *idl.InterfaceType, set *signature.TypeSet,
	method idl.Method, methodName string) (jen.Code, error) {

	// RegisterEvent and UnregisterEvent are implemented in basicObject.
	if methodName == "RegisterEvent" && method.ID == 0x0 {
		return jen.Id(`RegisterEvent(msg *net.Message, c Channel) error`), nil
	} else if methodName == "UnregisterEvent" && method.ID == 0x1 {
		return jen.Id(`UnregisterEvent(msg *net.Message, c Channel) error`), nil
	}

	tuple := method.Tuple()
	ret := method.Return
	// resolve MetaObject references for Generic objet.
	if ret.Signature() == signature.MetaObjectSignature {
		ret = signature.NewMetaObjectType()
	}

	tuple.RegisterTo(set)
	ret.RegisterTo(set)

	if ret.Signature() == "v" {
		return jen.Id(methodName).Add(tuple.Params()).Params(jen.Error()), nil
	}
	return jen.Id(methodName).Add(tuple.Params()).Params(ret.TypeName(),
		jen.Error()), nil
}

func generateSignalDef(itf *idl.InterfaceType, set *signature.TypeSet,
	tuple *signature.TupleType, signalName string) (jen.Code, error) {

	tuple.RegisterTo(set)
	return jen.Id(signalName).Add(tuple.Params()).Error(), nil
}

func methodBodyBlock(itf *idl.InterfaceType, method idl.Method,
	methodName string) (*jen.Statement, error) {

	writing := make([]jen.Code, 0)
	params := make([]jen.Code, 0)
	code := jen.Id("buf").Op(":=").Qual(
		"bytes", "NewBuffer",
	).Call(jen.Id("msg.Payload"))
	if len(method.Params) != 0 {
		writing = append(writing, code)
	}

	for _, param := range method.Params {
		params = append(params, jen.Id(param.Name))
		code = jen.List(jen.Id(param.Name), jen.Err()).Op(":=").Add(
			param.Type.Unmarshal("buf"),
		)
		writing = append(writing, code)
		code = jen.Id(`if err != nil {
		    return c.SendError(msg, fmt.Errorf("cannot read ` + param.Name + `: %s", err))
	        }`)
		writing = append(writing, code)
	}
	ret := method.Return

	// resolve MetaObject references for Generic objet.
	if ret.Signature() == signature.MetaObjectSignature {
		ret = signature.NewMetaObjectType()
	}

	// if has not return value
	if ret.Signature() == "v" {
		code = jen.Id("callErr := p.impl").Dot(methodName).Call(params...)
	} else {
		code = jen.Id("ret, callErr := p.impl").Dot(methodName).Call(params...)
	}
	writing = append(writing, code)
	code = jen.Id(`
	// do not respond to post messages.
	if msg.Header.Type == net.Post {
	    return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer`)
	writing = append(writing, code)
	if ret.Signature() != "v" {
		code = jen.Id("errOut").Op(":=").Add(ret.Marshal("ret", "&out"))
		writing = append(writing, code)

		code = jen.Id(`if errOut != nil {
		    return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	        }`)
		writing = append(writing, code)
	}
	code = jen.Id(`return c.SendReply(msg, out.Bytes())`)
	writing = append(writing, code)

	return jen.Block(
		writing...,
	), nil

}

func generateMethodMarshal(file *jen.File, itf *idl.InterfaceType,
	method idl.Method, methodName string) error {

	// RegisterEvent and UnregisterEvent are implemented in basicObject.
	if methodName == "RegisterEvent" && method.ID == 0x0 {
		file.Id(`
 func (p *stubObject) RegisterEvent(msg *net.Message, c Channel) error {
       return p.impl.RegisterEvent(msg, c)
 }`)
		return nil
	} else if methodName == "UnregisterEvent" && method.ID == 0x1 {
		file.Id(`
 func (p *stubObject) UnregisterEvent(msg *net.Message, c Channel) error {
       return p.impl.UnregisterEvent(msg, c)
 }`)
		return nil
	}

	body, err := methodBodyBlock(itf, method, methodName)
	if err != nil {
		return fmt.Errorf("create method body: %s", err)
	}

	file.Func().Params(jen.Id("p").Op("*").Id(stubName(itf.Name))).Id(methodName).Params(
		jen.Id("msg *net.Message"),
		jen.Id("c").Qual("github.com/lugu/qiloop/bus", "Channel"),
	).Params(
		jen.Id("error"),
	).Add(body)
	return nil
}

func propertyBodyBlock(itf *idl.InterfaceType, property idl.Property,
	signalName string) (*jen.Statement, error) {

	writing := make([]jen.Code, 0)
	code := jen.Var().Id("buf").Qual("bytes", "Buffer")
	writing = append(writing, code)

	for _, param := range property.Params {
		code = jen.If(jen.Err().Op(":=").Add(
			param.Type.Marshal(param.Name, "&buf"),
		).Op(";").Err().Op("!=").Nil()).Block(
			jen.Id(`return fmt.Errorf("serialize ` +
				param.Name + `: %s", err)`),
		)
		writing = append(writing, code)
	}
	code = jen.Id("err := p.signal.UpdateProperty").Call(
		jen.Lit(int(property.ID)),
		jen.Lit(property.Type().Signature()),
		jen.Id("buf.Bytes()"),
	)
	writing = append(writing, code)
	code = jen.Id(`
	if err != nil {
	    return fmt.Errorf("update ` +
		signalName + `: %s", err)
	}
	return nil`)
	writing = append(writing, code)
	return jen.Block(
		writing...,
	), nil
}

func signalBodyBlock(itf *idl.InterfaceType, signal idl.Signal,
	signalName string) (*jen.Statement, error) {

	writing := make([]jen.Code, 0)
	code := jen.Var().Id("buf").Qual("bytes", "Buffer")
	writing = append(writing, code)

	for _, param := range signal.Params {
		code = jen.If(jen.Err().Op(":=").Add(
			param.Type.Marshal(param.Name, "&buf"),
		).Op(";").Err().Op("!=").Nil()).Block(
			jen.Id(`return fmt.Errorf("serialize ` +
				param.Name + `: %s", err)`),
		)
		writing = append(writing, code)
	}
	// if has not return value
	code = jen.Id("err := p.signal.UpdateSignal").Call(
		jen.Lit(int(signal.ID)),
		jen.Id("buf.Bytes()"),
	)
	writing = append(writing, code)
	code = jen.Id(`
	if err != nil {
	    return fmt.Errorf("update ` +
		signalName + `: %s", err)
	}
	return nil`)
	writing = append(writing, code)
	return jen.Block(
		writing...,
	), nil
}

func generateSignalHelper(file *jen.File, itf *idl.InterfaceType,
	signal idl.Signal, signalName string) error {
	tuple := signal.Tuple()

	body, err := signalBodyBlock(itf, signal, signalName)
	if err != nil {
		return fmt.Errorf("create signal helper body: %s", err)
	}
	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id(signalName).Add(tuple.Params()).Error().Add(body)
	return nil
}

func generatePropertyHelper(file *jen.File, itf *idl.InterfaceType,
	property idl.Property, propertyName string) error {
	tuple := property.Tuple()

	body, err := propertyBodyBlock(itf, property, propertyName)
	if err != nil {
		return fmt.Errorf("create property helper body: %s", err)
	}
	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id(propertyName).Add(tuple.Params()).Error().Add(body)
	return nil
}

func generateStubMethods(file *jen.File, itf *idl.InterfaceType) error {

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		err := generateMethodMarshal(file, itf, method, methodName)
		if err != nil {
			return fmt.Errorf("create method marshall %s of %s: %s",
				method.Name, itf.Name, err)
		}
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = "Signal" + signalName
		err := generateSignalHelper(file, itf, signal, signalName)
		if err != nil {
			return fmt.Errorf("create signal marshall %s of %s: %s",
				signal.Name, itf.Name, err)
		}
		return nil
	}
	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]
		propertyName = "Update" + propertyName
		err := generatePropertyHelper(file, itf, property,
			propertyName)
		if err != nil {
			return fmt.Errorf("create property marshall %s of %s: %s",
				property.Name, itf.Name, err)
		}
		return nil
	}

	meta := itf.MetaObject()
	if err := meta.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("generate interface object %s: %s",
			itf.Name, err)
	}
	return nil
}
func generateStubMetaObject(file *jen.File, itf *idl.InterfaceType) error {
	metaMethods := func(d jen.Dict) {
		for _, method := range itf.Methods {
			d[jen.Lit(int(method.ID))] = jen.Values(jen.Dict{
				jen.Id("Uid"): jen.Lit(int(method.ID)),
				jen.Id("ReturnSignature"): jen.Lit(
					method.Return.Signature(),
				),
				jen.Id("Name"): jen.Lit(method.Name),
				jen.Id("ParametersSignature"): jen.Lit(
					method.Tuple().Signature(),
				),
			})
		}
	}
	metaSignals := func(d jen.Dict) {
		for _, signal := range itf.Signals {
			d[jen.Lit(int(signal.ID))] = jen.Values(jen.Dict{
				jen.Id("Uid"):  jen.Lit(int(signal.ID)),
				jen.Id("Name"): jen.Lit(signal.Name),
				jen.Id("Signature"): jen.Lit(
					signal.Type().Signature(),
				),
			})
		}
	}
	metaProperties := func(d jen.Dict) {
		for _, property := range itf.Properties {
			d[jen.Lit(int(property.ID))] = jen.Values(jen.Dict{
				jen.Id("Uid"):  jen.Lit(int(property.ID)),
				jen.Id("Name"): jen.Lit(property.Name),
				jen.Id("Signature"): jen.Lit(
					property.Type().Signature(),
				),
			})
		}
	}
	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id("metaObject").Params().Params(
		jen.Qual("github.com/lugu/qiloop/type/object", "MetaObject"),
	).Block(
		jen.Return().Qual(
			"github.com/lugu/qiloop/type/object", "MetaObject",
		).Values(
			jen.Dict{
				jen.Id("Description"): jen.Lit(itf.Name),
				jen.Id("Methods"): jen.Map(jen.Uint32()).Qual(
					"github.com/lugu/qiloop/type/object",
					"MetaMethod",
				).Values(
					jen.DictFunc(metaMethods),
				),
				jen.Id("Signals"): jen.Map(jen.Uint32()).Qual(
					"github.com/lugu/qiloop/type/object",
					"MetaSignal",
				).Values(
					jen.DictFunc(metaSignals),
				),
				jen.Id("Properties"): jen.Map(jen.Uint32()).Qual(
					"github.com/lugu/qiloop/type/object",
					"MetaProperty",
				).Values(
					jen.DictFunc(metaProperties),
				),
			},
		),
	)
	return nil
}

func generateStubObject(file *jen.File, itf *idl.InterfaceType) error {
	writing := make([]jen.Code, 0)
	if itf.Name == "Object" {
		code := jen.Id(`p.signal.Activate(activation)`)
		writing = append(writing, code)
		code = jen.Id(`p.session = activation.Session`)
		writing = append(writing, code)
		code = jen.Id(`err := p.impl.Activate(activation, p)
		if err != nil {
			p.signal.OnTerminate()
			return err
		}`)
		writing = append(writing, code)
		code = jen.Id(`err = p.obj.Activate(activation)
		if err != nil {
			p.impl.OnTerminate()
			p.signal.OnTerminate()
			return err
		}
		return nil`)
		writing = append(writing, code)
	} else {
		code := jen.Id(`p.session = activation.Session`)
		writing = append(writing, code)
		code = jen.Id(`p.service = activation.Service`)
		writing = append(writing, code)
		code = jen.Id(`p.serviceID = activation.ServiceID`)
		writing = append(writing, code)
		code = jen.Id(`return p.impl.Activate(activation, p)`)
		writing = append(writing, code)
	}

	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id("Activate").Params(
		jen.Id("activation").Qual(
			"github.com/lugu/qiloop/bus",
			"Activation",
		),
	).Params(jen.Error()).Block(writing...)

	writing = make([]jen.Code, 0)
	if itf.Name == "Object" {
		code := jen.Id(`p.obj.OnTerminate()`)
		writing = append(writing, code)
		code = jen.Id(`p.impl.OnTerminate()`)
		writing = append(writing, code)
		code = jen.Id(`p.signal.OnTerminate()`)
		writing = append(writing, code)
	} else {
		code := jen.Id(`p.impl.OnTerminate()`)
		writing = append(writing, code)
	}

	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id("OnTerminate").Params().Block(writing...)

	err := generateReceiveMethod(file, itf)
	if err != nil {
		return err
	}

	return generateStubPropertyCallback(file, itf)
}

func generateReceiveMethod(file *jen.File, itf *idl.InterfaceType) error {
	writing := make([]jen.Code, 0)

	prelude := jen.Comment("action dispatch")
	if itf.Name == "Object" {
		prelude = jen.Id(`from = p.impl.Tracer(msg, from)`)
	}

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		code := jen.Case(jen.Lit(int(method.ID)))
		writing = append(writing, code)

		code = jen.Return().Id("p").Dot(methodName).Call(
			jen.Id("msg"),
			jen.Id("from"),
		)
		writing = append(writing, code)
		return nil
	}

	signal := func(m object.MetaSignal, signalName string) error {
		return nil
	}

	property := func(p object.MetaProperty, propertyName string) error {
		return nil
	}

	meta := itf.MetaObject()
	err := meta.ForEachMethodAndSignal(method, signal, property)
	if err != nil {
		return fmt.Errorf("generate receive: %s: %s",
			itf.Name, err)
	}
	code := jen.Default()
	writing = append(writing, code)
	if itf.Name == "Object" {
		code = jen.Id(`return p.obj.Receive(msg, from)`)
		writing = append(writing, code)
	} else {
		code = jen.Return().Id("from").Op(".").Id("SendError").Call(
			jen.Id("msg"),
			jen.Qual(
				"github.com/lugu/qiloop/bus",
				"ErrActionNotFound",
			),
		)
		writing = append(writing, code)
	}

	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id("Receive").Params(
		jen.Id("msg").Op("*").Qual("github.com/lugu/qiloop/bus/net", "Message"),
		jen.Id("from").Qual("github.com/lugu/qiloop/bus", "Channel"),
	).Params(jen.Error()).Block(
		prelude,
		jen.Switch(jen.Id("msg.Header.Action")).Block(
			writing...,
		),
	)
	return nil
}

func generateStubPropertyCallback(file *jen.File, itf *idl.InterfaceType) error {

	writing := make([]jen.Code, 0)

	method := func(m object.MetaMethod, methodName string) error {
		return nil
	}

	signal := func(m object.MetaSignal, signalName string) error {
		return nil
	}

	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]
		code := jen.Case(jen.Lit(p.Name))
		writing = append(writing, code)
		code = jen.Id("buf").Op(":=").Qual(
			"bytes", "NewBuffer",
		).Call(jen.Id("data"))
		writing = append(writing, code)
		code = jen.List(jen.Id("prop"), jen.Err()).Op(":=").Add(
			property.Type().Unmarshal("buf"),
		)
		writing = append(writing, code)
		code = jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("cannot read "+propertyName+": %s"),
				jen.Err(),
			),
		)
		writing = append(writing, code)
		code = jen.Id(`return p.impl.On` + propertyName + `Change(prop)`)
		writing = append(writing, code)
		return nil
	}

	meta := itf.MetaObject()
	err := meta.ForEachMethodAndSignal(method, signal, property)
	if err != nil {
		return fmt.Errorf("failed at onPropertyChange for %s: %s",
			itf.Name, err)
	}
	code := jen.Id(`default:`).Return().Qual("fmt", "Errorf").Call(
		jen.Lit("unknown property %s"),
		jen.Id("name"),
	)
	writing = append(writing, code)

	file.Func().Params(
		jen.Id("p").Op("*").Id(stubName(itf.Name)),
	).Id("onPropertyChange").Params(
		jen.Id("name").String(),
		jen.Id("data []byte"),
	).Params(
		jen.Error(),
	).Block(
		jen.Switch().Id("name").Block(writing...),
	)
	return nil
}

func generateStubConstructor(file *jen.File, itf *idl.InterfaceType) error {
	writing := make([]jen.Code, 0)
	code := jen.Var().Id("stb").Id(stubName(itf.Name))
	writing = append(writing, code)
	code = jen.Id("stb.impl = impl")
	writing = append(writing, code)
	if itf.Name == "Object" {
		code = jen.Id("stb.signal").Op("=").Qual(
			"github.com/lugu/qiloop/bus",
			"newSignalHandler",
		).Call()
		writing = append(writing, code)
		code = jen.Return().Op("&").Id("stb")
		writing = append(writing, code)
	} else {
		code = jen.Id("obj").Op(":=").Qual(
			"github.com/lugu/qiloop/bus",
			"NewBasicObject",
		).Call(
			jen.Id("&stb"),
			jen.Id("stb.metaObject()"),
			jen.Id("stb.onPropertyChange"),
		)
		writing = append(writing, code)
		code = jen.Id(`stb.signal = obj`)
		writing = append(writing, code)
		code = jen.Id(`return obj`)
		writing = append(writing, code)
	}

	file.Commentf("%s returns an object using %s", itf.Name+"Object",
		implName(itf.Name))
	file.Func().Id(itf.Name+"Object").Params(
		jen.Id("impl").Id(implName(itf.Name)),
	).Qual(
		"github.com/lugu/qiloop/bus", "Actor",
	).Block(writing...)
	return generateObjectConstructor(file, itf)
}

func generateStubType(file *jen.File, itf *idl.InterfaceType) error {
	file.Commentf("%s implements server.Actor.", stubName(itf.Name))

	writing := make([]jen.Code, 0)
	code := jen.Id("impl").Id(implName(itf.Name))
	writing = append(writing, code)
	code = jen.Id("session").Qual("github.com/lugu/qiloop/bus", "Session")
	writing = append(writing, code)
	code = jen.Id("service").Qual("github.com/lugu/qiloop/bus", "Service")
	writing = append(writing, code)
	code = jen.Id("serviceID").Uint32()
	writing = append(writing, code)
	if itf.Name == "Object" {
		code = jen.Id("signal").Op("*").Qual(
			"github.com/lugu/qiloop/bus", "signalHandler")
		writing = append(writing, code)
		code = jen.Id("obj").Qual("github.com/lugu/qiloop/bus", "Actor")
		writing = append(writing, code)
	} else {
		code = jen.Id("signal").Qual("github.com/lugu/qiloop/bus",
			"SignalHandler")
		writing = append(writing, code)
	}

	file.Type().Id(stubName(itf.Name)).Struct(writing...)
	return nil
}

func generateObjectInterface(file *jen.File, set *signature.TypeSet,
	itf *idl.InterfaceType) error {

	// Proxy and stub shall generate the name method name: reuse
	// the MetaObject method ForEachMethodAndSignal to get an
	// ordered list of the method with uniq name.
	definitions := make([]jen.Code, 0)
	signalDefinitions := make([]jen.Code, 0)
	activate := jen.Id("Activate").Params(
		jen.Id("activation").Qual("github.com/lugu/qiloop/bus",
			"Activation"),
		jen.Id("helper").Id(itf.Name+"SignalHelper"),
	).Params(
		jen.Error(),
	)
	comment := jen.Comment("Activate is called before any other method.")
	definitions = append(definitions, comment)
	comment = jen.Comment("It shall be used to initialize the interface.")
	definitions = append(definitions, comment)
	comment = jen.Comment("activation provides runtime informations.")
	definitions = append(definitions, comment)
	comment = jen.Comment("activation.Terminate() unregisters the object.")
	definitions = append(definitions, comment)
	comment = jen.Comment("activation.Session can access other services.")
	definitions = append(definitions, comment)
	comment = jen.Comment("helper enables signals and properties updates.")
	definitions = append(definitions, comment)
	comment = jen.Comment("Properties must be initialized using helper,")
	definitions = append(definitions, comment)
	comment = jen.Comment("during the Activate call.")
	definitions = append(definitions, comment)
	definitions = append(definitions, activate)
	terminate := jen.Id("OnTerminate()")
	definitions = append(definitions, terminate)

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		def, err := generateMethodDef(itf, set, method, methodName)
		if err != nil {
			return fmt.Errorf("render method definition %s of %s: %s",
				method.Name, itf.Name, err)
		}
		definitions = append(definitions, def)
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = "Signal" + signalName
		def, err := generateSignalDef(itf, set, signal.Tuple(), signalName)
		if err != nil {
			return fmt.Errorf("render %s of %s: %s",
				s.Name, itf.Name, err)
		}
		signalDefinitions = append(signalDefinitions, def)
		return nil
	}
	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]

		comment := "On" + propertyName +
			`Change is called when the property is updated.`
		definitions = append(definitions, jen.Comment(comment))
		comment = `Returns an error if the property value is not allowed`
		definitions = append(definitions, jen.Comment(comment))

		callback := jen.Id("On" + propertyName + "Change").Add(
			property.Tuple().Params(),
		).Error()
		definitions = append(definitions, callback)

		signalName := "Update" + propertyName
		def, err := generateSignalDef(itf, set, property.Tuple(), signalName)
		if err != nil {
			return fmt.Errorf("render %s of %s: %s",
				p.Name, itf.Name, err)
		}

		signalDefinitions = append(signalDefinitions, def)
		return nil
	}

	meta := itf.MetaObject()
	if err := meta.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("generate interface object %s: %s",
			itf.Name, err)
	}

	if itf.Name == "Object" {
		code := jen.Id(`Tracer(msg *net.Message, from Channel) Channel`)
		definitions = append(definitions, code)
	}

	file.Commentf("%s interface of the service implementation",
		implName(itf.Name))
	file.Type().Id(implName(itf.Name)).Interface(
		definitions...,
	)

	file.Commentf("%s provided to %s a companion object",
		itf.Name+"SignalHelper", itf.Name)
	file.Type().Id(itf.Name + "SignalHelper").Interface(
		signalDefinitions...,
	)
	return nil
}

func generateObjectConstructor(file *jen.File, itf *idl.InterfaceType) error {
	if itf.Name == "Object" || itf.Name == "ServiceZero" {
		return nil
	}
	file.Comment("New" + itf.Name + " registers a new object to a service")
	file.Comment("and returns a proxy to the newly created object")
	file.Func().Params(
		jen.Id("c").Id("Constructor"),
	).Id("New"+itf.Name).Params(
		jen.Id("service").Qual("github.com/lugu/qiloop/bus", "Service"),
		jen.Id("impl").Id(implName(itf.Name)),
	).Params(
		jen.Id(itf.Name+"Proxy"),
		jen.Error(),
	).Block(
		jen.Id(`obj := `+itf.Name+`Object(impl)
		        objectID, err := service.Add(obj)
		        if err != nil {
		                return nil, err
		        }
		        stb := &stub`+itf.Name+`{}
		        meta := object.FullMetaObject(stb.metaObject())`),
		jen.Id("client").Op(":=").Qual(
			"github.com/lugu/qiloop/bus", "DirectClient",
		).Call(jen.Id("obj")),
		jen.Id("proxy").Op(":=").Qual(
			"github.com/lugu/qiloop/bus", "NewProxy",
		).Call(
			jen.Id("client"),
			jen.Id("meta"),
			jen.Id("service.ServiceID()"),
			jen.Id("objectID"),
		),
		jen.Id(`return Make`+itf.Name+`(c.session, proxy), nil`),
	)
	return nil
}
