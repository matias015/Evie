package postgresLib

import (
	"database/sql"
	"evie/values"
	"reflect"
	"strings"
)

type PostgreSQLConnection struct {
	Value *sql.DB
}

func (a PostgreSQLConnection) GetNumber() float64 {
	return 1
}
func (p PostgreSQLConnection) GetType() values.ValueType {
	return values.CustomType
}
func (p PostgreSQLConnection) GetBool() bool {
	return true
}

func (p PostgreSQLConnection) GetString() string {
	return "Postgres Connection object"
}

func (db PostgreSQLConnection) GetProp(v *values.RuntimeValue, name string) (values.RuntimeValue, error) {
	props := map[string]values.RuntimeValue{}

	props["ping"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		val := (*v).(*PostgreSQLConnection)

		err := val.Value.Ping()

		if err != nil {
			return values.BoolValue{Value: false}
		}

		return values.BoolValue{Value: true}
	}}

	props["close"] = values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		val := (*v).(*PostgreSQLConnection)
		val.Value.Close()

		return values.BoolValue{Value: true}
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
					return values.BoolValue{Value: true}
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
				return values.BoolValue{Value: true}
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
					dict.Value[colName] = values.BoolValue{Value: colVal.(bool)}
				}

			}

			// Add row
			arr.Value = append(arr.Value, &dict)
		}

		return &arr

	}}

	return props[name], nil
}
