package osLib

import (
	environment "evie/env"
	"evie/values"
)

func LoadProcessStruct(env *environment.Environment) {
	st := values.StructValue{}
	st.Name = "Process"
	st.Properties = make([]string, 0)

	methods := make(map[string]values.RuntimeValue, 0)

	st.Methods = methods
}
