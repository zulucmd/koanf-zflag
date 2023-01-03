# Koanf zflag Provider

[![GoDoc](https://godoc.org/github.com/zulucmd/koanf-zflag?status.svg)](https://godoc.org/github.com/zulucmd/koanf-zflag)
[![Go Report Card](https://goreportcard.com/badge/github.com/zulucmd/koanf-zflag)](https://goreportcard.com/report/github.com/zulucmd/koanf-zflag)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/zulucmd/koanf-zflag?sort=semver)](https://github.com/zulucmd/koanf-zflag/releases)
[![Build Status](https://github.com/zulucmd/koanf-zflag/actions/workflows/validate.yml/badge.svg)](https://github.com/zulucmd/koanf-zflag/actions/workflows/validate.yml)

<!-- toc -->

- [Introduction](#introduction)
- [Installation](#installation)
- [Documentation](#documentation)
  - [Reading from command line](#reading-from-command-line)

<!-- /toc -->

## Introduction

koanf-zflag is a Koanf provider to retrieve configuration from [zulucmd/zflag](https://github.com/zulucmd/zflag).

## Installation

koanf-zflag is available using the standard `go get` command.

Install by running:

```bash
go get github.com/zulucmd/koanf-zflag
```

## Documentation

### Reading from command line

The following example shows the use of kzflag.Provider, a wrapper over the spf13/pflag library, an advanced commandline
lib. For Go's built in flag package, use basicflag.Provider.

```go
package main

import (
  "fmt"
  "log"
  "os"

  kzflag "github.com/zulucmd/koanf-zflag"
  "github.com/zulucmd/zflag"
  "github.com/knadh/koanf"
  "github.com/knadh/koanf/parsers/json"
  "github.com/knadh/koanf/providers/file"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func main() {
  // Use the POSIX compliant zflag lib instead of Go's flag lib.
  f := zflag.NewFlagSet("config", zflag.ContinueOnError)
  f.Usage = func() {
    fmt.Println(f.FlagUsages())
    os.Exit(0)
  }
  // Path to one or more config files to load into koanf along with some config params.
  f.StringSlice("conf", []string{"example/conf.json"}, "path to one or more .toml config files")
  f.String("time", "2020-01-01", "a time string")
  f.String("type", "xxx", "type of the app")
  f.Parse([]string{"--type", "hello world"})

  // Load the config files provided in the commandline.
  cFiles, _ := f.GetStringSlice("conf")
  for _, c := range cFiles {
    if err := k.Load(file.Provider(c), json.Parser()); err != nil {
      log.Fatalf("error loading file: %v", err)
    }
  }

  // "time" and "type" may have been loaded from the config file, but
  // they can still be overridden with the values from the command line.
  // The bundled kzflag.Provider takes a flagset from the gorwarden/zflag lib.
  // Passing the Koanf instance to kzflag helps it deal with default command
  // line flag values that are not present in conf maps from previously loaded
  // providers.
  if err := k.Load(kzflag.Provider(f, ".", k), nil); err != nil {
    log.Fatalf("error loading config: %v", err)
  }

  fmt.Println("time is = ", k.String("time"))
  fmt.Println("type is = ", k.String("type"))

  // Output:
  // time is =  2019-01-01
  // type is =  hello world
}
```
