package instance

import (
	"strings"

	"wasm.go/binary"
)

// name(params)->(results)
func parseNameAndSig(nameAndSig string) (string, binary.FuncType) {
	idxOfLPar := strings.IndexByte(nameAndSig, '(')
	name := nameAndSig[:idxOfLPar]
	sig := nameAndSig[idxOfLPar:]
	return name, parseSig(sig)
}

func parseSig(sig string) binary.FuncType {
	paramsAndResults := strings.SplitN(sig, "->", 2)
	return binary.FuncType{
		ParamTypes:  parseValTypes(paramsAndResults[0]),
		ResultTypes: parseValTypes(paramsAndResults[1]),
	}
}

func parseValTypes(list string) []binary.ValType {
	list = strings.TrimSpace(list)
	list = list[1 : len(list)-1] // remove ()

	var valTypes []binary.ValType
	for _, t := range strings.Split(list, ",") {
		switch strings.TrimSpace(t) {
		case "i32":
			valTypes = append(valTypes, binary.ValTypeI32)
		case "i64":
			valTypes = append(valTypes, binary.ValTypeI64)
		case "f32":
			valTypes = append(valTypes, binary.ValTypeF32)
		case "f64":
			valTypes = append(valTypes, binary.ValTypeF64)
		}
	}
	return valTypes
}
