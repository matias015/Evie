package lib

import (
	environment "evie/env"
	fsLib "evie/lib/fs"
)

func GetLibMap() map[string]func(*environment.Environment) {

	m := make(map[string]func(*environment.Environment), 1)

	m["fs"] = fsLib.Load
	// m["os"] = osLib.Load
	// m["crypt"] = cryptLib.Load
	// m["http"] = htppLib.Load
	// m["json"] = jsonLib.Load
	// m["postgres"] = postgresLib.Load
	return m
}
