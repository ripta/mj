package main

import "flag"
import "fmt"
import "log"
import "os"

var logger = log.New(os.Stderr, "", 0)

func usage() {
	fmt.Fprintf(os.Stderr, "%s v%s built %s\n\n", os.Args[0], BuildVersion, BuildDate)
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [options...] <key=value...>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintln(os.Stderr, `Key-value pairs:
  mj foo=bar		{"foo":"bar"}
  mj foo.bar=baz	{"foo":{"bar":"baz"}}
  mj foo=bar baz=quux	{"baz":"quux","foo":"bar"}
			# Keys are automatically sorted
  mj a=b a=c		# Error: Key path was already assigned
  mj foo:bar=baz	{"foo:bar":"baz"}
  mj -s=: foo:bar=baz	{"foo":"bar=baz"}
  mj -p=: foo:bar=baz	{"foo":{"bar":"baz"}}
  mj -- -really=why	{"-really":"why"}
`)
}

func main() {
	var kvSeparator, pathSeparator string
	var showVersion bool

	flag.Usage = usage
	flag.StringVar(&kvSeparator, "s", "=", "Separator between key and value")
	flag.StringVar(&pathSeparator, "p", ".", "Separator between key-path components")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		logger.Printf("mj v%s (built %s)", BuildVersion, BuildDate)
		os.Exit(0)
	}

	p := &Processor{
		input:         Struct{},
		kvSeparator:   kvSeparator,
		pathSeparator: pathSeparator,
	}

	for i, arg := range flag.Args() {
		err := p.Process(arg)
		if err != nil {
			logger.Fatalf("%s: while processing argument #%d:\n\t%s\nError: %s\n", os.Args[0], i, arg, err)
		}
	}

	fmt.Println(p.Output())
}
