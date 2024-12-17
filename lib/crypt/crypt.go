package cryptLib

import (
	"crypto/md5"
	"encoding/hex"
	environment "evie/env"
	"evie/values"
)

func Load(env *environment.Environment) {

	structValue := values.StructValue{}
	structValue.Properties = []string{"crypt"}

	structValue.Methods = make(map[string]values.RuntimeValue)

	structValue.Methods["md5"] = values.NativeFunctionValue{Value: ToMd5}
	obj := values.ObjectValue{}
	obj.Struct = structValue
	obj.Value = make(map[string]values.RuntimeValue)

	env.Variables["crypt"] = obj

}

func ToMd5(args []values.RuntimeValue) values.RuntimeValue {
	hash := md5.Sum([]byte(args[0].GetStr()))
	text := hex.EncodeToString(hash[:])
	return values.StringValue{Value: text}
}
