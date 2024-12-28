package postgresLib

import (
	"database/sql"
	environment "evie/env"
	"evie/values"
	"log"

	_ "github.com/lib/pq"
)

func Load(env *environment.Environment) {
	ns := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	ns.Value["connect"] = values.NativeFunctionValue{Value: Connect}

	env.DeclareVar("postgres", ns)
}

func Connect(args []values.RuntimeValue) values.RuntimeValue {

	argslen := len(args)

	usern := "postgres"
	pw := "admin"
	dbname := "postgres"
	host := "127.0.0.1"

	if argslen > 0 {
		usernStr, ok := args[0].(values.StringValue)
		if !ok {
			return values.ErrorValue{Value: "Expected username to be a string"}
		}
		usern = usernStr.Value
	}
	if argslen > 1 {
		pwString, ok := args[1].(values.StringValue)
		if !ok {
			return values.ErrorValue{Value: "Expected password to be a string"}
		}
		pw = pwString.Value
	}
	if argslen > 2 {
		dbnameStr, ok := args[2].(values.StringValue)
		if !ok {
			return values.ErrorValue{Value: "Expected dbname to be a string"}
		}
		dbname = dbnameStr.Value
	}
	if argslen > 3 {
		hostStr, ok := args[3].(values.StringValue)
		if !ok {
			return values.ErrorValue{Value: "Expected host to be a string"}
		}
		host = hostStr.Value
	}

	connStr := "user=" + usern + " password=" + pw + " dbname=" + dbname + " host=" + host + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &PostgreSQLConnection{Value: db}

}
