package lib

import (
	environment "evie/env"
	cryptLib "evie/lib/crypt"
	fsLib "evie/lib/fs"
	osLib "evie/lib/os"
)

func GetLibMap() map[string]func(*environment.Environment) {
	var m map[string]func(*environment.Environment)

	m = make(map[string]func(*environment.Environment))

	m["fs"] = fsLib.Load
	m["os"] = osLib.Load
	m["crypt"] = cryptLib.Load
	return m
}
