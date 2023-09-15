// Package kzflag implements a koanf.Provider that reads commandline
// parameters as conf maps using zulucmd/zflag, a POSIX compliant
// alternative to Go's stdlib flag package.
package kzflag

import (
	"github.com/knadh/koanf/v2"

	"github.com/zulucmd/zflag/v2"
)

type Option func(f *KZFlag)

func applyOptions(f *KZFlag, options ...Option) {
	for _, option := range options {
		option(f)
	}
}

// WithCallback adds a callback that takes a (key, value) with the variable
// name and value and allows their modification.
// This is useful for cases where complex types like slices separated by
// custom separators.
// Returning "" for the key causes the particular flag to be disregarded.
func WithCallback(cb func(key string, value string) (string, interface{})) Option {
	return func(f *KZFlag) {
		f.cb = cb
	}
}

// WithFlagCallback takes a callback with *zflag.Flag and applies the callback
// to all items in the flagset in KZFlag. It does not parse
// *zflag.Flag values and expects the callback to process the keys and values
// from *zflag.Flag. FlagVal() can be used in the callback to avoid
// repeating the type-switch block for parsing values.
// Returning "" for the key causes the particular flag to be disregarded.
//
// Example:
//
//	p := koanf_zflag.Provider(flagset, ".", ko, WithFlagCallback(func(f *zflag.Flag) (string, interface{})) {
//	    // Transform the key in whatever manner.
//	    key := f.Name
//
//	    // Use FlagVal() and then transform the value, or don't use it at all
//	    // and add custom logic to parse the value.
//	    val := koanf_zflag.FlagVal(flagset, f)
//
//	    return key, val
//	 }))
func WithFlagCallback(cb func(f *zflag.Flag) (string, interface{})) Option {
	return func(f *KZFlag) {
		f.flagCB = cb
	}
}

// WithKoanf takes an optional (but recommended) Koanf instance to see if
// the flags defined have been set from other providers, for instance,
// a config file. If they are not, then the default values of the flags
// are merged. If they do exist, the flag values are not merged but only
// the values that have been explicitly set in the command line are merged.
func WithKoanf(ko *koanf.Koanf) Option {
	return func(f *KZFlag) {
		f.ko = ko
	}
}
