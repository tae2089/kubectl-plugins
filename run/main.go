package main

import (
	"github.com/spf13/cobra"
)

var rooCmd cobra.Command

func main() {
	if rooCmd.Use == "" {
		panic("root command not initialized")
	}
	rooCmd.Execute()
}
