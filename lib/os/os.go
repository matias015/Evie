package osLib

import (
	environment "evie/env"
)

func Load(env *environment.Environment) {

	// ns := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	// ns.Value["exec"] = values.NativeFunctionValue{Value: Exec}

	// env.DeclareVar("os", ns)

}

// func Exec(args []values.RuntimeValue) values.RuntimeValue {

// 	if len(args) == 0 {
// 		return values.ErrorValue{Value: "Missing arguments"}
// 	}

// 	strArgs := make([]string, len(args))

// 	for i, arg := range args {
// 		strArgs[i] = arg.(values.StringValue).Value
// 	}

// 	cmd := exec.Command(strArgs[0], strArgs[1:]...)

// 	output, err := cmd.Output()

// 	if err != nil {
// 		return values.ErrorValue{Value: err.Error()}
// 	}

// 	return values.StringValue{Value: string(output)}
// }
