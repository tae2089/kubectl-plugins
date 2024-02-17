//go:build kdesc

package main

import (
	"github.com/tae2089/kubectl-custom-cli/cmd/kdesc"
)

func init() {
	rooCmd = *kdesc.CreateRootCmd()
}
