package cryptLib

import (
	"crypto/md5"
	"encoding/hex"
	environment "evie/env"
	"evie/values"

	"golang.org/x/crypto/bcrypt"
)

func Load(env *environment.Environment) {

	structValue := values.StructValue{}
	structValue.Properties = []string{"crypt"}

	structValue.Methods = make(map[string]values.RuntimeValue)

	structValue.Methods["md5"] = values.NativeFunctionValue{Value: ToMd5}
	structValue.Methods["bcrypt"] = values.NativeFunctionValue{Value: BCrypt}

	obj := values.ObjectValue{}
	obj.Struct = structValue
	obj.Value = make(map[string]values.RuntimeValue)

	env.DeclareVar("crypt", obj)

}

func ToMd5(args []values.RuntimeValue) values.RuntimeValue {
	hash := md5.Sum([]byte(args[0].GetStr()))
	text := hex.EncodeToString(hash[:])
	return values.StringValue{Value: text}
}

func BCrypt(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing arguments for bcrypt"}
	}

	cost := bcrypt.DefaultCost

	if len(args) > 1 {
		cost = int(args[1].GetNumber())
	}

	arg := args[0].(values.StringValue).Value

	argBytes := []byte(arg)

	output, err := bcrypt.GenerateFromPassword(argBytes, cost)

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	return values.StringValue{Value: string(output)}

}
