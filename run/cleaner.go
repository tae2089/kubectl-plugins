//go:build cleaner

package main

import "github.com/tae2089/kubectl-custom-cli/cmd/cleaner"

func init() {
	rooCmd = *cleaner.CreateRootCmd()
}
