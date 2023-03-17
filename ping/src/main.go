package main

import (
	"github.com/Kong/go-pdk/server"
	"github.com/jabardigitalservice/kong-plugins/ping/src/internal"
)

var (
	Version  = "1.1"
	Priority = 3000
)

func main() {
	_ = server.StartServer(internal.New, Version, Priority)
}
