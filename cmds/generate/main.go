package main

import (
	"github.com/samxiaogo/trace_go"
	"github.com/samxiaogo/trace_go/cmds/generate/cmd"
)

func main() {
	defer trace_go.Print()()
	cmd.Execute()
}
