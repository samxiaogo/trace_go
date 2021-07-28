package cmd

import (
	"fmt"
	"github.com/samxiaogo/trace_go/pkg/ast"
	"github.com/samxiaogo/trace_go/pkg/path"
	"github.com/spf13/cobra"
	"os"
)

var (
	generatePath = "./"
)

var rootCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate trace to the project",
	Long:  `generate trace to the function and then run the project can see the trace data`,
	Run: func(cmd *cobra.Command, args []string) {
		path.Walk(generatePath, func(filePath string) error {
			err := path.ReWriteFile(filePath, ast.ParseAndAdd(filePath))
			if err != nil {
				panic(err)
			}
			return nil
		})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
