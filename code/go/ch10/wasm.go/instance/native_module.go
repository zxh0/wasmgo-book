package instance

var _ Module = (*nativeModule)(nil)

type nativeModule struct {
	exported map[string]interface{}
}

func NewNativeInstance() nativeModule {
	return nativeModule{
		exported: map[string]interface{}{},
	}
}

func (nm nativeModule) RegisterFunc(nameAndSig string, f GoFunc) {
	name, sig := parseNameAndSig(nameAndSig)
	nm.exported[name] = nativeFunction{t: sig, f: f}
}

func (nm nativeModule) Register(name string, x interface{}) {
	nm.exported[name] = x
}

func (nm nativeModule) GetMember(name string) interface{} {
	return nm.exported[name]
}

func (nm nativeModule) InvokeFunc(name string, args ...WasmVal) ([]WasmVal, error) {
	return nm.exported[name].(Function).Call(args...) // TODO
}

func (nm nativeModule) GetGlobalVal(name string) (WasmVal, error) {
	return nm.exported[name].(Global).Get(), nil // TODO
}

func (nm nativeModule) SetGlobalVal(name string, val WasmVal) error {
	nm.exported[name].(Global).Set(val) // TODO
	return nil
}
