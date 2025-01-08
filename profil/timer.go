package profil

import (
	"fmt"
	"sync"
	"time"
)

type Timer struct {
	Values map[string]int64
}

var instance *Timer
var once sync.Once

func NewTimer() *Timer {
	timer := Timer{Values: make(map[string]int64)}

	timer.Values["test"] = 0
	timer.Values["env_lookup"] = 0
	timer.Values["env_declare"] = 0
	timer.Values["env_set"] = 0
	timer.Values["get_prop"] = 0
	timer.Values["exec_fn"] = 0
	timer.Values["creating_scope"] = 0
	timer.Values["destroing_scope"] = 0
	timer.Values["reading_file"] = 0
	timer.Values["native_fn"] = 0
	timer.Values["evaluate_args"] = 0
	timer.Values["eval_calee"] = 0
	timer.Values["bin_exp"] = 0
	timer.Values["if_stmt"] = 0
	return &timer

}

func (t Timer) Display() {
	var total int64 = 0
	for k, v := range t.Values {
		if v == 0 {
			continue
		}
		total += v
		fmt.Println(k, v/1000, "ms")
	}

	fmt.Println("total", total/1000, "ms")
}

func (t *Timer) Add(name string, initial time.Time) {
	t.Values[name] += time.Since(initial).Microseconds()
}

func (t Timer) Init() time.Time {
	return time.Now()
}

func ObtenerInstancia() *Timer {
	once.Do(func() {
		instance = NewTimer()
	})
	return instance
}
