package cryptLib

import (
	"crypto/md5"
	"encoding/hex"
	environment "evie/env"
	"evie/values"

	"golang.org/x/crypto/bcrypt"
)

func Load(env *environment.Environment) {

	nm := values.RuntimeValue{Type: values.NamespaceType}

	fns := make(map[string]values.RuntimeValue)

	fns["md5"] = values.RuntimeValue{Type: values.NativeFunctionType, Value: ToMd5}
	fns["bcrypt"] = values.RuntimeValue{Type: values.NativeFunctionType, Value: BCrypt}

	nm.Value = fns

	env.DeclareVar("crypt", nm)

}

func ToMd5(args []values.RuntimeValue) values.RuntimeValue {
	hash := md5.Sum([]byte(args[0].Value.(string)))
	text := hex.EncodeToString(hash[:])
	return values.RuntimeValue{Type: values.StringType, Value: text}
}

func BCrypt(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.RuntimeValue{Type: values.ErrorType, Value: "Missing arguments for bcrypt"}
	}

	cost := bcrypt.DefaultCost

	if len(args) > 1 {
		cost = int(args[1].Value.(float64))
	}

	arg := args[0].Value.(string)

	argBytes := []byte(arg)

	output, err := bcrypt.GenerateFromPassword(argBytes, cost)

	if err != nil {
		return values.RuntimeValue{Type: values.ErrorType, Value: err.Error()}
	}

	return values.RuntimeValue{Type: values.StringType, Value: string(output)}

}
