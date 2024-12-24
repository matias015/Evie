package postgresLib

import (
	"database/sql"
	environment "evie/env"
	"evie/values"
	"log"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

type PostgreSQLConnection struct {
	Value *sql.DB
}

func (p PostgreSQLConnection) GetType() string {
	return "PostgreSQLConnection"
}

func (p PostgreSQLConnection) GetStr() string {
	return ""
}

func (p PostgreSQLConnection) GetNumber() float64 {
	return 0
}

func (p PostgreSQLConnection) GetBool() bool {
	return false
}

func (db PostgreSQLConnection) GetProp(v *values.RuntimeValue, name string) values.RuntimeValue {
	props := map[string]values.RuntimeValue{}

	props["ping"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		val := (*v).(*PostgreSQLConnection)

		err := val.Value.Ping()

		if err != nil {
			return values.BooleanValue{Value: false}
		}

		return values.BooleanValue{Value: true}
	}}

	props["close"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		val := (*v).(*PostgreSQLConnection)
		val.Value.Close()

		return values.BooleanValue{Value: true}
	}}

	props["query"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {

		// Arguments check
		if len(args) == 0 {
			return values.ErrorValue{Value: "No query specified, expected 1 argument"}
		}

		queryStr, ok := args[0].(values.StringValue)

		if !ok {
			return values.ErrorValue{Value: "Expected query to be a string"}
		}

		query := strings.TrimSpace(queryStr.Value)

		if !strings.HasPrefix(query, "SELECT") {
			result, err := db.Value.Exec(query)

			if err != nil {
				return values.ErrorValue{Value: err.Error()}
			}

			if strings.HasPrefix(query, "UPDATE") {
				rowsAffected, err := result.RowsAffected()
				if err != nil {
					return values.BooleanValue{Value: true}
				}
				return values.NumberValue{Value: float64(rowsAffected)}
				// LAST INSERTED ID IS NOT SUPPORTED FOR PQ
				// } else if strings.HasPrefix(query, "INSERT") {
				// 	lastInsertId, err := result.LastInsertId()
				// 	if err != nil {
				// 		return values.BooleanValue{Value: true}
				// 	}
				// 	return values.NumberValue{Value: float64(lastInsertId)}
			} else {
				return values.BooleanValue{Value: true}
			}

		}

		// Execute query
		rows, err := db.Value.Query(query)

		if err != nil {
			return values.ErrorValue{Value: err.Error()}
		}

		defer rows.Close()

		// Get column names
		columnNames, err := rows.Columns()

		if err != nil {
			return values.ErrorValue{Value: err.Error()}
		}

		// Make a list to hold the row values
		// to load dynamicly the columns in runtime structures
		rowsList := make([]interface{}, len(columnNames))
		rowsListPtrs := make([]interface{}, len(columnNames))

		for i := range rowsList {
			rowsListPtrs[i] = &rowsList[i]
		}

		// Returned value will be an array of dictionaries
		arr := values.ArrayValue{}
		arr.Value = make([]values.RuntimeValue, 0)

		// Load vals in dict
		for rows.Next() {

			dict := values.DictionaryValue{}
			dict.Value = make(map[string]values.RuntimeValue)

			err = rows.Scan(rowsListPtrs...)
			// now, rowsList has the returned literals values of the columns

			if err != nil {
				return values.ErrorValue{Value: err.Error()}
			}

			for i, colName := range columnNames {

				colVal := rowsList[i]

				switch reflect.TypeOf(colVal) {
				case reflect.TypeOf(""):
					dict.Value[colName] = values.StringValue{Value: colVal.(string)}
				case reflect.TypeOf(0):
					dict.Value[colName] = values.NumberValue{Value: colVal.(float64)}
				case reflect.TypeOf(true):
					dict.Value[colName] = values.BooleanValue{Value: colVal.(bool)}
				}

			}

			// Add row
			arr.Value = append(arr.Value, dict)
		}

		return &arr

	}}

	return props[name]
}

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
