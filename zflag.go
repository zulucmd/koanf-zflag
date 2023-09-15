// Package kzflag implements a koanf.Provider that reads commandline
// parameters as conf maps using zulucmd/zflag, a POSIX compliant
// alternative to Go's stdlib flag package.
package kzflag

import (
	"errors"

	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/v2"

	"github.com/zulucmd/zflag/v2"
)

// KZFlag implements a zflag command line provider.
type KZFlag struct {
	delim   string
	flagset *zflag.FlagSet
	ko      *koanf.Koanf
	cb      func(key string, value string) (string, interface{})
	flagCB  func(f *zflag.Flag) (string, interface{})
}

// Provider returns a commandline flags provider that returns
// a nested map[string]interface{} of environment variable where the
// nesting hierarchy of keys are defined by delim. For instance, the
// delim "." will convert the key `parent.child.key: 1`
// to `{parent: {child: {key: 1}}}`.
//
// Though optional, it is recommended to pass in a koanf.Koanf instance.
// You can add additional options using the With* functions calls.
func Provider(f *zflag.FlagSet, delim string, options ...Option) *KZFlag {
	k := &KZFlag{
		flagset: f,
		delim:   delim,
	}
	applyOptions(k, options...)
	return k
}

// Read reads the flag variables and returns a nested conf map.
func (p *KZFlag) Read() (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	p.flagset.VisitAll(func(f *zflag.Flag) {
		var (
			key = f.Name
			val interface{}
		)

		// If there is a (key, value) callback set, pass the key and string
		// value from the flagset to it and use the results.
		if p.cb != nil {
			key, val = p.cb(key, f.Value.String())
		} else if p.flagCB != nil {
			// If there is a zflag.Flag callback set, pass the flag as-is
			// to it and use the results from the callback.
			key, val = p.flagCB(f)
		} else {
			// There are no callbacks set. Use the in-built flag value parser.
			val = FlagVal(f)
		}

		if key == "" {
			return
		}

		// If the default value of the flag was never changed by the user,
		// it should not override the value in the conf map (if it exists in the first place).
		if !f.Changed && (p.ko == nil || p.ko.Exists(key)) {
			return
		}

		// No callback. Use the key and value as-is.
		mp[key] = val
	})

	return maps.Unflatten(mp, p.delim), nil
}

// ReadBytes is not supported by the zflag provider.
func (p *KZFlag) ReadBytes() ([]byte, error) {
	return nil, errors.New("zflag provider does not support this method")
}

// Watch is not supported.
func (p *KZFlag) Watch(cb func(event interface{}, err error)) error {
	return errors.New("koanf-zflag provider does not support this method")
}

// FlagVal examines a zflag.Flag and returns a typed value as an interface{}
// from the types that zflag supports. If it is of a type that isn't known
// for any reason, the value is returned as a string.
func FlagVal(f *zflag.Flag) interface{} {
	var (
		val interface{}
	)

	switch v := f.Value.(type) {
	case zflag.Getter:
		val = v.Get()
	case zflag.Value:
		val = v.String()
	}

	return val
}
