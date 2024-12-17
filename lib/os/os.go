package osLib

import (
	environment "evie/env"
	"evie/values"
)

func Load(env *environment.Environment) {
	structValue := values.StructValue{}
	structValue.Properties = []string{"os"}

	structValue.Methods = make(map[string]values.RuntimeValue)

	fn := values.NativeFunctionValue{}

	fn.Value = func(args []values.RuntimeValue) values.RuntimeValue {
		return values.StringValue{Value: "Hola modulo os"}
	}

	structValue.Methods["hello"] = fn

	obj := values.ObjectValue{}
	obj.Struct = structValue
	obj.Value = make(map[string]values.RuntimeValue)

	env.Variables["os"] = obj

}
