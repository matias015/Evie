package jsonLib

import (
	"evie/common"
	environment "evie/env"
	"evie/values"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

/*

JSON LIB TODO:

 -> Detect format error at parsin raw json to objects
 -> Allow parse objects
 -> Allow scape characters and specials, like \n

*/

func Load(env *environment.Environment) {
	namespace := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	namespace.Value["encode"] = values.NativeFunctionValue{Value: Encode}
	namespace.Value["decode"] = values.NativeFunctionValue{Value: Decode}

	env.DeclareVar("json", namespace)
}

func Encode(args []values.RuntimeValue) values.RuntimeValue {
	obj := args[0]

	output := ""

	switch obj.GetType() {
	case values.StringType:
		return values.StringValue{Value: "\"" + obj.GetString() + "\""}
	case values.NumberType:
		return values.StringValue{Value: obj.GetString()}
	case values.ArrayType:
		output = "["
		for _, v := range obj.(*values.ArrayValue).Value {
			output = output + Encode([]values.RuntimeValue{v}).GetString() + ", "
		}
		length := len(output)
		output = output[:length-2]
		output = output + "]"
	case values.DictionaryType:
		output = "{"
		for k, v := range obj.(*values.DictionaryValue).Value {
			output = output + "\"" + k + "\"" + ": " + Encode([]values.RuntimeValue{v}).GetString() + ", "
		}
		length := len(output)
		output = output[:length-2]
		output = output + "}"
	}

	return values.StringValue{Value: output}
}

func Decode(args []values.RuntimeValue) values.RuntimeValue {

	arg := args[0]

	raw := arg.(values.StringValue).Value

	iter := common.RuneIterator{Items: []rune(raw)}

	return DecodeValue(&iter)

}

func DecodeValue(iter *common.RuneIterator) values.RuntimeValue {

	for {
		if iter.Get() != ' ' {
			break
		}

		if iter.IsOutOfBounds() {
			return nil
		}

		iter.Eat()
	}

	if iter.Get() == '{' {
		return DecodeDictionary(iter)
	} else if iter.Get() == '[' {
		return DecodeArray(iter)
	} else {
		return DecodeSimpleValue(iter)
	}
}

func DecodeSimpleValue(iter *common.RuneIterator) values.RuntimeValue {

	word := ""

	for iter.HasNext() && !iter.IsOutOfBounds() && iter.Get() != '}' {

		if isNumber(iter.Get()) {
			word = word + string(iter.Get())
			iter.Eat()
			for isNumber(iter.Get()) || iter.Get() == '.' {
				word = word + string(iter.Get())
				iter.Eat()
			}

			parsed, err := strconv.ParseFloat(word, 64)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			return values.NumberValue{Value: parsed}
		}

		if iter.Get() == '"' {

			iter.Eat()

			for iter.HasNext() && iter.Get() != '"' {
				word += string(iter.Eat())
			}

			iter.Eat()

			return values.StringValue{Value: word}

		}

		iter.Eat()
	}
	return values.StringValue{Value: word}
}

func DecodeArray(iter *common.RuneIterator) values.RuntimeValue {
	arr := values.ArrayValue{Value: make([]values.RuntimeValue, 0)}

	iter.Eat()

	for iter.HasNext() && !iter.IsOutOfBounds() && iter.Get() != ']' {

		if iter.Get() == ',' {
			iter.Eat()
			continue
		}

		arr.Value = append(arr.Value, DecodeValue(iter))

	}

	return &arr
}

func DecodeDictionary(iter *common.RuneIterator) values.RuntimeValue {
	output := values.DictionaryValue{Value: make(map[string]values.RuntimeValue, 0)}

	key := ""

	iter.Eat() // this is the rbrace

	for iter.HasNext() && !iter.IsOutOfBounds() && iter.Get() != '}' {

		if iter.Get() == ' ' {
			iter.Eat()
			continue
		}

		if iter.Get() == ':' {
			iter.Eat()
			output.Value[key] = DecodeValue(iter)
			key = ""
			continue
		}

		if IsAlpha(iter.Get()) {
			key = key + string(iter.Get())
			iter.Eat()
			for IsAlpha(iter.Get()) || isNumber(iter.Get()) || iter.Get() == '_' {
				key = key + string(iter.Get())
				iter.Eat()
			}
			continue
		}
		iter.Eat()
	}

	return &output

}

func IsAlpha(char rune) bool {
	return unicode.IsLetter(char)
}

func isNumber(char rune) bool {
	return unicode.IsDigit(char)
}
