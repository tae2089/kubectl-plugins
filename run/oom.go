//go:build restarted

package main

import "github.com/tae2089/kubectl-custom-cli/cmd/restart"

func init() {
	rooCmd = *restart.CreateRootCmd()
}
