package main

import "flag"
import "fmt"
import "log"
import "os"
import "strings"

var logger = log.New(os.Stderr, "", 0)

var kvSeparator string
var showVersion bool

func usage() {
	fmt.Fprintf(os.Stderr, "%s v%s built %s\n\n", os.Args[0], versionString(), buildString())
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [options...] <key=value...>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	flag.Usage = usage
	flag.StringVar(&kvSeparator, "s", "=", "Separator between key and value")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		logger.Printf("mj v%s (built %s)", versionString(), buildString())
		os.Exit(0)
	}

	input := Struct{}

	for _, arg := range flag.Args() {
		substrings := strings.SplitN(arg, kvSeparator, 2)
		if len(substrings) != 2 {
			logger.Fatalf("Missing separator (%q) in %q\n", kvSeparator, arg)
		}

		err := input.Set(substrings[0], substrings[1])
		if err != nil {
			logger.Fatalf("%v\n", err)
		}
	}

	fmt.Println(input.String())
}
