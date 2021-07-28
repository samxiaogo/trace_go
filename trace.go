package trace_go

import (
	"fmt"
	"github.com/modern-go/gls"
	"runtime"
	"sync"
)

type funcDirection string

const (
	funcStart funcDirection = "->"
	funcEnd   funcDirection = "<-"
)

var (
	mu       sync.Mutex
	funcStep = make(map[int64]int)
)

// Print do print the trace method
func Print(args ...interface{}) func() {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic(ErrorNoCallFrame)
	}
	fn := runtime.FuncForPC(pc)
	methodName := fn.Name()
	gid := gls.GoID()

	mu.Lock()
	v := funcStep[gid]
	funcStep[gid] = v + 1
	mu.Unlock()
	trace(gid, methodName, file, line, funcStart, v+1, args...)

	return func() {
		mu.Lock()
		v := funcStep[gid]
		funcStep[gid] = v - 1
		mu.Unlock()
		trace(gid, methodName, file, line, funcEnd, v)
	}

}

func trace(gid int64, methodName, file string, line int, direction funcDirection, tabCount int, args ...interface{}) {
	var (
		tabString string
		argString string
	)
	for i := 0; i < tabCount; i++ {
		tabString += "\t"
	}

	for i, v := range args {
		argString += fmt.Sprintf("arg_%d:%#v ;", i+1, v)
	}
	fmt.Printf("t_g[%02d]:%s%s%s %s %s:%d\n", gid, tabString, direction, methodName, argString, file, line)
}
