package evruntime

import (
	"fmt"
)

type CallStackItem struct {
	ModuleName string
	Line       int
}

func (cs *CallStackItem) String() string {
	return "Module: " + cs.ModuleName + " -> Line: " + fmt.Sprint(cs.Line)
}

type CallStack struct {
	Items []CallStackItem
}

func (cs *CallStack) Add(line int, moduleName string) {
	// if len(cs.Items) > 1000 {
	// 	fmt.Println("MAX CALL STACK EXCEDED")
	// 	os.Exit(1)
	// }
	cs.Items = append(cs.Items, CallStackItem{Line: line, ModuleName: moduleName})
}

func (cs *CallStack) Remove() {
	cs.Items = cs.Items[0 : len(cs.Items)-1]
}
