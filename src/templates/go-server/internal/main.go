package main

import "github.com/Elbujito/2112/src/templates/go-server/internal/cmd"

var VERSION string = "0.0.1"

func main() {
	cmd.Version = VERSION
	cmd.Execute()
}
