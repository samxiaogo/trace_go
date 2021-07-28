package trace_go

import (
	"sync"
	"testing"
)

func TestTraceData(t *testing.T) {
	Demo()
}

func Demo() {
	defer Print()()
	var (
		n  = 2
		wg sync.WaitGroup
	)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			Demo1()
		}()
	}
	wg.Wait()
}

func Demo1() {
	defer Print()()
	demo2()
}

func demo2() {
	defer Print()()
	demo3(1, 2)
}

func demo3(c, d int) {
	defer Print(c, d)()
	demo4(&DemoStruct{
		Name: "213",
	}, []byte{1, 2, 3})
}

type DemoStruct struct {
	Name string
}

func demo4(d *DemoStruct, v []byte) {
	defer Print(d, v)()
}
