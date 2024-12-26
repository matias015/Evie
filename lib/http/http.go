package htppLib

import (
	environment "evie/env"
)

// func GetRequestStructValue() values.StructValue {
// 	req := values.StructValue{}
// 	req.Properties = []string{"method"}
// 	req.Methods = make(map[string]values.RuntimeValue)
// 	return req
// }

// LOAD
func Load(env *environment.Environment) {

	// // BASE NAMESPACE
	// namespace := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	// namespace.Value["route"] = values.NativeFunctionValue{Value: AddRoute}
	// namespace.Value["listen"] = values.NativeFunctionValue{Value: ListenAndServe}

	// env.DeclareVar("http", namespace)

	// // REQUEST VALUE STRUCT
	// env.DeclareVar("Request", GetRequestStructValue())
}

// func AddRoute(args []values.RuntimeValue) values.RuntimeValue {

// 	pattern := args[0].GetStr()
// 	function := args[1]
// 	fn := function.(values.FunctionValue)
// 	env := fn.Environment.(*environment.Environment)
// 	env.DeclareVar("Request", values.StringValue{Value: "Request he"})

// 	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

// 		reqObject := values.ObjectValue{}
// 		reqObject.Value = make(map[string]values.RuntimeValue)
// 		reqObject.Value["method"] = values.StringValue{Value: r.Method}
// 		reqObject.Struct = GetRequestStructValue()

// 		args := make([]interface{}, 0)

// 		args = append(args, reqObject)

// 		ret := fn.Evaluator.ExecuteCallback(fn, env, args)

// 		switch returned := ret.(type) {
// 		case values.ErrorValue:
// 			w.Write([]byte("500 Internal Server Error"))
// 		case *values.StringValue:
// 			w.Write([]byte(returned.GetStr()))
// 		case values.StringValue:
// 			w.Write([]byte(returned.GetStr()))
// 		}
// 	})
// 	return values.BooleanValue{Value: true}
// }

// func ListenAndServe(args []values.RuntimeValue) values.RuntimeValue {

// 	if len(args) == 0 {
// 		return values.ErrorValue{Value: "No port specified"}
// 	}

// 	port := args[0].GetStr()
// 	http.ListenAndServe(":"+port, nil)
// 	return values.BooleanValue{Value: true}
// }
