package main

import (
	"github.com/dswcpp/lazygit/pkg/app"
)

// These values may be set by the build script via the LDFLAGS argument
var (
	commit      string
	date        string
	version     = "1.1.3"
	buildSource = "unknown"
)

func main() {
	ldFlagsBuildInfo := &app.BuildInfo{
		Commit:      commit,
		Date:        date,
		Version:     version,
		BuildSource: buildSource,
	}

	app.Start(ldFlagsBuildInfo, nil)
}
