package main

import (
	"github.com/Elbujito/2112/src/app-service/internal/cmd"
)

// VERSION application version
var VERSION string = "0.0.1"

func main() {

	cmd.Version = VERSION
	cmd.Execute()
}
