package main

import "fmt"
import "log"
import "os"
import "strings"

var logger = log.New(os.Stderr, "", 0)
var sep string = "="

func main() {
  if len(os.Args) == 2 && os.Args[1] == "-version" {
    logger.Printf("MJ v%s (built %s)", versionString(), buildString())
    os.Exit(0)
  }

  input := Struct{}

  for _, arg := range os.Args[1:] {
    substrings := strings.SplitN(arg, sep, 2)
    err := input.Set(substrings[0], substrings[1])
    if err != nil {
      logger.Printf("%v\n", err)
      os.Exit(1)
    }
  }

  fmt.Println(input.String());
}
