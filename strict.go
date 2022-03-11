package gocuke

import "flag"

var flagStrict = flag.Bool("gocuke.strict", true, "will turn pending or undefined tests into a test failure (default true)")
