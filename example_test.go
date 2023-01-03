// Package kzflag implements a koanf.Provider that reads commandline
// parameters as conf maps using zulucmd/zflag, a POSIX compliant
// alternative to Go's stdlib flag package.
package kzflag_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	kzflag "github.com/zulucmd/koanf-zflag"
	"github.com/zulucmd/zflag"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func TestExample(t *testing.T) {
	t.Parallel()

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
	_ = f.Parse([]string{"--type", "hello world"})

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
	if err := k.Load(kzflag.Provider(f, ".", kzflag.WithKoanf(k)), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	fmt.Println("time is = ", k.String("time"))
	fmt.Println("type is = ", k.String("type"))

	// Output:
	// time is =  2019-01-01
	// type is =  hello world
}
