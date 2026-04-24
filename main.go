package main

import "github.com/iamllcoolray/stone-cli/cmd"

var Version = "dev"

func main() {
	cmd.Execute(Version)
}
