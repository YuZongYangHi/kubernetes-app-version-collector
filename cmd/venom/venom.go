package main

import (
	"github.com/YuZongYangHi/kubernetes-app-version-collector/cmd/venom/app"
	"k8s.io/component-base/cli"
	"os"
)

func main() {
	command := app.NewVenomCommand()
	code := cli.Run(command)
	os.Exit(code)
}
