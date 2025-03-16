package main

import (
	"github.com/johankristianss/evrium/internal/cli"
	"github.com/johankristianss/evrium/pkg/build"
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
