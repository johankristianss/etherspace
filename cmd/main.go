package main

import (
	"github.com/johankristianss/etherspace/internal/cli"
	"github.com/johankristianss/etherspace/pkg/build"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

func main() {
	build.BuildVersion = BuildVersion
	build.BuildTime = BuildTime
	cli.Execute()
}
