package jsonLib

import (
	environment "evie/env"
	"evie/utils"
	"evie/values"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func Load(env *environment.Environment) {
	namespace := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	namespace.Value["encode"] = values.NativeFunctionValue{Value: Encode}
	namespace.Value["decode"] = values.NativeFunctionValue{Value: Decode}

	env.Variables["json"] = namespace
}

func Encode(args []values.RuntimeValue) values.RuntimeValue {
	obj := args[0]

	output := ""

	switch val := obj.(type) {
	case values.StringValue:
		return values.StringValue{Value: "\"" + val.GetStr() + "\""}
	case values.NumberValue:
		return values.StringValue{Value: val.GetStr()}
	case *values.ArrayValue:
		output = "["
		for _, v := range val.Value {
			output = output + Encode([]values.RuntimeValue{v}).GetStr() + ", "
		}
		length := len(output)
		output = output[:length-2]
		output = output + "]"
	case values.DictionaryValue:
		output = "{"
		for k, v := range val.Value {
			output = output + "\"" + k + "\"" + ": " + Encode([]values.RuntimeValue{v}).GetStr() + ", "
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

	iter := utils.RuneIterator{Items: []rune(raw)}

	return DecodeValue(&iter)

}

func DecodeValue(iter *utils.RuneIterator) values.RuntimeValue {

	fmt.Println("VALUE TO DECODE VALUE: " + string(iter.Get()))

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

func DecodeSimpleValue(iter *utils.RuneIterator) values.RuntimeValue {

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

		iter.Eat()
	}
	return values.StringValue{Value: word}
}

func DecodeArray(iter *utils.RuneIterator) values.RuntimeValue {
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

func DecodeDictionary(iter *utils.RuneIterator) values.RuntimeValue {
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

	return output

}

func IsAlpha(char rune) bool {
	return unicode.IsLetter(char)
}

func isNumber(char rune) bool {
	return unicode.IsDigit(char)
}
