package mj

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
)

var logger = log.New(os.Stderr, "", 0)

func usage() {
	_ = showVersionInfo()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  %s [options...] <key=value...>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintln(os.Stderr, `Key-value pairs:
  mj foo=bar			{"foo":"bar"}
  mj foo.bar=baz		{"foo":{"bar":"baz"}}
  mj foo=bar baz=quux		{"baz":"quux","foo":"bar"}
  mj a=b a=c			# Error: Key path was already assigned
  mj foo:bar=baz		{"foo:bar":"baz"}
  mj -s=: foo:bar=baz		{"foo":"bar=baz"}
  mj -p=: foo:bar=baz		{"foo":{"bar":"baz"}}
  mj -- -really=why		{"-really":"why"}
  mj foo[]=abc foo[]=def	{"foo":["abc","def"]}
  
  echo hello-world > bar.txt
  mj foo=@bar.txt		{"foo":"@bar.txt"}
  mj -r=@ foo=@bar.txt		{"foo":"hello-world\n"}
  mj -r=@ foo=@<(date)		{"foo":"Mon Apr 25 00:00:34 PDT 2022\n"}

  mj foo=bar | mj baz=quux	{"baz":"quux"}
  mj foo=bar | mj -m=- baz=quux	{"baz":"quux","foo":"bar"}
  `)
}

func Run() int {
	var kvSeparator, pathSeparator string
	var readFilePrefix, readFrom string
	var showVersion bool

	flag.Usage = usage
	flag.StringVar(&readFrom, "m", "", "Merge input from file; use '-' for STDIN (default empty)")
	flag.StringVar(&kvSeparator, "s", "=", "Separator between key and value")
	flag.StringVar(&pathSeparator, "p", ".", "Separator between key-path components")
	flag.StringVar(&readFilePrefix, "r", "", "Prefix (for values) that indicate reading from a local file; reading value from a file is disabled when this flag is empty (default empty)")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		return showVersionInfo()
	}

	in := Struct{}
	if readFrom != "" {
		if readFrom == "-" {
			bs, err := io.ReadAll(os.Stdin)
			if err != nil {
				logger.Printf("%s: encountered error reading from STDIN: %v", os.Args[0], err)
				return 1
			}
			if err := json.Unmarshal(bs, &in); err != nil {
				logger.Printf("%s: encountered error unmarshaling payload from STDIN: %v", os.Args[0], err)
				return 1
			}
		} else {
			bs, err := os.ReadFile(readFrom)
			if err != nil {
				logger.Printf("%s: encountered error reading from %s: %v", os.Args[0], readFrom, err)
				return 1
			}
			if err := json.Unmarshal(bs, &in); err != nil {
				logger.Printf("%s: encountered error unmarshaling payload from %s: %v", os.Args[0], readFrom, err)
				return 1
			}
		}
	}

	p := &Processor{
		Input:             in,
		KeyValueSeparator: kvSeparator,
		KeyPathSeparator:  pathSeparator,
		ReadFilePrefix:    readFilePrefix,
	}

	for i, arg := range flag.Args() {
		err := p.Process(arg)
		if err != nil {
			logger.Printf("%s: encountered error while processing argument #%d: %q\n\t%v", os.Args[0], i, arg, err)
			return 1
		}
	}

	fmt.Println(p.Output())
	return 0
}

func showVersionInfo() int {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintf(os.Stderr, "%s [error: no build information]\n", os.Args[0])
		return 1
	}

	version := bi.Main.Version
	dirty := false
	if version == "(devel)" {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				version = s.Value
			}
			if s.Key == "vcs.modified" && s.Value == "true" {
				dirty = true
			}
		}
	}
	if dirty {
		version = version + "-dirty"
	}

	fmt.Fprintf(os.Stderr, "%s version %s\n", bi.Path, version)
	return 0
}
