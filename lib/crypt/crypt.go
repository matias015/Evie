package cryptLib

import (
	"crypto/md5"
	"encoding/hex"
	environment "evie/env"
	"evie/values"

	"golang.org/x/crypto/bcrypt"
)

func Load(env *environment.Environment) {

	nm := values.NamespaceValue{}

	fns := make(map[string]values.RuntimeValue)

	fns["md5"] = values.NativeFunctionValue{Value: ToMd5}
	fns["bcrypt"] = values.NativeFunctionValue{Value: BCrypt}

	nm.Value = fns

	env.DeclareVar("crypt", nm)

}

func ToMd5(args []values.RuntimeValue) values.RuntimeValue {
	hash := md5.Sum([]byte(args[0].(values.StringValue).Value))
	text := hex.EncodeToString(hash[:])
	return values.StringValue{Value: text}
}

func BCrypt(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing arguments for bcrypt"}
	}

	cost := bcrypt.DefaultCost

	if len(args) > 1 {
		cost = int(args[1].(values.NumberValue).Value)
	}

	arg := args[0].(values.StringValue).Value

	argBytes := []byte(arg)

	output, err := bcrypt.GenerateFromPassword(argBytes, cost)

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	return values.StringValue{Value: string(output)}

}
