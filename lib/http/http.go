package htppLib

import (
	environment "evie/env"
	"evie/values"
	"net/http"

	"github.com/sanity-io/litter"
)

// SERVER DATA TYPE
type HttpServer struct {
	Value *http.ServeMux
}

func (s HttpServer) SetValue(value values.RuntimeValue) {}

func (s HttpServer) GetStr() string     { return "" }
func (s HttpServer) GetNumber() float64 { return 1 }
func (s HttpServer) GetBool() bool {
	return false
}
func (s HttpServer) GetType() string {
	return "FileValue"
}

func (s *HttpServer) GetProp(v *values.RuntimeValue, name string) values.RuntimeValue {
	props := map[string]values.RuntimeValue{}

	props["listen"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		val := (*v).(*HttpServer)
		port := args[0].GetStr()
		http.ListenAndServe(":"+port, val.Value)
		return values.BooleanValue{Value: true}
	}}

	props["route"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		val := (*v).(*HttpServer)
		pattern := args[0].GetStr()
		function := args[1]
		fn := function.(values.FunctionValue)
		env := fn.Environment.(*environment.Environment)
		env.Variables["Request"] = values.StringValue{Value: "Request he"}
		val.Value.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

			ret := fn.Evaluator.ExecuteCallback(fn, env)

			switch returned := ret.(type) {
			case values.ErrorValue:
				litter.Dump(val.Value)
			case *values.StringValue:
				w.Write([]byte(returned.GetStr()))
			}
		})
		return values.BooleanValue{Value: true}
	}}

	return props[name]
}

// LOAD
func Load(env *environment.Environment) {

	// BASE NAMESPACE
	namespace := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}
	namespace.Value["createServer"] = values.NativeFunctionValue{Value: CreateServer}

	namespace.Value["route"] = values.NativeFunctionValue{Value: AddRoute}
	namespace.Value["listen"] = values.NativeFunctionValue{Value: ListenAndServe}

	env.Variables["http"] = namespace

	// REQUEST VALUE STRUCT
	req := values.StructValue{}
	req.Properties = []string{"method"}
	req.Methods = make(map[string]values.RuntimeValue)

	env.Variables["Request"] = req
}

func AddRoute(args []values.RuntimeValue) values.RuntimeValue {

	pattern := args[0].GetStr()
	function := args[1]
	fn := function.(values.FunctionValue)
	env := fn.Environment.(*environment.Environment)
	env.Variables["Request"] = values.StringValue{Value: "Request he"}
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		ret := fn.Evaluator.ExecuteCallback(fn, env)

		switch returned := ret.(type) {
		case values.ErrorValue:
			w.Write([]byte("500 Internal Server Error"))
		case *values.StringValue:
			w.Write([]byte(returned.GetStr()))
		case values.StringValue:
			w.Write([]byte(returned.GetStr()))
		}
	})
	return values.BooleanValue{Value: true}
}

func ListenAndServe(args []values.RuntimeValue) values.RuntimeValue {

	port := args[0].GetStr()
	http.ListenAndServe(":"+port, nil)
	return values.BooleanValue{Value: true}
}

func CreateServer(args []values.RuntimeValue) values.RuntimeValue {
	server := http.NewServeMux()

	httpServer := HttpServer{Value: server}

	return &httpServer
}
