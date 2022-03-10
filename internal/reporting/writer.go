package reporting

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var flagOut = flag.String("gocuke.out", "", "an output file for gocuke messages")

func getWriter() io.WriteCloser {
	if flagOut != nil && *flagOut != "" {
		w, err := os.OpenFile(*flagOut, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(fmt.Errorf("error opening file %s: %v", *flagOut, err))
		}
		return w
	}
	return nil
}
