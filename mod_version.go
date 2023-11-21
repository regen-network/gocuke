package gocuke

import (
	"runtime/debug"
)

// Version of package - based on Semantic Versioning 2.0.0 http://semver.org/
var Version = "v0.0.0-dev"

func init() {
	if info, available := debug.ReadBuildInfo(); available {
		if Version == "v0.0.0-dev" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
